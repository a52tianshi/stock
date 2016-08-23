package main

import (
	"flag"
	"fmt"
	//"io/ioutil"
	//"net/http"
	"os"
	//"regexp"
	//"time"
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
		if info.IsDir() == false && info.Size() < 1000 {
			os.Remove(path)
			a++
		}
		return nil
	})
	fmt.Println("总计删除", a, "个文件")
}
