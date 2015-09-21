<!DOCTYPE html>
<html>
<head lang="zh">
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../asserts/3rd/bootstrap/css/bootstrap.css">
    <link rel="stylesheet" href="../../asserts/3rd/datetimepicker/jquery.datetimepicker.css">
    <link rel="stylesheet" href="../../asserts/3rd/uploadify/uploadify.css" />
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
<script src="../../asserts/3rd/uploadify/jquery.uploadify.js"></script>
<script src="../../asserts/3rd/datetimepicker/jquery.datetimepicker.js"></script>
<script src="../../asserts/3rd/jquery-ui/jquery-ui.min.js"></script>
<script src="../../asserts/js/validateExtend.js"></script>
<script src="../../asserts/js/ui.js"></script>

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
    function setProductValues(item){
//        console.log(item)
        $("#product_key" ).val(item.keyword);
        $("#product_name" ).val(item.name);
        $("#product" ).val(item.sn);
        $("#productprice").val(item.price)
    }
    function clearProductValues(){
        $( "#product_name" ).val("");
        $( "#product" ).val("");
        $("#productprice").val("")
    }
    $(function () {
        $("#product_key").autocomplete({
            source: "/item/autocomplete/product",
            autoFocus:true,
            focus: function( event, ui ) {
                setProductValues(ui.item)
                return false;
            },
            minLength: 1,
            select: function( event, ui) {
                setProductValues(ui.item)
                return false;
            },
            change: function( event, ui ) {
                if(!ui.item){
                    clearProductValues()
                }
            }
        }).autocomplete( "instance" )._renderItem = function( ul, item ) {
            return $( "<li>" )
                    .append(item.keyword + "(" + item.name + ")")
                    .appendTo( ul );
        };

        $("#placedate").datetimepicker({timepicker:false,format:'Y.m.d',lang:'zh',value:new Date()})
        $("#requireddate").datetimepicker({timepicker:false,format:'Y.m.d',lang:'zh'})
    });
</script>
</body>
</html>