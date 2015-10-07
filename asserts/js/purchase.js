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