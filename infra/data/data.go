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
	NodeImage   string `json:"nodeImage"`
	FlavorVcpu  int    `json:"flavorVcpu"`
	FlavorRam   int    `json:"flavorRam"`
	FlavorDisk  int    `json:"flavorDisk"`
	MasterCount int    `json:"masterCount"`
	WorkerCount int    `json:"workerCount"`
}

type DeleteClusterReq struct {
	ClusterName string `json:"clusterName"`
}
