package main

import "github.com/literaen/simple_project/tasks/internal/web"

func main() {
	web := web.NewWeb()
	web.Init()
	web.Run()
}
