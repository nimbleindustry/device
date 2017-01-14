package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/nimbleindustry/device/common"
	"github.com/nimbleindustry/device/define"

	"github.com/fsnotify/fsnotify"
	"github.com/nimbleindustry/suture"
)

// ConfigService manages asset, device and connectivity configuration state
type ConfigService struct {
	common.Service

	StartDelay time.Duration // Duration to delay prior to starting the service
	LogFunc    func(string)  // Destination for logging

	stop chan bool
}

// Serve is called by this service's supervisor—it should not be called directly.
// Exiting or panicing from this function will force the supervisor to attempt to restart.
func (svc *ConfigService) Serve() {
	svc.ServiceState = suture.ServiceNotRunning

	// timeout, if spec'd, can be used ease initialization and avoid race conditions
	if svc.StartDelay > 0 {
		svc.LogFunc(fmt.Sprintf("%s delays start for %s\n", svc.Name, svc.StartDelay))
		time.Sleep(svc.StartDelay)
	}

	// if services that this service depends on were specified, wait for them to start
	if !svc.WaitForServices() {
		svc.LogFunc("One or more dependent services not found")
		return
	}

	// this channel used to interrupt for/select loop
	svc.stop = make(chan bool)

	// watch the config directory for configuration changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		svc.ServiceState = suture.ServicePaused
		svc.LogFunc("ConfigService: cannot create fsnotify/watcher")
		return
	}
	defer watcher.Close()
	watcher.Add(define.ConfigPath)

	svc.loadConfigFiles()

	svc.LogFunc(fmt.Sprintf("%s running normally", svc.Name))
	for {
		// important to set the state here for dependent services
		svc.ServiceState = suture.ServiceNormal
		select {
		case <-svc.stop:
			svc.ServiceState = suture.ServiceNotRunning
			// Clean up resources here, know that Serve will get called again
			return
		case event := <-watcher.Events:
			svc.LogFunc(fmt.Sprintf("%s detects config change %s, %s", svc.Name, event.Name, event.String()))
			switch event.Name {
			case define.AssetConfigPath:
				svc.loadAssetConfigFile()
				common.SendBusMessage(define.AssetConfigUpdated, event)
			case define.EquipmentConfigPath:
				svc.loadEquipmentConfigFile()
				common.SendBusMessage(define.EquipmentConfigUpdated, event)
			case define.ConnectivityConfigPath:
				svc.loadConnectionsConfigFile()
				common.SendBusMessage(define.ConnectivityConfigUpdated, event)
			}
		case err := <-watcher.Errors:
			svc.ServiceState = suture.ServicePaused
			svc.LogFunc(fmt.Sprintf("ConfigService fsnotify/watcher error %s", err))
			return
		}
	}
}

// Stop is called by a supervisor to signal that the service should be stopped. Every
// effort should be made to clean up resources and put the service in a state in which
// it could be restarted.
func (svc *ConfigService) Stop() {
	svc.LogFunc(fmt.Sprintf("%s stops as directed by supervisor", svc.Name))
	svc.stop <- true
}

// State returns the state of service
func (svc *ConfigService) State() int {
	return svc.ServiceState
}

func loadJSON(path string, object interface{}) error {
	if len(path) == 0 {
		return errors.New("Path empty")
	}
	var configSource []byte
	configSource, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(configSource), object)
}

func (svc *ConfigService) loadConfigFiles() {
	configLoads := []struct {
		path   string
		object interface{}
	}{
		{define.AssetConfigPath, &common.AssetConfig},
		{define.EquipmentConfigPath, &common.EquipmentConfig},
		{define.ConnectivityConfigPath, &common.ConnectionConfig},
	}
	for _, v := range configLoads {
		err := loadJSON(v.path, v.object)
		if err != nil {
			svc.LogFunc(fmt.Sprintf("ConfigService warns. Error loading %s, %s", v.path, err))
		}
	}
}

func (svc *ConfigService) loadAssetConfigFile() {
	err := loadJSON(define.AssetConfigPath, &common.AssetConfig)
	if err != nil {
		svc.LogFunc(fmt.Sprintf("ConfigService: warns. Error loading asset config file, %s", err))
	}
}

func (svc *ConfigService) loadConnectionsConfigFile() {
	err := loadJSON(define.ConnectivityConfigPath, &common.ConnectionConfig)
	if err != nil {
		svc.LogFunc(fmt.Sprintf("ConfigService: warns. Error loading connections config file, %s", err))
	}
}

func (svc *ConfigService) loadEquipmentConfigFile() {
	err := loadJSON(define.EquipmentConfigPath, &common.EquipmentConfig)
	if err != nil {
		svc.LogFunc(fmt.Sprintf("ConfigService: warns. Error loading equipment config file, %s", err))
	}
}
