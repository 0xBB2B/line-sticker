package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

var waitGroup = new(sync.WaitGroup)

func download(i int) {
	url := fmt.Sprintf("https://stickershop.line-scdn.net/stickershop/v1/sticker/%d/ANDROID/sticker.png", i)
	fmt.Printf("开始下载:%s\n", url)
	res, errHTTP := http.Get(url)
	if errHTTP != nil || res.StatusCode != 200 {
		fmt.Printf("下载失败:%s\n", res.Request.URL)
	}

	fmt.Printf("开始读取文件内容,url=%s\n", url)
	data, errRead := ioutil.ReadAll(res.Body)
	if errRead != nil {
		fmt.Println("读取数据失败")
	}

	img, _, errDecode := image.Decode(bytes.NewReader(data))
	if errDecode != nil {
		fmt.Println(errDecode)
	}

	// 缩略图的大小
	dstImage512 := imaging.Resize(img, 512, 0, imaging.Lanczos)
	errSave := imaging.Save(dstImage512, fmt.Sprintf("./2b/%d.png", i))
	if errSave != nil {
		fmt.Println(errSave)
	}

	//计数器-1
	waitGroup.Done()
}

func main() {
	//创建多个协程，同时下载多个图片
	now := time.Now()
	n := 145198350
	os.MkdirAll("./images", 0755)

	for i := 0; i < 16; i++ {
		//计数器+1
		waitGroup.Add(1)
		go download(n + i)
	}

	//等待所有协程操作完成
	waitGroup.Wait()
	fmt.Printf("下载总时间:%v\n", time.Now().Sub(now))
}
