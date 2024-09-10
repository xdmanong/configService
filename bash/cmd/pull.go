package cmd

import (
	"configService/service"
	"fmt"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull configFile from zk",
	Long: `pull configFile from zk, first arg is procName, second arg is file path to save configFile
            examples: ` + programName + ` pull procService`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		procName := args[0]
		var cs service.ConfigService
		err := cs.PullConfigFile(Zk, procName)
		if err != nil {
			fmt.Println("pull config failed!")
			return
		}
		fmt.Println("pull config completed!")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
