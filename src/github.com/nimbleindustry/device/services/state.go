package services

import (
	"fmt"
	"time"

	"github.com/nimbleindustry/device/common"

	"github.com/nimbleindustry/device/define"
	"github.com/nimbleindustry/device/integrations"

	"github.com/nimbleindustry/suture"
)

// StateService is responsible for reporting the status of the Nimble (field) Device
// to one or more integrated service implementations (/integrations). Three explicit
// data points are reported: percentage of (RAM) consumed, percentage of disk consumed,
// and the load factor for the last minute. Implicitly, the existence of this record
// at a particular time interval indicates the Device was running at that moment.
type StateService struct {
	common.Service

	StartDelay time.Duration // Duration to delay prior to starting the service
	LogFunc    func(string)  // Destination for logging

	integrations []integrations.Integration
	stop         chan bool
}

// Serve is called by this service's supervisor—it should not be called directly.
// Exiting or panicing from this function will force the supervisor to attempt to restart.
func (svc *StateService) Serve() {
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

	common.SendBusMessage(define.TopicStateReport, common.GetSystemState())

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
			return
		case <-time.After(timeout):
			common.SendBusMessage(define.TopicStateReport, common.GetSystemState())
		}
	}
}

// Stop is called by a supervisor to signal that the service should be stopped. Every
// effort should be made to clean up resources and put the service in a state in which
// it could be restarted.
func (svc *StateService) Stop() {
	svc.LogFunc(fmt.Sprintf("%s stops as directed by supervisor", svc.Name))
	svc.stop <- true
}

// State returns the state of service
func (svc *StateService) State() int {
	return svc.ServiceState
}
