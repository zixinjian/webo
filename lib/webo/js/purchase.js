/**
 * Created by rick on 15/10/6.
 */

function rowStyleOvertime(row, index) {
    if(row.godowndate == ""){
        godownDate = moment()
    }else{
        godownDate = moment(row.godowndate, "YYYY.MM.DD")
    }
    requiredDate = moment(row.requireddate, "YYYY.MM.DD")
    if (requiredDate.diff(godownDate, "days") < 0){
        return {classes: "danger"};
    }
    if (requiredDate.diff(godownDate, "days") < 2 ){
        return {classes: "warning"};
    }
    return {}
}

function supplierFormatter(value, row){
    return wbSprintf('<span title="%s(%s)">%s</span>', value, row.suppliername, value)
}
function modelFormatter(value, row){
    if (row.product == ""){
        return value
    }
    return wbSprintf('<a class="text-info-dker" href="/static/files/product/%s" target="_blank" title="附件" data-toggle="poplayer" data-placement="bottom" data-url="/static">%s</a>', row.product, value)
}


