package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
	//	"strings"
	"path/filepath"

	"github.com/golang/glog"
)

//var timestart time.Time = time.Date(2016, time.January, 1, 0, 0, 0, 0, time.Local) ///////////////重要

const base_url = "http://market.finance.sina.com.cn/downxls.php?date=2016-08-19&symbol=sz002223"
const base_url1 = "http://market.finance.sina.com.cn/downxls.php?date="
const base_url2 = "&symbol="

var channel chan string
var work chan int
var datadir string

func main() {
	channel = make(chan string, 10000)
	work = make(chan int, 300)
	fmt.Println("")
	flag.Set("logtostderr", "true")
	flag.Parse()
	datadir = os.Getenv("STOCKDATA")
	glog.Infoln(datadir)
	//	glog.Infoln("开始时间", timestart)
	//	channel = make(chan string, 1000)
	//	work = make(chan int, 500)
	glog.Infoln("开始")
	f, e1 := os.Open(datadir + "/codelist")
	if e1 != nil {
		glog.Infoln(e1)
		return
	}
	buflist, _ := ioutil.ReadAll(f)
	f.Close()
	reg, _ := regexp.Compile("\\S+")
	codelist := reg.FindAllString(string(buflist), -1)
	os.Mkdir(datadir+"/raw1/", 0777)
	for _, code := range codelist {
		os.Mkdir(datadir+"/raw1/"+code, 0777)
	}
	glog.Infoln("下载")
	var requestnumber int
	var count int
	//	for date := time.Now().Add(-time.Second * 86400 * 0); date.After(timestart); date = date.Add(-time.Second * 86400) {
	//		if date.Weekday() == time.Sunday || date.Weekday() == time.Saturday {
	//			continue
	//		}
	//		dd := date.Format("2006-01-02")
	//		for _, code := range codelist {
	//			requestnumber++
	//			time.Sleep(time.Second * 0)
	//			go DownStock(code, dd)
	//		}
	//	}
	codereg := regexp.MustCompile("s[hz][0-9]{6,6}")
	datereg := regexp.MustCompile("[-0-9]{10,10}")
	filepath.Walk(datadir+"/raw2/", func(path string, info os.FileInfo, err error) error {
		if info.Size() < 10 {
			return nil
		}
		ff, _ := os.Open(path)
		str, _ := ioutil.ReadAll(ff)
		codeone := codereg.FindString(string(path))
		dates := datereg.FindAllString(string(str), -1)
		for _, dd := range dates {
			if Exist(datadir + "/raw1/" + codeone + "/" + dd) {
				requestnumber++
				count++
				continue
			}
			requestnumber++
			time.Sleep(time.Second * 0)
			go DownStock(codeone, dd)
		}

		ff.Close()
		return nil
	})
	//var requestnumber int = len(codelist) * (int(time.Now().Unix()-timestart.Unix()) / 86400) //总计请求数目

	glog.Infoln(count, requestnumber)
	for {
		<-channel
		count++
		glog.Infoln(count, requestnumber)
		if count >= int(requestnumber) {
			break
		}
	}
}
func DownStock(code string, date string) {
	if Exist(datadir + "/raw1/" + code + "/" + date) {
		channel <- "1"
		return
	}
	work <- 1

	if resp, err := http.Get(base_url1 + date + base_url2 + code); err == nil {
		glog.Infoln(resp, err, base_url1+date+base_url2+code)
		buf, err2 := ioutil.ReadAll(resp.Body)
		if len(buf) > 10 && err2 == nil && resp.StatusCode == http.StatusOK {
			os.Remove(datadir + "/raw1/" + code + "/" + date)
			f, _ := os.Create(datadir + "/raw1/" + code + "/" + date)
			f.Write(buf)
			f.Close()
			resp.Body.Close()
		} else {

		}

	}
	channel <- "1"
	<-work
}
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
