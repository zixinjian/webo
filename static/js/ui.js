/**
 * Created by rick on 15/7/19.
 */
function hideAlert(){
    $(".alert").hide()
}
function showSuccess(tip){
    showAlert("success", tip)
}
function showError(tip){
    showAlert("danger", tip)
}
function showAlert(type, tip){
    $(".alert").addClass("alert-"+type)
    $(".alert").text(tip)
    $(".alert").show()
}