package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func setUpRoute(r *gin.Engine) {
	r.Use(XSSMiddleware)

	// 判断状态
	r.GET("/", authMiddleware, func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/index")
	})

	// 登录
	r.GET("/login", authMiddlewareLogin, func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	// 注册
	r.GET("/register", authMiddlewareLogin, func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})

	// 主页
	r.GET("/index", authMiddleware, func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// 用户主页 , authMiddleware
	r.GET("/user", authMiddleware, func(c *gin.Context) {
		c.HTML(http.StatusOK, "user.html", gin.H{
			"admin": userInfo(c).UserType,
		})
	})

	r.GET("/edit-post/:post_id", authMiddleware, func(c *gin.Context) {
		postID := c.Param("post_id")
		p, err := postInfo(postID)
		if err != nil || p.UserID != userInfo(c).UserID {
			c.Abort()
			return
		}

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"mod":     0,
			"id":      postID,
			"content": p.Content,
		})
	})

	r.GET("/create-post", authMiddleware, func(c *gin.Context) {
		c.HTML(http.StatusOK, "edit.html", gin.H{
			"mod": 1,
		})
	})

	r.POST("/api/user/reg", authMiddlewareLogin, func(c *gin.Context) {
		user := User{}

		// 获取数据并报错
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, &Response{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
				Data: nil,
			})
			return
		}

		c.JSON(register(user))
	})

	r.GET("/post/:post_id", func(c *gin.Context) {
		postID := c.Param("post_id")

		p, err := postInfo(postID)

		if err != nil {
			c.JSON(http.StatusBadRequest, &Response{
				Code: http.StatusBadRequest,
				Msg:  "post id error",
			})
		}

		c.JSON(http.StatusOK, Response{
			Code: http.StatusOK,
			Data: p,
		})
	})

	r.GET("/report-post/:post_id", func(c *gin.Context) {
		postID := c.Param("post_id")

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"mod": 2,
			"id":  postID,
		})
	})

	r.GET("/api/user/info", authMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, userInfo(c))
	})

	r.POST("/api/user/login", authMiddlewareLogin, func(c *gin.Context) {
		c.JSON(login(c))
	})

	r.POST("/api/student/post", authMiddleware, func(c *gin.Context) {
		c.JSON(sendPost(c))
	})

	r.GET("/api/student/post", authMiddleware, func(c *gin.Context) {
		n, e := c.GetQuery("amount")
		t, e1 := c.GetQuery("time")
		if !(e && e1) {
			c.JSON(getPost("-1", time.Now()))
		} else {
			t, _ := time.Parse("2006-01-02 15:04:05", t)
			c.JSON(getPost(n, t))
		}

	})

	r.DELETE("/api/student/post", authMiddleware, func(c *gin.Context) {
		c.JSON(delPost(c))
	})

	r.PUT("/api/student/post", authMiddleware, func(c *gin.Context) {
		c.JSON(revisePost(c))
	})

	r.POST("/api/student/report-post", authMiddleware, func(c *gin.Context) {
		c.JSON(reportPost(c))
	})

	r.GET("/api/student/report-post", authMiddleware, func(c *gin.Context) {
		c.JSON(getReportPost(c))
	})

	r.GET("/api/admin/report", authMiddlewareAdmin, func(c *gin.Context) {
		c.JSON(AdminGetReportPost())
	})

	r.POST("/api/admin/report", authMiddlewareAdmin, func(c *gin.Context) {
		c.JSON(AdminReportPost(c))
	})
	// 404 重定向
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})
}

func authMiddleware(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")

	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	} else if !SessionExist(sessionID) {
		c.SetCookie("session_id", "", -1, "/", "", false, false)
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	c.Next()
}

func authMiddlewareAdmin(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")

	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	} else if !SessionExist(sessionID) {
		c.SetCookie("session_id", "", -1, "/", "", false, false)
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	} else if SessionGet(sessionID).UserType != 2 {
		c.Redirect(http.StatusForbidden, "/index")
		c.Abort()
		return
	}

	c.Next()
}

func authMiddlewareLogin(c *gin.Context) {
	_, err := c.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		c.Next()
	} else {
		c.Redirect(http.StatusFound, "/index")
	}
}

func XSSMiddleware(c *gin.Context) {
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Next()
}
