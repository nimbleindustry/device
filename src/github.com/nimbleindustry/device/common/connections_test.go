package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadModbusTCPDefinitions(t *testing.T) {
	var connections Connections
	err := json.Unmarshal([]byte(connectionFixture1), &connections)
	assert.Nil(t, err, "unmarshall failed")
	record := connections.GetMachineConnection(ModbusTCP)
	assert.NotEmpty(t, record.Type, "Expected to find non-empty modbusTCP entry")
}

func TestUnfoundConnectionDefinitions(t *testing.T) {
	var connections Connections
	err := json.Unmarshal([]byte(connectionFixture1), &connections)
	assert.Nil(t, err, "unmarshall failed")
	record := connections.GetMachineConnection("WillNotFind")
	assert.Empty(t, record.Type, "Expected to find no entries")
}

const connectionFixture1 = `{
    "deviceId": "0d80005e",
    "deviceState": [
        {
            "provider": "InitialState",
            "protocol": "REST",
            "providerKey": "m6XSk21RY8wHU00GHuh2NUBx2790n80o",
            "endpoint": "https://groker.initialstate.com/api"
        }
    ],
    "machineOperationsAndTelemetry": [
        {
            "provider": "InitialState",
            "protocol": "REST",
            "providerKey": "m6XSk21RY8wHU00GHuh2NUBx2790n80o",
            "endpoint": "https://groker.initialstate.com/api"
        }
    ],
    "machineIntegration": [
        {
            "type": "modbusTCP",
            "endpoint": "10.0.1.30"
        },
        {
            "type": "modbusRTU",
            "endpoint": "/dev/tty00"
        },
        {
            "type": "modbusTCP",
            "endpoint": "10.0.1.31"
        }        
    ],
    "historian": [
        {
            "type": "influxdb",
            "endpoint": "127.0.0.1",
            "port": 3001
        }
    ]
}`
