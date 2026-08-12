package main

import (
	"fmt"
	"log"
	"net/http"

	"infinite-calc/web"
)

func main() {
	web.RegisterRoutes()

	addr := ":8080"
	fmt.Println("========================================")
	fmt.Println("  无限规则通用计算器 已启动")
	fmt.Println("  模块：血量 · 掷骰 · 防御预设")
	fmt.Println("  请在浏览器打开: http://localhost" + addr)
	fmt.Println("========================================")
	log.Fatal(http.ListenAndServe(addr, nil))
}
