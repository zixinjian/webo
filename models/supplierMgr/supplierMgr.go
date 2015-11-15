package supplierMgr

import (
	"strings"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/svc"
	"webo/models/t"
	"webo/models/wbo"
	"github.com/astaxie/beego"
)

func Get(sn string) (map[string]interface{}, bool) {
	params := t.Params{
		s.Sn: sn,
	}
	status, retMap := svc.Get(s.Supplier, params)
	return retMap, strings.EqualFold(stat.Success, status)
}

func GetSupplierListByProductSn(sn string)(string, []map[string]interface{}){
	query := `select supplier.* from supplier as supplier, supplier_product as rel WHERE supplier.sn = rel.supplier AND rel.product = ?`
	return wbo.DbQueryValues(query, sn)
}

func AddSupplierRels(product string, suppliers []string) string{
	beego.Debug("AddSuppliers suppliers:", suppliers)
	for _, supplier := range suppliers{
		params := t.Params{s.Supplier:supplier, s.Product:product}
		beego.Debug("AddSuppliers", params)
		svc.Add("supplier_product", params)
	}
	return stat.Success
}
func UpdateSupplierRels(product string, suppliers []string)string{
	beego.Debug("UpdateSupplierRels:", product, " ", suppliers)
	query := `SELECT * FROM supplier_product WHERE product = ?`
	st, vs := wbo.DbQueryValues(query, product)
	if st != stat.Success{
		return st
	}
	oldRelMaps := make(map[string]map[string]interface{})
	for _, v := range vs{
		supplier, _ := v[s.Supplier]
		beego.Debug("UpdateSupplierRels oldRels:", supplier)
		oldRelMaps[supplier.(string)] = v
	}
	for _, supplier := range suppliers{
		if _, ok:= oldRelMaps[supplier];ok{
			delete(oldRelMaps, supplier)
		}else{
			// 新增
			beego.Debug("UpdateSupplierRels", t.Params{s.Supplier:supplier, s.Product:product})
			svc.Add("supplier_product", t.Params{s.Supplier:supplier, s.Product:product})
		}
	}
	for _, oldRel := range oldRelMaps{
		if sn, ok := oldRel[s.Sn];ok{
			DelSupplierRel(sn.(string))
		}
	}
	return stat.Success
}

func DelSupplierRel(sn string){
	beego.Debug("DelSupplierRel sn:", sn)
	deleteSql := `DELETE FROM supplier_product WHERE sn=?`
	wbo.DbDelete(deleteSql, sn)
}