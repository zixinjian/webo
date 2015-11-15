/**
 * Created by rick on 15/11/11.
 */

var $power = $("#power")

function addSupplier(item){
    suppliers = []
    $("#supplierList").find('[data-wb-c="supplierSn"]').each(function(idx, sitem){
        suppliers.push($(sitem).val())
    });
    if(suppliers.indexOf(item.sn) >=0){
        return
    }
    supplierBtn = '<div class="btn-group dropdown supplierBtnGroup">'
    + wbSprintf('<a class="btn m-b-xs btn-sm btn-default" data-toggle="dropdown" aria-expanded="false">%s<span class="caret"></span>', item.keyword)
    + wbSprintf('<input type="hidden" data-wb-c="supplierSn" name="supplierlist", value="%s">', item.sn)
    + wbSprintf('</a><ul class="dropdown-menu"><li><a id="deleteSupplier_%s">删除</a></li></ul></div>', item.sn)
    $("#supplierList").append(supplierBtn)
    $("#deleteSupplier_" + item.sn).on("click", function(evt){
        $(evt.target).parents(".supplierBtnGroup").remove()
    })
}
$(function() {
    $("#power").wrapAll('<div class="input-group"></div>')
    $("#power").after('<span class="input-group-addon">KW</span>')
    $('#file_upload').uploadify({
        'swf'      : '../../lib/3rd/uploadify/uploadify.swf',
        'uploader' : '/item/upload/product?sn=' + $("#sn").val(),
        'cancelImg': '../../lib/3rd/uploadify/uploadify-cancel.png',
        'fileObjName':'uploadFile'
    });

    $("#supplier_key").autocomplete({
        source: "/item/autocomplete/supplier",
        autoFocus:false,
        focus: function( event, ui ) {
            return false;
        },
        minLength: 1,
        select: function(event, ui) {
            $("#supplier_key").val("")
            if(ui && ui.item){
                addSupplier(ui.item)
            }
            return true;
        }
    }).autocomplete( "instance" )._renderItem = function( ul, item ) {
        return $( "<li>" )
            .append(item.keyword + "(" + item.name + ")")
            .appendTo( ul );
    };
});

function calRetailPrice(){
    clearInputError("#price")
    clearInputError("#profitrat")
    price = $("#price").val()
    rat = $("#profitrat").val()
    if (isNaN(price) || price == ""){
        showInputError("#price", "参考价格填写错误")
        return
    }
    if (isNaN(rat) || rat == ""){
        showInputError("#profitrat", "调价比例填写错误")
        return
    }
    rat = wbToMoney(rat)
    retailprice = wbToMoney(parseFloat(price) * (100+rat)/100)
    $("#retailprice").val(retailprice)
}