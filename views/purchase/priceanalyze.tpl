<!DOCTYPE html>
<html>
<head lang="zh">
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../../lib/bootstrap/css/bootstrap.css">
    <link rel="stylesheet" href="../../lib/jquery/jquery-ui/jquery-ui.min.css">
    <link rel="stylesheet" href="../../lib/webo/css/ui.css">
</head>
<body>
<div class="container-fluid" style="padding:20px;">
    <div class="form-inline">
        <div class="form-group">
            <label for="product_key">产品关键字</label>
            <input type="text" class="form-control" id="product_key" placeholder="请输入关键字">
            <label for="product_name">名称</label>
            <input type="text" class="form-control" id="product_name" placeholder="自动联想" readonly>
            <input type="hidden" id="product">
        </div>
        <button type="button" id="analyzeBtn" class="btn btn-default btn-primary">分析</button>
    </div>
    <div class="alert" role="alert" style="display: none;margin-top:20px"></div>
    <div id="analyzeCharts" style="height: 500px;margin-top:20px"></div>
</div>

<script src="../../lib/jquery/jquery/jquery.js"></script>
<script src="../../lib/bootstrap/js/bootstrap.min.js"></script>
<script src="../../lib/echart/echarts-all.js"></script>
<script src="../../lib/jquery/jquery-ui/jquery-ui.min.js"></script>
<script src="../../lib/webo/js/ui.js"></script>
<script>
//    PNotify.prototype.options.styling = "bootstrap3";
    var analyzChart = echarts.init(document.getElementById("analyzeCharts"));
    var option = {
        tooltip: {
            show: true
        },
        legend: {
            data:['单价', "参考价格"]
        },
        xAxis : [
            {
                type : 'category',
                data : []
            }
        ],
        yAxis : [
            {
                type : 'value'
            }
        ],
        series : [
            {
                "name":"单价",
                "type":"line",
                "data":[]
            },
            {
                "name":"参考价格",
                "type":"line",
                "data":[]
            }
        ]
    };
    function setChartData(rows){
        option.xAxis[0].data =[]
        option.series[0].data = []
        for (i in rows){
            row = rows[i]
            if(row.unitprice == 0){
                continue
            }
            option.xAxis[0].data.push(row.sn)
            option.series[0].data.push(row.unitprice)
            option.series[1].data.push(row.productprice)
        }
        analyzChart.setOption(option)
    }
    $(function () {
        $("#analyzeBtn").on("click", function(){
            $.post("/item/list/purchase",
                {
                    product:$("#product").val()
                },
                function(data,status){
                    if (status != "success" || data.Status != "success" || data.rows.length <= 0){
                        analyzChart.clear()
                        showError("未查到数据！")
                        return
                    }
                    hideAlert()
                    setChartData(data.rows)
                });
        });
        $("#product_key").autocomplete({
            source: "/item/autocomplete/product",
            autoFocus:true,
            focus: function( event, ui ) {
                $( "#product_key" ).val(ui.item.keyword );
                $( "#product_name" ).val(ui.item.name);
                $( "#product" ).val(ui.item.sn);
                return false;
            },
            minLength: 1,
            select: function( event, ui) {
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
        $("#analyzeCharts").height(getHeight)
        function getHeight() {
            return $(window).height()-100;
        }
    })
</script>
</body>
</html>