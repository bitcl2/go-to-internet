package main

import (
    "flag"
    "fmt"
    "gitee.com/Luna-CY/go-to-internet/config"
    "gitee.com/Luna-CY/go-to-internet/proxy"
    "golang.org/x/net/http2"
    "log"
    "net/http"
    "os"
)

// Usage 打印控制台Usage信息
func Usage() {
    _, _ = fmt.Fprintln(flag.CommandLine.Output(), "server -H Hostname -c CRT -k KEY [options]")

    flag.PrintDefaults()
}

func main() {
    c := &config.Config{}

    flag.StringVar(&c.Hostname, "H", "", "域名，该域名应该与证书的域名一致")
    flag.IntVar(&c.Port, "p", 443, "监听端口号")

    flag.StringVar(&c.SSLCrtFile, "c", "", "SSL CRT文件路径")
    flag.StringVar(&c.SSLKeyFile, "k", "", "SSL KEY文件路径")

    flag.BoolVar(&c.Authorize, "auth", false, "是否开启用户身份验证，默认不启用")
    flag.Usage = Usage
    flag.Parse()

    if "" == c.Hostname || "" == c.SSLCrtFile || "" == c.SSLKeyFile {
        flag.Usage()

        os.Exit(0)
    }

    server := &http.Server{Addr: fmt.Sprintf(":%d", c.Port)}
    server.Handler = &proxy.Proxy{}

    if err := http2.ConfigureServer(server, &http2.Server{}); nil != err {
        log.Fatal("配置http/2服务器失败")
    }

    fmt.Printf("启动监听 %v:%d ...\n", c.Hostname, c.Port)
    if err := server.ListenAndServeTLS(c.SSLCrtFile, c.SSLKeyFile); nil != err {
        log.Fatal("启动http/2服务器失败")
    }
}
