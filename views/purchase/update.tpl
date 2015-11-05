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
    <div class="alert" role="alert" style="display: none">添加成功！</div>
    <form class="form-horizontal" id="item_form">
        {{str2html .Form}}
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
</body>
</html>