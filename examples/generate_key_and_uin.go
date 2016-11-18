package main

import (
	"log"
	"net/http"
    "net/url"
    "bufio"
    "fmt"
    "os"
	// "github.com/sundy-li/wechat_spider"
     "github.com/eromoe/wechat_spider"
	"github.com/elazarl/goproxy"
)

func main() {
	var port = "8899"
	proxy := goproxy.NewProxyHttpServer()
	//open it see detail logs
	// wechat_spider.Verbose = true
	proxy.OnResponse().DoFunc(
		wechat_spider.ProxyHandle(&CustomProcessor{}),
	)
	log.Println("server will at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, proxy))

}

//Just to implement Output Method of interface{} Processor
type CustomProcessor struct {
	wechat_spider.BaseProcessor
}

func (c *CustomProcessor) Output() {
	//Just print the length of result urls
	println("result urls size =>", len(c.Urls()))
    uf := c.Urls()[0]
    u, _ := url.Parse(uf)
    m, _ := url.ParseQuery(u.RawQuery)
    biz_id := m["__biz"][0]

    println("=======================")
    println("=======================")
    println("biz", biz_id)
    println("=======================")
    println("=======================")


    biz_dir := fmt.Sprintf("./biz/%s", biz_id)
    if _, err := os.Stat(biz_dir); os.IsNotExist(err) {
        os.Mkdir(biz_dir, 0644)
    }
    urls_path := fmt.Sprintf("%s/urls.txt", biz_dir)
    f, _ := os.Create(urls_path)
    w := bufio.NewWriter(f)

    for _, url := range c.Urls() {
        w.WriteString(url+"\n")
    }
    w.Flush()
    f.Close()
}
