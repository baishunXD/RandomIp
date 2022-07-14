package UseProxy

import (
    "fmt"
    "github.com/gocolly/colly"
    "github.com/tidwall/gjson"
    "github.com/elazarl/goproxy"
    "time"
    "net/http"
    "crypto/tls"
    "net/url"
)

type IPlib struct {
    ProxyIp string
    Timeone int
}

var IPStruct IPlib


func RequestIP(ProxyType string) {
    c := colly.NewCollector()
    url := "http://192.168.40.137:5010/get?type="+ProxyType //proxy_pool链接
    fmt.Println(url)
    var body string
    c.WithTransport(&http.Transport{
        MaxIdleConnsPerHost:   10,
        ResponseHeaderTimeout: time.Second * time.Duration(5),
        TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
    })
    c.OnRequest(func(r *colly.Request) {
        r.Headers.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.1 Safari/605.1.15")
        r.Headers.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
        r.Headers.Set("Accept-Encoding","gzip, deflate, br")
        r.Headers.Set("Accept-Language","zh-CN,zh;q=0.9")
        r.Headers.Set("Connection","keep-alive")
    })


    c.OnResponse(func(r *colly.Response) { 
        body = string(r.Body)
    })
    c.Visit(url)
    IPStruct.Timeone = int(time.Now().UnixNano())
    IPStruct.ProxyIp = gjson.Get(body, "proxy").String()
    fmt.Printf("%s\n", IPStruct.ProxyIp)
}


func ProxyStart(ProxyType string) {
    RequestIP(ProxyType)
    proxy := goproxy.NewProxyHttpServer()
    proxy.Tr = &http.Transport{
        Proxy: func(req *http.Request) (*url.URL, error) {
            return url.Parse("http://"+IPStruct.ProxyIp)
    }}

    proxy.OnRequest().DoFunc(
    func(r *http.Request,ctx *goproxy.ProxyCtx)(*http.Request,*http.Response) {
        st,_ :=time.ParseDuration("-5s")  //每次IP更换时间
        ts := time.Now().Add(st)
        if int(ts.UnixNano()) > IPStruct.Timeone {
            fmt.Println("准备更换IP")
            RequestIP(ProxyType)
            
        }
        return r,nil
    })
    fmt.Println("这次换的IP是："+IPStruct.ProxyIp)

    
    proxy.ConnectDial = proxy.NewConnectDialToProxy("http://"+IPStruct.ProxyIp)
    proxy.Verbose = true
    http.ListenAndServe(":8088", proxy)
}