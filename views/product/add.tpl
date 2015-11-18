<!DOCTYPE html>
<html>
<head lang="zh">
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1" />
    <link rel="stylesheet" href="../../lib/app/css/app.min.css"/>
    <link rel="stylesheet" href="../../lib/font-awesome/css/font-awesome.min.css" type="text/css" />
    <link rel="stylesheet" href="../../lib/simple-line-icons/css/simple-line-icons.css" type="text/css" />
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
<div class="container-fluid bg-white">
    <form class="form-horizontal" id="item_form" enctype="multipart/form-data">
        {{str2html .Form_sn}}
        <div class="form-group">
            <label class="col-sm-2 control-label">类别</label>
            <div class="col-sm-8">
                <select class="input-block-level form-control" data-validate="{required: true, messages:{required:'请输入类别'}}" name="category" id="category" autocomplete="off" value="cate_engine" >
                    {{str2html .CategoryOptions}}
                </select>
            </div>
        </div>
        {{str2html .Form_name}}
        {{str2html .Form_brand}}
        {{str2html .Form_model}}
        {{str2html .Form_power}}
        {{str2html .Form_detail}}
        <div class="form-group">
            <label class="col-sm-2 control-label">附件</label>
            <div class="col-sm-8">
                <input type="file" name="fileUpload" id="file_upload" />
            </div>
        </div>
        <div class="form-group">
            <label class="col-sm-2 control-label">供应商</label>
            <div class="col-sm-8">
                <input type="text" class=" form-control" id="supplier_key">
                <span id="supplierList" class="help-block" style="margin-bottom: 0"></span>
            </div>
        </div>
        {{str2html .Form_size}}
        {{str2html .Form_freight}}
        {{str2html .Form_price}}
        {{str2html .Form_profitrat}}
        <div class="form-group">
            <label class="col-sm-2 control-label">卖价</label>
            <div class="col-sm-8">
                <div class="input-group"><span class="input-group-addon">￥</span><input type="text" class="input-block-level form-control" data-validate="{required: false, number:true, messages:{required:'请输入正确的卖价!'}}" name="retailprice" id="retailprice" autocomplete="off" value="" />
                    <a class="btn btn-sm input-group-addon" id="calc">计算</a>
                </div>
            </div>
        </div>
    </form>
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
<script src="../../lib/webo/js/util.js"></script>
<script src="../../lib/webo/js/catagory.js"></script>
<script src="../../lib/webo/js/product.js"></script>
<script>
    function showResponse(resp) {
        if(resp.ret == "success"){
            top.hideTopModal()
            refreshContent()
        }else{
            if(resp.ret == "duplicated_value"){
                if (resp.result == "product.model, product.name"){
                    showError("添加失败! 重复的" + "型号" +  "。")
                }else{
                    showError("添加失败! 重复的" + resp.result +  "。")
                }
            }else{
                showError("添加失败! " + resp.result +  "。")
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
            url: "/item/product/add",
            success: showResponse
        });
    }
    $(function(){
        initCategory($("#name"))
        $("#calc").click(calRetailPrice)
    })
</script>
</body>
</html>