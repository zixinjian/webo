package ctrl

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"time"
	"wb/cc"
	"wb/cs"
	"wb/ii"
	"wb/lg"
	"wb/om"
	"wb/usr"
	"wb/ut"
)

type BaseController struct {
	beego.Controller
	CtxVersion string
	CtxUser    usr.Usr
}

type NestPreparer interface {
	NestPrepare()
}

func (this *BaseController) Prepare() {
	user, ok := this.GetSession(cc.SessionUser).(usr.Usr)
	if !ok {
		this.CtxUser = usr.NoUser
	} else {
		this.CtxUser = user
		this.Data[cc.CtxUserName] = this.CtxUser.GetName()
		this.Data[cc.CtxUserUserName] = this.CtxUser.GetUserName()
		this.Data[cc.CtxUserSn] = this.CtxUser.GetSn()
		this.Data[cc.CtxRole] = this.CtxUser.GetRole()
		for k, v := range this.CtxUser.GetExtMap() {
			this.Data["Ctx"+k] = v
		}
	}
	this.CtxVersion = beego.AppConfig.String(cc.AppConfVersion)
	this.Data[cc.CtxVersion] = this.CtxVersion

	// 调用子类的NestPrepare
	if app, ok := this.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}
func (c *BaseController) SendJson(jsonObject interface{}) {
	c.Data["json"] = jsonObject
	c.ServeJSON()
}
func (this *BaseController) GetSessionUser() usr.Usr {
	if user, ok := this.GetSession(cc.SessionUser).(usr.Usr); ok {
		return user
	}
	lg.Critical(usr.CodeErrorNotLogin)
	return usr.NoUser
}
func (c *BaseController) GetCurRole() string {
	return c.GetSessionUser().GetRole()
}
func (c *BaseController) GetCurUserSn() string {
	return c.GetSessionUser().GetSn()
}
func (c *BaseController) GetCurTime() string {
	return ut.GetCurDbTime()
}
func (c *BaseController) GetMonthTime() string {
	return ut.GetMonthTime()
}
func (c *BaseController) GetSecondMonthTime() string {
	return ut.GetSecondMonthTime()
}
func (c *BaseController) GetYesterdayDate() string {
	return ut.GetYesterdayDbDate()
}
func (c *BaseController) GetTodayDbDate() string {
	return ut.GetTodayDbDate()
}
func (c *BaseController) GetTomorrowDbDate() string {
	return ut.GetTomorrowDbDate()
}
func (this *BaseController) GetRequestMapFromJson() (map[string]interface{}, error) {
	requestBody := this.Ctx.Input.RequestBody
	var requestMap map[string]interface{}
	if err := json.Unmarshal(requestBody, &requestMap); err == nil {
		lg.Debug("BaseController.GetRequestMapFromJson requestMap: ", requestMap)
		return requestMap, err
	} else {
		lg.Error("BaseController.GetRequestMapFromJson requestMap: ", requestBody)
		return make(map[string]interface{}), err
	}
}
func (this *BaseController) getRequestRowsFromJson() ([]map[string]interface{}, error) {
	requestBody := this.Ctx.Input.RequestBody
	var rows []map[string]interface{}
	err := json.Unmarshal(requestBody, &rows)
	lg.Debug("BaseController.getRequestRowsFromJson rows: ", rows)
	return rows, err
}

// 替换特殊字符 curuser curdate curtime 返回值 string 替换后的字符，bool 是否替换
func (this *BaseController) replaceSpecialValues(value string) (string, bool) {
	switch value {
	case cc.CurUser:
		return this.GetCurUserSn(), true
	case cc.CurDate, cc.CurTime:
		return ut.GetCurDbTime(), true
	}
	if strings.HasPrefix(value, "Ctx_User_") {
		key := strings.Replace(value, "Ctx_User_", "", -1)
		return this.CtxUser.GetString(key), true
	}
	return value, false
}
func (this *BaseController) getFilteredQueryParamFromJsonMap(oItemInfo ii.ItemInfo, requestMap map[string]interface{}) map[string]interface{} {
	queryParams := make(om.Params, 0)
	fieldMap := oItemInfo.GetFieldMap()
	for k, v := range requestMap {
		if field, ok := fieldMap[k]; ok {
			if fv, fok := field.GetQueryValue(v); fok {
				switch field.Model {
				case ii.MEnum:
					queryParams["%"+k] = field.GetEnumKey(fv)
				case ii.MTime, ii.MCurTime, ii.MDate:
					if sv, ok := fv.(string); ok {
						if strings.Contains(sv, "~") {
							ts := strings.Split(sv, "~")
							if len(ts) < 2 {
								continue
							}
							if t0, err := ut.GetTimeFromStr(ts[0]); err == nil {
								queryParams[">"+k] = t0.Add(-1 * time.Second).Format(ut.DbTimeFormatter)
							} else {
								lg.Error("getFilteredQueryParamFromJsonMap get time0 error:", err.Error())
							}
							t1, _ := ut.GetDbTimeFromInput(ts[1])
							if t1 != "" {
								queryParams["<"+k] = t1
							}
						} else {
							queryParams["%"+k] = strings.Replace(sv, "-", "", -1)
						}
					}
				default:
					queryParams["%"+k] = fv
				}
			} else {
				lg.Error(fmt.Sprintf("Check param[%s]value %v error", k, v))
			}
		} else {
			lg.Error(fmt.Sprintf("Check param[%s]value %v error no such field", k, v))
		}
	}
	return queryParams
}

func (this *BaseController) getQueryParamFromJsonMap(requestMap map[string]interface{}) map[string]interface{} {
	queryParams := make(om.Params, len(requestMap))
	for k, v := range requestMap {
		queryParams["%"+k] = v
	}
	return queryParams
}

func (this *BaseController) getLimitParamFromJsonMap(requestMap map[string]interface{}) om.LimitParams {
	limitParams := make(map[string]int64, 0)
	if k, ok := requestMap["limit"]; ok {
		limitParams["limit"] = int64(k.(float64))
	}
	if k, ok := requestMap["offset"]; ok {
		limitParams["offset"] = int64(k.(float64))
	}
	return limitParams
}
func (this *BaseController) getOrderParamFromJsonMap(requestMap map[string]interface{}) om.OrderParams {
	orderByParams := make(om.OrderParams, 0)
	if sort, ok := requestMap["sort"]; ok {
		sortStr := strings.TrimSpace(sort.(string))
		if sortStr != "" {
			order := "asc"
			if o, ok := requestMap["order"]; ok {
				if strings.TrimSpace(o.(string)) == "desc" {
					order = "desc"
				}
			}
			orderByParams[sortStr] = order
		}
	}
	return orderByParams
}
func (this *BaseController) GetTrimString(key string) string {
	return strings.TrimSpace(this.GetString(key))
}
func (this *BaseController) SetDataWithoutDefaultStr(field, defaultValue, value string) {
	if defaultValue == "" {
		this.Data[field] = value
	} else {
		this.Data[field] = defaultValue
	}
}
func (this *BaseController) FillFormElement(elementMap map[string]string) {
	pre := "Form_"
	for k, v := range elementMap {
		this.Data[pre+k] = v
	}
}
func (this *BaseController) FillValue(oldValueMap cs.MObject) {
	pre := "V_"
	for k, v := range oldValueMap {
		this.Data[pre+k] = v
	}
}
func (this *BaseController) SendZeroTableResult(status string) {
	this.SendJson(&cs.TableResult{status, int64(0), make([]cs.MObject, 0)})
}
func (this *BaseController) SendJsonResult(status, desc string) {
	this.SendJson(&cs.JsonResult{status, desc})
}
func (this *BaseController) SendJsonData(ret, msg string, data interface{}) {
	this.SendJson(&cs.JsonData{Ret: ret, Msg: msg, Data: data})
}
