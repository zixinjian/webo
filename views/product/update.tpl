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
        <div class="form-group">
            <label class="col-sm-3 control-label">附件</label>
            <div class="col-sm-6">
                <input type="file" name="fileUpload" id="file_upload" />
            </div>
        </div>
    </form>
</div>

<script src="../../lib/3rd/jquery/jquery.js"></script>
<script src="../../lib/3rd/bootstrap/js/bootstrap.min.js"></script>
<script src="../../lib/3rd/jquery/jquery.form.js"></script>
<script src="../../lib/3rd/jquery/validate/jquery.metadata.js"></script>
<script src="../../lib/3rd/jquery/validate/jquery.validate.js"></script>
<script src="../../lib/webo/js/validateExtend.js"></script>
<script src="../../lib/3rd/uploadify/jquery.uploadify.js"></script>
<script src="../../lib/3rd/datetimepicker/jquery.datetimepicker.js"></script>
<script src="../../lib/3rd/jquery-ui/jquery-ui.min.js"></script>
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
    $(function() {
        $('#file_upload').uploadify({
            'swf'      : '../../asserts/3rd/uploadify/uploadify.swf',
            'uploader' : '/item/upload/product?sn=' + $("#sn").val(),
            'cancelImg': '../../asserts/3rd/uploadify/uploadify-cancel.png',
            'fileObjName':'uploadFile'
        });
    });
</script>
{{str2html .Onload}}
</body>
</html>