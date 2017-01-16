/*
   Copyright (C) 2016  Industrial Internet Systems, LLC dba NimbleIndustry

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nimbleindustry/device/define"
	"github.com/nimbleindustry/device/services"
	"github.com/nimbleindustry/device/services/fieldbus"
	"github.com/nimbleindustry/suture"
)

import _ "net/http/pprof"
import (
	"net/http"
)

var masterSupervisor *suture.Supervisor
var defaultServiceSpec = suture.Spec{Log: localLog, FailureThreshold: 3}

func main() {
	sysLogPtr := flag.Bool("syslog", false, "Send logging output to syslog")
	profileFlag := flag.Bool("profile", false, "Enable cpu/mem profiling")
	flag.Parse()
	if *sysLogPtr {
		logWriter, err := syslog.New(syslog.LOG_NOTICE, define.SystemName)
		if err == nil {
			log.SetOutput(logWriter)
		}
	}
	if *profileFlag {
		go func() {
			// Run the Profiler on port 6789
			log.Println(http.ListenAndServe(":6789", nil))
		}()
	}

	masterSupervisor = suture.New(define.MasterSupervisorName, defaultServiceSpec)
	startServices(masterSupervisor, localLog)

	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		for {
			sig := <-sigs
			localLog(sig.String())
			switch sig {
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGTERM:
				fallthrough
			case syscall.SIGSTOP:
				fallthrough
			case syscall.SIGQUIT:
				masterSupervisor.Stop()
				time.Sleep(time.Second * 1)
				done <- true
			}
		}
	}()

	masterSupervisor.Serve()
	<-done
}

func localLog(s string) {
	log.Println(s)
}

/*
	Services:

	MasterSupervisor
	- Config [1]                Manages asset, device and connectivity configuration
	- State [1]                 Responsible for sending state of the (computing) device
	- Integration [1]			Responsible for managing integrations to IIoT services
	- FieldbusSupervisor        Supervisor for all fieldbus integration services
		- ModbusTCP [0-1]       Service to manage I/O to Modbus TCP
        - ModbusRTU [0-1]       Service to manage I/O to Modbus RTU (serial) (PLANNED)
		- OPCUA [0-1]           Service to manage I/O to OPC/UA (PLANNED)
        - Serial [0-*]          Service to broker I/O to serial controllers (PLANNED)
	- Prediction [1]            Service for on-machine predictive analytics (PLANNED)
    - REST [1]                  Service for REST access to device (PLANNED)
*/

func startServices(supervisor *suture.Supervisor, logFunc func(string)) {

	configService := &services.ConfigService{LogFunc: logFunc,
		StartDelay: time.Duration(500 * time.Millisecond)}
	configService.Name = define.ConfigServiceName

	integrationsService := &services.IntegrationsService{LogFunc: logFunc,
		StartDelay: time.Duration(100 * time.Millisecond)}
	integrationsService.Name = define.IntegrationsServiceName
	integrationsService.AddServiceDependentUpon(define.ConfigServiceName)

	stateService := &services.StateService{LogFunc: logFunc,
		StartDelay: time.Duration(1000 * time.Millisecond)}
	stateService.Name = define.StateServiceName
	stateService.AddServiceDependentUpon(define.IntegrationsServiceName)

	fieldBusSupervisor := suture.New(define.FieldbusSupervisorName, defaultServiceSpec)
	supervisor.Add(fieldBusSupervisor)

	modbusTCPService := &fieldbus.ModbusTCPService{LogFunc: logFunc,
		StartDelay: time.Duration(1000 * time.Millisecond)}
	modbusTCPService.Name = define.ModbusTCPServiceName
	modbusTCPService.AddServiceDependentUpon(define.ConfigServiceName)
	fieldBusSupervisor.Add(modbusTCPService)

	supervisor.Add(configService)
	supervisor.Add(stateService)
	supervisor.Add(integrationsService)
}
