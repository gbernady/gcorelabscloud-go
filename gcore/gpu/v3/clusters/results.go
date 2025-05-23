package clusters

import (
	"encoding/json"
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a cluster resource.
func (r commonResult) Extract() (*Cluster, error) {
	var s Cluster
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Image.
type GetResult struct {
	commonResult
}

type ExternalInterface struct {
	Name     *string      `json:"name"`
	Type     string       `json:"type"`
	IPFamily IPFamilyType `json:"ip_family"`
}

type FloatingIP struct {
	Source string `json:"source"`
}

type SubnetInterface struct {
	NetworkID  string      `json:"network_id"`
	Name       *string     `json:"name"`
	Type       string      `json:"type"`
	SubnetID   string      `json:"subnet_id"`
	FloatingIP *FloatingIP `json:"floating_ip"`
}

type AnySubnetInterface struct {
	NetworkID  string       `json:"network_id"`
	Name       *string      `json:"name"`
	Type       string       `json:"type"`
	IPFamily   IPFamilyType `json:"ip_family"`
	IPAddress  *string      `json:"ip_address"`
	FloatingIP *FloatingIP  `json:"floating_ip"`
}

type InterfaceUnion struct {
	ExternalInterface  *ExternalInterface
	SubnetInterface    *SubnetInterface
	AnySubnetInterface *AnySubnetInterface
}

func (i *InterfaceUnion) InterfaceType() string {
	if i.ExternalInterface != nil {
		return i.ExternalInterface.Type
	}
	if i.SubnetInterface != nil {
		return i.SubnetInterface.Type
	}
	if i.AnySubnetInterface != nil {
		return i.AnySubnetInterface.Type
	}
	return ""
}

// ExtractInterfaceType extracts the interface type from the data, returning the type if it is valid.
// If the type is not valid, an error is returned. If the type is not found, an error is returned.
func (i *InterfaceUnion) ExtractInterfaceType(data []byte) (string, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return "", err
	}

	interfaceType, ok := raw["type"]
	if !ok {
		return "", fmt.Errorf("interface type not specified, unable to unmarshal interface")
	}
	allTypes := []string{"external", "subnet", "any_subnet"}
	for _, t := range allTypes {
		if interfaceType == t {
			return t, nil
		}
	}
	return "", fmt.Errorf("invalid interface type: %s", interfaceType)
}

func (i *InterfaceUnion) MarshalJSON() ([]byte, error) {
	if i.ExternalInterface != nil {
		return json.Marshal(i.ExternalInterface)
	}
	if i.SubnetInterface != nil {
		return json.Marshal(i.SubnetInterface)
	}
	if i.AnySubnetInterface != nil {
		return json.Marshal(i.AnySubnetInterface)
	}
	return nil, fmt.Errorf("no valid interface type")
}

func (i *InterfaceUnion) UnmarshalJSON(data []byte) error {
	interfaceType, err := i.ExtractInterfaceType(data)
	if err != nil {
		return err
	}
	if interfaceType == "external" {
		var iface ExternalInterface
		if err := json.Unmarshal(data, &iface); err != nil {
			return err
		}
		i.ExternalInterface = &iface
		return nil
	}
	if interfaceType == "subnet" {
		var iface SubnetInterface
		if err := json.Unmarshal(data, &iface); err != nil {
			return err
		}
		i.SubnetInterface = &iface
		return nil
	}
	if interfaceType == "any_subnet" {
		var iface AnySubnetInterface
		if err := json.Unmarshal(data, &iface); err != nil {
			return err
		}
		i.AnySubnetInterface = &iface
		return nil
	}
	return fmt.Errorf("no valid interface type: %s", interfaceType)
}

// Volume represents a volume structure.
type Volume struct {
	Size                 int        `json:"size"`
	Type                 VolumeType `json:"type"`
	DeletedOnTermination bool       `json:"deleted_on_termination"`
	Name                 *string    `json:"name"`
	BootIndex            *int       `json:"boot_index"`
	ImageID              *string    `json:"image_id"`
	SnapshotID           *string    `json:"snapshot_id"`
	Tags                 []Tag      `json:"tags"`
}

type ClusterServerSettings struct {
	Interfaces     []InterfaceUnion `json:"interfaces"`
	SecurityGroups []string         `json:"security_groups"`
	Volumes        []Volume         `json:"volumes"`
	UserData       string           `json:"user_data"`
	KeypairName    *string          `json:"keypair_name"`
}

type Cluster struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	Status          ClusterStatusType        `json:"status"`
	FlavorID        string                   `json:"flavor_id"`
	ServersCount    int                      `json:"servers_count"`
	CreatedAt       gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt       *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	ServersIDs      *[]string                `json:"servers_ids"`
	ServersSettings ClusterServerSettings    `json:"servers_settings"`
	Tags            []Tag                    `json:"tags"`
}

// Tag represents a cluster tag
type Tag struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	ReadOnly bool   `json:"read_only"`
}

// ClusterPage is the page returned by a pager when traversing over a collection of clusters.
type ClusterPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of AI Clusters has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ClusterPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a AIClusterPage struct is empty.
func (r ClusterPage) IsEmpty() (bool, error) {
	is, err := ExtractClusters(r)
	return len(is) == 0, err
}

// ExtractClusters accepts a Page struct, specifically a ClusterPage struct,
// and extracts the elements into a slice of Cluster structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractClusters(r pagination.Page) ([]Cluster, error) {
	var s []Cluster
	err := ExtractClustersInto(r, &s)
	return s, err
}

func ExtractClustersInto(r pagination.Page, v interface{}) error {
	return r.(ClusterPage).Result.ExtractIntoSlicePtr(v, "results")
}

type ClusterTaskResult struct {
	AIClusters []string `mapstructure:"ai_clusters"`
	// etc
}

func ExtractClusterIDFromTask(task *tasks.Task) (string, error) {
	var result ClusterTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode GPU cluster information in task structure: %w", err)
	}
	if len(result.AIClusters) == 0 {
		return "", fmt.Errorf("cannot decode GPU cluster information in task structure: empty list")
	}
	return result.AIClusters[0], nil
}
