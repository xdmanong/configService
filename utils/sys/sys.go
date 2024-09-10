package sys

import (
	"fmt"
	"net"
	"os/exec"
)

func GetLocalIpv4() (string, error) {
	inter, err := net.InterfaceByName("eth0")
	if err != nil {
		return "", err
	}

	addrs, err := inter.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ip := addr.(*net.IPNet)
		if ip.IP.DefaultMask() != nil {
			return ip.IP.String(), nil
		}
	}
	return "", fmt.Errorf("IP address not found for eth0")
}

func SendFileToHost(filepath string, hostName string, ip string) {
	cmdStr := "scp " + filepath + " " + hostName + ":" + filepath
	cmd := exec.Command("scp", filepath, hostName+":"+filepath)
	// 设置标准输入输出为空，这样子进程就不会继承当前进程的输入输出
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	//windows
	//cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
	err := cmd.Start()
	fmt.Println("starting send: ", cmdStr)
	if err != nil {
		fmt.Println("send file:", filepath, " to host:", hostName, " wrong! ", err)
		fmt.Println("try use ip:", ip)
		cmdStr = "scp " + filepath + " " + ip + ":" + filepath
		fmt.Println("starting send: ", cmdStr)
		cmd = exec.Command("scp", filepath, ip+":"+filepath)
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil
		err = cmd.Start()
		if err != nil {
			fmt.Println("send file:", filepath, " to ip:", ip, " wrong! ", err)
		} else {
			go func() {
				err := cmd.Wait()
				if err != nil {
					fmt.Println("send file:", filepath, " to ip:", ip, " wrong!, err: ", err)
				}
				fmt.Println("send file:", filepath, " to ip:", ip, " success!")
			}()
		}

	} else {
		go func() {
			err := cmd.Wait()
			if err != nil {
				fmt.Println("send file:", filepath, " to host:", hostName, " wrong!, err: ", err)
			}
			fmt.Println("send file:", filepath, " to host:", hostName, " success!")
		}()
	}
}
