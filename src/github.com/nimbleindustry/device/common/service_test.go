package common

import (
	"log"
	"testing"
	"time"

	"math"

	"github.com/nimbleindustry/suture"
)

type TestService struct {
	Service

	timeoutToStart time.Duration
	shutdown       chan bool
	successChannel chan bool
}

func (svc *TestService) Serve() {
	defer func() { /*log.Printf("Service %s returning\n", s)*/ }()
	svc.ServiceState = suture.ServiceNotRunning
	svc.WaitForServices()
	svc.shutdown = make(chan bool)
	log.Printf("%s sleeping for %s\n", svc.Name, svc.timeoutToStart)
	time.Sleep(svc.timeoutToStart)
	timeout := time.Duration(math.MaxInt32 * time.Second)
	if svc.Name == "serviceB" {
		timeout = time.Duration(3 * time.Second)
	}
	log.Println("Timeout is: ", time.Now().Add(timeout))
	for {
		svc.ServiceState = suture.ServiceNormal
		log.Printf("%s running normally", svc.Name)
		select {
		case <-svc.shutdown:
			svc.ServiceState = suture.ServiceNotRunning
			return
		case <-time.After(timeout):
			svc.successChannel <- true
			log.Printf("%s timeout, returning", svc.Name)
			svc.ServiceState = suture.ServiceNotRunning
			return
		}
	}
}

func (svc *TestService) Stop() {
	svc.shutdown <- true
}

func TestDependentService(t *testing.T) {

	supervisor := suture.New("testSupervisor", suture.Spec{FailureDecay: 100})
	serviceA := &TestService{timeoutToStart: time.Duration(3 * time.Second), successChannel: make(chan bool, 1)}
	serviceA.Name = "serviceA"
	serviceC := &TestService{timeoutToStart: time.Duration(2 * time.Second), successChannel: make(chan bool, 1)}
	serviceC.Name = "serviceC"

	serviceB := &TestService{timeoutToStart: time.Duration(1 * time.Second), successChannel: make(chan bool, 1)}
	serviceB.Name = "serviceB"
	serviceB.AddServiceDependentUpon("serviceA")
	serviceB.AddServiceDependentUpon("serviceC")

	supervisor.Add(serviceA)
	supervisor.Add(serviceB)
	supervisor.Add(serviceC)

	go supervisor.Serve()
	success := <-serviceB.successChannel

	if !success {
		t.Fail()
	}
}
