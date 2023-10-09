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
			flavor_disk INT NOT NULL,
			status TEXT NOT NULL,
			bastion_ip TEXT
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

func (p *postgresHandler) GetAllClusters() ([]data.ClusterSpec, error) {
	rows, err := p.db.Query("SELECT cluster_name, worker_count, master_count, node_image, flavor_vcpu, flavor_ram, flavor_disk, status, bastion_ip FROM clusters")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clusters []data.ClusterSpec
	for rows.Next() {
		var cluster data.ClusterSpec
		var bastionIp sql.NullString

		if err := rows.Scan(&cluster.ClusterName, &cluster.WorkerCount, &cluster.MasterCount, &cluster.NodeImage, &cluster.FlavorVcpu, &cluster.FlavorRam, &cluster.FlavorDisk, &cluster.Status, &bastionIp); err != nil {
			return nil, err
		}

		cluster.BastionIP = bastionIp.String

		clusters = append(clusters, cluster)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clusters, nil
}

func (p *postgresHandler) AddCluster(cluster data.ClusterSpec) {
	fmt.Println(cluster)
	_, err := p.db.Exec(`INSERT INTO clusters(cluster_name, worker_count, master_count, node_image, flavor_vcpu, flavor_ram, flavor_disk, status) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		cluster.ClusterName, cluster.WorkerCount, cluster.MasterCount, cluster.NodeImage, cluster.FlavorVcpu, cluster.FlavorRam, cluster.FlavorDisk, "Creating")
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
	row := p.db.QueryRow(`SELECT cluster_name, worker_count, master_count, node_image, flavor_vcpu, flavor_ram, flavor_disk, status, bastion_ip
	FROM clusters WHERE cluster_name = $1`, clusterName)

	var cluster data.ClusterSpec
	err := row.Scan(&cluster.ClusterName, &cluster.WorkerCount, &cluster.MasterCount, &cluster.NodeImage, &cluster.FlavorVcpu, &cluster.FlavorRam, &cluster.FlavorDisk, &cluster.Status, &cluster.BastionIP)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func (p *postgresHandler) UpdateClusterStatusAndIP(clusterName, status, bastionIP string) error {
	_, err := p.db.Exec(`UPDATE clusters SET status = $1, bastion_ip = $2 WHERE cluster_name = $3`,
		status, bastionIP, clusterName)
	if err != nil {
		log.Printf("Failed to update the cluster status and IP for: %s, Error: %v", clusterName, err)
		return err
	}
	return nil
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
