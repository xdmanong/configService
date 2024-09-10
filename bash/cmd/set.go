/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set a path and value in zookeeper",
	Long:  `examples: ` + programName + ` set /myConfig abc `,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		zNode := args[0]
		data := args[1]
		err := Zk.CreateNodeOrUpdateIfExist(zNode, []byte(data))
		if err != nil {
			fmt.Println("wrong zk path")
			return
		}
		fmt.Println("set path: ", zNode, " value: ", data, "success!")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
