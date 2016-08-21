package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
)

const base_url = "http://market.finance.sina.com.cn/downxls.php?date=2016-08-19&symbol=sz002223"
const base_url2 = "http://qt.gtimg.cn/q=sh600000"

var channel chan string
var work chan int

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	datadir := os.Getenv("STOCKDATA")
	glog.Infoln(datadir)
	channel = make(chan string, 1000)
	work = make(chan int, 500)
	fmt.Println("开始")
	fmt.Printf("%06d", 0)
	for i := 0; i < 604000; i++ {
		time.Sleep(time.Nanosecond * 0)
		if i < 3000 {
			go GetSZ(i)
		} else if i < 301000 && i >= 300000 {
			go GetSZ(i)
		} else if i >= 600000 {
			go GetSH(i)
		}
	}
	var buf []string
	var count int
	for {
		bufs := <-channel
		count++
		if len(bufs) > 1 {
			buf = append(buf, bufs)
		}
		glog.Infoln(count, len(buf), buf[len(buf)-1])
		if count >= 8000 {
			break
		}
	}
	os.Remove(datadir + "/codelist")
	file, err := os.Create(datadir + "/codelist")
	glog.Infoln(file, err)
	for _, s := range buf {
		file.WriteString(fmt.Sprintln(s))
	}
	file.Close()
	glog.Infoln("结束")
}
func GetSH(code int) {
	work <- 1
	var SH string = "http://qt.gtimg.cn/q=sh"
	codestr := fmt.Sprintf("%06d", code)
	if resp, err := http.Get(SH + codestr); err == nil {
		buf, _ := ioutil.ReadAll(resp.Body)
		if len(buf) > 80 {
			channel <- "sh" + codestr
		} else {
			channel <- "1"
		}
	}
	<-work
}
func GetSZ(code int) {
	work <- 1
	var SZ string = "http://qt.gtimg.cn/q=sz"
	codestr := fmt.Sprintf("%06d", code)
	if resp, err := http.Get(SZ + codestr); err == nil {
		buf, _ := ioutil.ReadAll(resp.Body)
		if len(buf) > 80 {
			channel <- "sz" + codestr
		} else {
			channel <- "1"
		}
	}
	<-work
}
