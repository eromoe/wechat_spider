# wechat_spider
微信公众号爬虫  (只需设置代理, 一键可以爬取指定公众号的所有历史文章)

常见问题[FAQ][3]

代理服务端: 通过Man-In-Middle 代理方式获取微信服务端返回, 自动模拟请求自动分页,抓取对应点击的所有历史文章

客户端:   暂时支持 win,macos,android三大平台,  iphone由于https证书问题后续再支持

#### 代理服务端
- 一个简单的Demo  [simple_server.go][1]

```
package main

import (
	"log"
	"net/http"

	"github.com/sundy-li/wechat_spider"

	"github.com/elazarl/goproxy"
)

func main() {
	var port = "8899"
	proxy := goproxy.NewProxyHttpServer()
	//open it see detail logs
	// wechat_spider.Verbose = true
	proxy.OnResponse().DoFunc(
		wechat_spider.ProxyHandle(wechat_spider.NewBaseProcessor()),
	)
	log.Println("server will at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, proxy))

}
```

- 自定义输出源,实现Processor接口的Output方法即可, [custom_output_server.go][2]


[1]: https://github.com/sundy-li/wechat_spider/blob/master/examples/simple_server.go
[2]: https://github.com/sundy-li/wechat_spider/blob/master/examples/custom_output_server.go
[3]: https://github.com/sundy-li/wechat_spider/blob/master/FAQ.md

- 微信会屏蔽频繁的请求,所以历史文章的翻页请求调用了Sleep()方法, 默认每个请求休眠50ms,可以根据实际情况自定义Processor覆盖此方法


#### 客户端使用:
  (确保客户端 能正常访问 代理服务端的服务)

- Android客户端使用方法:
  运行后, 设置手机的代理为 本机ip 8899端口,  打开微信客户端, 点击任一公众号查看历史文章按钮, 即可爬取该公众号的所有历史文章(已经支持自动翻页爬取)
- win/mac客户端,设置下全局代理对应 代理服务端的服务和端口,同理点击任一公众号查看历史文章按钮
- 自动化批量爬取所有公众号:  Windows客户端获取批量公众号所有历史文章方法
  1. 要求安装windows +  微信pc版本 + ActivePython3 + autogui, 设置windows下全局代理对应 代理服务端的服务和端口
  2. 修改 win_client.py 中的 bizs参数, 通过pyautogui.position() 瞄点设置 first_ret, rel_link 坐标
  3. 执行 python win_client.py 将自动生成链接,模拟点击

#### 增加
biz, uin, key 导出文件功能