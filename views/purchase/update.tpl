<!DOCTYPE html>
<html>
<head lang="zh">
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../lib/app/css/app.min.css"/>
    <link rel="stylesheet" href="../../lib/jquery/datetimepicker/jquery.datetimepicker.css">
    <link rel="stylesheet" href="../../lib/uploadify/uploadify/uploadify.css" />
    <link rel="stylesheet" href="../../lib/jquery/jquery-ui/jquery-ui.min.css">
    <link rel="stylesheet" href="../../lib/webo/css/ui.css">
    <!-- Le HTML5 shim, for IE6-8 support of HTML5 elements -->
    <!--[if lt IE 9]>
    <script src="../../lib/html5shiv.min.js"></script>
    <![endif]-->
</head>
<body>
<div class="container-fluid">
    <form class="form-horizontal" id="item_form">
        {{str2html .Form_sn}}
        <div class="form-group">
            <label class="col-sm-2 control-label">类别</label>
            <div class="col-sm-8">
                <select class="input-block-level form-control" data-validate="{required: true, messages:{required:'请输入类别'}}" name="category" id="category" autocomplete="off" value="" disabled>
                    {{str2html .CategoryOptions}}
                </select>
            </div>
        </div>
        {{str2html .Form_productname}}
        {{str2html .Form_product}}
        {{str2html .Form_brand}}
        {{str2html .Form_model}}
        {{str2html .Form_power}}
        {{str2html .Form_num}}
        {{str2html .Form_placedate}}
        {{str2html .Form_requireddate}}
        {{str2html .Form_requireddepartment}}
        {{str2html .Form_unitprice}}
        {{str2html .Form_productprice}}
        {{str2html .Form_totalprice}}
        {{str2html .Form_freightprice}}
        {{if .NeedSupplier}}
        <div class="form-group">
            <label class="col-sm-2 control-label">供应商关键字</label>
            <div class="col-sm-8">
                <select class="input-block-level form-control" name="supplier" id="supplier" autocomplete="off" value="{{.supplier}}" >
                    {{str2html .SupplierOptions}}
                </select>
                <input type="hidden" name="suppliername" id="suppliername2" value="{{.suppliername}}">
            </div>
        </div>
        {{else}}
        <div class="form-group">
            <label class="col-sm-2 control-label">供应商关键字</label>
            <div class="col-sm-8">
            <input type="text" class="input-block-level form-control" data-validate="{required: true, messages:{required:'请输入正确的供应商!'}}" name="suppliername" id="suppliername" autocomplete="off" value="{{.suppliername}}">
            <input type="hidden" name="supplier" id="supplier2" value="{{.supplier}}">
            </div>
        </div>
        {{end}}
        {{str2html .Form_buyer}}
        {{str2html .Form_orderdate}}
        {{str2html .Form_predictdeliverydate}}
        {{str2html .Form_actualdeliverydate}}
        {{str2html .Form_arrivaldate}}
        {{str2html .Form_godowndate}}
        {{str2html .Form_paymentamount}}
        {{str2html .Form_paymentdate}}
        {{str2html .Form_mark}}
        {{str2html .Form_creater}}
    </form>
</div>
<div class="modal fade" id="top_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content" style="width: 800px;">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title" id="top_modal_title">付款</h4>
            </div>
            <div class="modal-body" id="top_modal_body">

            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" id="top_modal_btn_cancel" data-dismiss="modal">取消</button>
                <button id="top_modal_btn_ok" type="button" class="btn btn-primary">确定</button>
            </div>
        </div>
    </div>
</div>
<script src="../../lib/app/js/app.min.js"></script>
<script src="../../lib/jquery/jquery/jquery.form.js"></script>
<script src="../../lib/jquery/jquery/validate/jquery.metadata.js"></script>
<script src="../../lib/jquery/jquery/validate/jquery.validate.js"></script>
<script src="../../lib/uploadify/uploadify/jquery.uploadify.js"></script>
<script src="../../lib/jquery/datetimepicker/jquery.datetimepicker.js"></script>
<script src="../../lib/jquery/jquery-ui/jquery-ui.min.js"></script>
<script src="../../lib/webo/js/validateExtend.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script>
    function showResponse(resp) {
        if(resp.ret == "success"){
            top.hideTopModal()
            refreshContent()
        }else{
            showError("更新失败!")
        }
    }
    var refreshContent
    function onTopModalOk(options){
        if(options.refreshContent){
            refreshContent = options.refreshContent
        }
        if (! $("#item_form").valid()){
            return "not"
        }
        $("#item_form").ajaxSubmit({
            type: "post",
            url: "{{.Service}}",
            success: showResponse
        });
        return "not"
    }
    $(function(){
        $("#paymentamount").wrapAll('<div class="input-group"></div>')
        $("#paymentamount").after('<a class="btn btn-sm input-group-addon" id="calc">付款</a>')
//        $("#paymentamount").click(calPayment)
        $("#paymentdate").wrapAll('<div class="input-group"></div>')
        $("#paymentdate").after('<a class="btn btn-sm input-group-addon" id="calc">明细</a>')
    });
</script>
{{str2html .Onload}}
{{if .NeedSupplier}}
<script>
    $(function(){
        $("#supplier").on("change", function(){
            supplierName = $("#supplier").find("option:selected").text()
            $("#suppliername2").val(supplierName)
        })
    })
</script>
{{else}}
<script>
    $(function(){
        $("#suppliername").autocomplete({
            source: "/item/autocomplete/supplier",
            autoFocus:true,
            autoFill:true,
            focus: function( event, ui ) {
                return true;
            },
            minLength: 1,
            select: function( event, ui) {
                $( "#suppliername" ).val(ui.item.keyword);
                $( "#supplier2" ).val(ui.item.sn);
                return false;
            },
            change: function( event, ui ) {
                console.log("ui.item", ui.item)
                if(!ui.item){
                    $( "#supplier2" ).val("");
                }else{
                    $( "#supplier2" ).val(ui.item.sn);
                }
            }
        }).autocomplete( "instance" )._renderItem = function( ul, item ) {
            return $( "<li>" )
                    .append(item.keyword + "(" + item.name + ")")
                    .appendTo( ul );
        };
    })
</script>
{{end}}
</body>
</html>