<!DOCTYPE html>
<html>
<meta charset="UTF-8">
<link rel="stylesheet" href="../../lib/app/css/app.min.css"/>
<link rel="stylesheet" href="../../lib/bootstrap-table/bootstrap-table.css">
<link rel="stylesheet" href="../../lib/3rd/bootstrap-editable/bootstrap3-editable/css/bootstrap-editable.css">
<link rel="stylesheet" href="../../lib/webo/css/ui.css">
<!-- Le HTML5 shim, for IE6-8 support of HTML5 elements -->
<!--[if lt IE 9]>
<script src="../../lib/html5shiv.min.js"></script>
<![endif]-->
</head>
<body>
<div>
    <table id="item_table"
           data-show-refresh="true"
           data-show-columns="true"
           data-search="true"
           data-page-size="25"
           data-sort-name="rat"
           data-sort-order="desc"
           data-toolbar=".toolbar">
        <thead>
        <tr>
            <th data-field="buyer"  data-sortable="true">采购人</th>
            <th data-field="intime"  data-sortable="true">延期数量</th>
            <th data-field="total"  data-sortable="true">总数量</th>
            <th data-field="rat"  data-sortable="true">及时率(%)</th>
        </tr>
        </thead>
    </table>
</div>
<script src="../../lib/app/js/app.min.js"></script>
<script src="../../lib/bootstrap-table/bootstrap-table.js"></script>
<script src="../../lib/bootstrap-table/locale/bootstrap-table-zh-CN.js"></script>
<script src="../../lib/webo/js/poplayer.js"></script>
<script src="../../lib/webo/js/util.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script>
    var $table = $("#item_table")
    $(function(){
        $table.bootstrapTable({url:"/purchase/list/buyertimely", method:"post", sidePagination:"server", pagination:true, height:getTableHeight()});
        $(window).resize(function () {
            $table.bootstrapTable('resetView', {
                height: getTableHeight()
            });
        });
    });
</script>
</body>
</html>