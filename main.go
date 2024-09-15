package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	db = DBInit()
	setUpRoute(r)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
