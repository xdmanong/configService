/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"configService/service"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var addFlag bool
var deleteFlag bool
var listFlag bool

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "manage the hostNodes and ip",
	Long: `examples: ` + programName + ` nodes --add/-A main1 172.16.17.100 
           ` + programName + ` nodes --delete/-D main1 
           ` + programName + ` nodes --list/-L`,
	Args: func(cmd *cobra.Command, args []string) error {
		if listFlag && len(args) != 0 {
			return errors.New("list dont need args")
		}
		if addFlag && len(args) != 2 {
			return errors.New("add need 2 args")
		}
		if deleteFlag && len(args) != 1 {
			return errors.New("delete need 1 args")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var cs service.ConfigService
		if addFlag && !deleteFlag && !listFlag {
			nodeName := args[0]
			ip := args[1]
			err := cs.CreateNodeInfo(Zk, nodeName, ip)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("create node success")
		} else if !addFlag && deleteFlag && !listFlag {
			nodeName := args[0]
			err := cs.DeleteNodesInfo(Zk, nodeName)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("delete node success")
		} else if !addFlag && !deleteFlag && listFlag {
			cs.ListNodes(Zk)
		} else {
			fmt.Println("wrong options")
		}
	},
}

func init() {
	rootCmd.AddCommand(nodesCmd)
	nodesCmd.Flags().BoolVarP(&addFlag, "add", "A", false, "Add a Node")
	nodesCmd.Flags().BoolVarP(&deleteFlag, "delete", "D", false, "Delete a Node")
	nodesCmd.Flags().BoolVarP(&listFlag, "list", "L", false, "List a Node")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
