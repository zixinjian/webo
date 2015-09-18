<!DOCTYPE html>
<html>
<head lang="zh">
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../asserts/3rd/bootstrap/css/bootstrap.css">
    <link rel="stylesheet" href="../../asserts/3rd/datetimepicker/jquery.datetimepicker.css">
    <link rel="stylesheet" href="../../asserts/3rd/uploadify/uploadify.css" />
    <link rel="stylesheet" href="../../asserts/3rd/jquery-ui/jquery-ui.min.css">
</head>
<body>
<div class="container-fluid">
    <div class="alert" role="alert" style="display: none">添加成功！</div>
    <form class="form-horizontal" id="item_form" enctype="multipart/form-data">
    {{str2html .Form}}
    </form>
</div>
<script src="../../asserts/3rd/jquery/jquery.js"></script>
<script src="../../asserts/3rd/bootstrap/js/bootstrap.min.js"></script>
<script src="../../asserts/3rd/jquery/jquery.form.js"></script>
<script src="../../asserts/3rd/jquery/validate/jquery.metadata.js"></script>
<script src="../../asserts/3rd/jquery/validate/jquery.validate.js"></script>
<script src="../../asserts/3rd/uploadify/jquery.uploadify.js"></script>
<script src="../../asserts/3rd/datetimepicker/jquery.datetimepicker.js"></script>
<script src="../../asserts/3rd/jquery-ui/jquery-ui.min.js"></script>
<script src="../../static/js/validateExtend.js"></script>
<script src="../../static/js/ui.js"></script>
<script src="../../asserts/webo/util.js"></script>
<script src="../../static/js/travel.js"></script>
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
    function beforeSubmit(a){
        console.log("before submit", a)
        for (idx in a){
            valueMap = a[idx]
            if (valueMap.name == "expayrat"){
                valueMap.value = valueMap.value/100
            }
        }
        return true
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
            url: "{{.Service}}",
            beforeSubmit:beforeSubmit,
            success: showResponse
        });
    }
    $(function(){
        $("#expayrat").wrapAll('<div class="input-group"></div>')
        $("#expayrat").after('<span class="input-group-addon">%</span>')
        $("#payment").wrapAll('<div class="input-group"></div>')
        $("#payment").after('<a class="btn btn-sm input-group-addon" id="calc">计算</a>')
        $("#calc").click(calPayment)
    });
</script>
{{str2html .Onload}}
</body>
</html>