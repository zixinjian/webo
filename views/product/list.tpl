<!DOCTYPE html>
<html>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../lib/app/css/app.min.css">
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
                    data-width="75px"
                    data-sortable="false"><span style="width: 150px">操作</span></th>
                <th data-field="sn" data-visible="true" ><span style="width: 150px">编号</span></th>
                <th data-field="category">类别</th>
                <th data-field="name">名称</th>
                <th data-field="brand">品牌</th>
                <th data-field="keyword" data-visible="false" >关键词</th>
                <th data-field="model">型号</th>
                <th data-field="power">功率</th>
                <th data-field="detail">产品详情</th>
                <th data-field="supplier" data-formatter="supplierFormatter">供应商</th>
                <th data-field="size">尺寸重量</th>
                <th data-field="freight">预计运费</th>
                <th data-field="price">参考价格</th>
                <th data-field="profitrat">外卖比例</th>
                <th data-field="retailprice">卖价</th>

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
    function supplierFormatter(value, row){
        suppliers = row.supplier_enum
        ks = []
        for (i in suppliers){
            supplier= suppliers[i]
            ks.push(supplier.keyword)
        }
        return ks.join(",")
    }
    function actionFormatter(value, row) {
        return [
            '<a class="update" href="javascript:" title="修改" style="margin-right: 5px;"><i class="icon-note text-primary-dker"></i></a>',
            wbSprintf('<a class="file" href="javascript:" title="附件"><i class="icon-tag text-primary-dker"></i></a>', row.sn),
        ].join('');
    }
    window.actionEvents = {
        'click .update': function (e, value, row) {
            top.showTopModal({url:"{{.updateUrl}}?sn=" + row.sn, refreshContent:refreshContent});
        },
        'click .file': function (e, value, row) {
            window.open("/static/files/product/" + row.sn)
        }
    }
</script>
</body>
</html>