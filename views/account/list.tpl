<!DOCTYPE html>
<html>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../lib/app/css/app.min.css">
    <link rel="stylesheet" href="../../lib/3rd/bootstrap-table/bootstrap-table.css">
    <link rel="stylesheet" href="../../lib/webo/css/overwrite.css">
</head>
<body>
<div>
    <p class="toolbar">
        <a id="add_item" class="create btn btn-primary">新建</a>
    </p>
    <table id="item_table"
           data-show-refresh="true"
           data-show-columns="true"
           data-search="true"
           data-page-size="25"
           data-query-params="queryParams"
           data-toolbar=".toolbar">
        <thead>
            <tr>
                <th data-field="action"
                    data-align="center"
                    data-formatter="actionFormatter"
                    data-events="actionEvents"
                    data-width="75px">  [ 操作 ]  </th>
                <th data-field="sn" data-visible="false" >编号</th>
                <th data-field="incident"  data-sortable="true">付款事由</th>
                <th data-field="supplier"  data-sortable="true">供应商</th>
                <th data-field="amount"  data-sortable="true">金额</th>
                <th data-field="payday"  data-sortable="true">付款日期</th>
                <th data-field="payername"  data-sortable="true">付款人</th>
                <th data-field="purchase"  data-sortable="true">订单号</th>
                <th data-field="paytype"  data-sortable="true">付款方式</th>
                <th data-field="mark"  data-sortable="true">备注</th>
            </tr>
        </thead>
    </table>
</div>
<script src="../../lib/app/js/app.min.js"></script>
<script src="../../lib/3rd/bootstrap-table/bootstrap-table.js"></script>
<script src="../../lib/3rd/bootstrap-table/locale/bootstrap-table-zh-CN.js"></script>
<script src="../../lib/webo/util.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script>
    var $table = $("#item_table")
    $(function(){
        $table.bootstrapTable({url:"{{.listUrl}}", method:"post", sidePagination:"server", pagination:true, height:getTableHeight()});
        $("#add_item").on("click", function(){
            top.showTopModal({url:"{{.addUrl}}", refreshContent:refreshContent});
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
    function queryParams(params){
        return params
    }
    function actionFormatter(value, row) {
        return [
            '<a class="update" href="javascript:" title="修改" style="margin-right: 5px;"><i class="glyphicon glyphicon-edit"></i></a>',
            wbSprintf('<a class="file" href="/static/files/{{.item}}/%s" target="_blank" title="附件" data-toggle="poplayer" data-placement="bottom" data-url="/static"><i class="glyphicon glyphicon-file"></i></a>', row.sn),
        ].join('');
    }
    window.actionEvents = {
        'click .update': function (e, value, row) {
            top.showTopModal({url:"{{.updateUrl}}?sn=" + row.sn, refreshContent:refreshContent});
        }
    }
</script>
</body>
</html>