package resource

type Node struct {
	Id        string            `json:"id,omitempty"`
	Os        string            `json:"os,omitempty"`
	Kernel    string            `json:"kernel,omitempty"`
	Arch      string            `json:"arch,omitempty"`
	BootTime  int64             `json:"boot_time,omitempty"`
	LocalTime string            `json:"local_time,omitempty"`
	Stats     map[string]string `json:"stats,omitempty"`
}
