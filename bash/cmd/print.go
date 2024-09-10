/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var recurPrintFlag bool

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "print zNodes of the path, has one arg as the path",
	Long:  `examples: ` + programName + ` create print [--recur/-R] /configService/path/procService`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		zNode := args[0]
		if !recurPrintFlag {
			Zk.PrintAllPath(zNode)
		} else {
			Zk.PrintZkNodes(zNode)
		}
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	printCmd.Flags().BoolVarP(&recurPrintFlag, "recur", "R", false, "if print the dir recurrently")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// printCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// printCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
