/**
 * Created by rick on 15/9/15.
 */

function calPayment(){
    actualexpenses = $("#actualexpenses").val()
    expenses = $("#expenses").val()
    rat = $("#expayrat").val()
    if (isNaN(expenses) || rat == ""){
        alert("预计费用填写错误")
        return
    }
    if (isNaN(actualexpenses) || rat == ""){
        alert("实际费用填写错误")
        return
    }
    if (isNaN(rat) || rat == ""){
        alert("额外报销比例填写错误")
        return
    }
    rat = wbToMoney(rat)
    payment = getPayment(actualexpenses, expenses, rat)
    $("#payment").val(payment)
}
function getExPayRat(){
    rat = $("#expayrat").val()

}
function getPayment(actualexpenses, expenses, rat){
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
        return wbToMoney(parseFloat(expenses) + parseFloat(diff * rat/100))
    }
    return wbToMoney(actualexpenses) - wbToMoney(diff * rat/100)
}