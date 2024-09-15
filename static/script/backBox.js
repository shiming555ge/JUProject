// 弹幕特效
$(document).ready(()=>{
    $(".backBox>.track").each(function () {
        $(this).css({"animation": "trackScroll "+ (Math.random()*6+4).toFixed(1).toString() +"s linear infinite"})
        $(this).css({"top": (Math.random()*95).toFixed(1) + "%"})
    })
})