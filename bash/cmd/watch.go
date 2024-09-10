package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var stopFlag bool
var allFlag bool

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch the changes of config file and sync to all hosts",
	Long: `examples: ` + programName + ` watch [--stop/-S] procService 
           ` + programName + ` watch --all/-A [--stop/-S] procService`,
	Args: func(cmd *cobra.Command, args []string) error {
		if allFlag && len(args) != 0 {
			return errors.New("to all flag dont need args")
		}
		if !allFlag && len(args) != 1 {
			return errors.New("just need 1 arg (procName)")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !stopFlag {
			if !allFlag {
				procName := args[0]
				res, err := Call(context.Background(), "WatchConfigFileUpdate", procName)
				if err != nil {
					fmt.Println("WatchConfigFileUpdate method call failed! ", err)
					return
				}
				fmt.Println(res)
			} else {
				res, err := Call(context.Background(), "WatchAllConfigFileUpdate", "")
				if err != nil {
					fmt.Println("WatchAllConfigFileUpdate method call failed!")
					return
				}
				fmt.Println(res)
			}
		} else {
			if !allFlag {
				procName := args[0]
				res, err := Call(context.Background(), "StopWatchConfigFileUpdate", procName)
				if err != nil {
					fmt.Println("StopWatchConfigFileUpdate method call failed!")
					return
				}
				fmt.Println(res)
			} else {
				res, err := Call(context.Background(), "StopWatchAllConfigFileUpdate", "")
				if err != nil {
					fmt.Println("StopWatchAllConfigFileUpdate method call failed!")
					return
				}
				fmt.Println(res)
			}
		}
	},
}

func init() {
	watchCmd.Flags().BoolVarP(&stopFlag, "stop", "S", false, "if stop the watching")
	watchCmd.Flags().BoolVarP(&allFlag, "all", "A", false, "if apply to all services in zookeeper")
	rootCmd.AddCommand(watchCmd)
}
