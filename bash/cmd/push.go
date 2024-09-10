package cmd

import (
	"configService/service"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var forceFlag bool
var pathFlag bool

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push configFile to zookeeper",
	Long: `push .yaml to zookeeper， the first arg is procName, the second arg is configFilePath, please check if path is valid 
           examples: ` + programName + ` push [--force/-F] procService procService,yaml
           ` + programName + ` push [--force/-F] --path/-P procService /home/nari/ns6000/config/procService,yaml`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		procName := args[0]
		yamlFilePath := args[1]
		if !pathFlag && !isValidYAMLFile(yamlFilePath) {
			fmt.Printf("The file %s is not a valid YAML file.\n", yamlFilePath)
			os.Exit(1)
		}

		var cs service.ConfigService
		err := cs.PushConfigFile(Zk, yamlFilePath, procName, forceFlag, pathFlag)
		if err != nil {
			fmt.Println("push data failed!")
			return
		}
		fmt.Println("push data success!")
	},
}

func init() {
	pushCmd.Flags().BoolVarP(&forceFlag, "force", "F", false, "if force creating the path ")
	pushCmd.Flags().BoolVarP(&pathFlag, "path", "P", false, "if push the Path")
	rootCmd.AddCommand(pushCmd)
}

func isValidYAMLFile(filePath string) bool {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}
