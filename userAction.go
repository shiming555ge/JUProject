package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 注册
func register(user User) (int, Response) {

	c := strings.Count(user.Username, "")
	// 用户名过长过短错误
	if c < 7 || c > 12 {
		return http.StatusUnauthorized, Response{
			Code: http.StatusUnauthorized,
			Msg:  "用户名过长或过短",
			Data: nil,
		}
	}

	symbols := "<>\"'\\/();="
	// 用户名存在特殊字符错误
	for _, char := range user.Username {
		if strings.Contains(symbols, string(char)) {
			return http.StatusUnauthorized, Response{
				Code: http.StatusUnauthorized,
				Msg:  "用户名存在特殊字符",
				Data: nil,
			}
		}
	}

	// 用户名已存在错误
	var count int64
	db.Model(&User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return http.StatusUnauthorized, Response{
			Code: http.StatusUnauthorized,
			Msg:  "用户名已存在",
			Data: nil,
		}
	}

	// 设置userid
	var lUser User
	result := db.Model(&User{}).Order("user_id desc").First(&lUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user.UserID = 1
	} else {
		user.UserID = lUser.UserID + 1
	}

	// 创建记录
	db.Create(&user)
	return http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  "注册成功",
		Data: nil,
	}
}

// 登录
func login(c *gin.Context) (int, Response) {
	var user, cUser User

	// 获取数据并报错
	err := c.BindJSON(&user)
	if err != nil {
		return http.StatusBadRequest, Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
			Data: nil,
		}
	}

	// 查询记录
	result := db.Where("username = ? AND password = ?", user.Username, user.Password).First(&user)
	// 记录不存在
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return http.StatusUnauthorized, Response{
			Code: http.StatusUnauthorized,
			Msg:  "密码错误或者账号不存在",
			Data: nil,
		}
	}

	// 登录成功
	sessionID := strings.ToLower(uuid.New().String())
	SessionSave(Session{
		SessionID: sessionID,
		UserID:    user.UserID,
		UserType:  user.UserType,
		Time:      time.Now(),
	})

	c.SetCookie("session_id", sessionID, 60*60*24*7, "/", "", false, false)

	return http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  "登陆成功",
		Data: &gin.H{
			"user_id":   cUser.UserID,
			"user_type": cUser.UserType,
		},
	}

}

// 发布帖子
func sendPost(c *gin.Context) (int, Response) {
	var p struct {
		UserID  uint64 `json:"user_id"`
		Content string `json:"content"`
	}
	err := c.BindJSON(&p)
	sessionID, _ := c.Cookie("session_id")
	p.UserID = SessionGet(sessionID).UserID

	if err != nil {
		return http.StatusBadRequest, Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
			Data: nil,
		}
	}

	var lPost Post
	result := db.Model(&Post{}).Order("post_id desc").First(&lPost)
	var PostID uint64
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		PostID = 1
	} else {
		PostID = lPost.PostID + 1
	}

	result = db.Create(&Post{
		PostID:  PostID,
		UserID:  p.UserID,
		Content: p.Content,
		Time:    time.Now(),
	})

	if result.Error != nil {
		return http.StatusBadRequest, Response{
			Code: http.StatusBadRequest,
			Msg:  result.Error.Error(),
			Data: nil,
		}
	}

	var u User
	result = db.Model(&User{}).Where("user_id = ?", p.UserID).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return http.StatusUnauthorized, Response{
			Code: http.StatusUnauthorized,
			Msg:  "user_id不存在",
		}
	}
	u.Written += fmt.Sprintf("%d;", PostID)
	db.Save(&u)

	return http.StatusOK, Response{
		Code: http.StatusOK,
		Msg:  "发送成功",
		Data: nil,
	}
}

// 获取帖子
func getPost(s string, t time.Time) (int, Response) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return http.StatusBadRequest, Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
			Data: nil,
		}
	}

	rows, err := db.Model(&Post{}).Order("time desc").Where("time < ? and counter >= 0", t).Select("post_id", "user_id", "content", "time").Limit(n).Rows()
	if err != nil {
		return http.StatusBadRequest, Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
			Data: nil,
		}
	}

	var posts []struct {
		PostID  uint64    `json:"post_id"`
		UserID  uint64    `json:"user_id"`
		Content string    `json:"content"`
		Time    time.Time `json:"time"`
	}

	defer rows.Close()

	for rows.Next() {
		var (
			postID  uint64
			userID  uint64
			content string
			t1      time.Time
		)
		err := rows.Scan(&postID, &userID, &content, &t1)
		if err != nil {
			return http.StatusBadRequest, Response{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			}
		}

		posts = append(posts, struct {
			PostID  uint64    `json:"post_id"`
			UserID  uint64    `json:"user_id"`
			Content string    `json:"content"`
			Time    time.Time `json:"time"`
		}{postID, userID, content, t1})
	}

	return http.StatusOK, Response{
		Code: http.StatusOK,
		Data: gin.H{
			"post_list": posts,
		},
		Msg: "读取成功",
	}
}

// 删除帖子
func delPost(c *gin.Context) (int, Response) {
	sessionID, _ := c.Cookie("session_id")
	u := SessionGet(sessionID)

	PostID, e := c.GetQuery("post_id")
	if !e {
		return http.StatusBadRequest, Response{
			Code: http.StatusBadRequest,
			Msg:  "缺少post_id",
			Data: nil,
		}
	}

	var (
		result *gorm.DB
		p      Post
	)

	if u.UserType == 1 {
		result = db.Model(&Post{}).Where("user_id = ? AND post_id = ?", u.UserID, PostID).First(&p)
	} else {
		result = db.Model(&Post{}).Where("post_id = ?", PostID).First(&p)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return http.StatusUnauthorized, Response{
			Code: http.StatusUnauthorized,
			Msg:  "post_id不存在",
			Data: nil,
		}
	}

	p.Counter = -1
	p.Reported = ""

	db.Save(&p)

	return http.StatusOK, Response{
		Code: http.StatusOK,
	}
}

func reportPost(c *gin.Context) (int, Response) {
	d := struct {
		UserID uint64 `json:"user_id"`
		PostID uint64 `json:"post_id"`
		Reason string `json:"reason"`
	}{}
	err := c.BindJSON(&d)
	sessionID, _ := c.Cookie("session_id")
	d.UserID = SessionGet(sessionID).UserID

	if err != nil {
		return http.StatusBadRequest, Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
			Data: nil,
		}
	}

	// 获取原帖
	var p Post
	result := db.Model(&Post{}).Where("post_id = ?", d.PostID).First(&p)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return http.StatusUnauthorized, Response{
			Code: http.StatusUnauthorized,
			Msg:  "post_id不存在",
		}
	}

	// 检测二次举报
	for _, i := range strings.Split(p.Reported, ";") {
		if strconv.Itoa(int(d.UserID)) == strings.Split(i, ":")[0] {
			return http.StatusUnauthorized, Response{
				Code: http.StatusUnauthorized,
				Msg:  "你已经举报过了",
			}
		}
	}

	// 添加到用户举报
	var u User
	result = db.Model(&User{}).Where("user_id = ?", d.UserID).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return http.StatusUnauthorized, Response{
			Code: http.StatusUnauthorized,
			Msg:  "user_id不存在",
		}
	}
	u.Reported += fmt.Sprintf("%d:%s;", d.PostID, d.Reason)
	db.Save(&u)

	p.Counter += 1
	p.Reported += fmt.Sprintf("%d:%s;", d.UserID, d.Reason)
	db.Save(&p)
	return http.StatusOK, Response{
		Code: http.StatusOK,
	}
}

func revisePost(c *gin.Context) (int, Response) {
	var d struct {
		UserID  uint64 `json:"user_id"`
		PostID  uint64 `json:"post_id"`
		Content string `json:"content"`
	}

	err := c.BindJSON(&d)
	sessionID, _ := c.Cookie("session_id")
	d.UserID = SessionGet(sessionID).UserID

	if err != nil {
		return http.StatusBadRequest, Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
			Data: nil,
		}
	}

	var p Post

	result := db.Model(&Post{}).Where("post_id = ? and user_id = ?", d.PostID, d.UserID).First(&p)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return http.StatusUnauthorized, Response{
			Code: http.StatusUnauthorized,
			Msg:  "帖子不存在或无权修改",
		}
	}

	p.Content = d.Content
	db.Save(&p)

	return http.StatusOK, Response{
		Code: http.StatusOK,
	}
}

func getReportPost(c *gin.Context) (int, Response) {
	UserID, e := c.GetQuery("user_id")
	if !e {
		return http.StatusBadRequest, Response{
			Code: http.StatusBadRequest,
			Msg:  "缺少必要参数user_id",
		}
	}

	var rl string

	var pl []struct {
		PostID  uint64 `json:"post_id"`
		Content string `json:"content"`
		Reason  string `json:"reason"`
		Status  int    `json:"status"`
	}

	var p Post
	result := db.Model(&User{}).Where("user_id = ?", UserID).Select("reported").Scan(&rl)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return http.StatusUnauthorized, Response{
			Code: http.StatusUnauthorized,
			Msg:  "user_id不存在",
		}
	}

	for _, i := range strings.Split(rl, ";") {
		if i == "" {
			break
		}
		t := strings.Split(i, ":")
		result := db.Model(&Post{}).Where("post_id = " + t[0]).First(&p)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			continue
		}

		if p.Counter == -1 {
			pl = append(pl, struct {
				PostID  uint64 `json:"post_id"`
				Content string `json:"content"`
				Reason  string `json:"reason"`
				Status  int    `json:"status"`
			}{
				PostID:  p.PostID,
				Content: p.Content,
				Reason:  t[1],
				Status:  1,
			})
		} else if p.Counter == 0 {
			pl = append(pl, struct {
				PostID  uint64 `json:"post_id"`
				Content string `json:"content"`
				Reason  string `json:"reason"`
				Status  int    `json:"status"`
			}{
				PostID:  p.PostID,
				Content: p.Content,
				Reason:  t[1],
				Status:  2,
			})
		} else {
			pl = append(pl, struct {
				PostID  uint64 `json:"post_id"`
				Content string `json:"content"`
				Reason  string `json:"reason"`
				Status  int    `json:"status"`
			}{
				PostID:  p.PostID,
				Content: p.Content,
				Reason:  t[1],
				Status:  0,
			})
		}
	}
	return http.StatusOK, Response{
		Code: http.StatusOK,
		Data: gin.H{
			"report_list": pl,
		},
	}
}

func AdminGetReportPost() (int, Response) {
	var pl []struct {
		Username string `json:"username"`
		PostID   uint64 `json:"post_id"`
		Content  string `json:"content"`
		Reason   string `json:"reason"`
	}

	rows, _ := db.Model(&Post{}).Where("counter > 0").Order("counter desc").Select("post_id", "content", "reported").Rows()

	defer rows.Close()

	for rows.Next() {
		var (
			postID   uint64
			username string
			content  string
			reason   string

			rl string
		)

		err := rows.Scan(&postID, &content, &rl)
		if err != nil {
			return http.StatusBadRequest, Response{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			}
		}

		for _, i := range strings.Split(rl, ";") {
			if i == "" {
				break
			}
			t := strings.Split(i, ":")
			username += t[0] + ";"
			reason += t[1] + ";"
		}

		pl = append(pl, struct {
			Username string `json:"username"`
			PostID   uint64 `json:"post_id"`
			Content  string `json:"content"`
			Reason   string `json:"reason"`
		}{
			Username: username,
			PostID:   postID,
			Content:  content,
			Reason:   reason,
		})

	}
	return http.StatusOK, Response{
		Code: http.StatusOK,
		Data: gin.H{
			"report_list": pl,
		},
	}
}

func AdminReportPost(c *gin.Context) (int, Response) {
	var d struct {
		UserID   uint64 `json:"user_id"`
		PostID   uint64 `json:"post_id"`
		Approval uint64 `json:"approval"`
	}

	err := c.BindJSON(&d)
	if err != nil {
		return http.StatusBadRequest, Response{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}
	}

	sessionID, _ := c.Cookie("session_id")
	u := SessionGet(sessionID)
	d.UserID = u.UserID

	var p Post

	db.Model(&Post{}).Where("post_id = ?", d.PostID).First(&p)

	switch d.Approval {
	case 1:
		p.Counter = -1
		//var u User
		//db.Model(&User{}).Where("user_id = ?", d.UserID).First(&u)
		//u.Reported += fmt.Sprintf("%d:管理权限;", d.PostID)
		//db.Save(&u)
		break
	case 2:
		p.Counter = 0
		p.Reported = ""
		break
	}
	db.Save(&p)

	return http.StatusOK, Response{
		Code: http.StatusOK,
	}
}

func userInfo(c *gin.Context) User {
	sessionID, _ := c.Cookie("session_id")

	var u User
	db.Model(&User{}).Where("user_id = ?", SessionGet(sessionID).UserID).First(&u)

	return u
}

func postInfo(postID string) (Post, error) {
	var p Post
	result := db.Model(&Post{}).Where("post_id = " + postID).First(&p)

	fmt.Println(p)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Post{}, result.Error
	}

	return p, nil

}
