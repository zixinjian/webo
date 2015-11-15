<!DOCTYPE html>
<html>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../lib/app/css/app.min.css">
    <link rel="stylesheet" href="../../lib/bootstrap-table/bootstrap-table.css">
    <link rel="stylesheet" href="../../lib/simple-line-icons/css/simple-line-icons.css">
    <link rel="stylesheet" href="../../lib/webo/css/ui.css">
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
                    data-sortable="false"
                    data-width="75px">  [ 操作 ]  </th>
                {{str2html .thlist}}
            </tr>
        </thead>
    </table>
</div>
<script src="../../lib/app/js/app.min.js"></script>
<script src="../../lib/bootstrap-table/bootstrap-table.js"></script>
<script src="../../lib/bootstrap-table/locale/bootstrap-table-zh-CN.js"></script>
<script src="../../lib/webo/js/util.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script>
    var $table = $("#item_table")
    $(function(){
        $table.bootstrapTable({url:"{{.listUrl}}", method:"post", sidePagination:"server", pagination:true, height:getTableHeight(),
            fixedColumns: true,fixedNumber:1});
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
            '<a class="update" href="javascript:" title="修改" style="margin-right: 5px;"><i class="icon-note text-primary-dker"></i></a>',
            wbSprintf('<a class="file" href="/static/files/{{.item}}/%s" target="_blank" title="附件"><i class="icon-tag text-primary-dker"></i></a>', row.sn),
        ].join('');
    }
    window.actionEvents = {
        'click .update': function (e, value, row) {
            top.showTopModal({url:"{{.updateUrl}}?sn=" + row.sn, refreshContent:refreshContent});
        },
        'click .remove': function (e, value, row) {
            if (confirm('你确定要删除本行吗?')) {
                $.ajax({
                    url: API_URL + row.id,
                    type: 'delete',
                    success: function () {
                        $table.bootstrapTable('refresh');
                        showAlert('Delete item successful!', 'success');
                    },
                    error: function () {
                        showAlert('Delete item error!', 'danger');
                    }
                })
            }
        }
    }
</script>
</body>
</html>