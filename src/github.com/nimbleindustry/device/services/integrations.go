package services

import (
	"fmt"
	"time"

	"github.com/nimbleindustry/device/common"

	"github.com/nimbleindustry/device/define"
	"github.com/nimbleindustry/device/integrations"

	"github.com/nimbleindustry/suture"
)

// IntegrationsService is responsible for maintaining connections
// to one or more integrated service implementations as found in the
// /integrations folder and defined in connections.json
//
// Ops (operations) Integrations receive machine operational and sensor telemetry data
// State Integrations receive information regarding the state and health of the Nimble Device
type IntegrationsService struct {
	common.Service

	StartDelay time.Duration // Duration to delay prior to starting the service
	LogFunc    func(string)  // Destination for logging

	opsIntegrations   []integrations.Integration
	stateIntegrations []integrations.Integration
	stop              chan bool
}

// Serve is called by this service's supervisor—it should not be called directly.
// Exiting or panicing from this function will force the supervisor to attempt to restart.
func (svc *IntegrationsService) Serve() {
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

	svc.loadIntegrations()

	// this channel used to interrupt for/select loop
	svc.stop = make(chan bool)

	// break the for loop below every 1 minute to send the state of this Device
	timeout := time.Duration(1 * time.Minute)
	svc.LogFunc(fmt.Sprintf("%s begins running normally", svc.Name))
	for {
		// important to set the state here for dependent services
		svc.ServiceState = suture.ServiceNormal
		select {
		case <-svc.stop:
			svc.ServiceState = suture.ServiceNotRunning
			// Clean up resources here, know that Serve will get called again
			svc.clean()
			return
		case <-common.BusChannel(define.AssetConfigUpdated):
			svc.LogFunc(fmt.Sprintf("%s advises that the asset config file was updated, no action taken", svc.Name))
		case <-common.BusChannel(define.ConnectivityConfigUpdated):
			svc.LogFunc(fmt.Sprintf("%s advises that the connections config file was updated, reloading integratations", svc.Name))
			svc.loadIntegrations()
		case <-common.BusChannel(define.EquipmentConfigUpdated):
			svc.LogFunc(fmt.Sprintf("%s advises that the equipment config file was updated, no action taken", svc.Name))
		case msg := <-common.BusChannel(define.TopicStateReport):
			for _, v := range svc.stateIntegrations {
				err := v.SendState(msg.(*common.SystemState))
				if err != nil {
					svc.LogFunc(fmt.Sprintf("%s warns: error sending state data to %s, %s", svc.Name, v.Record().Endpoint, err))
				}
			}
		case msg := <-common.BusChannel(define.TopicOpsReport):
			for _, v := range svc.opsIntegrations {
				err := v.SendData(msg.(map[string]interface{}))
				if err != nil {
					svc.LogFunc(fmt.Sprintf("%s warns: error sending ops data to %s, %s", svc.Name, v.Record().Endpoint, err))
				}
			}
		case <-time.After(timeout):
			// Maybe test/tickle the connections every timeout?
		}
	}
}

// Stop is called by a supervisor to signal that the service should be stopped. Every
// effort should be made to clean up resources and put the service in a state in which
// it could be restarted.
func (svc *IntegrationsService) Stop() {
	svc.LogFunc(fmt.Sprintf("%s stops as directed by supervisor", svc.Name))
	svc.stop <- true
}

// State returns the state of service
func (svc *IntegrationsService) State() int {
	return svc.ServiceState
}

func (svc *IntegrationsService) clean() {
	for _, v := range svc.opsIntegrations {
		v.Close()
	}
	for _, v := range svc.stateIntegrations {
		v.Close()
	}
}

func (svc *IntegrationsService) loadIntegrations() {
	svc.clean()
	svc.stateIntegrations = integrations.GetDeviceStateIntegrations()
	for _, v := range svc.stateIntegrations {
		err := v.Connect()
		if err != nil {
			svc.LogFunc(fmt.Sprintf("%s reports an error connecting to integration %s: %s", svc.Name, v.Record().Provider, err))
		} else {
			svc.LogFunc(fmt.Sprintf("%s reports connection to %s integration at %s", svc.Name, v.Record().Provider, v.Record().Endpoint))
		}
	}
	svc.opsIntegrations = integrations.GetDeviceOperationsAndTelemetryIntegrations()
	for _, v := range svc.opsIntegrations {
		err := v.Connect()
		if err != nil {
			svc.LogFunc(fmt.Sprintf("%s reports an error connecting to integration %s: %s", svc.Name, v.Record().Provider, err))
		} else {
			svc.LogFunc(fmt.Sprintf("%s reports connection to %s integration at %s", svc.Name, v.Record().Provider, v.Record().Endpoint))
		}
	}
}
