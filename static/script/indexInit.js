const now = new Date();

const year = now.getFullYear();
const month = ('0' + (now.getMonth() + 1)).slice(-2);
const day = ('0' + now.getDate()).slice(-2);
const hours = ('0' + now.getHours()).slice(-2);
const minutes = ('0' + now.getMinutes()).slice(-2);
const seconds = ('0' + now.getSeconds()).slice(-2);

const formattedTime = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}+08:00`

$(document).ready(()=>{
    $.ajax({
        url: "/api/student/post",
        type: "get",
        data: JSON.stringify({
            amount: 10,
            time: formattedTime
        }),
        success: function (data) {
            pl = eval("("+JSON.stringify(data)+").data.post_list")

            for (let i of pl) {
                $(".mainPage").append("<div class=\"column\" id='"+i.post_id+"'data-content=\""+i.content+"\">\n" +
                    "<div class=\"title\">"+i.content.split('\n')[0]+"</div>\n" +
                    "<div class=\"time\">"+i.time.slice(0,19)+"</div>\n" +
                    "<div class=\"author\">uid "+i.user_id+"</div>\n" +
                    "</div>")
            }

            $(".column").click(function (){
                page = 0
                pages = [$(this).attr("data-content")]

                $(".columnSelected").removeClass("columnSelected")
                $(this).addClass("columnSelected")

                jump(page)
            })

            pages.push($(".mainPage>.column:first-child").attr("data-content"))
            $(".mainPage>.column:first-child").addClass("columnSelected")
            PageReady()
        }
    })
    $(".operator>.subscribe").click(()=>{
        alert("下次一定")
    })
    $(".operator>.report").click(()=>{
        window.location.href = "/report-post/"+$(".columnSelected:first").attr("id")
    })
})