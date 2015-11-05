
<!DOCTYPE html>
<html>
<head lang="zh">
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../lib/3rd/bootstrap/css/bootstrap.css">
    <link rel="stylesheet" href="../../lib/3rd/datetimepicker/jquery.datetimepicker.css">
    <link rel="stylesheet" href="../../lib/3rd/uploadify/uploadify.css" />
    <link rel="stylesheet" href="../../lib/3rd/jquery-ui/jquery-ui.min.css">
</head>
<body>
<div class="container-fluid">
    <form class="form-horizontal" id="item_form" enctype="multipart/form-data">
        <input type="hidden" id="sn" name="sn" value="2015103117335307">
        <div class="form-group">
        <label class="col-sm-3 control-label">付款事由</label>
        <div class="col-sm-6">
            <input type="text" class="input-block-level form-control" data-validate="{required: true, messages:{required:'请输入正确的付款事由!'}}" name="incident" id="incident" autocomplete="off" value="" />
        </div>
    </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">供应商关键字</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" id="supplier_key" value="" />
                <label>供应商名称</label><input type="text" class="input-block-level form-control" readonly="true" id="supplier_name" name="supplier_name" data-validate="{required: false, messages:{required:'请输入正确的供应商!'}}" value="" placeholder="自动联想">
                <input type="hidden" id="supplier" name="supplier" value="">
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">收款方</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" data-validate="{required: true, messages:{required:'请输入正确的收款方!'}}" name="payee" id="payee" autocomplete="off" value="" />
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">金额</label>
            <div class="col-sm-6">
                <div class="input-group"><span class="input-group-addon">￥</span><input type="text" class="input-block-level form-control" data-validate="{required: true, number:true, messages:{required:'请输入正确的金额!'}}" name="amount" id="amount" autocomplete="off" value="" /></div>
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">付款日期</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control datetimepicker" data-validate="{required: true, messages:{required:'请输入付款日期!'}}" name="payday" id="payday" autocomplete="off" value="curdate" />
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">付款人</label>
            <div class="col-sm-6">
                <select class="input-block-level form-control" data-validate="{required: true, messages:{required:'请输入付款人'}}" name="payer" id="payer" autocomplete="off" value="" >

                </select>
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">订单号</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" data-validate="{required: false, messages:{required:'请输入正确的订单号!'}}" name="purchase" id="purchase" autocomplete="off" value="" />
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">付款方式</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" data-validate="{required: true, messages:{required:'请输入正确的付款方式!'}}" name="paytype" id="paytype" autocomplete="off" value="" />
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-3 control-label">备注</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" data-validate="{required: false, messages:{required:'请输入正确的备注!'}}" name="mark" id="mark" autocomplete="off" value="" />
            </div>
        </div>

    </form>
</div>

<script src="../../lib/3rd/jquery/jquery.js"></script>
<script src="../../lib/3rd/bootstrap/js/bootstrap.min.js"></script>
<script src="../../lib/3rd/jquery/jquery.form.js"></script>
<script src="../../lib/3rd/jquery/validate/jquery.metadata.js"></script>
<script src="../../lib/3rd/jquery/validate/jquery.validate.js"></script>
<script src="../../lib/3rd/uploadify/jquery.uploadify.js"></script>
<script src="../../lib/3rd/datetimepicker/jquery.datetimepicker.js"></script>
<script src="../../lib/3rd/jquery-ui/jquery-ui.min.js"></script>
<script src="../../lib/webo/js/validateExtend.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script>
    function showResponse(resp) {
        if(resp.ret == "success"){
            top.hideTopModal()
            refreshContent()
        }else{
            if(resp.ret == "duplicated_value"){
                showError("添加失败! 重复的" + resp.result +  "。")
            }
        }
    }
    var refreshContent
    function onTopModalOk(options){
        if(options.refreshContent){
            refreshContent = options.refreshContent
        }


        if (! $("#item_form").valid()){
            return
        }
        $("#item_form").ajaxSubmit({
            type: "post",
            url: "\/item\/add\/account",
            success: showResponse
        });
    }
</script>
<script>$(function(){
    $("#supplier_key").autocomplete({
        source: "/item/autocomplete/supplier",
        autoFocus:true,
        focus: function( event, ui ) {
            $( "#supplier_key" ).val(ui.item.keyword);
            $( "#supplier_name" ).val(ui.item.name);
            $( "#supplier" ).val(ui.item.sn);
            return false;
        },
        minLength: 1,
        select: function( event, ui) {
            $( "#supplier_key" ).val(ui.item.keyword);
            $( "#supplier_name" ).val(ui.item.name);
            $( "#supplier" ).val(ui.item.sn);
            return false;
        },
        change: function( event, ui ) {
            if(!ui.item){
                $( "#supplier_name" ).val("");
                $( "#supplier" ).val("");
            }
        }
    })
            .autocomplete( "instance" )._renderItem = function( ul, item ) {
        return $( "<li>" )
                .append(item.keyword + "(" + item.name + ")")
                .appendTo( ul );
    };

    $("#payday").datetimepicker({timepicker:false,format:'Y.m.d',scrollMonth:false, lang:'zh',new date()})
});</script>

</body>
</html>