package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var recurFlag bool

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a path in zookeeper, the value is nil",
	Long:  `examples: ` + programName + ` create [--recur/-R] /configService/path/procService`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		zNode := args[0]
		if recurFlag {
			err := Zk.EnsureDir(zNode)
			if err != nil {
				fmt.Println("wrong zk path")
				return
			}
		} else {
			err := Zk.CreateNode(zNode, nil)
			if err != nil {
				fmt.Println("wrong zk path")
				return
			}
		}
		fmt.Println("create success")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVarP(&recurFlag, "recur", "R", false, "if create the dir recurrently")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
