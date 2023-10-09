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
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
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
		w.Header().Set("Access-Control-Allow-Origin", "http://192.168.10.20:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Allow", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
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
	deleteClusterReq := data.DeleteClusterReq{
		ClusterName: r.URL.Query().Get("clusterName"),
	}

	if deleteClusterReq.ClusterName == "" {
		http.Error(w, "Cluster name is required", http.StatusBadRequest)
		return
	}

	params := structToMap(deleteClusterReq)
	triggerJenkinsJob("deleteCluster", params)

	a.db.DeleteCluster(deleteClusterReq.ClusterName)

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

func (a *AppHandler) getImageInfo(w http.ResponseWriter, r *http.Request) {
	authOpts := getAuthOpts()

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		http.Error(w, fmt.Sprintf("error setting authentication: %v", err), http.StatusInternalServerError)
		return
	}

	opts := gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")}
	imageClient, err := openstack.NewImageServiceV2(provider, opts)
	if err != nil {
		http.Error(w, fmt.Sprintf("error setting image client: %v", err), http.StatusInternalServerError)
		return
	}

	allPages, err := images.List(imageClient, images.ListOpts{}).AllPages()
	if err != nil {
		http.Error(w, fmt.Sprintf("error listing images: %v", err), http.StatusInternalServerError)
		return
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		http.Error(w, fmt.Sprintf("error extracting images: %v", err), http.StatusInternalServerError)
		return
	}

	var imageInfos []data.ImageInfo
	for _, img := range allImages {
		info := data.ImageInfo{
			Name: img.Name,
			UUID: img.ID,
		}
		imageInfos = append(imageInfos, info)
	}

	rd.JSON(w, http.StatusOK, imageInfos)
}

func (a *AppHandler) getAllClusters(w http.ResponseWriter, r *http.Request) {
	clusters, err := a.db.GetAllClusters()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving all clusters: %v", err), http.StatusInternalServerError)
		return
	}

	if len(clusters) == 0 {
		rd.JSON(w, http.StatusNotFound, map[string]string{"message": "No clusters found"})
		return
	}

	rd.JSON(w, http.StatusOK, clusters)
}

func (a *AppHandler) jenkinsCompletionNotify(w http.ResponseWriter, r *http.Request) {
	var notif data.JenkinsCompletionNotification

	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode Jenkins completion notification: %v", err), http.StatusBadRequest)
		return
	}

	err = a.db.UpdateClusterStatusAndIP(notif.ClusterName, "Active", notif.BastionIP)
	if err != nil {
		http.Error(w, fmt.Sprintf("error updating cluster status: %v", err), http.StatusInternalServerError)
		return
	}

	rd.JSON(w, http.StatusOK, map[string]string{"message": "Cluster updated successfully!"})
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
	r.HandleFunc("/images", a.getImageInfo).Methods("GET")
	r.HandleFunc("/cluster/count", a.getClusterCount).Methods("GET")
	r.HandleFunc("/clusters", a.getAllClusters).Methods("GET")
	r.HandleFunc("/cluster/spec", a.getClusterSpec).Methods("GET")
	r.HandleFunc("/cluster/create", a.createCluster).Methods("POST", "OPTIONS")
	r.HandleFunc("/cluster/delete", a.deleteCluster).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/jenkins/notify-completion", a.jenkinsCompletionNotify).Methods("POST")

	return a
}
