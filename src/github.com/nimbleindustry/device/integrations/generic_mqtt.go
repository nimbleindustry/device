package integrations

import (
	"encoding/json"
	"errors"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"github.com/nimbleindustry/device/common"
)

// GenericMQTT implements an integration with MQTT brokers
type GenericMQTT struct {
	client MQTT.Client
	record common.ConnectionRecord
}

// SetRecord associates the passed connection record
func (i *GenericMQTT) SetRecord(record common.ConnectionRecord) {
	i.record = record
}

// Record returns the associated connection information
func (i *GenericMQTT) Record() *common.ConnectionRecord {
	return &i.record
}

// Connect attempts a connection to the GenericMQTT broker
func (i *GenericMQTT) Connect() (err error) {
	if len(i.record.Endpoint) == 0 {
		err = errors.New("GenericMQTT endpoint setting unexpectedly nil")
		return
	}
	opts := MQTT.NewClientOptions()
	opts.AddBroker(i.record.Endpoint)
	opts.SetClientID(common.AssetConfig.MachineID)
	i.client = MQTT.NewClient(opts)
	if token := i.client.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
	}
	return
}

// Close the connection to the MQTT broker
func (i *GenericMQTT) Close() error {
	i.client.Disconnect(1000)
	return nil
}

type mqttDataMessage struct {
	Timestamp time.Time     `json:"ts"`
	Tags      *common.Asset `json:"tags"`
	Body      interface{}   `json:"body"`
}

// SendData transmits telemetry and operations data to the MQTT broker
func (i *GenericMQTT) SendData(data map[string]interface{}) error {
	msg := &mqttDataMessage{}
	msg.Timestamp = time.Now()
	msg.Tags = &common.AssetConfig
	msg.Body = data
	if false {
		debugData("mqtt data send:", msg)
		return nil
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	topic := "/ops"
	token := i.client.Publish(topic, 1, false, bytes)
	token.Wait()
	return nil
}

// SendState transmits Device state to the MQTT broker
func (i *GenericMQTT) SendState(data *common.SystemState) error {
	msg := &mqttDataMessage{}
	msg.Timestamp = time.Now()
	msg.Tags = &common.AssetConfig
	msg.Body = data
	if false {
		debugData("mqtt state send:", msg)
		return nil
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	topic := "/state"
	token := i.client.Publish(topic, 1, false, bytes)
	token.Wait()
	return nil
}

// ReceiveData is not implemented
func (i *GenericMQTT) ReceiveData(data interface{}) error {
	return nil
}
