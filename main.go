package main

import (
	"randomip/UseProxy"
	"os"
	"fmt"
)



func main(){
	fmt.Printf("使用工具参数 main.exe https 或者 http\n")
	UseProxy.ProxyStart(os.Args[1])
}