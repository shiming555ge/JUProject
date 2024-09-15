$(document).ready(()=>{
    $("#register").click(()=>{
        let username = $("#username").val()
        let password = $("#password").val()
        let confirm = $("#confirm").val()

        if (password != confirm) {
            alert("两次密码不匹配")
            $("#password").val("")
            $("#confirm").val("")
        }

        $.ajax({
            url: "/api/user/reg",
            type: "post",
            data: JSON.stringify({
                username: username,
                password: password,
                name: username,
                user_type: 1
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