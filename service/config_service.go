package service

import (
	"configService/conf"
	"configService/utils"
	"configService/utils/sys"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

type ConfigService struct {
	ZkWatchClientMap sync.Map
	lock             sync.Mutex
}

var cfg *ConfigService

func GetConfigService() *ConfigService {
	return cfg
}

func (*ConfigService) PushConfigFile(s *utils.ZkClient, filepath string, procName string, forceFlag bool, pathFlag bool) error {
	// 假设你有一个YAML文件路径
	var data []byte
	var err error
	var zNode string
	if !pathFlag {
		yamlFilePath := filepath
		zNode = utils.GetConfigZNode(procName)
		// 读取YAML文件
		data, err = os.ReadFile(yamlFilePath)
		if err != nil {
			fmt.Printf("unable to read YAML file: %v\n", err)
			return err
		}
	} else {
		zNode = utils.GetPathZNode(procName)
		data = []byte(filepath)
	}

	if forceFlag {
		fmt.Printf("creating Znode: %s\n", zNode)
		err := s.EnsureDir(zNode)
		if err != nil {
			fmt.Printf("Create Znode: %s error !\n", zNode)
			return err
		}
	}

	err = s.CreateNodeOrUpdateIfExist(zNode, data)
	if err != nil {
		return err
	}

	fmt.Println("YAML file content has been stored to ZooKeeper.")
	return nil
}

func (*ConfigService) CreateNodeInfo(s *utils.ZkClient, nodeName string, ip string) error {
	zNode := utils.GetNodesZNode(nodeName)
	ip = utils.TrimString(ip)
	if !utils.IsValidIPv4(ip) {
		return errors.New("error ipv4 format")
	}
	data := []byte(ip)
	err := s.CreateNodeOrUpdateIfExist(zNode, data)
	return err
}

func (*ConfigService) ListNodes(s *utils.ZkClient) {
	s.PrintAllValue(conf.NodesPreFix)
}

func (*ConfigService) DeleteNodesInfo(s *utils.ZkClient, nodeName string) error {
	zNode := utils.GetNodesZNode(nodeName)
	err := s.Conn.Delete(zNode, 0)
	return err
}

func (*ConfigService) PullConfigFile(s *utils.ZkClient, procName string) error {
	conn := s.Conn
	zNode := utils.GetConfigZNode(procName)
	pNode := utils.GetPathZNode(procName)
	// 获取节点内容
	fileData, _, err := conn.Get(zNode)
	if err != nil {
		fmt.Printf("Get node %s error: %v\n", zNode, err)
		return err
	}

	data, _, err := conn.Get(pNode)
	if err != nil {
		fmt.Printf("Get node %s error: %v\n", pNode, err)
		return err
	}
	filePath := string(data)
	writeErr := updateLocalConfigFile(fileData, filePath)
	if writeErr != nil {
		return writeErr
	}
	fmt.Println(string(fileData))
	return nil
}

func (s *ConfigService) StopWatchConfigFileUpdate(procName string) {
	fmt.Println("stop watching ", procName)
	zk, exists, err := s.getZkWatchClientIfExist(procName)
	if err != nil {
		fmt.Println(err)
	}
	if exists {
		zk.Stop()
	}
}

func (s *ConfigService) WatchConfigFileUpdate(procName string) error {
	zk, exist, err := s.getZkWatchClient(procName)
	if zk == nil || err != nil {
		fmt.Println("get zkClient of ", procName, "wrong ,", err)
		return err
	}
	if exist && atomic.LoadInt32(&zk.Status) == 1 {
		fmt.Println("already in watch: ", procName)
		return errors.New("already in watch: " + procName)
	}
	if !exist {
		zk.SetCallback(UpdateConfigFile)
	}
	atomic.StoreInt32(&zk.Status, 1)
	zk.CancelCtx, zk.Cancel = context.WithCancel(context.Background())
	go func(procName string) {
		fmt.Println("start watching: ", procName)
		err = zk.ZkWatch(procName)
		defer func() {
			fmt.Println("watching ", procName, "exited!")
		}()
		if err != nil {
			fmt.Println(err)
		}
	}(procName)
	return nil
}

func (s *ConfigService) getServicesList() []string {
	zk, err := utils.NewClient(conf.ZkConfSlice, "/", 5)
	if err != nil {
		fmt.Println(err)
	}
	confServiceList, _ := zk.GetAllPath(strings.TrimRight(conf.ConfigPreFix, "/"))
	fmt.Println("confList:", confServiceList)
	pathServiceList, _ := zk.GetAllPath(strings.TrimRight(conf.PathPreFix, "/"))
	fmt.Println("pathList:", pathServiceList)
	serviceList := utils.Intersect(confServiceList, pathServiceList)
	return serviceList
}

func (s *ConfigService) WatchAllConfigFileUpdate() {
	serviceList := s.getServicesList()
	fmt.Println(serviceList)
	for _, procName := range serviceList {
		go func(procName string) {
			err := s.WatchConfigFileUpdate(procName)
			if err != nil {
				fmt.Println("start watching service: "+procName+"wrong! ", err)
			}
		}(procName)
	}
}

func (s *ConfigService) StopWatchAllConfigFileUpdate() {
	serviceList := s.getServicesList()
	for _, procName := range serviceList {
		s.StopWatchConfigFileUpdate(procName)
	}
}

func (s *ConfigService) getZkWatchClientIfExist(procName string) (*utils.ZkClient, bool, error) {
	if ret, found := s.ZkWatchClientMap.Load(procName); found {
		zkClient, ok := ret.(*utils.ZkClient)
		if !ok {
			fmt.Println("convert to ZkClient failed!")
			return nil, false, errors.New("convert to ZkClient failed")
		}
		return zkClient, true, nil
	} else {
		return nil, false, errors.New("cannot found ZkClient")
	}
}

func (s *ConfigService) getZkWatchClient(procName string) (*utils.ZkClient, bool, error) {
	path := utils.GetConfigZNode(procName)
	s.lock.Lock()
	defer s.lock.Unlock()
	tmpRet, exists := s.ZkWatchClientMap.Load(procName)
	if exists && tmpRet == nil {
		s.ZkWatchClientMap.Delete(procName)
	} else if exists && tmpRet != nil {
		zkClient, ok := tmpRet.(*utils.ZkClient)
		if !ok {
			fmt.Println("convert to ZkClient failed!")
			s.ZkWatchClientMap.Delete(procName)
			return nil, false, errors.New("convert to ZkClient failed")
		}
		return zkClient, true, nil
	}
	zk, err := utils.NewClient(conf.ZkConfSlice, path, 5)
	if zk == nil || err != nil {
		fmt.Println("creat zkConnection error ", err)
		return nil, false, err
	}
	s.ZkWatchClientMap.Store(procName, zk)
	return zk, false, nil

}

func (s *ConfigService) PrintPath(zk *utils.ZkClient, procName string) {
	zNode := utils.GetPathZNode(procName)
	zk.PrintZNodeValue(zNode)
}

func (s *ConfigService) PrintConfig(zk *utils.ZkClient, procName string) {
	zNode := utils.GetConfigZNode(procName)
	zk.PrintZNodeValue(zNode)
}

func (s *ConfigService) SetModel(jsonStr string) {
	var result map[string]any
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	procName := result["proc_name"]
	model := result["model"]
	fmt.Println(procName)
	fmt.Println(model)
}

func (s *ConfigService) CreateService(zk *utils.ZkClient, procName string, filePath string, configFilePath string) error {
	err := s.PushConfigFile(zk, filePath, procName, true, false)
	if err != nil {
		fmt.Println("create configNode of: ", procName, "wrong! ", err)
		return err
	}
	err = s.PushConfigFile(zk, configFilePath, procName, true, true)
	if err != nil {
		fmt.Println("create pathNode of: ", procName, "wrong! ", err)
		return err
	}
	return nil
}

func (s *ConfigService) DeleteService(zk *utils.ZkClient, procName string) error {
	ret := 0
	configPath := utils.GetConfigZNode(procName)
	path := utils.GetPathZNode(procName)
	err := zk.Conn.Delete(configPath, 0)
	if err != nil {
		ret = 2
		fmt.Println("delete pathNode of: ", procName, "wrong! ", err)
	}
	err = zk.Conn.Delete(path, 0)
	if err != nil {
		ret++
		fmt.Println("delete configNode of: ", procName, "wrong! ", err)
	}
	if ret == 0 {
		return nil
	} else if ret == 2 {
		return errors.New("delete pathNode of: " + procName + "failed! ")
	} else if ret == 1 {
		return errors.New("delete configNode of: " + procName + "failed! ")
	} else if ret == 3 {
		return errors.New("delete both node failed")
	}
	return nil
}

func (s *ConfigService) ListServices(zk *utils.ZkClient) {
	zk.PrintAllPath(strings.TrimRight(conf.ConfigPreFix, "/"))
}

func updateLocalConfigFile(data []byte, filepath string) error {
	err := os.WriteFile(filepath, data, 0644)
	if err != nil {
		fmt.Printf("无法写入文件：%s err: %v\n", filepath, err)
		return err
	}
	fmt.Println("send file: ", filepath, " to host: ", conf.HostName, " success!")
	return nil
}

func UpdateConfigFile(data []byte, filepath string, nodesMap map[string]string) error {
	err := updateLocalConfigFile(data, filepath)
	if err != nil {
		return err
	}
	copyFileToRemote(filepath, nodesMap)
	return nil
}

func copyFileToRemote(filepath string, nodesMap map[string]string) {
	localIp, _ := sys.GetLocalIpv4()
	for hostName, ip := range nodesMap {
		if hostName == conf.HostName || ip == localIp {
			continue
		} else {
			sys.SendFileToHost(filepath, hostName, ip)
		}
	}
}

func init() {
	cfg = &ConfigService{
		ZkWatchClientMap: sync.Map{},
	}
}
