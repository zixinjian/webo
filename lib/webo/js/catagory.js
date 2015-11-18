/**
 * Created by rick on 15/11/13.
 */


cateDict = {}
function autoFillName($name, changePower){
    var selectCate = $('#category').val()
    if (!(selectCate in cateDict)){
        return
    }
    option = cateDict[selectCate]
    if(cateDict[selectCate][0]){
        $name.val(cateDict[selectCate][1]);
        $name.attr("readonly", true)
        if(changePower){
            $power.val("")
        }
    }else{
        $name.attr("readonly", false)
        $name.val("")
        if(changePower){
            $power.val("0")
        }
    }
}
function initCategory($name){
    $('#category').find("option").each(function () {
        $option = $(this)
        value = $option.val()
        label = $option.text()
        cateDict[value]= [$option.attr("data-wb-a-flag")== "yes", label]
    })
    autoFillName($name)
    $('#category').change(function(){
        autoFillName($name, true)
    })
}