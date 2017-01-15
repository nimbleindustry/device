
NimbleIndustry: Device
=========
A fault-tolerant, extensible implementation of an IIoT gateway written in Go.

[![Build Status](https://travis-ci.org/nimbleindustry/device.svg?branch=master)](https://travis-ci.org/nimbleindustry/device) 

*Device* is our implemenation of an IIoT gateway. Typically IIoT gateways are deployed on/near industrial machines in order to shuttle certain [field bus](https://en.wikipedia.org/wiki/Fieldbus) messages over to modern IIoT platforms such as [Predix](https://www.predix.com/) or [SightMachine](http://sightmachine.com/). 

Overview
--------
This project contains source for a fully functional, standalone IIoT gateway. *Device* is typically deployed on a Linux-based headless CPU (we like, and have tested on, the Intel NUC line). The job of any IIoT gateway is to allow the transfer of operational data and sensor telemetry from industrial equipment to IIoT platforms. This project abides by providing access to one or more field bus networks such as Modbus as well as access to one or more IIoT platforms.

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
Industrial settings are inhospital places for computers. *Device* aims to be highly fault-tolerant. For instance if the serial connection to a fieldbus interface is interrupted, *Device* gracefully attempts to reconnect and uses exponential backoff techniques.

### Extensible
This project is specifically designed to easily add new support for fieldbus and IIoT integrations.  

Building
--------

Deploying
---------
- configuration files

Testing
-------

Hardware Configuration
-------------------

Roadmap
-------



