<!DOCTYPE html>
<html>
<meta charset="UTF-8">
<link rel="stylesheet" href="../../lib/bootstrap/css/bootstrap.min.css">
<link rel="stylesheet" href="../../lib/bootstrap-table/bootstrap-table.css">
<link rel="stylesheet" href="../../lib/3rd/bootstrap-editable/bootstrap3-editable/css/bootstrap-editable.css">
<link rel="stylesheet" href="../../lib/webo/css/ui.css">
</head>
<body>
<div>
    <table id="item_table"
           data-show-refresh="true"
           data-show-columns="true"
           data-page-size="25"
           data-sort-name="rat"
           data-sort-order="desc"
           data-toolbar=".toolbar">
        <thead>
        <tr>
            <th data-field="supplier"  data-sortable="true">供应商</th>
            <th data-field="intime"  data-sortable="true">达标数</th>
            <th data-field="total"  data-sortable="true">总数</th>
            <th data-field="rat"  data-sortable="true">及时率(%)</th>
        </tr>
        </thead>
    </table>
</div>
<script src="../../lib/jquery/jquery/jquery.js"></script>
<script src="../../lib/bootstrap/js/bootstrap.min.js"></script>
<script src="../../lib/bootstrap-table/bootstrap-table.js"></script>
<script src="../../lib/bootstrap-table/locale/bootstrap-table-zh-CN.js"></script>
<script src="../../lib/webo/js/poplayer.js"></script>
<script src="../../lib/webo/js/util.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script>
    var $table = $("#item_table")
    $(function(){
        $table.bootstrapTable({url:"/purchase/list/suppliertimely", method:"post", sidePagination:"server", pagination:true, height:getTableHeight()});
        $(window).resize(function () {
            $table.bootstrapTable('resetView', {
                height: getTableHeight()
            });
        });
    });
</script>
</body>
</html>