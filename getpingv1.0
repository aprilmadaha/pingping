package main

import (
	//	"fmt"

	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error
var remarkArray []Remark //备注结构体数组
var remarkIp []string    //备注结构体的IP
var device []Device
var getping []Getping

type Remark struct { //备注结构体
	gorm.Model
	Desc string //描述
	Ip   string //ip地址
}

type Getping struct { //ip结构体
	gorm.Model
	Tim  string //运行程序服务器/容器时间
	Ip   string //目的地址
	Xmt  string //发送次数
	Rcv  string //接收到的回包
	Loss string //丢包率
	Min  string //最小值
	Avg  string //平均值
	Max  string //最大值
}

type Alive struct {
	gorm.Model
	Tim  string  //运行程序服务器/容器时间
	Tar  string  //targets
	Ali  string  //alive
	Unr  string  //unreachable
	Uadd string  //unknown addresses
	Per  float32 //percentage百分比
}

type Delay struct {
	gorm.Model
	Tim string //运行程序服务器/容器时间
	Min string //min round trip time
	Avg string
	Max string
	Ela string //elapsed real time
}

type Device struct {
	Desc string //描述
	Ip   string //ip地址
	// Tim  string //运行程序服务器/容器时间
	// //	Dst  string //目的地址
	// Xmt  string //发送次数
	// Rcv  string //接收到的回包
	// Loss string //丢包率
	// Min  string //最小值
	Avg string //平均值
	//Max string //最大值
}

func main() {

	db, err = gorm.Open("mysql", "root:123456@(mysql)/ping?charset=utf8mb4&parseTime=True&loc=Local")
	//	db, err = gorm.Open("sqlite3", "../sql/ping.db") //打开数据库
	if err != nil {
		log.Fatal("db connect error")
	}
	defer db.Close()
	// ticker := time.NewTicker(time.Second * 60)
	// for {
	// 	select {
	// 	case <-ticker.C:

	//db, err := gorm.Open("mysql", "root:123456@(localhost)/ping?charset=utf8mb4&parseTime=True&loc=Local")
	currentTime := time.Now()
	time := currentTime.Format("2006-01-02 15:04")
	// db, err := gorm.Open("mysql", "root:123456@(mysql)/ping?charset=utf8mb4&parseTime=True&loc=Local")

	db.AutoMigrate(&Remark{}, &Getping{}, &Alive{}, &Delay{}, &Device{})
	db.Find(&remarkArray)              //获取数据
	for _, nume := range remarkArray { //过滤出数据库的ip列的数据
		remarkIp = append(remarkIp, nume.Ip)
	}

	argArray := make([]string, 0) //整合参数
	argArray = append(argArray, "-a")
	argArray = append(argArray, "-c")
	argArray = append(argArray, "3")
	argArray = append(argArray, "-s")
	argArray = append(argArray, remarkIp...)
	// argArray := []string{"-a", "-c", "3", "-g", "192.168.1.0/24", "-s"}

	cmd := exec.Command("fping", argArray...) //运行命令
	cmdValue, _ := cmd.CombinedOutput()       //输出运行的命令
	fmt.Println("string(cmdValue)", string(cmdValue))
	reIp := regexp.MustCompile(`(.*) +: xmt/rcv/%loss = (.*)/(.*)/(.*), min/avg/max = (.*)/(.*)/(.*)`) //过滤100%丢包的数据
	ipValueArray := reIp.FindAllStringSubmatch(string(cmdValue), -1)
	for _, ipValue := range ipValueArray { //把运行的结果放到数据库里面
		yyy := strings.TrimSpace(ipValue[1])
		db.Create(&Getping{Tim: time,
			Ip:   yyy,
			Xmt:  ipValue[2],
			Rcv:  ipValue[3],
			Loss: ipValue[4],
			Min:  ipValue[5],
			Avg:  ipValue[6],
			Max:  ipValue[7]})
		fmt.Println("ipvalue", ipValue, "time", time)
	}
	fmt.Println("time", time)

	//	db.Table("remarks").Where("getpings.tim = ?", time).Select("remarks.desc, remarks.ip, getpings.avg").Joins("JOIN getpings on getpings.ip = remarks.ip").Row(&device)

	reLive := regexp.MustCompile(`(.*) targets\n(.*) alive\n(.*) unreachable\n(.*) unknown addresses`) //提取存活个数
	liveValueArray := reLive.FindAllStringSubmatch(string(cmdValue), -1)
	for _, liveValue := range liveValueArray { //把运行的结果放到数据库里面
		liy := strings.Replace(liveValue[1], " ", "", -1)
		liy1 := strings.Replace(liveValue[2], " ", "", -1)
		liveValueTarFloat, _ := strconv.ParseFloat(liy, 32)
		liveValueAliFloat, _ := strconv.ParseFloat(liy1, 32)
		liveValuePer := float32(liveValueAliFloat / liveValueTarFloat)

		fmt.Println("liveValue[1],liveValue[2]", liveValue[1], liveValue[2])
		fmt.Println("liy,liy1", liy, liy1)
		fmt.Println("liveValuePer,liveValueTarFloat,liveValueAliFloat", liveValuePer, liveValueTarFloat, liveValueAliFloat)
		db.Create(&Alive{Tim: time,
			Tar:  liveValue[1],
			Ali:  liveValue[2],
			Unr:  liveValue[3],
			Uadd: liveValue[4],
			Per:  liveValuePer})
	}

	reDelay := regexp.MustCompile(`(.*) ms [(]min round trip time[)]\n(.*) ms [(]avg round trip time[)]\n(.*) ms [(]max round trip time[)]\n(.*) sec [(]elapsed real time[)]`) //提取延迟时间和程序运行时间
	delayValueArray := reDelay.FindAllStringSubmatch(string(cmdValue), -1)
	for _, delayValue := range delayValueArray { //把运行的结果放到数据库里面
		db.Create(&Delay{Tim: time,
			Min: delayValue[1],
			Avg: delayValue[2],
			Max: delayValue[3],
			Ela: delayValue[4]})
	}

	//db.Table("remarks").Select("remarks.desc, remarks.ip, getpings.avg").Joins("JOIN getpings on getpings.ip = remarks.ip").Create(&device)
	//	db.Table("remarkArray").Select("remarkArray.Desc, remarkArray.Ip,getping.Avg").Joins("left join getping on getping.Ip = remarkArray.Ip").Scan(&device)
	//fmt.Println("device", device)
	//	db.Table("remarks").Where("getpings.tim = ?", time).Select("remarks.desc, remarks.ip, getpings.avg").Joins("left JOIN getpings on getpings.ip = remarks.ip").Find(&device)
	//db.Table("remarks").Select("remarks.desc, remarks.ip, getpings.avg").Joins("left JOIN getpings on getpings.ip = remarks.ip").Scan(&device)
	//	db.Table("remarks").Where("getpings.tim = ?", time).Select("remarks.desc, remarks.ip, getpings.avg").Joins("JOIN getpings on getpings.ip = remarks.ip").Find(&device)
	//db.Table("remarks").Where("getpings.tim = ?", time).Select("remarks.desc, remarks.ip, getpings.avg").Joins("JOIN getpings on getpings.ip = remarks.ip").Create(&device)
	db.Table("remarks").Where("getpings.tim = ?", time).Select("remarks.desc, remarks.ip, getpings.avg").Joins("left join getpings on getpings.ip = remarks.ip").Scan(&device)
	fmt.Println("device", device)
	//db.Create(&device)
	for _, devicenum := range device {
		db.Create(&Device{
			Desc: devicenum.Desc,
			Ip:   devicenum.Ip,
			Avg:  devicenum.Avg,
		})
	}
}

// }

// }
