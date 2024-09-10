package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

var programName string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   programName,
	Short: "a cli of zookeeper, also a terminal to manager configs",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	execPath := os.Args[0]
	// 使用path包的Base函数获取路径的最后一部分，即程序名
	programName = path.Base(execPath)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bash.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
