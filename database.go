package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

// 数据库初始化
func DBInit() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("default.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库打开失败")
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil
	}

	err = db.AutoMigrate(&Post{})
	if err != nil {
		return nil
	}

	err = db.AutoMigrate(&Session{})
	if err != nil {
		return nil
	}

	return db
}

// 数据结构定义
type User struct {
	UserID   uint64 `gorm:"primary_key" json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	UserType uint8  `json:"user_type"`

	Subscribed string `json:"subscribed"`
	Reported   string `json:"reported"`
	Written    string `json:"written"`
}

type Response struct {
	Code uint16      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Session struct {
	SessionID string `gorm:"primary_key" json:"session_id"`
	UserID    uint64 `json:"user_id"`
	UserType  uint8  `json:"user_type"`
	Time      time.Time
}

type Post struct {
	PostID     uint64    `gorm:"primary_key" json:"post_id"`
	UserID     uint64    `json:"user_id"`
	Content    string    `json:"content"`
	Time       time.Time `json:"time"`
	Counter    int16     `json:"counter"`
	Subscribed string    `json:"subscribed"`
	Reported   string    `json:"reported"`
}
