package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	storagequotasets "github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	computequotasets "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gorilla/mux"
	"github.com/jaehanbyun/infra/data"
	"github.com/jaehanbyun/infra/model"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	rd *render.Render
)

const (
	JenkinsURL      = "http://192.168.10.20:9090"
	JenkinsUser     = "jaehan"
	JenkinsAPIToken = "11b790f6e2ec50b9555fa4685c9bafa050"
)

type AppHandler struct {
	http.Handler
	db model.DBHandler
}

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token ,Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")

			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func (a *AppHandler) getQuotaUsage(w http.ResponseWriter, r *http.Request) {
	authOpts := getAuthOpts()

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		http.Error(w, fmt.Sprintf("error setting authentification: %v", err), http.StatusInternalServerError)
	}

	opts := gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")}

	computeClient, err := openstack.NewComputeV2(provider, opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("error setting computeclient: %v", err), http.StatusInternalServerError)
	}

	storageClient, err := openstack.NewBlockStorageV3(provider, opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("error setting storageclient: %v", err), http.StatusInternalServerError)
	}

	computeQuota, err := computequotasets.Get(computeClient, os.Getenv("OS_USER_DOMAIN_NAME")).Extract()
	if err != nil {
		http.Error(w, fmt.Sprintf("error extracting computequota: %v", err), http.StatusInternalServerError)
	}

	storageQuota, err := storagequotasets.Get(storageClient, os.Getenv("OS_USER_DOMAIN_NAME")).Extract()
	if err != nil {
		http.Error(w, fmt.Sprintf("error extracting storagequota: %v", err), http.StatusInternalServerError)
	}

	allServers, err := servers.List(computeClient, servers.ListOpts{TenantID: os.Getenv("OS_USER_DOMAIN_ID")}).AllPages()
	if err != nil {
		http.Error(w, fmt.Sprintf("error listing servers: %v", err), http.StatusInternalServerError)
		return
	}

	instances, err := servers.ExtractServers(allServers)
	if err != nil {
		http.Error(w, fmt.Sprintf("error extracting servers: %v", err), http.StatusInternalServerError)
		return
	}

	totalCPUUsage := 0
	totalRAMUsage := 0
	totalStorageUsage := 0

	for _, instance := range instances {
		flavorID := instance.Flavor["id"].(string)
		flavor, err := flavors.Get(computeClient, flavorID).Extract()
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting flavor details: %v", err), http.StatusInternalServerError)
			return
		}
		totalCPUUsage += flavor.VCPUs
		totalRAMUsage += flavor.RAM
		totalStorageUsage += flavor.Disk
	}

	response := data.QuotaUsage{
		InstanceUsage: len(instances),
		InstanceLimit: computeQuota.Instances,
		CpuUsage:      totalCPUUsage,
		CpuLimit:      computeQuota.Cores,
		RamUsage:      totalRAMUsage / 1024,
		RamLimit:      computeQuota.RAM / 1024,
		StorageUsage:  totalStorageUsage,
		StorageLimit:  storageQuota.Gigabytes,
	}

	rd.JSON(w, http.StatusOK, response)
}

func (a *AppHandler) getClusterSpec(w http.ResponseWriter, r *http.Request) {
	clusterName := r.URL.Query().Get("clusterName")
	if clusterName == "" {
		http.Error(w, "Cluster name is required", http.StatusBadRequest)
		return
	}

	cluster, err := a.db.GetClusterSpec(clusterName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving cluster spec: %v", err), http.StatusInternalServerError)
		return
	}

	if cluster == nil {
		http.Error(w, "Cluster not found", http.StatusNotFound)
		return
	}

	rd.JSON(w, http.StatusOK, cluster)
}

func (a *AppHandler) createCluster(w http.ResponseWriter, r *http.Request) {
	var spec data.ClusterSpec

	err := json.NewDecoder(r.Body).Decode(&spec)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode cluster spec: %v", err), http.StatusBadRequest)
		return
	}

	a.db.AddCluster(spec)

	params := structToMap(spec)
	triggerJenkinsJob("createCluster", params)

	rd.JSON(w, http.StatusCreated, map[string]string{"message": "Cluster created successfully!"})
}

func (a *AppHandler) deleteCluster(w http.ResponseWriter, r *http.Request) {
	clusterName := data.DeleteClusterReq{
		ClusterName: r.URL.Query().Get("clusterName"),
	}

	if clusterName.ClusterName == "" {
		http.Error(w, "Cluster name is required", http.StatusBadRequest)
		return
	}

	params := structToMap(clusterName)
	triggerJenkinsJob("deleteCluster", params)

	// a.db.DeleteCluster(clusterName)

	rd.JSON(w, http.StatusCreated, map[string]string{"message": "Cluster deleted successfully!"})
}

func (a *AppHandler) getClusterCount(w http.ResponseWriter, r *http.Request) {
	count, err := a.db.GetClusterCount()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving cluster count: %v", err), http.StatusInternalServerError)
		return
	}

	rd.JSON(w, http.StatusOK, map[string]int{"count": count})
}

func (a *AppHandler) deleteTerraformStatefile(w http.ResponseWriter, r *http.Request) {
	authOpts := getAuthOpts()
	clusterName := "terraform-" + r.URL.Query().Get("clusterName") + "-tfstate"
	if clusterName == "" {
		http.Error(w, "Cluster name is required", http.StatusBadRequest)
		return
	}

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		http.Error(w, fmt.Sprintf("error setting authentication: %v", err), http.StatusInternalServerError)
		return
	}

	objectStoreClient, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")})
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating object storage client: %v", err), http.StatusInternalServerError)
		return
	}

	objectList, err := objects.List(objectStoreClient, clusterName, nil).AllPages()
	if err != nil {
		http.Error(w, fmt.Sprintf("error listing objects in container: %v", err), http.StatusInternalServerError)
		return
	}

	objectNames, err := objects.ExtractNames(objectList)
	if err != nil {
		http.Error(w, fmt.Sprintf("error extracting object names: %v", err), http.StatusInternalServerError)
		return
	}

	for _, objectName := range objectNames {
		_, err := objects.Delete(objectStoreClient, clusterName, objectName, nil).Extract()
		if err != nil {
			http.Error(w, fmt.Sprintf("error deleting object %s: %v", objectName, err), http.StatusInternalServerError)
			return
		}
	}

	_, err = containers.Delete(objectStoreClient, clusterName).Extract()
	if err != nil {
		http.Error(w, fmt.Sprintf("error deleting container: %v", err), http.StatusInternalServerError)
		return
	}

	rd.JSON(w, http.StatusOK, map[string]string{"message": "Container deleted successfully!"})
}

func MakeHandler() *AppHandler {
	rd = render.New()
	r := mux.NewRouter()
	r.Use(enableCORS)

	neg := negroni.Classic()
	neg.UseHandler(r)

	a := &AppHandler{
		Handler: neg,
		db:      model.NewDBHandler(),
	}

	err := a.db.Init()
	if err != nil {
		panic(err)
	}

	r.HandleFunc("/quota/usage", a.getQuotaUsage).Methods("GET")
	r.HandleFunc("/cluster/count", a.getClusterCount).Methods("GET")
	r.HandleFunc("/cluster/spec", a.getClusterSpec).Methods("GET")
	r.HandleFunc("/cluster/create", a.createCluster).Methods("POST")
	r.HandleFunc("/cluster/delete", a.deleteCluster).Methods("DELETE")
	r.HandleFunc("/terraform/statefile/delete", a.deleteTerraformStatefile).Methods("DELETE")

	return a
}
