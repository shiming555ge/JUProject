package main

import (
	"errors"
	"gorm.io/gorm"
)

func SessionSave(session Session) {
	db.Create(&session)

	//  可插入session回收代码
}

func SessionExist(sessionID string) bool {
	var session Session
	result := db.Model(&Session{}).Where("session_id = ?", sessionID).First(&session)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func SessionGet(sessionID string) Session {
	var session Session
	db.Model(&Session{}).Where("session_id = ?", sessionID).First(&session)
	return session
}
