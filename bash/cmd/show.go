package cmd

import (
	"configService/service"
	"github.com/spf13/cobra"
)

var serviceConfigFlag bool
var servicePathFlag bool

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show services path -P, show services config -C",
	Long:  `examples: ` + programName + ` show [--config/-C] [--path/-P] procService`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		procName := args[0]
		var cs service.ConfigService
		if servicePathFlag {
			cs.PrintPath(Zk, procName)
		}
		if serviceConfigFlag {
			cs.PrintConfig(Zk, procName)
		}
	},
}

func init() {
	servicesCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVarP(&serviceConfigFlag, "config", "C", false, "show the service config content")
	showCmd.Flags().BoolVarP(&servicePathFlag, "path", "P", false, "show the service config file path")
}
