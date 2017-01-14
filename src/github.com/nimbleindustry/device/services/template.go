package services

import (
	"fmt"
	"math"
	"time"

	"github.com/nimbleindustry/device/common"

	"github.com/nimbleindustry/suture"
)

// TestService provides an example of how services should be coded.
type TestService struct {
	common.Service

	StartDelay time.Duration // Duration to delay prior to starting the service
	LogFunc    func(string)  // Destination for logging

	stop chan bool
}

// Serve is called by this service's supervisor—it should not be called directly.
// Exiting or panicing from this function will force the supervisor to attempt to restart.
func (svc *TestService) Serve() {
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

	// this timeout can be used to interrupt for/select loop, note the
	// very lengthy duration here (which in effect is infinite)
	timeout := time.Duration(math.MaxInt32 * time.Second)
	svc.LogFunc(fmt.Sprintf("%s begins running normally", svc.Name))
	for {
		// important to set the state here for dependent services
		svc.ServiceState = suture.ServiceNormal
		select {
		case <-svc.stop:
			svc.ServiceState = suture.ServiceNotRunning
			// Clean up resources here, know that Serve will get called again
			return
		case msg := <-common.BusChannel("someChannel"):
			// perform some operation on msg
			svc.LogFunc(msg.(string))
		case <-time.After(timeout):
			svc.ServiceState = suture.ServiceNotRunning
			return
		}
	}
}

// Stop is called by a supervisor to signal that the service should be stopped. Every
// effort should be made to clean up resources and put the service in a state in which
// it could be restarted.
func (svc *TestService) Stop() {
	svc.LogFunc(fmt.Sprintf("%s stops as directed by supervisor", svc.Name))
	svc.stop <- true
}

// State returns the state of service
func (svc *TestService) State() int {
	return svc.ServiceState
}
