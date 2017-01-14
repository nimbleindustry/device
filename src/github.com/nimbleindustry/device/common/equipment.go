package common

// MLMap provides i18n support
type MLMap map[string]string

// MachineIntegration houses fieldbus integration configuration information
type MachineIntegration struct {
	ModbusEntries []ModbusEntry `json:"modbus,omitempty"`
}

// ModbusEntry defines a single modbus port's configuration information
type ModbusEntry struct {
	RegisterName string `json:"registerName"`
	Address      int    `json:"address"`
	Class        string `json:"class"`
	Desc         MLMap  `json:"desc"`
	Functions    []int  `json:"functions"`
}

// ErrorCode maps codes to (i18n) descriptions
type ErrorCode struct {
	Code string `json:"code"`
	Desc MLMap  `json:"desc"`
}

// Maintenance defines a maintenance entry for specific equipment
type Maintenance struct {
	ID       string `json:"id"`
	Title    MLMap  `json:"title"`
	Desc     MLMap  `json:"desc"`
	Interval string `json:"interval"`
	URI      string `json:"uri"`
}

// Equipment collects configuration information about a machien
type Equipment struct {
	Ref                 string             `json:"ref"`
	URI                 string             `json:"uri"`
	Entity              string             `json:"entity"`
	Models              []string           `json:"models"`
	Versions            []string           `json:"versions"`
	Desc                MLMap              `json:"desc"`
	MachineIntegrations MachineIntegration `json:"machineIntegration"`
	ErrorCodes          []ErrorCode        `json:"errorCodes"`
	MaintenanceOps      []Maintenance      `json:"maintenance"`
}

// FindModbusEntriesByFunction returns all ModbusEntry 's defined in the equipment configuration object
// using the modbus function type as a search parameter. If no entries are found, an empty slice is returned.
func (machineIntegration MachineIntegration) FindModbusEntriesByFunction(functions []int) (entries []ModbusEntry) {
	for _, v := range functions {
		for _, entry := range machineIntegration.ModbusEntries {
			for _, function := range entry.Functions {
				if v == function {
					entries = append(entries, entry)
				}
			}
		}
	}
	return
}

// FindModbusEntriesByClass returns all ModbusEntry 's defined in the equipment configuration using the class
// of entry as a search parameter. If no entries are found, an empty slice is returned.
func (machineIntegration MachineIntegration) FindModbusEntriesByClass(class string) (entries []ModbusEntry) {
	for _, entry := range machineIntegration.ModbusEntries {
		if entry.Class == class {
			entries = append(entries, entry)
		}
	}
	return
}
