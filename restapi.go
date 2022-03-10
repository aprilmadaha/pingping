package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type Remark struct { //备注结构体
	//gorm.Model
	Desc string `json:"desc"`                   //描述
	Ip   string ` gorm:"primary_key" json:"ip"` //ip地址
}

func main() {

	db, err = gorm.Open("sqlite3", "./test.db")
	//	db, err = gorm.Open(sqlite.Open("test.db", &gorm.Conifig{})
	//	db, err := gorm.Open(sqlite.Open("test.db"), nil)
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("db connect error")
	}
	defer db.Close()
	db.AutoMigrate(&Remark{})

	r := gin.Default()
	r.GET("/users", index)
	r.GET("/users/:ip", show)
	r.POST("/ip", store)
	r.PUT("/users/:ip", update)
	r.DELETE("/users/:ip", destroy)
	_ = r.Run("8090")

}

func index(c *gin.Context) {
	var remark []Remark
	db.Find(&remark)
	c.JSON(200, remark)
}

func show(c *gin.Context) {
	ip := c.Params.ByName("ip")
	//ip := c.Param("ip")
	var remark []Remark
	//s := db.First(&remark, "ip=?", ip`)
	//db.Where("ip=?", ip).First(&remark)
	db.Where("ip LIKE ?", "%"+ip+"%").Find(&remark)
	// if remark.Ip == "" {
	// 	c.JSON(404, gin.H{"message": "ip not found"})
	// }
	// c.JSON(200, remark)
	fmt.Println("ip", ip, "remark", remark)
	//fmt.Println("s", s)
	c.JSON(200, remark)
	//	db.Select("ip").Find(&remark, ip)
	//c.JSON(200, remark)
}

func store(c *gin.Context) {
	var remark []Remark
	_ = c.BindJSON(&remark)
	// fmt.Println("remark", remark)
	//	db.Save(&remark)
	c.JSON(200, remark)
	//_ = c.BindJSON(&remark)
	//c.JSON(200, remark)
}

func update(c *gin.Context) {
	ip := c.Params.ByName("ip")
	var remark Remark
	db.First(&remark, ip)
	if remark.Ip == "" {
		c.JSON(404, gin.H{"message": "user not found"})
		return
	} else {
		_ = c.BindJSON(&remark)
		db.Save(&remark)
		c.JSON(200, remark)
	}
}

func destroy(c *gin.Context) {
	ip := c.Params.ByName("ip")
	var remark Remark
	db.First(&remark, ip)
	if remark.Ip == "" {
		c.JSON(404, gin.H{"message": "user not found"})
		return
	} else {
		_ = c.BindJSON(&remark)
		db.Delete(&remark)
		c.JSON(200, gin.H{"message": "delete success"})
	}
}
