package controller

import (
	"configService/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

var cs *service.ConfigService

type ConfigController struct{}

func (pc *ConfigController) WatchConfigFileUpdate(ctx context.Context, procName string, res *string) error {
	err := cs.WatchConfigFileUpdate(procName)
	if err != nil {
		*res = "start Watching error"
		return err
	}
	*res = "start watching success"
	return nil
}

func (pc *ConfigController) WatchAllConfigFileUpdate(ctx context.Context, procName string, res *string) error {
	fmt.Println("start watching all config")
	go cs.WatchAllConfigFileUpdate()
	*res = "start all watching"
	return nil
}

func (pc *ConfigController) StopWatchConfigFileUpdate(ctx context.Context, procName string, res *string) error {
	cs.StopWatchConfigFileUpdate(procName)
	*res = "stop watching"
	return nil
}

func (pc *ConfigController) StopWatchAllConfigFileUpdate(ctx context.Context, procName string, res *string) error {
	fmt.Println("stop watching all config")
	cs.StopWatchAllConfigFileUpdate()
	*res = "stop all watching"
	return nil
}

func (pc *ConfigController) SetModel(ctx context.Context, jsonStr string, res *string) error {
	var result map[string]any
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	procName := result["proc_name"]
	model := result["model"]
	fmt.Println(procName)
	fmt.Println(model)
	if procName == nil || model == nil {
		fmt.Println("In SetModel: wrong args")
		return errors.New("wrong args")
	}
	//cs.SetModel(jsonStr)
	return nil
}

func init() {
	cs = service.GetConfigService()
}
