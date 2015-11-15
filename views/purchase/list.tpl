<!DOCTYPE html>
<html>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1" />
<link rel="stylesheet" href="../../lib/simple-line-icons/css/simple-line-icons.css" type="text/css" />
<link rel="stylesheet" href="../../lib/app/css/app.min.css" type="text/css" />
<link rel="stylesheet" href="../../lib/bootstrap-table/bootstrap-table.css">
<link rel="stylesheet" href="../../lib/webo/css/ui.css">
</head>
<body>
<div>
    <div class="toolbar" style="line-height: 20px">
        <div class="form-group">
            {{str2html .buyers}}
        </div>
    </div>
    <table id="item_table"
           data-show-refresh="true"
           data-show-columns="true"
           data-search="true"
           data-page-size="25"
           data-query-params="queryParams"
           data-row-style="rowStyleOvertime"
           data-toolbar=".toolbar">
        <thead>
        <tr>
            <th data-field="action"
                data-align="center"
                data-formatter="actionFormatter"
                data-events="actionEvents"
                data-sortable="false"
                data-width="75px">  [ 操作 ]  </th>
            <th data-field="sn"  data-sortable="true" data-visible="false">[ 编号 ]</th>
            <th data-field="category"  data-sortable="true">类别</th>
            <th data-field="productname"  data-sortable="false">商品名称</th>
            <th data-field="brand"  data-sortable="false">品牌</th>
            <th data-field="model"  data-sortable="false" data-formatter="modelFormatter">型号</th>
            <th data-field="power"  data-sortable="false">功率(KW)</th>
            <th data-field="num"  data-sortable="true">数量</th>
            <th data-field="placedate"  data-sortable="true">下单日期</th>
            <th data-field="requireddate"  data-sortable="true">需用日期</th>
            <th data-field="requireddepartment"  data-sortable="true">申请部门</th>
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
            <th data-field="paymentamount" data-visible="false" data-sortable="true">付款金额</th>
            <th data-field="paymentdate" data-visible="false" data-sortable="true">付款日期</th>
            <th data-field="file" data-visible="false" data-sortable="true">附件</th>
            <th data-field="changelog"  data-sortable="true">变更情况</th>
            <th data-field="mark"  data-sortable="true">备注</th>
        </tr>
        </thead>
    </table>
</div>
<script src="../../lib/app/js/app.min.js"></script>
<script src="../../lib/bootstrap-table/bootstrap-table.js"></script>
<script src="../../lib/bootstrap-table/locale/bootstrap-table-zh-CN.js"></script>
<script src="../../lib/moment/moment.js"></script>
<script src="../../lib/webo/js/util.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script src="../../lib/webo/js/purchase.js"></script>
{{str2html .queryParams}}
<script>
    var $table = $("#item_table")
    $(function(){
        $table.bootstrapTable({url:"{{.listUrl}}", method:"post", sidePagination:"server", pagination:true,
            height:getTableHeight(), sortName:"placedate", sortOrder:"{{.sortOrder}}",
            fixedColumns: true,fixedNumber:1});
        $("#add_item").on("click", function(){
            top.showTopModal({url:"{{.addUrl}}", refreshContent:refreshContent});
        })
        $table.on("post-body.bs.table", function(){
//            console.log("post-body.bs.table")
            $('[data-toggle="tooltip"]').tooltip()
        })
        $(window).resize(function () {
            $table.bootstrapTable('resetView', {
                height: getTableHeight()
            });
        });
    });
    function refreshContent(options){
        top.hideTopModal()
        $table.bootstrapTable("refresh")
    }
    function actionFormatter(value, row) {
        return [
            '<a class="update" href="javascript:" title="修改" style="margin-right: 5px;"><i class="icon-note text-primary-dker"></i></a>',
            wbSprintf('<a class="file" href="/static/files/purchase/%s" target="_blank" title="附件"><i class="icon-tag text-primary-dker"></i></a>', row.sn),
        ].join('');
    }
    window.actionEvents = {
        'click .update': function (e, value, row) {
            top.showTopModal({url:"{{.updateUrl}}?sn=" + row.sn, refreshContent:refreshContent});
        }
    }
    $(function(){
        $("[data-model='buyers']").on("change", function(){
            var selectedValue = $("input[name='buyers']:checked").val();
            if (selectedValue == "all"){
                $table.bootstrapTable("refresh")
            }else{
                $table.bootstrapTable("refresh", {query: {buyer: selectedValue}})
            }
        })
    });
    function rowStyle(row, index) {
        if(row.godowndate == ""){
            now = moment()
            requiredDate = moment(row.requireddate, "YYYY.MM.DD")
            if (requiredDate.diff(now, "days") < 0){
                return {classes: "danger"};
            }
            if (requiredDate.diff(now, "days") < 2 ){
                return {classes: "warning"};
            }
        }
        return {}
    }
</script>
</body>
</html>