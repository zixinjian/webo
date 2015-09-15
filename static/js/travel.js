/**
 * Created by rick on 15/9/15.
 */

function calPayment(){
    actualexpenses = $("#actualexpenses").val()
    expenses = $("#expenses").val()
    payment = getPayment(actualexpenses, expenses)
    $("#payment").val(payment)
    //console.log("blur", actualexpenses, expenses, payment)
}

function getPayment(actualexpenses, expenses){
    if (actualexpenses == 0 && expenses == 0){
        return 0
    }
    if (!actualexpenses || !expenses){
        return ""
    }
    diff = wbToMoney(parseFloat(actualexpenses) - parseFloat(expenses))
    if (diff == 0){
        return wbToMoney(expenses)
    }
    if (diff > 0){
        return wbToMoney(parseFloat(expenses) + parseFloat(diff * 0.3))
    }
    return wbToMoney(actualexpenses) - wbToMoney(diff * 0.3)
}