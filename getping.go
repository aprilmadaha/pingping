package main

import (
	"log"
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os/exec"
	"regexp"
)

var db *gorm.DB
var err error
var remarkArray []Remark //备注结构体数组
var remarkIp []string    //备注结构体的IP

type Remark struct { //备注结构体
	gorm.Model
	Desc string //描述
	Ip   string //ip地址
}

type Getping struct { //ip结构体
	gorm.Model
	Tim  string //运行程序服务器/容器时间
	Dst  string //目的地址
	Xmt  string //发送次数
	Rcv  string //接收到的回包
	Loss string //丢包率
	Min  string //最小值
	Avg  string //平均值
	Max  string //最大值
}

type Alive struct {//ping结果汇总结构体
	gorm.Model
	Tim  string //运行程序服务器/容器时间
	Tar  string //targets 总共有多少IP
	Ali  string //alive 存活的IP主机数
	Unr  string //unreachable 不能到达的
	Uadd string //unknown addresses 不能识别的
}

type Delay struct {//延迟的汇总结构体
	gorm.Model
	Tim string //运行程序服务器/容器时间
	Min string //min round trip time 最小延迟
	Avg string
	Max string
	Ela string //elapsed real time 运行程序的时间fping
}

func main() {
	ticker := time.NewTicker(time.Second * 60)//每60秒执行一次
	for {
		select {
		case <-ticker.C:

			//db, err := gorm.Open("mysql", "root:123456@(localhost)/ping?charset=utf8mb4&parseTime=True&loc=Local")
			currentTime := time.Now()
			db, err := gorm.Open("mysql", "root:123456@(ping-mysql)/ping?charset=utf8mb4&parseTime=True&loc=Local")
			defer db.Close()
			if err != nil {
				log.Fatal("db connect error")
			}
			db.AutoMigrate(&Remark{}, &Getping{}, &Alive{}, &Delay{})
			argArray := []string{"-a", "-c", "3", "-g", "172.26.18.32/28", "-s"}

			cmd := exec.Command("fping", argArray...) //运行命令
			cmdValue, _ := cmd.CombinedOutput()       //输出运行的命令
			//fmt.Println("string(cmdValue)", string(cmdValue))
			reIp := regexp.MustCompile(`(.*) +: xmt/rcv/%loss = (.*)/(.*)/(.*), min/avg/max = (.*)/(.*)/(.*)`) //过滤100%丢包的数据
			ipValueArray := reIp.FindAllStringSubmatch(string(cmdValue), -1)
			for _, ipValue := range ipValueArray { //把运行的结果放到数据库里面
				db.Create(&Getping{Tim: currentTime.Format("2006-01-02 15:04:05"),
					Dst:  ipValue[1],
					Xmt:  ipValue[2],
					Rcv:  ipValue[3],
					Loss: ipValue[4],
					Min:  ipValue[5],
					Avg:  ipValue[6],
					Max:  ipValue[7]})
			}

			reLive := regexp.MustCompile(`(.*) targets\n(.*) alive\n(.*) unreachable\n(.*) unknown addresses`) //提取存活个数
			liveValueArray := reLive.FindAllStringSubmatch(string(cmdValue), -1)
			for _, liveValue := range liveValueArray { //把运行的结果放到数据库里面
				db.Create(&Alive{Tim: currentTime.Format("2006-01-02 15:04:05"),
					Tar:  liveValue[1],
					Ali:  liveValue[2],
					Unr:  liveValue[3],
					Uadd: liveValue[4]})
			}

			reDelay := regexp.MustCompile(`(.*) ms [(]min round trip time[)]\n(.*) ms [(]avg round trip time[)]\n(.*) ms [(]max round trip time[)]\n(.*) sec [(]elapsed real time[)]`) //提取延迟时间和程序运行时间
			delayValueArray := reDelay.FindAllStringSubmatch(string(cmdValue), -1)
			for _, delayValue := range delayValueArray { //把运行的结果放到数据库里面
				db.Create(&Delay{Tim: currentTime.Format("2006-01-02 15:04:05"),
					Min: delayValue[1],
					Avg: delayValue[2],
					Max: delayValue[3],
					Ela: delayValue[4]})
			}
		}
	}

}
