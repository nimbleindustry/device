package define

// Global Definitions
const (
	SystemName = "nimble-device"
)

// Protocol/Fieldbus Definitions
const (
	Modbus    = "modbus"
	ModbusTCP = "modbusTCP"
	ModbusRTU = "modbusRTU"
	OPCUA     = "OPCUA"
)

// Supervisor and Service identifiers
const (
	MasterSupervisorName    = "MasterSupervisor"
	ConfigServiceName       = "ConfigService"
	IntegrationsServiceName = "IntegrationsService"
	StateServiceName        = "StateService"
	FieldbusSupervisorName  = "FieldbusSupervisor"
	ModbusRTUServiceName    = "ModbusRTUService"
	ModbusTCPServiceName    = "ModbusTCPService"
	OPCUAServiceName        = "OPCUAService"
	SerialServiceName       = "SerialService"
	GatewaySupervisorName   = "GatewaySupervisor"
	MQTTServiceName         = "MQTTService"
	RESTGatewayServiceName  = "RESTGatewayService"
	PredictionServiceName   = "PredictionService"
	RESTServiceName         = "RESTService"
)

// Internal messaging topics
const (
	// Files have been updated
	EquipmentConfigUpdated    = "EquipmentConfigUpdated"
	AssetConfigUpdated        = "AssetConfigUpdated"
	ConnectivityConfigUpdated = "ConnectivityConfigUpdate"

	// Messages from the state service
	TopicStateReport = "TopicStateReport"

	// Messages from field bus integrations
	TopicOpsReport = "TopicOpsReport"
)

// Service Providers
const (
	InitialState = "InitialState"
	GenericMQTT  = "GenericMQTT"
	Predix       = "Predix"
	AWS          = "AWS"
	SightMachine = "SightMachine"
)
