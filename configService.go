package main

import (
	"configService/conf"
	"configService/controller"
	"configService/utils/logger"
	"context"
	"flag"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	goVersion string // Go 编译版本
	buildTime string // 构建时间
	buildDir  string // 构建目录

)

var termCh = make(chan int)

func init() {
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("RegisterName: %s\n", conf.Cfg.Server.RegisterName)
		fmt.Printf("GoVersion: %s\n", goVersion)
		fmt.Printf("BuildTime: %s\n", buildTime)
		fmt.Printf("BuildPath: %s\n", buildDir)
		os.Exit(1)
	}
}
func main() {
	cfg := conf.Cfg
	addr := getServerHost(conf.Cfg.Server.Addr)
	fmt.Println(cfg)

	// 初始化日志功能
	logger.Init("configService", addr)
	defer logger.Close()
	logger.Infof("-----configService start-----")
	s := server.NewServer()
	addRegistryPlugin(s, addr, cfg.Zookeeper.Addr, cfg.Server.BasePath)
	if err := s.RegisterName(cfg.Server.RegisterName, new(controller.ConfigController), ""); err != nil {
		fmt.Println(err)
		return
	}
	err := s.Serve("tcp", conf.Cfg.Server.Addr)
	if err != nil {
		return
	}
	<-termCh
}

// 服务退出(ctrl+c和kill有效, kill -9强制退出不会执行term)
func term(s *server.Server) {
	signalChan := make(chan os.Signal, 1)
	// 收到syscall.SIGINT, syscall.SIGTERM, os.Kill信号时转发到signalChan
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Kill)

	sig := <-signalChan
	logger.Infof("服务收到退出信号: %+v", sig)

	// 执行服务器退出处理
	s.Shutdown(context.Background()) //关闭监听,注销zk记录,关闭空闲连接,打印过程中的错误
	logger.Infof("服务%v退出", conf.Cfg.Server.Addr)

	close(termCh)
}

// addRegistryPlugin 添加zookeeper注册插件
func addRegistryPlugin(s *server.Server, addr string, zkAddrs []string, basePath string) {
	logger.Info("启动", "服务在zookeeper的注册地址:", addr)

	r := &serverplugin.ZooKeeperRegisterPlugin{
		ServiceAddress:   "tcp@" + addr,
		ZooKeeperServers: zkAddrs,
		BasePath:         basePath,
		Metrics:          metrics.NewRegistry(),
		UpdateInterval:   time.Minute,
	}

	err := r.Start()
	if err != nil {
		panic(err)
	}
	s.Plugins.Add(r)
}

// getServerHost 获取本地服务器地址
func getServerHost(addr string) string {
	var IP string
	if !strings.HasPrefix(addr, ":") {
		return addr
	}
	serverAddr := os.Getenv("SERVERADDR")
	if serverAddr == "" {
		ipList := make([]string, 0)
		netInterfaces, err := net.Interfaces()
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(netInterfaces); i++ {
			if (netInterfaces[i].Flags & net.FlagUp) != 0 {
				addrs, _ := netInterfaces[i].Addrs()
				for _, address := range addrs {
					if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
						if ipNet.IP.To4() != nil {
							ipList = append(ipList, ipNet.IP.String())
						}
					}
				}
			}
		}
		if len(ipList) <= 0 {
			panic("获取本地IP失败")
		}
		IP = ipList[0] + addr
	} else {
		IP = serverAddr
	}
	return IP
}

// checkZookeeperService 检查是否存在已经注册的服务
func checkZookeeperService(zkAddrs []string, basePath, servicePath string) {
	registerPath := basePath + "/" + servicePath
	if conn, _, err := zk.Connect(zkAddrs, time.Second*5); err != nil {
		logger.Fatal("zookeeper服务连接失败，程序退出")
	} else {
		defer conn.Close()
		children, _, _ := conn.Children(registerPath)
		if len(children) > 0 {
			logger.Fatal("已存在注册的zookeeper服务:", children[0], ",程序退出")
		}
	}
}
