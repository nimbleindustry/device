package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nimbleindustry/device/common"
)

const (
	apiVersion         = "0.4.0"
	mediaType          = "application/json"
	accessKeyHeader    = "X-IS-AccessKey"
	buckerKeyHeader    = "X-IS-BucketKey"
	maxIdleConnections = 10
	requestTimeout     = 5
	events             = "events"
	buckets            = "buckets"
	versions           = "versions"
	iso8601            = "2006-01-02T15:04:05Z"
)

// InitialState implements an integration with the InitialState service (initialstate.com)
//
// Bucket keys:
// -- for Device state: state|<entity>/<location>/<machineId>
//    The value pairs in this bucket are fixed to memory and disk consumption along with
//    system load status
// -- for Device telemetry and ops: ops|<entity>/<location>/<machineId>
//    The value pairs in this bucket are dependent on the registerName entries for each
//    field bus entry in the machineIntegration section of the equipment configuration object.
type InitialState struct {
	client *http.Client
	record common.ConnectionRecord
}

// SetRecord sets the passed connection record
func (i *InitialState) SetRecord(record common.ConnectionRecord) {
	i.record = record
}

// Record returns the associated connection information
func (i *InitialState) Record() *common.ConnectionRecord {
	return &i.record
}

// Connect attempts a connection to the InitialState service
// Returns an error if the connection experiences a problem
func (i *InitialState) Connect() (err error) {
	i.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: time.Duration(requestTimeout) * time.Second,
	}
	if _, err = i.sendRequest("GET", versions, "", nil); err != nil {
		return
	}
	// create the state and telemetry buckets (null op if already created)
	if _, err = i.sendRequest("POST", buckets, i.getStateBucketName(), ""); err != nil {
		return
	}
	_, err = i.sendRequest("POST", buckets, i.getOpsAndTelemetryBucketName(), "")
	return
}

// Close the connection
func (i *InitialState) Close() error {
	i.client = nil
	return nil
}

// SendData transmits telemetry and operations data to the Initial State service
func (i *InitialState) SendData(data map[string]interface{}) error {
	timestamp := time.Now()
	body := i.transformOpsData(timestamp, data)
	if false {
		debugData("data send:", body)
		return nil
	}
	_, err := i.sendRequest("POST", events, i.getOpsAndTelemetryBucketName(), body)
	return err
}

// SendState transmits Device state to the Initial State Service
func (i *InitialState) SendState(data *common.SystemState) error {
	body := i.transformStateData(data)
	if false {
		debugData("state send:", body)
		return nil
	}
	_, err := i.sendRequest("POST", events, i.getStateBucketName(), body)
	return err
}

// ReceiveData is not implemented
func (i *InitialState) ReceiveData(data interface{}) error {
	return nil
}

func debugData(prefix string, data interface{}) {
	var prettyJSON bytes.Buffer
	b, _ := json.Marshal(data)
	json.Indent(&prettyJSON, b, "", "\t")
	fmt.Println(prefix, string(prettyJSON.Bytes()))
}

func (i *InitialState) getStateBucketName() string {
	return fmt.Sprintf("state|%s/%s/%s", common.AssetConfig.Entity, common.AssetConfig.Location, common.AssetConfig.MachineID)
}

func (i *InitialState) getOpsAndTelemetryBucketName() string {
	return fmt.Sprintf("ops|%s/%s/%s", common.AssetConfig.Entity, common.AssetConfig.Location, common.AssetConfig.MachineID)
}

func (i *InitialState) sendRequest(verb string, resource string, bucket string, body interface{}) (reply []byte, err error) {
	var req *http.Request
	uri := i.record.Endpoint + "/" + resource

	if body != nil {
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return
		}
		req, err = http.NewRequest(verb, uri, buf)
		if err != nil {
			return
		}
	} else {
		req, err = http.NewRequest(verb, uri, nil)
		if err != nil {
			return
		}
	}
	req.Header.Set("Content-Type", mediaType)
	req.Header.Set(accessKeyHeader, i.record.ProviderKey)
	if len(bucket) > 0 {
		req.Header.Set(buckerKeyHeader, bucket)
	}
	response, err := i.client.Do(req)
	if err != nil && response == nil {
		return
	}
	defer response.Body.Close()
	reply, err = ioutil.ReadAll(response.Body)
	return
}

func (i *InitialState) transformStateData(record *common.SystemState) interface{} {
	timestamp := record.Timestamp.Format(iso8601)
	state := []struct {
		Iso8601 string  `json:"iso8601"`
		Key     string  `json:"key"`
		Value   float64 `json:"value"`
	}{
		{timestamp, "memoryConsumed", record.MemoryConsumed},
		{timestamp, "diskConsumed", record.DiskConsumed},
		{timestamp, "loadAverage", record.LoadAverage},
	}
	return state
}

func (i *InitialState) transformOpsData(timestamp time.Time, data map[string]interface{}) interface{} {
	ts := timestamp.Format(iso8601)
	type telemetry struct {
		Iso8601 string      `json:"iso8601"`
		Key     string      `json:"key"`
		Value   interface{} `json:"value"`
	}
	var telemetryArray []telemetry
	for k, v := range data {
		record := telemetry{}
		record.Iso8601 = ts
		record.Key = k
		record.Value = v
		telemetryArray = append(telemetryArray, record)
	}
	return telemetryArray
}
