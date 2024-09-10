package cmd

import (
	"configService/service"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var addSFlag bool
var deleteSFlag bool
var listSFlag bool

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "manage the servicesNodes ",
	Long: `examples: ` + programName + ` services --add/-A procService ./procService.yaml /home/nari/ns6000/config/procService.yaml  
           ` + programName + ` services --delete/-D procService 
           ` + programName + ` services --list/-L`,
	Args: func(cmd *cobra.Command, args []string) error {
		if listSFlag && len(args) != 0 {
			return errors.New("list dont need args")
		}
		if addSFlag && len(args) != 3 {
			return errors.New("add need 2 args")
		}
		if deleteSFlag && len(args) != 1 {
			return errors.New("delete need 1 args")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var cs service.ConfigService
		if addSFlag && !deleteSFlag && !listSFlag {
			nodeName := args[0]
			filePath := args[1]
			configFilePath := args[2]
			err := cs.CreateService(Zk, nodeName, filePath, configFilePath)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("create service success")
		} else if !addSFlag && deleteSFlag && !listSFlag {
			nodeName := args[0]
			err := cs.DeleteService(Zk, nodeName)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("delete service success!")
		} else if !addSFlag && !deleteSFlag && listSFlag {
			cs.ListServices(Zk)
		} else {
			fmt.Println("wrong options")
		}
	},
}

func init() {
	rootCmd.AddCommand(servicesCmd)
	servicesCmd.Flags().BoolVarP(&addSFlag, "add", "A", false, "Add a Node")
	servicesCmd.Flags().BoolVarP(&deleteSFlag, "delete", "D", false, "Delete a Node")
	servicesCmd.Flags().BoolVarP(&listSFlag, "list", "L", false, "List a Node")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
