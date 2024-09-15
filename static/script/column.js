// 文章被选中效果
$(document).ready(()=>{
    $(".mainPage>.column:first").addClass("columnSelected")
    $(".mainPage>.column").click(function () {
        $(".mainPage>.column").removeClass("columnSelected")
        $(this).addClass("columnSelected")
    })
})