<!DOCTYPE html>
<html>
<head lang="zh">
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../lib/3rd/bootstrap/css/bootstrap.css">
    <link rel="stylesheet" href="../../lib/3rd/datetimepicker/jquery.datetimepicker.css">
    <link rel="stylesheet" href="../../lib/3rd/uploadify/uploadify.css" />
    <link rel="stylesheet" href="../../lib/3rd/jquery-ui/jquery-ui.min.css">
    <style>
        .ui-autocomplete-loading {
            background: white url("../../lib/webo/images/ui-anim_basic_16x16.gif") right center no-repeat;
        }
    </style>
</head>
<body>
<div class="container-fluid">
    <div class="alert" role="alert" style="display: none">添加成功！</div>
    <form class="form-horizontal" id="item_form">
        <input type="hidden" id="sn" name="sn" value="{{.sn}}">
        <div class="form-group">
        <label class="col-sm-3 control-label">类别</label>
        <div class="col-sm-6">
            <select class="input-block-level form-control" name="category" id="category" autocomplete="off" value="" >
                <option value="cate_engine">柴油机</option>
                <option value="cate_generator">电机</option>
                <option value="cate_waterbox">水箱</option>
                <option value="cate_epart">电器件</option>
                <option value="cate_parts">配件</option>
                <option value="cate_other">其他</option>
                <option value="cate_newpdt">新产品</option>
            </select>
        </div>
    </div>
        <input type="hidden" id="product" name="product" value="">
        <input type="hidden" id="supplier" name="supplier" value="">
        <div class="form-group">
            <label class="col-sm-3 control-label">商品名称</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" id="productname" name="productname" value=""
                       data-rule-required="true"/>
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">型号</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" name="model" id="model" autocomplete="off" value=""
                       data-rule-required="true"/>
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">品牌</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" name="brand" id="brand" autocomplete="off" value="" />
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">功率</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" name="power" id="power" autocomplete="off" value="" readonly
                       data-rule-required="true" data-rule-number="true" data-msg-number="请输入正确的功率!" />
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">参考价</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" name="productprice" id="productprice" autocomplete="off" value="" readonly='true'/>
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">采购人</label>
            <div class="col-sm-6">
                <select class="input-block-level form-control" name="buyer" id="buyer" autocomplete="off" value=""
                        data-rule-required="true">
                    {{range .Buyers}}
                        <option value="{{.Sn}}">{{.Name}}</option>
                    {{end}}
                </select>
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">数量</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" data-validate="{required: true, messages:{required:'请输入正确的数量!'}}" name="num" id="num" autocomplete="off" value="1" />
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">下单日期</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control datetimepicker" name="placedate" id="placedate" autocomplete="off" value="curtime"
                       data-rule-required="true"/>
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">需用日期</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control datetimepicker" name="requireddate" id="requireddate" autocomplete="off" value=""
                       data-rule-required="true" data-rule-date="true"/>
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">申请部门</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" name="requireddepartment" id="requireddepartment" autocomplete="off" value=""
                       data-rule-required="true"/>
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">备注</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" name="mark" id="mark" autocomplete="off" value="" />
            </div>
        </div>
    </form>
</div>

<script src="../../lib/3rd/jquery/jquery.js"></script>
<script src="../../lib/3rd/bootstrap/js/bootstrap.min.js"></script>
<script src="../../lib/3rd/jquery/jquery.form.js"></script>
<script src="../../lib/3rd/jquery/validate/jquery.validate.min.js"></script>
<script src="../../lib/3rd/uploadify/jquery.uploadify.js"></script>
<script src="../../lib/3rd/datetimepicker/jquery.datetimepicker.js"></script>
<script src="../../lib/3rd/jquery-ui/jquery-ui.min.js"></script>
<script src="../../lib/webo/js/validateExtend.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script src="../../lib/webo/util.js"></script>

<script>
    cateNoName = {
        cate_engine:"柴油机",
        cate_generator:"电机",
        cate_waterbox:"水箱"
    }
    cateNameValues = wbGetMapValue(cateNoName)
    var $productName = $('#productname')
    var $power = $("#power")
    var $category = $("#category")
    var $form = $("#item_form")
    function showResponse(resp) {
        if(resp.ret == "success"){
            top.hideTopModal()
            refreshContent()
        }else{
            showError("添加失败!")
        }
    }
    var refreshContent
    function onTopModalOk(options){
        if(options.refreshContent){
            refreshContent = options.refreshContent
        }
        if (! $form.valid()){
            return
        }
        $form.ajaxSubmit({
            type: "post",
            url: "{{.Service}}",
            success: showResponse
        });
    }
    function setProductValues(item){
        $("#product" ).val(item.sn);
        $("#productprice").val(item.price)
        $("#brand").val(item.brand)
        $("#supplier").val(item.supplier)
        $power.val(item.power)
    }
    function clearProductValues(){
        $("#product").val("");
        $("#productprice").val("")
        $("#supplier").val("")
        $("#brand").val("")
        $("#model").val("")
        $power.val("0")
    }
    $(function () {
        $("#power").wrapAll('<div class="input-group"></div>')
        $("#power").after('<span class="input-group-addon">KW</span>')

        var selectCate = $('#category').val()
        if(selectCate in cateNoName && $productName.val() == ""){
            $productName.val(cateNoName[selectCate]);
            $productName.attr("readonly", true)
            wbGetParentFromGroup("#productname").hide()
        }
        $('#category').change(function(){
            clearProductValues()
            var selectCate = $('#category').val()
            $("#model").autocomplete("option", "source", "/item/autocomplete/product?category=" + selectCate)
            if(selectCate in cateNoName){
                $productName.val(cateNoName[selectCate]);
                $productName.attr("readonly", true)
                $power.val("")
                wbGetParentFromGroup("#productname").hide()
                wbGetParentFromGroup("#power").show()

            }else{
                $productName.attr("readonly", false)
                $productName.val("")
                $power.val("0")
                wbGetParentFromGroup("#productname").show()
                wbGetParentFromGroup("#power").hide()
            }
        })
        $("#model").autocomplete({
            source: "/item/autocomplete/product?category=" + $("#category").val(),
            autoFocus:false,
            focus: function( event, ui ) {
                setProductValues(ui.item)
                return false;
            },
            minLength: 1,
            select: function( event, ui) {
                $("#model").val(ui.item.model)
                setProductValues(ui.item)
                return false;
            },
            change: function( event, ui ) {
                if(!ui.item){
                    $("#product").val("");
                    $("#productprice").val("")
                    $("#supplier").val("")
                    $("#brand").val("")
                    $power.val("")
                }
            }
        }).autocomplete( "instance" )._renderItem = function( ul, item ) {
            return $( "<li>" )
                    .append(item.keyword + "(" + item.name + ")")
                    .appendTo( ul );
        };

        $("#placedate").datetimepicker({timepicker:false,format:'Y.m.d',lang:'zh',value:new Date(), scrollMonth:false})
        $("#requireddate").datetimepicker({timepicker:false,format:'Y.m.d',lang:'zh', scrollMonth:false})
    });
</script>
</body>
</html>