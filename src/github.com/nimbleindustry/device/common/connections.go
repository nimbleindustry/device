package common

// Protocol/Fieldbus Definitions
const (
	Modbus    = "modbus"
	ModbusTCP = "modbusTCP"
	ModbusRTU = "modbusRTU"
	OPCUA     = "OPCUA"
)

// ConnectionRecord defines fieldbus and IIoT integration specifics
type ConnectionRecord struct {
	Provider    string `json:"provider"`
	Type        string `json:"type"`
	Endpoint    string `json:"endpoint"`
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	ProviderKey string `json:"providerKey, omitempty"`
	Baudrate    int    `json:"baudRate, omitempty"`
}

// Connections defines the arrays of ConnectionRecords defined for the system
type Connections struct {
	DeviceID                          string             `json:"deviceId"`
	DeviceStateConnections            []ConnectionRecord `json:"deviceState"`
	HistorianConnections              []ConnectionRecord `json:"historian"`
	MachineConnections                []ConnectionRecord `json:"machineIntegration"`
	OperationsAndTelemetryConnections []ConnectionRecord `json:"machineOperationsAndTelemetry"`
}

// GetMachineConnection returns a ConnectionRecord from the stored MachineConnections
// that matches the supplied class
func (conn Connections) GetMachineConnection(class string) (record ConnectionRecord) {
	if len(conn.MachineConnections) == 0 {
		return
	}
	for _, v := range conn.MachineConnections {
		if v.Type == class {
			return v
		}
	}
	return
}
