package model

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jaehanbyun/infra/data"
	_ "github.com/lib/pq"
)

type postgresHandler struct {
	db *sql.DB
}

func (h *postgresHandler) Init() error {
	_, err := h.db.Exec(`
		CREATE TABLE IF NOT EXISTS clusters (
			id SERIAL PRIMARY KEY,
			cluster_name TEXT NOT NULL,
			worker_count INT NOT NULL,
			master_count INT NOT NULL,
			node_image TEXT NOT NULL,
			flavor_vcpu INT NOT NULL,
			flavor_ram INT NOT NULL,
			flavor_disk INT NOT NULL
		)
	`)
	return err
}

func (p *postgresHandler) Close() {
	p.db.Close()
}

func (p *postgresHandler) GetClusterCount() (int, error) {
	row := p.db.QueryRow("SELECT COUNT(*) FROM clusters")
	var count int
	err := row.Scan(&count)
	return count, err
}

func (p *postgresHandler) AddCluster(cluster data.ClusterSpec) {
	_, err := p.db.Exec(`INSERT INTO clusters(cluster_name, worker_count, master_count, node_image, flavor_vcpu, flavor_ram, flavor_disk) 
	VALUES($1, $2, $3, $4, $5, $6, $7)`,
		cluster.ClusterName, cluster.WorkerCount, cluster.MasterCount, cluster.NodeImage, cluster.FlavorVcpu, cluster.FlavorRam, cluster.FlavorDisk)
	if err != nil {
		log.Printf("failed to insert a cluster: %v", err)
	}
}

func (p *postgresHandler) DeleteCluster(clusterName string) error {
	_, err := p.db.Exec(`DELETE FROM clusters WHERE cluster_name = $1`, clusterName)
	if err != nil {
		log.Printf("Failed to delete the cluster with name: %s, Error: %v", clusterName, err)
		return err
	}
	return nil
}

func (p *postgresHandler) GetClusterSpec(clusterName string) (*data.ClusterSpec, error) {
	row := p.db.QueryRow(`SELECT cluster_name, worker_count, master_count, node_image, flavor_vcpu, flavor_ram, flavor_disk 
	FROM clusters WHERE cluster_name = $1`, clusterName)

	var cluster data.ClusterSpec
	err := row.Scan(&cluster.ClusterName, &cluster.WorkerCount, &cluster.MasterCount, &cluster.NodeImage, &cluster.FlavorVcpu, &cluster.FlavorRam, &cluster.FlavorDisk)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func newPostgresHandler() DBHandler {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "cluster",
	)

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	return &postgresHandler{database}
}
