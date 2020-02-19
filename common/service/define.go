package service

type ServiceId string

const (
	DataHub   ServiceId = "datahub"
	BootMan   ServiceId = "bootman"
	ApiServer ServiceId = "apiserver"
	DeviceMan ServiceId = "deviceman"
	PipeTask  ServiceId = "pipetask"
)

type ServiceFunction string

const (
	FuncAny   ServiceFunction = "#"
	FuncGet   ServiceFunction = "get"
	FuncList  ServiceFunction = "list"
	FuncState ServiceFunction = "state" //TODO
)
