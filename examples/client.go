package main

import (
	"fmt"
	"gdiamond/client/client"
)

//implement  listener.ManagerListener
type myListener struct {
}

func (receiver myListener) ReceiveConfigInfo(content string) {
	fmt.Println("config changed:", content)
}

func main() {
	//client
	cli := client.NewClient()

	//publish
	ok := cli.PublishConfig("linna", "DEFAULT_GROUP", "Who is linna?")
	fmt.Println(ok)

	//get
	content := cli.GetConfig("linna", "DEFAULT_GROUP", 100)
	fmt.Println(content)

	//get and set listener
	content = cli.GetConfigAndSetListener("linna", "DEFAULT_GROUP", 100, myListener{})
	fmt.Println(content)
}
