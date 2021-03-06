package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/golang/glog"
)

var timestart int = -1 + 4*2015 + 1
var timeend int = -1 + 4*2016 + 3

const base_url = "http://vip.stock.finance.sina.com.cn/corp/go.php/vMS_FuQuanMarketHistory/stockid/002223.phtml?year=2016&jidu=3"
const base_url1 = "http://vip.stock.finance.sina.com.cn/corp/go.php/vMS_FuQuanMarketHistory/stockid/"
const base_url2 = ".phtml?year="

var channel chan string
var work chan int
var datadir string
var subdir string = "/raw2/"

func main() {
	channel = make(chan string, 10000)
	work = make(chan int, 100)
	fmt.Println("")
	flag.Set("logtostderr", "true")
	flag.Parse()
	datadir = os.Getenv("STOCKDATA")
	glog.Infoln(datadir)
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
	fmt.Println(len(codelist))
	os.Mkdir(datadir+subdir, 0777)
	for _, code := range codelist {
		os.Mkdir(datadir+subdir+code, 0777)
	}
	glog.Infoln("下载")
	var requestnumber int
	for date := timeend; date >= timestart; date-- {

		dd := strconv.Itoa(date/4) + "&jidu=" + strconv.Itoa(date%4+1)
		for _, code := range codelist {
			requestnumber++
			time.Sleep(time.Second * 0)
			if date == timeend {
				go DownRecoveryFactor(code, dd, true)
			} else {
				go DownRecoveryFactor(code, dd, false)
			}
		}
	}
	//var requestnumber int = len(codelist) * (int(time.Now().Unix()-timestart.Unix()) / 86400) //总计请求数目
	var count int
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
func DownRecoveryFactor(code string, date string, cover bool) {
	regdate := regexp.MustCompile("&date=([-0-9]{10,10})'>")
	regvalue := regexp.MustCompile("<td class=\"tdr\"><div align=\"center\">([0-9\\.]+)</div></td>\\s+ </tr>")

	//fmt.Println(1)
	if Exist(datadir+subdir+code+"/"+date) && cover == false {
		//glog.Infoln("pass")
		channel <- "1"
		return
	}
	work <- 1
	//glog.Infoln(base_url1 + code[2:] + base_url2 + date)
	if resp, err := http.Get(base_url1 + code[2:] + base_url2 + date); err == nil {
		//glog.Infoln(base_url1 + code[2:] + base_url2 + date)
		buf, err2 := ioutil.ReadAll(resp.Body)
		if len(buf) > 100 && err2 == nil && resp.StatusCode == http.StatusOK {
			os.Remove(datadir + subdir + code + "/" + date)
			f, _ := os.Create(datadir + subdir + code + "/" + date)
			v_date := regvalue.FindAllStringSubmatch(string(buf), -1)
			v_factor := regdate.FindAllStringSubmatch(string(buf), -1)
			//glog.Warningln(len(v_factor))
			for i := 0; i < len(v_factor); i++ {
				f.Write([]byte(v_factor[i][1]))
				f.Write([]byte(","))
				f.Write([]byte(v_date[i][1]))
				f.Write([]byte("\n"))
			}
			f.Close()
			resp.Body.Close()
		} else {
			glog.Infoln(base_url1 + code[2:] + base_url2 + date)
			glog.Infoln("c")
			if cover == true {
				panic("a")
			}
		}

	} else {
		glog.Infoln("d")
		if cover == true {
			panic("b")
		}
	}
	channel <- "1"
	<-work
}
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
