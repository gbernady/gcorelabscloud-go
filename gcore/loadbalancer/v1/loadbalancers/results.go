package loadbalancers

import (
	"fmt"
	"net"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/ddos/v1/ddos"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/lbflavors"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a loadbalancer resource.
func (r commonResult) Extract() (*LoadBalancer, error) {
	var s LoadBalancer
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a LoadBalancer.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a LoadBalancer.
type UpdateResult struct {
	commonResult
}

// LoadBalancer represents a loadbalancer structure.
type LoadBalancer struct {
	Name                  string                          `json:"name"`
	ID                    string                          `json:"id"`
	ProvisioningStatus    types.ProvisioningStatus        `json:"provisioning_status"`
	OperationStatus       types.OperatingStatus           `json:"operating_status"`
	VipAddress            net.IP                          `json:"vip_address"`
	VipPortID             string                          `json:"vip_port_id"`
	Listeners             []gcorecloud.ItemID             `json:"listeners"`
	CreatorTaskID         *string                         `json:"creator_task_id"`
	TaskID                *string                         `json:"task_id"`
	CreatedAt             gcorecloud.JSONRFC3339Z         `json:"created_at"`
	UpdatedAt             *gcorecloud.JSONRFC3339Z        `json:"updated_at"`
	ProjectID             int                             `json:"project_id"`
	RegionID              int                             `json:"region_id"`
	Region                string                          `json:"region"`
	Tags                  []string                        `json:"tags"`
	Flavor                lbflavors.Flavor                `json:"flavor"`
	Metadata              []metadata.Metadata             `json:"metadata"`
	DdosProfile           *ddos.Profile                   `json:"ddos_profile"`
	VrrpIPs               []NetworkPortFixedIP            `json:"vrrp_ips"`
	VipIPFamilyType       *types.IPFamilyType             `json:"vip_ip_family"`
	AdditionalVips        []NetworkPortFixedIP            `json:"additional_vips"`
	FloatingIPs           []instances.FloatingIP          `json:"floating_ips"`
	Logging               *Logging                        `json:"logging"`
	PreferredConnectivity types.PreferredConnectivityType `json:"preferred_connectivity"`
}

// NetworkPortFixedIP represents VRRP entry structure.
type NetworkPortFixedIP struct {
	IpAddress net.IP `json:"ip_address"`
	SubnetID  string `json:"subnet_id"`
}

func (lb LoadBalancer) IsDeleted() bool {
	return lb.ProvisioningStatus == types.ProvisioningStatusDeleted
}

// Logging represents logging configuration for a loadbalancer.
type Logging struct {
	Enabled             bool             `json:"enabled"`
	TopicName           string           `json:"topic_name"`
	DestinationRegionID int              `json:"destination_region_id"`
	RetentionPolicy     *RetentionPolicy `json:"retention_policy"`
}

// RetentionPolicy represents log retention policy configuration for a loadbalancer.
type RetentionPolicy struct {
	Period int `json:"period"`
}

// LoadBalancerPage is the page returned by a pager when traversing over a
// collection of loadbalancers.
type LoadBalancerPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of loadbalancers has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r LoadBalancerPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a LoadBalancerPage struct is empty.
func (r LoadBalancerPage) IsEmpty() (bool, error) {
	is, err := ExtractLoadBalancers(r)
	return len(is) == 0, err
}

// ExtractLoadBalancer accepts a Page struct, specifically a LoadBalancerPage struct,
// and extracts the elements into a slice of LoadBalancer structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractLoadBalancers(r pagination.Page) ([]LoadBalancer, error) {
	var s []LoadBalancer
	err := ExtractLoadBalancersInto(r, &s)
	return s, err
}

func ExtractLoadBalancersInto(r pagination.Page, v interface{}) error {
	return r.(LoadBalancerPage).Result.ExtractIntoSlicePtr(v, "results")
}

type LoadBalancerTaskResult struct {
	LoadBalancers []string `json:"loadbalancers"`
}

func ExtractLoadBalancerIDFromTask(task *tasks.Task) (string, error) {
	var result LoadBalancerTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode loadbalancer information in task structure: %w", err)
	}
	if len(result.LoadBalancers) == 0 {
		return "", fmt.Errorf("cannot decode loadbalancer information in task structure: %w", err)
	}
	return result.LoadBalancers[0], nil
}
