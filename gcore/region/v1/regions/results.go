package regions

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/region/v1/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a region resource.
func (r commonResult) Extract() (*Region, error) {
	var s Region
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Region.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Region.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Region.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	commonResult
}

// Keystone represents a keystone structure.
type Keystone struct {
	ID                        int                       `json:"id"`
	URL                       gcorecloud.URL            `json:"url"`
	State                     types.KeystoneState       `json:"state"`
	KeystoneFederatedDomainID string                    `json:"keystone_federated_domain_id"`
	CreatedOn                 gcorecloud.JSONRFC3339NoZ `json:"created_on"`
	AdminPassword             string                    `json:"admin_password"`
}

// Region represents a region structure.
type Region struct {
	ID                int                       `json:"id"`
	DisplayName       string                    `json:"display_name"`
	KeystoneName      string                    `json:"keystone_name"`
	State             types.RegionState         `json:"state"`
	CreatorTaskID     *string                   `json:"creator_task_id"`
	EndpointType      types.EndpointType        `json:"endpoint_type"`
	ExternalNetworkID string                    `json:"external_network_id"`
	SpiceProxyURL     gcorecloud.URL            `json:"spice_proxy_url"`
	CreatedOn         gcorecloud.JSONRFC3339NoZ `json:"created_on"`
	KeystoneID        int                       `json:"keystone_id"`
	Keystone          Keystone                  `json:"keystone"`
}

// RegionPage is the page returned by a pager when traversing over a
// collection of regions.
type RegionPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of regions has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r RegionPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a RegionPage struct is empty.
func (r RegionPage) IsEmpty() (bool, error) {
	is, err := ExtractRegions(r)
	return len(is) == 0, err
}

// ExtractRegion accepts a Page struct, specifically a RegionPage struct,
// and extracts the elements into a slice of Region structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractRegions(r pagination.Page) ([]Region, error) {
	var s []Region
	err := ExtractRegionsInto(r, &s)
	return s, err
}

func ExtractRegionsInto(r pagination.Page, v interface{}) error {
	return r.(RegionPage).Result.ExtractIntoSlicePtr(v, "results")
}