
NimbleIndustry: Device
=========
A fault-tolerant, extensible implementation of an IIoT gateway written in Go.

[![Build Status](https://travis-ci.org/nimbleindustry/device.svg?branch=master)](https://travis-ci.org/nimbleindustry/device) 
[![Go Report Card](https://goreportcard.com/badge/github.com/nimbleindustry/device)](https://goreportcard.com/report/github.com/nimbleindustry/device)

*Device* is our implemenation of an IIoT gateway. Typically IIoT gateways are deployed on/near industrial machines in order to shuttle certain [field bus](https://en.wikipedia.org/wiki/Fieldbus) messages over to modern IIoT platforms such as [Predix](https://www.predix.com/) or [SightMachine](http://sightmachine.com/). 

Overview
--------
This project contains source for a fully functional, standalone IIoT gateway. *Device* is typically deployed on a Linux-based headless computer (we like, and have tested on, the Intel NUC line outfitted with Ubuntu). The job of any IIoT gateway is to allow the transfer of operational data and sensor telemetry from industrial equipment to IIoT platforms. This project abides by providing the means to access one or more field bus networks such as Modbus; then automagically forwarding that data to one or more IIoT platforms.

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
Industrial settings are inhospitable places for computers. *Device* aims to be highly fault-tolerant. For instance if the serial connection to a field bus interface is interrupted, *Device* gracefully attempts to reconnect and uses exponential backoff techniques in respect of system resources. In this example, other field bus or IIoT connections would be unaffected.

Device uses a hierarchical services architecture based on [supervisor trees](https://github.com/nimbleindustry/suture).

### Extensible
Want to add your field bus or IIoT system? This project is specifically designed to easily add support for new fieldbus and IIoT integrations. See [this page](here) for an example of adding a new IIoT integration to *Device*. You can also [contact us](mailto:info@nimbleindustry.com) if you'd like to discuss custom integrations.

### Simplicity
*Device* is written in golang. Once built, the binary image has no external dependencies and can run on a Linux computer as a defined service. The design employs concurrency yet consumes a minimum of system resources. For instance, in our lab an outfitted Intel NUC running *Device* which is attached to a Modbus-based PLC and the Initial State service and an MQTT broker (Mosquitto) has been running for months with 100% uptime.

The Device image contains available built-in diagnostics (sampled profiling) to allow monitoring of system resource consumption. The Device also supports reporting of its own state to IIoT integrations.
  

Building
--------
The only requirement is [go](https://golang.org), version 1.6 or better (Note, this project uses golang *vendoring* for its external dependencies and these dependencies are committed to this repo).

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
Deploying *Device* is as simple as installing the compiled image on your chosen gateway computer. Ideally you should configure *Device* as a System-V or upstart service. For Ubuntu, see [this reference](https://help.ubuntu.com/community/UbuntuBootupHowto). 

### Process Flags

- -syslog: send the log advisory/error output to the default system log, usually syslog on Linux
- -profile: enable remote profiling inspection via HTTP port 6789

### Configuration Files
Three separate configuration files bind the *Device* to its specfic installation—the system is configured completely using these files.

The *config service* monitors these files and restarts affected services automatically. The files are expected to be found in ```/etc/opt/nimble``` on Linux. When running (testing) on a Mac or Windows computer, the files are expected to be in the relative, local directory named ```./conf```. Examples of these files can be seen in the ```conf``` directory that is part of this project.

##### Equipment Configuration
```equipment.json``` defines information about the industrial equipment onto which the *Device* is being integrated. It contains the equipment's model number for instance as well as the field bus or other integration points. Theoretically, all machines from the same manufacturer, with the same model number, and the same firmware should be able to share this configuration file.

We are building a web application (machineconfig.com) that allows equipment manufacturers to edit, store and manage their equipment configurations in a centralized repository.

##### Asset Configuration
```asset.json``` identifies a piece of equipment in place. This file contains information such as the equipment identifier (e.g. robot9), the entity (the company using the equipment), the class of equipment (see above), and, for example, the assembly line in which the equipment is operating.

##### Connections Configuration
```connections.json``` defines the field bus and IIoT integrations that the Device should attempt to manage along with the endpoints and relevant connection access keys (if applicable).


Hardware Configuration
-------------------
The Device can run on pretty much any hardware/OS combination (x86, ARM, Linux, Windows, OSX). We can recommend (and have tested heavily on) the [Intel NUC](http://www.intel.com/content/www/us/en/nuc/overview.html) line of headless computers running Ubuntu 14.04.

IIoT Gateways often access sensor and operational data via serial communication lines or even direct-wire mechanisms. Because of this they are typically installed on or near the equipment they are monitoring. But that is definitely not a requirement.

If you'd prefer to get hardware that is configured for your equipment *and* preloaded with the *Device* image directly from us, please feel free to [reach out](mailto:info@nimbleindustry.com) for more information. 

Contributing
------------
Contributions to this project should follow the standard git model: 

	fork ➝ feature-branch ➝ pull-request 
	

Please base your pull request against our master branch and be sure to include details/reason for the change as well as the testing that was performed against the change.

Roadmap
-------
As time allows, we're working on improvements to this project and related technology by order of importance:

- more testing
- json schemas for configuration files
- additional fieldbus integrations (Modbus RTU, OPC/UA, CAN bus)
- additional IIoT platform integrations (Predix, AWS IoT, etc)
- equipment configuration editor (for equipment manufacturers and integrators) as an online app at http://machineconfig.com
- GPIO/ADC integrations
- onboard data historian (maybe using InfluxDB)
- onboard, sub-second predictions of sensor telemetry and operational data

License
-------
To encourage contributions, especially in the form of new field bus and IIoT integrations, this project is licensed under the GNU LGPLv3. The interfaces (as defined by the license) shall be considered those "Bus" mechanisms within that allow you to add functionality specfic and not needed by general use. 

In other words, if you add an OPC/UA integration, please share. If you add functionality that is specific to your company, your larger work can be distributed using different licensing terms.

