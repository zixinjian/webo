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

function layoutAutoHeight(){
    $.each($("[layout-auto-height]"), function(){
        var outHeight = $(this).attr("layout-auto-height")
        //console.log("outHeight", outHeight)
        $(this).height($(window).height() + parseInt(outHeight))
    });
}

//$(function(){
//    $(window).resize(function () {
//        if ($table){
//            $table.bootstrapTable('resetView', {
//                height: getHeight()
//            });
//        }
//    });
//});

function getTableHeight() {
    return $(window).height();
}