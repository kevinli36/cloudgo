package main

import (
	"os"

	"github.com/kevinli36/cloudgo/service"
	flag "github.com/spf13/pflag"
)

const (
	PORT string = "8080"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}

	//使用pflag，在启动服务端程序时设置监听端口
	pPort := flag.StringP("port", "p", PORT, "PORT for httpd listening")
	flag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}

	//创建新的服务端，并在设置的端口启动服务
	server := service.NewServer()
	server.Run(":" + port)
}
