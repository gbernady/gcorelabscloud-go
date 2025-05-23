package clusters

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

const (
	clustersPath = "clusters"
	actionPath   = "action"
)

// ClustersURL returns URL for GPU clusters operations
func ClustersURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL(clustersPath)
}

// ClusterURL returns URL for specific GPU cluster operations
func ClusterURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return c.ServiceURL(clustersPath, clusterID)
}

// ClusterActionURL returns URL for performing an action on a GPU cluster
func ClusterActionURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return c.ServiceURL(clustersPath, clusterID, actionPath)
}
