package main

import (
	"flag"
	"fmt"
	//	"io"
	"io/ioutil"
	//"net/http"
	"os"
	//"regexp"
	//"time"
	"strings"
	//	"strings"
	"path/filepath"

	"github.com/golang/glog"
)

func main() {
	fmt.Println("")
	flag.Set("logtostderr", "true")
	flag.Parse()
	datadir := os.Getenv("STOCKDATA")
	var a int
	glog.Infoln(datadir)
	filepath.Walk(datadir+"/raw1/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false && (info.Size() < 100 || info.Size() == 615) {
			if info.Size() == 615 {
				//f, _ := os.OpenFile(path)
				b, _ := ioutil.ReadFile(path)
				//glog.Infoln(len(strings.Split(string(b), "\t")), path)
				if len(strings.Split(string(b), "\t")) > 4 {
					glog.Infoln(len(strings.Split(string(b), "\t")), path, "!!!!!!!!!!!")
				}
				//f.Close()
			}
			//glog.Infoln(path)
			os.Remove(path)
			a++
		}
		return nil
	})
	glog.Infoln("总计删除", a, "个文件")
}
