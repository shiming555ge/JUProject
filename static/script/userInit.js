function hover (t) {
    $(t).children(".reason").show()
    $(t).children(".title").hide()
}
function out (t) {
    $(t).children(".title").show()
    $(t).children(".reason").hide()
}

$(document).ready(()=>{
    // 切换项目
    $(".optBox").hide()
    $(".myPosts").show()

    $(".optList>*").click(function (){
        $(".optBox").hide()
        $($(this).attr("data-page")).show()
        $(".optList>.columnSelected").removeClass("columnSelected")
        $(this).addClass("columnSelected")
    })
    //
    $(".exitBtn").click(()=>{
        document.cookie= "session_id=; expires=Thu, 01 Jan 1970 00:00:00 GMT"
        window.location.reload()
    })

// 初始化页面信息
$.ajax({
    url: "/api/user/info",
    type: "get",
    async: false,
    success: function (info) {
        $.ajaxSettings.async = false;
        // 我的帖子
        wr = info.written

        for (let i of wr.split(";")) {
            if (i == "") break
            $.get("/post/" + i, function (data) {
                $(".myPosts").append("<div class=\"column\" id='"+data.data.post_id+"' data-content='"+data.data.content+"'>\n" +
                    "<div class=\"title\">"+data.data.content.split('\n')[0]+"</div>\n" +
                    "<div class=\"time\">"+data.data.time.slice(0,19)+"</div>\n" +
                    "<img class=\"postBtn\" id='editBtn' src=\"/static/src/edit.png\" title='修改该帖' style='right: 2.5rem'>\n" +
                    "<img class=\"postBtn\" id='delBtn' src=\"/static/src/delete.png\" title='删除该帖'>\n" +
                    "</div>")
            })
        }
        // 我的举报
        re = info.reported
        for (let i of re.split(";")) {
            if (i == "") break
            $.get("/post/" + i.split(":")[0], (data) => {
                $(".myReport").append("<div class=\"column\" onmouseover='hover(this)' onmouseout='out(this)' id='"+data.data.post_id+"' data-content='"+data.data.content+"'>\n" +
                    "<div class=\"title\">"+data.data.content.split('\n')[0]+"</div>\n" +
                    "<div class=\"reason\">"+i.split(':')[1]+"</div>\n" +
                    "<div class=\"author\">uid "+data.data.user_id+"</div>\n" +
                    "</div>")
            })
        }

        // 管理举报
        if (admin == 2) {
        $.get("/api/admin/report",  (data)=>{
            pl = data.data.report_list
            for (let i of pl) {
                $(".adminReport").append("<div class=\"column\" onmouseover='hover(this)' onmouseout='out(this)' id='"+i.post_id+"' data-content='"+i.content+"'>\n" +
                    "<div class=\"title\">"+i.content.split('\n')[0]+"</div>\n" +
                    "<div class=\"reason\">"+i.reason+"</div>\n" +
                    "<img class=\"postBtn\" id='igBtn' src=\"/static/src/delete.png\" title='不认可' style='right: 2rem'>\n" +
                    "<img class=\"postBtn\" id='passBtn' src=\"/static/src/pass.png\" title='认可'>\n" +
                    "</div>")
            }
        })
        }

        // 我的收藏
        re = info.subscribed
        for (let i of re.split(";")) {
            if (i == "") break
            $.get("/post/" + i,  (data) => {
                $(".mySubscribe").append("<div class=\"column\" onmouseover='hover(this)' onmouseout='out(this)' id='"+data.data.post_id+"' data-content='"+data.data.content+"'>\n" +
                    "<div class=\"title\">"+data.data.content.split('\n')[0]+"</div>\n" +
                    "<div class=\"author\">uid "+data.data.user_id+"</div>\n" +
                    "</div>")

            })
        }
        $.ajaxSettings.async = true;

        // 我的信息
        $(".myInfo>.user_id>.title").append(info.user_id)
        $(".myInfo>.username>.title").append(info.username)
        $(".myInfo>.name>.title").append(info.name)
        $(".myInfo>.password>.title").append(info.password)
    }
})
    // 初始化卡片
    pages.push($(".myPosts>.column:first-child").attr("data-content"));
    page = 0
    PageReady()

    $(".myPosts>.column:first-child").addClass("columnSelected")
    $(".optBox>.columnSelected>.postBtn").show()

    $(".optBox:not(.myInfo)>.column").click(function (){
        pages = [$(this).attr("data-content")]
        page = 0
        jump(0)

        $(".optBox>.columnSelected>.postBtn").hide()
        $(".optBox>.columnSelected").removeClass("columnSelected")
        $(this).addClass("columnSelected")
        $(".optBox>.columnSelected>.postBtn").show()
    })

    $(".adminReport>.column>#editBtn").click(function (){
        window.location.href = "/edit-post/" + $(this).parent().attr("id")
    })
    $(".adminReport>.column>#delBtn").click(function (){
        $.ajax({
            url: "/api/student/post?post_id=" + $(this).parent().attr("id"),
            type: "delete",
            success: function(data){
                alert("删除成功")
                window.location.reload()
            },
            error: function(data){
                alert(data.msg)
            }
        })
    })
    $(".adminReport>.column>#igBtn").click(function (){
        $.post("/api/admin/report", JSON.stringify({
            user_id: 0,
            post_id: parseInt($(this).parent().attr("id")),
            approval: 2
        }), function (data){
            window.location.reload()
        },"JSON")
    })
    $(".adminReport>.column>#passBtn").click(function (){
        $.post("/api/admin/report", JSON.stringify({
            user_id: 0,
            post_id: parseInt($(this).parent().attr("id")),
            approval: 1
        }), function (data){
            window.location.reload()
        },"JSON")
    })
})