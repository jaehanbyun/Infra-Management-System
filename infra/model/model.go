package model

import "github.com/jaehanbyun/infra/data"

type DBHandler interface {
	GetClusterCount() (int, error)
	AddCluster(data.ClusterSpec)
	GetClusterSpec(clusterName string) (*data.ClusterSpec, error)
	Close()
	Init() error
}

func NewDBHandler() DBHandler {
	return newPostgresHandler()
}
