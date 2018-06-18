package kubepod

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kolide/osquery-go/plugin/table"
)

// Pod type
type pod struct {
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		UID       string `json:"uid"`
	} `json:"metadata"`
	Spec struct {
		NodeName string `json:"nodeName"`
	} `json:"spec"`
	Status struct {
		Phase     string    `json:"phase"`
		HostIP    string    `json:"hostIP"`
		PodIP     string    `json:"podIP"`
		StartTime time.Time `json:"startTime"`
	} `json:"status"`
}

// Podlist type
type podlist struct {
	Items []pod `json:"items"`
}

// Return pertinent pod info
func (c *pod) GetPodInfo() map[string]string {
	return map[string]string{
		"name":      c.Metadata.Name,
		"namespace": c.Metadata.Namespace,
		"node":      c.Spec.NodeName,
		"status":    c.Status.Phase,
		"host_ip":   c.Status.HostIP,
		"pod_ip":    c.Status.PodIP,
	}
}

// Generate list of pod info
func (c *podlist) GenerateTable() []map[string]string {
	results := make([]map[string]string, 0)
	for _, p := range c.Items {
		results = append(results, p.GetPodInfo())
	}
	return results
}

// FoobarColumns returns the columns that our table will return.
func KubePodColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("namespace"),
		table.TextColumn("node"),
		table.TextColumn("status"),
		table.TextColumn("host_ip"),
		table.TextColumn("pod_ip"),
	}
}

// FoobarGenerate will be called whenever the table is queried. It should return
// a full table scan.
func KubePodGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {

	url := "http://127.0.0.1:10255/pods"
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "osquery")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println("Failed to query the kubelet API:", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return nil, err
	}

	pods := podlist{}
	jsonerr := json.Unmarshal(body, &pods)
	if jsonerr != nil {
		log.Println("Failed to parse PodList result:", jsonerr)
		return nil, jsonerr
	}

	return pods.GenerateTable(), nil

}
