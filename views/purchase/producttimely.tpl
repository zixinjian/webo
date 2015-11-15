<!DOCTYPE html>
<html>
<meta charset="UTF-8">
<link rel="stylesheet" href="../../lib/bootstrap/css/bootstrap.min.css">
<link rel="stylesheet" href="../../lib/bootstrap-table/bootstrap-table.css">
<link rel="stylesheet" href="../../lib/3rd/bootstrap-editable/bootstrap3-editable/css/bootstrap-editable.css">
<link rel="stylesheet" href="../../lib/jquery/jquery-ui/jquery-ui.min.css">
<link rel="stylesheet" href="../../lib/webo/css/ui.css">
</head>
<body>
<div>
    <div style="line-height: 20px; padding:20px;">
        <div class = "form-inline">
            <div class="form-group">
                <label for="product_key">产品关键字</label>
                <input type="text" class="form-control" id="product_key" placeholder="请输入关键字">
                <label for="product_name">名称</label>
                <input type="text" class="form-control" id="product_name" placeholder="自动联想" readonly>
                <input type="hidden" id="product">
            </div>
            <button type="button" id="analyzeBtn" class="btn btn-default btn-primary">统计</button>
            <div class="alert" role="alert" style="display: none; margin: 10px 0 0 0"></div>
        </div>
    </div>
    <table id="item_table"
           data-row-style="rowStyleOvertime"
           data-page-size="25">
        <thead>
        <tr>
            <th data-field="sn"  data-sortable="true" data-visible="false">编号</th>
            <th data-field="placedate"  data-sortable="true">下单日期</th>
            <th data-field="requireddate"  data-sortable="true">需用日期</th>
            <th data-field="requireddepartment"  data-sortable="true">申请部门</th>
            <th data-field="num"  data-sortable="true">数量</th>
            <th data-field="unitprice"  data-sortable="true">单价</th>
            <th data-field="productprice"  data-sortable="true">参考价</th>
            <th data-field="totalprice"  data-sortable="true">总价</th>
            <th data-field="freightprice"  data-sortable="true">运费</th>
            <th data-field="supplierkey"  data-sortable="true" data-formatter="supplierFormatter">供应商关键词</th>
            <th data-field="buyer"  data-sortable="true">采购人</th>
            <th data-field="orderdate"  data-sortable="true" data-order="desc">订货日期</th>
            <th data-field="predictdeliverydate"  data-sortable="true">预计发货日期</th>
            <th data-field="actualdeliverydate"  data-sortable="true">实际发货日</th>
            <th data-field="arrivaldate"  data-sortable="true">到货日期</th>
            <th data-field="godowndate"  data-sortable="true">入库日期</th>
            <th data-field="changelog"  data-sortable="true">变更情况</th>
            <th data-field="mark"  data-sortable="true">备注</th>
        </tr>
        </thead>
    </table>
</div>
<script src="../../lib/jquery/jquery/jquery.js"></script>
<script src="../../lib/bootstrap/js/bootstrap.min.js"></script>
<script src="../../lib/bootstrap-table/bootstrap-table.js"></script>
<script src="../../lib/bootstrap-table/locale/bootstrap-table-zh-CN.js"></script>
<script src="../../lib/jquery/jquery-ui/jquery-ui.min.js"></script>
<script src="../../lib/moment/moment.js"></script>
<script src="../../lib/webo/js/util.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script src="../../lib/webo/js/purchase.js"></script>
<script>
    var $table = $("#item_table")
    function responseHandler(res){
        return res.rows
    }
    function supplierFormatter(value, row){
        return wbSprintf('<span title="%s(%s)">%s</span>', value, row.suppliername, value)
    }
    $(function(){
        $table.bootstrapTable({url:"/item/list/purchase?product", method:"post", sidePagination:"server", pagination:true, height:getTableHeight() -150});
        $(window).resize(function () {
            $table.bootstrapTable('resetView', {
                height: getTableHeight() - 150
            });
        });
        $("#analyzeBtn").on("click", function(){
            productSn = $("#product").val()
            if (!productSn){
                return
            }
            $.post("/purchase/calc/producttimely",
                {
                    product:productSn
                },
                function(data, status){
                    if (status != "success"){
                        showError("未查到数据！")
                        return
                    }
                    calcRet = wbSprintf("延期:%s, 总数:%s, 及时率: %s", data.delay, data.total, data.rat)
//                    console.log("calcRet", calcRet)
                    showSuccess(calcRet)
                });
            $table.bootstrapTable('refresh', {
                url:"/item/list/purchase",
                query:{"product":productSn}
            });
        });
        $("#product_key").autocomplete({
            source: "/item/autocomplete/product",
            autoFocus:true,
            focus: function( event, ui ) {
                if (!ui.item){
                    return
                }
                $( "#product_key" ).val(ui.item.keyword );
                $( "#product_name" ).val(ui.item.name);
                $( "#product" ).val(ui.item.sn);
                return false;
            },
            minLength: 1,
            select: function( event, ui) {
                if (!ui.item){
                    return
                }
                $( "#product_key" ).val(ui.item.keyword);
                $( "#product_name" ).val(ui.item.name);
                $( "#product" ).val(ui.item.sn);
                return false;
            },
            change: function( event, ui ) {
                if(!ui.item){
                    $( "#product_key" ).val("");
                    $( "#product_name" ).val(ui.item.name);
                    $( "#product" ).val("");
                }
            }
        }).autocomplete( "instance" )._renderItem = function( ul, item ) {
            return $( "<li>" )
                    .append(item.keyword + "(" + item.name + ")")
                    .appendTo( ul );
        };
    });
</script>
</body>
</html>