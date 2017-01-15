
NimbleIndustry: Device
=========
A fault-tolerant, extensible implementation of an IIoT gateway written in Go.

[![Build Status](https://travis-ci.org/nimbleindustry/device.svg?branch=master)](https://travis-ci.org/nimbleindustry/device) 

*Device* is our implemenation of an IIoT gateway. Typically IIoT gateways are deployed on/near industrial machines in order to shuttle certain [field bus](https://en.wikipedia.org/wiki/Fieldbus) messages over to modern IIoT platforms such as [Predix](https://www.predix.com/) or [SightMachine](http://sightmachine.com/). 

Overview
--------
This project contains source for a fully functional, standalone IIoT gateway. *Device* is typically deployed on a Linux-based headless computer (we like, and have tested on, the Intel NUC line outfitted with Ubuntu). The job of any IIoT gateway is to allow the transfer of operational data and sensor telemetry from industrial equipment to IIoT platforms. This project abides by providing access to one or more field bus networks such as Modbus as well as access to one or more IIoT platforms.

#### Field Bus Integration
- Modbus TCP
- Modbus RTU (in development)
- OPC/UA (planned)
- CAN Bus (planned)
- Generic Serial (planned)
- Direct Wire (planned)

#### IIoT Integration
- Initial State (http://initialstate.com)
- Generic MQTT
- Amazon IoT (in development)
- SightMachine (http://sightmachine.com, planned)
- Predix (planned)


### Fault Tolerance
Industrial settings are inhospital places for computers. *Device* aims to be highly fault-tolerant. For instance if the serial connection to a fieldbus interface is interrupted, *Device* gracefully attempts to reconnect and uses exponential backoff techniques in respect of system resources. Other fieldbus or IIoT connections would be unaffected.

Device uses a hierarchical services architecture based on [supervisor trees](https://github.com/nimbleindustry/suture).

### Extensible
Want to add your field bus or IIoT system? This project is specifically designed to easily add support for new fieldbus and IIoT integrations. See [this page](here) for an example of adding a new IIoT integration to *Device*. You can also [contact us](mailto:info@nimbleindustry.com) if you'd like to discuss custom integrations.

### Simplicity
*Device* is written in golang. Once built, the binary image has no external dependencies and can run on a Linux computer as a defined service. The design employs concurrency yet consumes very little system resources. For instance, an outfitted Intel NUC running *Device* which is attached to a Modbus-based PLC and the Initial State service and an MQTT broker (Mosquitto) has been running for months with 100% uptime.

Built-in diagnostics are available (sampled profiling) to allow monitoring of system resource consumption.
  

Building
--------
Requirements

- go, version 1.6 or better (Note, this project uses golang *vendoring* for its external dependencies. These dependencies are committed to this repo 

```bash
$ git clone http://github.com/nimbleindustry/device
$ cd device
$ export GOPATH=`pwd`
$ go get github.com/nimbleindustry/suture
$ cd src/github.com/nimbleindustry/device
$ go build -o device main.go
```

To build a Linux image (if you're building on a Mac or Windows computer), build this way

```bash
$ GOOS=linux go build -o device main.go
``` 

You can also [contact us](mailto:info@nimbleindustry.com) if you'd like to purchase prebuilt images installed and configured on Intel NUC node computers.

Testing
-------
```bash
$ cd ${GOPATH}/src/github.com/nimbleindustry/device
$ go get github.com/stretchr/testify/assert
$ ./test.sh
```


Deploying
---------
Deploying *Device* is as simple as installing the compiled image on your chosen gateway computer. Ideally you should configure *Device* as a System-V or upstart service, for Ubuntu, see [this reference](https://help.ubuntu.com/community/UbuntuBootupHowto). 

### Process Flags

- -syslog: send the log output to the default system log
- -profile: enable remote profiling inspection via HTTP port 6789

### Configuration Files



Hardware Configuration
-------------------

Contributing
------------

Roadmap
-------



