$(document).ready(()=>{
    $("#login").click(()=>{
        let username = $("#username").val()
        let password = $("#password").val()
        let check = $(this).children("check").val()
        $.ajax({
            url: "/api/user/login",
            type: "post",
            data: JSON.stringify({
                username: username,
                password: password
            }),
            dataType: "JSON",
            contentType: "application/json",
            success: function () {
                window.location.href = "/index"
            },
            error: function (data) {
                let d = eval("("+data.responseText+")")
                alert(d.msg)
            }
        })
    })
})