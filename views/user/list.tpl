<!DOCTYPE html>
<html>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../lib/app/css/app.min.css" type="text/css">
    <link rel="stylesheet" href="../../lib/bootstrap-table/bootstrap-table.css">
    <link rel="stylesheet" href="../../lib/simple-line-icons/css/simple-line-icons.css">
    <link rel="stylesheet" href="../../lib/webo/css/ui.css">
    <!-- Le HTML5 shim, for IE6-8 support of HTML5 elements -->
    <!--[if lt IE 9]>
    <script src="../../lib/html5shiv.min.js"></script>
    <![endif]-->
    </head>
<body>
<div>
    <p class="toolbar">
        <a id="add_item" class="create btn btn-primary">新建</a>
    </p>
    <table id="item_table"
           data-show-refresh="true"
           data-show-columns="true"
           data-page-size="25"
           data-sortable="false"
           data-toolbar=".toolbar">
        <thead>
            <tr>
                <th data-field="action"
                    data-align="center"
                    data-formatter="actionFormatter"
                    data-events="actionEvents"
                    data-width="75px">  [ 操作 ]  </th>
                {{str2html .thlist}}
            </tr>
        </thead>
    </table>
</div>
<script src="../../lib/app/js/app.min.js"></script>
<script src="../../lib/bootstrap-table/bootstrap-table.js"></script>
<script src="../../lib/bootstrap-table/locale/bootstrap-table-zh-CN.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script>
    var $table = $("#item_table")
    $(function(){
        $table.bootstrapTable({url:"{{.listUrl}}", method:"post", sidePagination:"server", pagination:true,
            height:getTableHeight(), rowStyle:rowStyle, sortName:"flag", sortOrder:"asc",fixedColumns: true,fixedNumber:1});
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
    function actionFormatter(value, row) {
        return [
            '<a class="update" href="javascript:" title="修改" style="margin-right: 5px;"><i class="icon-note text-primary-dker"></i></a>',
        ].join('');
    }
    window.actionEvents = {
        'click .update': function (e, value, row) {
            top.showTopModal({url:"{{.updateUrl}}?sn=" + row.sn, refreshContent:refreshContent});
//        },
//        'click .disable': function (e, value, row) {
//            if (confirm('你确定要禁用本用户吗?')) {
//                $.ajax({
//                    url: "/user/disable?sn=" + row.sn,
//                    type: 'Get',
//                    success: function () {
//                        $table.bootstrapTable('refresh');
//                        showAlert('已禁用!', 'success');
//                    },
//                    error: function () {
//                        showAlert('Delete item error!', 'danger');
//                    }
//                })
//            }
        }
    }
    function rowStyle(row, index) {
        if(row.flag != "可用"){
            return {classes: "warning"};
        }
        return {}
    }
</script>
</body>
</html>