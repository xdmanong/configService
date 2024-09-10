/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"bash/cmd"
	"configService/conf"
	"configService/utils"
	"fmt"
)

func main() {
	var err error
	cmd.Zk, err = utils.NewClient(conf.ZkConfSlice, "/", 5)
	if err != nil {
		fmt.Println("init ZkClient error, ", err)
		return
	}
	defer cmd.Zk.Close()
	cmd.Init(conf.ZkConfSlice, "realtime/scada/configService/v1")
	cmd.Execute()
}
