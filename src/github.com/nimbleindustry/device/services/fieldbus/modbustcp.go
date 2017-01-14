package fieldbus

import (
	"errors"
	"fmt"
	"time"

	"github.com/nimbleindustry/device/common"
	"github.com/nimbleindustry/device/define"
	"github.com/nimbleindustry/device/integrations"

	"github.com/goburrow/modbus"
	"github.com/nimbleindustry/suture"
)

const (
	modbusConnectionTimeout = 5 * time.Second
	modbusSampleFrequency   = 5 * time.Second
)

// ModbusTCPService provides access to configured modbus tcp interfaces.
type ModbusTCPService struct {
	common.Service
	GenericModbusService

	StartDelay time.Duration // Duration to delay prior to starting the service
	LogFunc    func(string)  // Destination for logging

	stop         chan bool
	integrations []integrations.Integration

	handler *modbus.TCPClientHandler
}

// Serve is called by this service's supervisor—it should not be called directly.
// Exiting or panicing from this function will force the supervisor to attempt to restart.
func (svc *ModbusTCPService) Serve() {
	svc.ServiceState = suture.ServiceNotRunning

	// timeout, if spec'd, can be used ease initialization and avoid race conditions
	if svc.StartDelay > 0 {
		svc.LogFunc(fmt.Sprintf("%s delays start for %s", svc.Name, svc.StartDelay))
		time.Sleep(svc.StartDelay)
	}

	// if services that this service depends on were specified, wait for them to start
	if !svc.WaitForServices() {
		svc.LogFunc(fmt.Sprintf("%s exits as one or more dependent services not found", svc.Name))
		return
	}

	// this channel used to interrupt for/select loop
	svc.stop = make(chan bool)

	// initialize the modbus connection and mappings
	configErr := svc.initConfigurations()
	fieldBusErr := svc.initBusIntegration()
	if configErr != nil || fieldBusErr != nil {
		// configurations not set for modbus, we should wait for config update messages
		svc.LogFunc(fmt.Sprintf("%s exits awaiting configuration updates", svc.Name))
		return
	}
	if configErr == nil && fieldBusErr == nil {
		svc.LogFunc(fmt.Sprintf("%s begins running normally", svc.Name))
	}
	timeout := time.Duration(modbusSampleFrequency)
	for {
		// important to set the state here for dependent services
		svc.ServiceState = suture.ServiceNormal
		select {
		case <-svc.stop:
			svc.ServiceState = suture.ServiceNotRunning
			// Clean up resources here, know that Serve will get called again
			svc.clean()
			return
		case <-time.After(timeout):
			// collect all modbus input values
			err := svc.initConnection()
			if err != nil {
				// there was a problem connecting to modbus, force backoff recovery
				svc.LogFunc(fmt.Sprintf("%s exits due to connection error: %s", svc.Name, err))
				return
			}
			m, err := svc.readAllInputs()
			if err != nil {
				svc.closeConnection()
				svc.LogFunc(fmt.Sprintf("%s exits, error reading input data: %s", svc.Name, err))
				svc.ServiceState = suture.ServiceNotRunning
				return
			}
			svc.closeConnection()
			common.SendBusMessage(define.TopicOpsReport, m)
		}
	}
}

// Stop is called by a supervisor to signal that the service should be stopped. Every
// effort should be made to clean up resources and put the service in a state in which
// it could be restarted.
func (svc *ModbusTCPService) Stop() {
	svc.LogFunc(fmt.Sprintf("%s stops as directed by supervisor", svc.Name))
	svc.stop <- true
}

// State returns the state of service
func (svc *ModbusTCPService) State() int {
	return svc.ServiceState
}

func (svc *ModbusTCPService) clean() {
	svc.closeConnection()
}

func (svc *ModbusTCPService) initConnection() error {
	if svc.connection.Type == "" {
		return errors.New("No modbusTCP connection records")
	}
	addr := svc.connection.Endpoint
	port := svc.connection.Port
	handler := modbus.NewTCPClientHandler(fmt.Sprintf("%s:%d", addr, port))
	handler.Timeout = modbusConnectionTimeout
	handler.SlaveId = 1
	err := handler.Connect()
	if err != nil {
		return err
	}
	svc.LogFunc(fmt.Sprintf("%s establishes connection to modbus slave %s", svc.Name, handler.Address))
	svc.handler = handler
	svc.client = modbus.NewClient(handler)
	return nil
}

func (svc *ModbusTCPService) closeConnection() (err error) {
	if svc.handler != nil {
		err = svc.handler.Close()
		if err != nil {
			svc.LogFunc(fmt.Sprintf("%s reports error disconnecting from modbus slave %s, %s", svc.Name, svc.connection.Endpoint, err))
		} else {
			svc.LogFunc(fmt.Sprintf("%s disconnects from modbus slave %s", svc.Name, svc.connection.Endpoint))
		}
	}
	return
}
