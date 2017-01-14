package common

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/nimbleindustry/suture"
)

// A Service is goroutine-based 'process' that attempts to create fault-tolerant, recoverable services
type Service struct {
	Name         string
	ServiceState int

	servicesDependentUpon map[string]bool
}

// Serve satisfies the interface defined by suture.
func (svc *Service) Serve() {
	fmt.Println("Base class service starting")
}

// Stop satisfies the interface defined by suture.
func (svc *Service) Stop() {
	fmt.Println("Base class service stopping")
}

// State satisfies the interface defined by suture.
func (svc *Service) State() int {
	return svc.ServiceState
}

func (svc *Service) String() string {
	return svc.Name
}

// WaitForServices will block the execution of Service until dependent Services are running normally.
func (svc *Service) WaitForServices() bool {
	if svc.servicesDependentUpon != nil {
		return suture.WaitForServices(svc.servicesDependentUpon, time.Duration(10*time.Second))
	}
	return true
}

// AddServiceDependentUpon adds the name of a service that this service is dependent upon.
// Multiple names can be added.
func (svc *Service) AddServiceDependentUpon(name string) {
	if svc.servicesDependentUpon == nil {
		svc.servicesDependentUpon = make(map[string]bool, 1)
	}
	svc.servicesDependentUpon[name] = true
}

// Detail returns information about the Service
func (svc *Service) Detail() string {
	var buffer bytes.Buffer
	buffer.WriteString("Service: ")
	buffer.WriteString(svc.Name)
	buffer.WriteString(" [")
	buffer.WriteString("Status: ")
	buffer.WriteString(strconv.Itoa(svc.State()))
	buffer.WriteString(", ")
	if svc.servicesDependentUpon != nil {
		buffer.WriteString("Dependent Upon: (")
		for k := range svc.servicesDependentUpon {
			buffer.WriteString(k + ", ")
		}
		buffer.Truncate(buffer.Len() - 2)
		buffer.WriteString(")")
	}
	buffer.WriteString("]")
	return buffer.String()
}
