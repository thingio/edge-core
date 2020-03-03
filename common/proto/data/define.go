package data

type TsColumn struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// a general TimeSeries dataset, TODO: should be compatible with popular tsdb
type TsDataSet struct {
	Columns TsColumn                 `json:"columns,omitempty"`
	Ts      []int64                  `json:"ts,omitempty"`
	Data    map[string][]interface{} `json:"data,omitempty"` // interface{} refers to float/int
}


type Alert struct {
	Id         string `json:"id,omitempty"`
	PipeTaskId string `json:"pipetask_id,omitempty"`
	Topic      string `json:"topic,omitempty"`
	Message    string `json:"message,omitempty"`
	Level      string `json:"level,omitempty"`
	Data       string `json:"data,omitempty"`  // better be []byte, but after json serialization, they will be base64 anyway
	Image      string `json:"image,omitempty"` // same as above
	Ts         string `json:"ts,omitempty"`
}