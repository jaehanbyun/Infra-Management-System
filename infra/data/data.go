package data

type QuotaUsage struct {
	InstanceUsage int `json:"instance_usage"`
	InstanceLimit int `json:"instance_limit"`
	CpuUsage      int `json:"cpu_usage"`
	CpuLimit      int `json:"cpu_limit"`
	RamUsage      int `json:"ram_usage"`
	RamLimit      int `json:"ram_limit"`
	StorageUsage  int `json:"storage_usage"`
	StorageLimit  int `json:"storage_limit"`
}

type ClusterSpec struct {
	ClusterName string `json:"clusterName"`
	WorkerCount int    `json:"workerCount"`
	MasterCount int    `json:"masterCount"`
	NodeImage   string `json:"nodeImage"`
	FlavorVcpu  int    `json:"flavorVcpu"`
	FlavorRam   int    `json:"flavorRam"`
	FlavorDisk  int    `json:"flavorDisk"`
	Status      string `json:"status,omitempty"`
	BastionIP   string `json:"bastionIP,omitempty"`
}

type DeleteClusterReq struct {
	ClusterName string `json:"clusterName"`
}

type ImageInfo struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type JenkinsCompletionNotification struct {
	ClusterName string `json:"clusterName"`
	BastionIP   string `json:"bastionIP"`
}
