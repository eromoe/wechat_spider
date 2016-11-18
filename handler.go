package wechat_spider

import (
    "bytes"
    "io"
    "io/ioutil"
    "bufio"
    "log"
    "net/http"
    "os"
    "reflect"
    "strings"
    "fmt"
    "net/url"
    "github.com/elazarl/goproxy"
)

var (
    Verbose = false
    Logger  = log.New(os.Stderr, "", log.LstdFlags)
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func writefile(path string, text string) {
    f, _ := os.Create(path)
    w := bufio.NewWriter(f)
    w.WriteString(text+"\n")
    w.Flush()
    f.Close()
}

func ProxyHandle(proc Processor) func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
    return func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
        if ctx.Req.URL.Path == `/mp/getmasssendmsg` && !strings.Contains(ctx.Req.URL.RawQuery, `f=json`) {
            fmt.Println("get history urls")
            var data []byte
            var err error
            data, resp.Body, err = copyReader(resp.Body)
            if err != nil {
                return resp
            }
            t := reflect.TypeOf(proc)
            v := reflect.New(t.Elem())
            p := v.Interface().(Processor)
            go func() {
                err = p.Process(ctx.Req, data)
                if err != nil {
                    Logger.Println(err.Error())
                }
                p.Output()
            }()
        }

        // write key
        if (ctx.Req.URL.Path == `/mp/getappmsgext`) {
            u, _ := url.Parse(ctx.Req.URL.RequestURI())
            m, _ := url.ParseQuery(u.RawQuery)

            fmt.Println("get biz, key and uin")

            biz_id := m["__biz"][0]
            biz_dir := fmt.Sprintf("./biz/%s", biz_id)

            if _, err := os.Stat(biz_dir); os.IsNotExist(err) {
                os.Mkdir(biz_dir, 0644)
            }

            key_path := fmt.Sprintf("./biz/%s/%s", biz_id, "key")
            uin_path := fmt.Sprintf("./biz/%s/%s", biz_id, "uin")

            writefile(key_path, m["key"][0])
            writefile(uin_path, m["uin"][0])
            // err := ioutil.WriteFile(key_path, m["key"][0], 0644)
            // check(err)
        }

        return resp
    }

}

// One of the copies, say from b to r2, could be avoided by using a more
func copyReader(b io.ReadCloser) (bs []byte, r2 io.ReadCloser, err error) {
    var buf bytes.Buffer
    if _, err = buf.ReadFrom(b); err != nil {
        return nil, b, err
    }
    if err = b.Close(); err != nil {
        return nil, b, err
    }
    return buf.Bytes(), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
