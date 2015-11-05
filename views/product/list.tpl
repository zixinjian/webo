<!DOCTYPE html>
<html>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../lib/3rd/bootstrap/css/bootstrap.min.css">
    <link rel="stylesheet" href="../../lib/3rd/bootstrap-table/bootstrap-table.css">
    <link rel="stylesheet" href="../../lib/3rd/bootstrap-editable/bootstrap3-editable/css/bootstrap-editable.css">
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
           data-query-params="queryParams"
           data-page-size="25"
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
<script src="../../lib/3rd/jquery/jquery.js"></script>
<script src="../../lib/3rd/bootstrap/js/bootstrap.min.js"></script>
<script src="../../lib/3rd/bootstrap-table/bootstrap-table.js"></script>
<script src="../../lib/3rd/bootstrap-table/locale/bootstrap-table-zh-CN.js"></script>
<script src="../../lib/3rd/bootstrap-table/extensions/editable/bootstrap-table-editable.js"></script>
<script src="../../lib/3rd/bootstrap-editable/bootstrap3-editable/js/bootstrap-editable.js"></script>
<script src="../../lib/webo/poplayer.js"></script>
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
            wbSprintf('<a class="file" href="/static/files/product/%s" target="_blank" title="附件" data-toggle="poplayer" data-placement="bottom" data-url="/static"><i class="glyphicon glyphicon-file"></i></a>', row.sn),
        ].join('');
    }
    window.actionEvents = {
        'click .update': function (e, value, row) {
            top.showTopModal({url:"{{.updateUrl}}?sn=" + row.sn, refreshContent:refreshContent});
        },
        'click .file': function (e, value, row) {
//            $e = $(e.currentTarget)
//            $e.poplayer({url:"/static"})
//            $e.poplayer('show')
//            console.log(e, value, row)
        }
    }
</script>
</body>
</html>