package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get the value of a path",
	Long:  `examples: ` + programName + ` create get /configService/path/procService`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		zNode := args[0]
		data, err := Zk.GetData(zNode)
		if err != nil {
			fmt.Println("get failed: ", err)
			return
		}
		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
