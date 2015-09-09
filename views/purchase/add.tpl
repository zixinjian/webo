<!DOCTYPE html>
<html>
<head lang="zh">
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../asserts/3rd/bootstrap/css/bootstrap.css">
    <link rel="stylesheet" href="../../asserts/3rd/datetimepicker/jquery.datetimepicker.css">
    <link rel="stylesheet" href="../../asserts/3rd/jquery-ui/jquery-ui.min.css">
    <style>
        .ui-autocomplete-loading {
            background: white url("../../asserts/webo/images/ui-anim_basic_16x16.gif") right center no-repeat;
        }
    </style>
</head>
<body>
<div class="container-fluid">
    <div class="alert" role="alert" style="display: none">添加成功！</div>
    <form class="form-horizontal" id="item_form">
        {{str2html .Form}}
    </form>
</div>

<script src="../../asserts/3rd/jquery/jquery.js"></script>
<script src="../../asserts/3rd/bootstrap/js/bootstrap.min.js"></script>
<script src="../../asserts/3rd/jquery/jquery.form.js"></script>
<script src="../../asserts/3rd/jquery/validate/jquery.metadata.js"></script>
<script src="../../asserts/3rd/jquery/validate/jquery.validate.js"></script>
<script src="../../asserts/3rd/datetimepicker/jquery.datetimepicker.js"></script>
<script src="../../asserts/3rd/jquery-ui/jquery-ui.min.js"></script>
<script src="../../static/js/ui.js"></script>
<script>
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
        if (! $("#item_form").valid()){
            return
        }
        $("#item_form").ajaxSubmit({
            type: "post",
            url: "{{.Service}}",
            success: showResponse
        });
    }
    $(function () {
        $("#supplier_key").autocomplete({
            source: "/item/autocomplete/supplier",
            autoFocus:true,
            focus: function( event, ui ) {
                $( "#supplier_key" ).val(ui.item.keyword);
                $( "#supplier_name" ).val(ui.item.name);
                $( "#supplier" ).val(ui.item.sn);
                return false;
            },
            minLength: 2,
            select: function( event, ui) {
                $( "#supplier_key" ).val(ui.item.keyword);
                $( "#supplier_name" ).val(ui.item.name);
                $( "#supplier" ).val(ui.item.sn);
                return false;
            },
            change: function( event, ui ) {
                console.log("ui", ui.item)
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
    })
</script>
{{str2html .Onload}}
</body>
</html>