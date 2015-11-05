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
           data-page-size="25"
           data-toolbar=".toolbar">
        <thead>
        <tr>
            <th data-field="action"
                data-align="center"
                data-formatter="actionFormatter"
                data-events="actionEvents"
                data-width="75px">  [ 操作 ]  </th>
            <th data-field="sn" data-visible="false" >编号</th>
            <th data-field="travelername"  data-sortable="true">出差人</th>
            <th data-field="task"  data-sortable="true">任务</th>
            <th data-field="starttime"  data-sortable="true">开始日期</th>
            <th data-field="endtime"  data-sortable="true">结束日期</th>
            <th data-field="path"  data-sortable="true">路线</th>
            <th data-field="requirement"  data-sortable="true">达标标准</th>
            <th data-field="expenses"  data-sortable="true">预计费用</th>
            <th data-field="actualexpenses"  data-sortable="true">实际费用</th>
            <th data-field="detail"  data-sortable="true">费用分项</th>
            <th data-formatter="diffexpenseFormat">节约/超支</th>
            <th data-field="diffrate" data-formatter="diffrateFormat">节约/超支报销比例</th>
            <th data-field="payment" data-sortable="true">实报费用</th>
            <th data-field="approvername"  data-sortable="true">审批人</th>
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
    function actionFormatter(value, row) {
        return [
            '<a class="update" href="javascript:" title="修改" style="margin-right: 5px;">审批</a>',
//            wbSprintf('<a class="file" href="/static/files/{{.item}}/%s" target="_blank" title="附件" data-toggle="poplayer" data-placement="bottom" data-url="/static"><i class="glyphicon glyphicon-file"></i></a>', row.sn),
        ].join('');
    }
    function diffexpenseFormat(value, row){
        if (row.actualexpenses == 0 && row.expenses == 0){
            return 0
        }
        if (!row.actualexpenses || !row.expenses){
            return ""
        }
        return getDiffexpenese(row)
    }
    function diffrateFormat(value, row){
        return row.expayrat *100 + "%"
    }
    function rowStyle(row, index) {
        if(getDiffexpenese() > 0){
            return {classes: "warning"};
        }
        return {}
    }
    function getDiffexpenese(row){
        return wbToMoney(parseFloat(row.actualexpenses) - parseFloat(row.expenses))
    }
    window.actionEvents = {
        'click .update': function (e, value, row) {
            top.showTopModal({url:"{{.updateUrl}}?sn=" + row.sn, refreshContent:refreshContent});
        }
    }
</script>
</body>
</html>