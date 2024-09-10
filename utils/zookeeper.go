package utils

import (
	"configService/conf"
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type ZkClient struct {
	zkServers []string
	Conn      *zk.Conn
	zkRoot    string                                        // 服务根节点，这里是/
	Status    int32                                         //client 状态 1 runing
	Cancel    context.CancelFunc                            //手动关闭ctx
	CancelCtx context.Context                               //手动关闭ctx
	callback  func([]byte, string, map[string]string) error //TODO
}

func (s *ZkClient) Callback() func([]byte, string, map[string]string) error {
	return s.callback
}

func (s *ZkClient) SetCallback(callback func([]byte, string, map[string]string) error) {
	s.callback = callback
}

func NewClient(zkServers []string, zkRoot string, timeout int) (*ZkClient, error) {
	client := new(ZkClient)
	client.zkServers = zkServers
	client.zkRoot = zkRoot
	client.CancelCtx, client.Cancel = context.WithCancel(context.Background())
	client.Status = 1
	// 连接服务器s
	// option := zk.WithEventCallback(client.EventCallback)
	// conn, _, err := zk.Connect(zkServers, time.Duration(timeout)*time.Second, option)
	conn, _, err := zk.Connect(zkServers, time.Duration(timeout)*time.Second, zk.WithLogInfo(false))
	if err != nil {
		fmt.Println("zk connect failed!")
		return nil, err
	}
	//log.Println("zookeeper connetion ok")
	client.Conn = conn
	// 创建服务根节点
	if err := client.ensureRoot(); err != nil {
		client.Close()
		return nil, err
	}
	//log.Println("ensureRoot ok")
	return client, nil
}

func (s *ZkClient) Stop() {
	s.Cancel()
	atomic.StoreInt32(&s.Status, 0)
}

func (s *ZkClient) Close() {
	s.Cancel()
	atomic.StoreInt32(&s.Status, 0)
	if s.Conn != nil {
		s.Conn.Close()
	}
}

func (s *ZkClient) ensureRoot() error {
	exists, _, err := s.Conn.Exists(s.zkRoot)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("ErrNodeExists")
	}
	return nil
}

func (s *ZkClient) CreateNode(zNode string, data []byte) error {
	_, err := s.Conn.Create(zNode, data, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Printf("cannot create ZNode: %v\n\n", err)
		return err
	}
	return nil
}

func (s *ZkClient) CreateNodeOrUpdateIfExist(zNode string, data []byte) error {
	// 检查ZooKeeper节点是否存在，如果不存在则创建
	conn := s.Conn
	exists, _, err := conn.Exists(zNode)
	if err != nil {
		//logger.Fatalf("cannot check if ZNode exists: %v", err)
		fmt.Printf("cannot check if ZNode exists: %v\n", err)
		return err
	}

	if exists == false {
		_, err = conn.Create(zNode, data, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			//logger.Fatalf("cannot create ZNode: %v", err)
			fmt.Printf("cannot create ZNode: %s\n", zNode)
			return err
		}

	} else {
		_, err = conn.Set(zNode, data, -1)
		if err != nil {
			//logger.Fatalf("cannot set ZNode data: %v", err)
			fmt.Printf("cannot set ZNode data: %v\n", err)
			return err
		}

	}
	return nil
}

func (s *ZkClient) ZkWatch(procName string) error {
	zNode := GetConfigZNode(procName)
	filepath, err := s.GetData(GetPathZNode(procName))
	if err != nil {
		fmt.Println("cannot find path of: ", procName)
		return err
	}
	fmt.Println("Watch path:", zNode)
Loop:
	for atomic.LoadInt32(&s.Status) == 1 {
		_, _, getCh, err := s.Conn.GetW(zNode)
		if err != nil {
			fmt.Println(err)
			time.Sleep(10 * time.Second) //5s 之后
			goto Loop
		}
		select {
		case chEvent := <-getCh:
			{
				fmt.Println("watch callback path:", chEvent.Path, "event_type:", chEvent.Type)
				if chEvent.Type == zk.EventNodeCreated {
					fmt.Println("EventNodeCreated")
					//TODO
				} else if chEvent.Type == zk.EventNodeDeleted {
					fmt.Println("EventNodeDeleted")
					//TODO
				} else if chEvent.Type == zk.EventNodeDataChanged {
					fmt.Println("EventNodeDeleted")
					data, err := s.GetData(chEvent.Path)
					if err == nil {
						nodesMap := make(map[string]string)
						s.GetHostNodes(nodesMap)
						s.callback(data, string(filepath), nodesMap)
					} else {
						//alter
						fmt.Println(err)
					}
				} else if chEvent.Type == zk.EventNodeChildrenChanged {
					fmt.Println("EventNodeChildrenChanged")
					//todo
				}
			}

		case <-s.CancelCtx.Done():
			break
		}
	}
	fmt.Println("STOP WATCH PROC: ", procName)
	return nil
}

func (s *ZkClient) GetData(path string) ([]byte, error) {
	data, _, err := s.Conn.Get(path)
	return data, err
}

func (s *ZkClient) GetAllPath(zNode string) ([]string, error) {
	conn := s.Conn
	// 列出根节点下的所有子节点
	children, _, err := conn.Children(zNode)
	if err != nil {
		fmt.Println("List children nodes error,", err)
		return nil, err
	}
	nodes := make([]string, 0)
	// 打印所有子节点的路径
	for _, child := range children {
		nodes = append(nodes, child)
	}
	return nodes, nil
}

func (s *ZkClient) PrintAllPath(path string) {
	conn := s.Conn
	// 列出根节点下的所有子节点
	children, _, err := conn.Children(path)
	if err != nil {
		fmt.Println("List children nodes error,", err)
	}

	// 打印所有子节点的路径
	for _, child := range children {
		fmt.Println(child)
	}
}

func (s *ZkClient) PrintAllValue(path string) {
	conn := s.Conn
	// 列出根节点下的所有子节点

	children, _, err := conn.Children(strings.TrimRight(path, "/"))
	if err != nil {
		fmt.Println("List children nodes error,", children)
	}

	// 打印所有子节点的路径
	for _, child := range children {
		zNode := path + child
		data, _, _ := conn.Get(zNode)
		fmt.Println(child + " : " + string(data))
	}
}

func (s *ZkClient) EnsureDir(dir string) error {
	conn := s.Conn
	// 分割目录路径
	parts := strings.Split(dir, "/")
	ensure := ""

	// 创建目录
	for _, part := range parts[1:] {
		ensure += "/" + part
		exists, _, err := conn.Exists(ensure)
		if err != nil {
			//logger.Fatalf("cannot check if ZNode exists: %v", err)
			fmt.Printf("cannot check if ZNode exists: %v\n", err)
			return err
		}

		if exists == false {
			_, err = conn.Create(ensure, nil, 0, zk.WorldACL(zk.PermAll))
			if err != nil {
				//logger.Fatalf("cannot create ZNode: %v", err)
				fmt.Printf("cannot create ZNode: %s\n", ensure)
				return err
			}

		}
	}

	return nil
}

func (s *ZkClient) GetHostNodes(nodesMap map[string]string) {
	conn := s.Conn
	// 列出根节点下的所有子节点

	children, _, err := conn.Children(strings.TrimRight(conf.NodesPreFix, "/"))
	if err != nil {
		fmt.Println("List children nodes error,", children)
	}

	// 打印所有子节点的路径
	for _, child := range children {
		zNode := conf.NodesPreFix + child
		data, _, _ := conn.Get(zNode)
		nodesMap[child] = string(data)
	}
}

func (s *ZkClient) PrintZkNodes(zNode string) {
	conn := s.Conn
	children, _, err := conn.Children(zNode)
	if err != nil {
		fmt.Printf("Error getting children of %s: %s\n", zNode, err)
	}

	for _, child := range children {
		childZNode := path.Join(zNode, child)
		s.PrintZkNodes(childZNode)
	}

	if len(children) == 0 {
		// 如果是叶子节点，直接打印
		fmt.Println(zNode)
	}
}

func (s *ZkClient) PrintZNodeValue(zNode string) {
	data, err := s.GetData(zNode)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

func GetConfigZNode(procName string) string {
	return conf.ConfigPreFix + procName
}

func GetPathZNode(procName string) string {
	return conf.PathPreFix + procName
}

func GetNodesZNode(procName string) string {
	return conf.NodesPreFix + procName
}
