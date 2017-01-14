package integrations

import (
	"fmt"

	"github.com/nimbleindustry/device/common"
	"github.com/nimbleindustry/device/define"
)

// Integration defines an interface for all external integrations with IIoT providers
type Integration interface {
	SetRecord(common.ConnectionRecord)
	Record() *common.ConnectionRecord
	Connect() error
	Close() error
	SendState(*common.SystemState) error
	SendData(map[string]interface{}) error
	ReceiveData(interface{}) error
}

// CreateIntegration instantiates a new Integration based on the supplied provider string
func CreateIntegration(provider string) Integration {
	switch provider {
	case define.InitialState:
		return new(InitialState)
	case define.GenericMQTT:
		return new(GenericMQTT)
	}
	return nil
}

// GetDeviceStateIntegrations loads all defined device state integrations from the 'Connections' configuration object
func GetDeviceStateIntegrations() []Integration {
	var integrations []Integration
	for _, v := range common.ConnectionConfig.DeviceStateConnections {
		integration := CreateIntegration(v.Provider)
		if integration == nil {
			fmt.Println("Warning, unable to create Integration object for ", v.Provider)
		}
		integration.SetRecord(v)
		integrations = append(integrations, integration)
	}
	return integrations
}

// GetDeviceOperationsAndTelemetryIntegrations loads all defined device operations and telemetry integrations from the
// 'Connections' configuration object
func GetDeviceOperationsAndTelemetryIntegrations() []Integration {
	var integrations []Integration
	for _, v := range common.ConnectionConfig.OperationsAndTelemetryConnections {
		integration := CreateIntegration(v.Provider)
		if integration == nil {
			fmt.Println("Warning, unable to create Integration object for ", v.Provider)
		}
		integration.SetRecord(v)
		integrations = append(integrations, integration)
	}
	return integrations
}
