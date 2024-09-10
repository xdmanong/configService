package cmd

import (
	"configService/utils"
	"context"
	"errors"
	"github.com/smallnest/rpcx/client"
)

var ConfigClient *client.XClientPool
var Zk *utils.ZkClient

// Call 自定义操作方法，提供自定义接口名及json字符串参数
func Call(ctx context.Context, fun string, args string) (string, error) {
	if ConfigClient == nil {
		return "", errors.New("未初始化的dbio客户端")
	}
	var res string
	err := ConfigClient.Get().Call(ctx, fun, args, &res)
	return res, err
}

// Init 获取一个rpc客户端, 程序启动时调用，设置zk地址和服务地址
func Init(zkAddr []string, svPath string) {
	// 设置zookeeper客户端配置
	dis, err := client.NewZookeeperDiscovery(
		"/rpcx_bcService", // 基地址
		svPath,            // 服务路由:runtime/scada/debug/v1
		zkAddr,            // zk服务器地址及端口
		nil)
	if err != nil {
		ConfigClient = nil
	} else {
		ConfigClient = client.NewXClientPool( // 长连接建议使用NewXClientPool池连接
			3,
			svPath, // 服务路由:runtime/scada/debug/v1
			client.Failfast,
			client.RandomSelect,
			dis,
			client.DefaultOption)
	}
}

// Close 销毁dbio客户端
func Close() {
	if ConfigClient != nil {
		ConfigClient.Close()
	}
	ConfigClient = nil
}
