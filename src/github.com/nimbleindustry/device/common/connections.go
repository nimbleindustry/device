package common

// Protocol/Fieldbus Definitions
const (
	Modbus    = "modbus"
	ModbusTCP = "modbusTCP"
	ModbusRTU = "modbusRTU"
	OPCUA     = "OPCUA"
)

type ConnectionRecord struct {
	Provider    string `json:"provider"`
	Type        string `json:"type"`
	Endpoint    string `json:"endpoint"`
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	ProviderKey string `json:"providerKey, omitempty"`
	Baudrate    int    `json:"baudRate, omitempty"`
}

type Connections struct {
	DeviceID                          string             `json:"deviceId"`
	DeviceStateConnections            []ConnectionRecord `json:"deviceState"`
	HistorianConnections              []ConnectionRecord `json:"historian"`
	MachineConnections                []ConnectionRecord `json:"machineIntegration"`
	OperationsAndTelemetryConnections []ConnectionRecord `json:"machineOperationsAndTelemetry"`
}

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
