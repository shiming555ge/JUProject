默认管理员
	账号 admin123
	密码 12345678

测试环境
	go1.22.4
	windows 11 64位

默认使用session_id存储登录信息，因为api要求做了user_id传参，实际user_id参数无意义，使用apifox测试时需要在cookie中新增session_id项，值在登陆后自动存储，可以手动从浏览器复制