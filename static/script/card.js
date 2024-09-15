// 跳转函数
function jump(i) {
    if (pages.length == 0 || pages[0] == "") {
        $(".operator>*").hide()
        return
    }

    // 检查下标越界
    while (i+1>pages.length){
        i--
    }
    while (i < 0){
        i++
    }

    $(".card>.textBox").hide()

    let tb = $(".textBox:first")
    let hei = $(".card:first").height()
    let rem = $(".operator>.leftA:first").height() / 2

    // 赋值
    tb.html(pages[i].replaceAll("\n","<br>"))

    // 判断是否需要修改字数
    if (i+1 == pages.length) {
        // 需要
        if (tb.height() + 3 * rem > hei){
            pages.push("")
        }
        // 修改字数
        while ( tb.height() + 3 * rem > hei ) {
            if (tb.html().slice(-4) == "<br>"){
                pages[pages.length-1] = "\n" + pages[pages.length-1]
                tb.html(tb.html().slice(0, -4))
            }
            else {
                pages[pages.length - 1] = tb.html().slice(-1) + pages[pages.length - 1]
                tb.html(tb.html().slice(0, -1))
            }
        }
        pages[i] = tb.html().replaceAll("<br>","\n")
    }

    // operator 修改
    $(".operator>*").show()
    $(".operator>.report").hide()
    $(".operator>.subscribe").hide()
    $(".commentBox").hide()

    if (i+1 == pages.length) {
        $(".operator>.rightA").hide()
    }
    if (i == 0) {
        $(".operator>.leftA").hide()
    }

    $(".card>.textBox").show()

    return i;
}

function PageReady() {
    // 初始化card内容
    jump(page)

    // 事件注册
    $(".operator>.rightA").click(()=>{
        page = jump(page+1)
    })
    $(".operator>.leftA").click(()=>{
        page = jump(page-1)
    })

    function upEnent(){
        $(".textBox").hide()
        $(".operator>*").hide()
        $(".operator>.report").show()
        $(".operator>.subscribe").show()
        $(".operator>.downA").click(()=>{
            jump(page)
            $(".operator>.downA").click(downEvent)
        }).show()
    }
    function downEvent() {
        $(".textBox").hide()
        $(".operator>*").hide()
        $(".card>.commentBox").show()
        $(".operator>.upA").click(()=>{
            jump(page)
            $(".operator>.upA").click(upEnent)
        }).show()
    }

    $(".operator>.upA").click(upEnent)
    $(".operator>.downA").click(downEvent)
}