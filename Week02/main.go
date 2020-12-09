package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	r := gin.Default()
	r.GET("/apis/v1/price", getUserInfo)
	log.Fatal(r.Run(":8080"))
}

func getUserInfo(ctx *gin.Context) {
	query := ctx.Query("code")
	_ = query
	prod, err := dao(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"user": prod.Price})
}

func dao(code string) (Product, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var p Product
	db.First(&p, "code = ?", code)
	return p, nil
}
