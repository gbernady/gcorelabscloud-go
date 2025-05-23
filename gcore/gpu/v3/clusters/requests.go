package clusters

import (
	"net/http"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/pagination"
)

// ClusterActionOptsBuilder allows extensions to add parameters to the action request.
type ClusterActionOptsBuilder interface {
	ToClusterActionMap() (map[string]interface{}, error)
}

// ClusterActionOpts represents options used to run an action on a cluster.
type ClusterActionOpts struct {
	Action       ClusterActionType `json:"action" required:"true" validate:"required,enum"`
	ServersCount *int              `json:"servers_count,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
}

// Validate checks if the ClusterActionOpts is valid.
func (opts ClusterActionOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToActionMap builds a request body from ClusterActionOpts.
func (opts ClusterActionOpts) ToClusterActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	mp, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// DeleteClusterOptsBuilder allows extensions to add parameters to delete cluster options.
type DeleteClusterOptsBuilder interface {
	ToClusterDeleteQuery() (string, error)
}

// DeleteClusterOpts specifies the parameters for the Delete method.
type DeleteClusterOpts struct {
	AllFloatingIPs      bool     `q:"all_floating_ips" validate:"omitempty,allowed_without=FloatingIPIDs"`
	AllReservedFixedIPs bool     `q:"all_reserved_fixed_ips" validate:"omitempty,allowed_without=ReservedFixedIPIDs"`
	AllVolumes          bool     `q:"all_volumes" validate:"omitempty,allowed_without=VolumeIDs"`
	FloatingIPIDs       []string `q:"floating_ip_ids" validate:"omitempty,allowed_without=AllFloatingIPs,dive,uuid4" delimiter:"comma"`
	ReservedFixedIPIDs  []string `q:"reserved_fixed_ip_ids" validate:"omitempty,allowed_without=AllReservedFixedIPs,dive,uuid4" delimiter:"comma"`
	VolumeIDs           []string `q:"volume_ids" validate:"omitempty,allowed_without=AllVolumes,dive,uuid4" delimiter:"comma"`
}

// Validate checks if the provided options are valid.
func (opts DeleteClusterOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToDeleteClusterActionMap builds a request body from DeleteInstanceOpts.
func (opts DeleteClusterOpts) ToClusterDeleteQuery() (string, error) {
	if err := opts.Validate(); err != nil {
		return "", err
	}
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// RenameClusterOptsBuilder allows extensions to add parameters to rename cluster options.
type RenameClusterOptsBuilder interface {
	ToRenameClusterActionMap() (map[string]interface{}, error)
}

// RenameClusterOpts specifies the parameters for the Rename method.
type RenameClusterOpts struct {
	Name string `json:"name" validate:"required"`
}

// Validate checks if the provided options are valid.
func (opts RenameClusterOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

// ToRenameClusterActionMap builds a request body from RenameInstanceOpts.
func (opts RenameClusterOpts) ToRenameClusterActionMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return gcorecloud.BuildRequestBody(opts, "")
}

type ServerCredentialsOpts struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	KeypairName string `json:"keypair_name,omitempty"`
}

type ServerSettingsOpts struct {
	Interfaces     []InterfaceOpts        `json:"interfaces"`
	SecurityGroups []string               `json:"security_groups,omitempty"`
	Volumes        []VolumeOpts           `json:"volumes"`
	UserData       *string                `json:"user_data,omitempty"`
	Credentials    *ServerCredentialsOpts `json:"credentials,omitempty"`
}

// VolumeOpts represents options used to create a volume.
type VolumeOpts struct {
	Source              VolumeSource      `json:"source" validate:"required,enum"`
	BootIndex           int               `json:"boot_index" validate:"required"`
	DeleteOnTermination bool              `json:"delete_on_termination,omitempty"`
	Name                string            `json:"name" validate:"required"`
	Size                int               `json:"size,omitempty" validate:"required"`
	ImageID             string            `json:"image_id,omitempty" validate:"rfe=Source:image,allowed_without=SnapshotID,omitempty,uuid4"`
	SnapshotID          string            `json:"snapshot_id,omitempty" validate:"rfe=Source:snapshot,allowed_without=ImageID,omitempty,uuid4"`
	Tags                map[string]string `json:"tags,omitempty"`
	Type                VolumeType        `json:"type,omitempty" validate:"required,enum"`
}

type InterfaceOpts interface {
	implInterfaceOpts()
}

func (ExternalInterfaceOpts) implInterfaceOpts()  {}
func (SubnetInterfaceOpts) implInterfaceOpts()    {}
func (AnySubnetInterfaceOpts) implInterfaceOpts() {}

type ExternalInterfaceOpts struct {
	Name     *string      `json:"name,omitempty"`
	Type     string       `json:"type" validate:"required"`
	IPFamily IPFamilyType `json:"ip_family,omitempty"`
}

type FloatingIPOpts struct {
	Source string `json:"source" validate:"required,enum"`
}

type SubnetInterfaceOpts struct {
	NetworkID  string          `json:"network_id" validate:"required"`
	Name       *string         `json:"name,omitempty"`
	Type       string          `json:"type" validate:"required"`
	SubnetID   string          `json:"subnet_id" validate:"required"`
	FloatingIP *FloatingIPOpts `json:"floating_ip,omitempty"`
}

type AnySubnetInterfaceOpts struct {
	NetworkID  string          `json:"network_id" validate:"required"`
	Name       *string         `json:"name,omitempty"`
	Type       string          `json:"type" validate:"required"`
	IPFamily   IPFamilyType    `json:"ip_family,omitempty"`
	IPAddress  *string         `json:"ip_address,omitempty"`
	FloatingIP *FloatingIPOpts `json:"floating_ip,omitempty"`
}

// CreateClusterOpts allows extensions to add parameters to create cluster options.
type CreateClusterOpts struct {
	Name            string             `json:"name" validate:"required"`
	Flavor          string             `json:"flavor" validate:"required"`
	Tags            map[string]string  `json:"tags,omitempty"`
	ServersCount    int                `json:"servers_count,omitempty"`
	ServersSettings ServerSettingsOpts `json:"servers_settings,omitempty"`
}

// Validate checks if the provided options are valid.
func (opts CreateClusterOpts) Validate() error {
	return gcorecloud.ValidateStruct(opts)
}

func (opts CreateClusterOpts) ToCreateClusterMap() (map[string]interface{}, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	mp, err := gcorecloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// CreateClusterOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateClusterOptsBuilder interface {
	ToCreateClusterMap() (map[string]interface{}, error)
}

// List returns a pager for listing GPU clusters
func List(client *gcorecloud.ServiceClient) pagination.Pager {
	url := ClustersURL(client)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListAll is a convenience function that returns all GPU clusters
func ListAll(client *gcorecloud.ServiceClient) ([]Cluster, error) {
	pages, err := List(client).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractClusters(pages)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// Get retrieves a specific GPU cluster by its ID.
func Get(client *gcorecloud.ServiceClient, clusterID string) (r GetResult) {
	url := ClusterURL(client, clusterID)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

// Delete removes a specific GPU cluster by its ID.
func Delete(client *gcorecloud.ServiceClient, clusterID string, opts DeleteClusterOptsBuilder) (r tasks.Result) {
	url := ClusterURL(client, clusterID)
	if opts != nil {
		query, err := opts.ToClusterDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	_, r.Err = client.DeleteWithResponse(url, &r.Body, nil) // nolint
	return
}

// Rename changes the name of a GPU cluster.
func Rename(client *gcorecloud.ServiceClient, clusterID string, opts RenameClusterOptsBuilder) (r GetResult) {
	b, err := opts.ToRenameClusterActionMap()
	if err != nil {
		r.Err = err
		return
	}

	url := ClusterURL(client, clusterID)
	_, r.Err = client.Patch(url, b, &r.Body, &gcorecloud.RequestOpts{ // nolint
		OkCodes: []int{http.StatusOK},
	})
	return
}

// Create creates a new GPU cluster.
func Create(client *gcorecloud.ServiceClient, opts CreateClusterOptsBuilder) (r tasks.Result) {
	b, err := opts.ToCreateClusterMap()
	if err != nil {
		r.Err = err
		return
	}

	url := ClustersURL(client)
	_, r.Err = client.Post(url, b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{http.StatusOK, http.StatusCreated},
	})
	return
}

// Action run an action on the GPU cluster.
func Action(client *gcorecloud.ServiceClient, clusterID string, opts ClusterActionOptsBuilder) (r tasks.Result) {
	b, err := opts.ToClusterActionMap()
	if err != nil {
		r.Err = err
		return
	}

	url := ClusterActionURL(client, clusterID)
	_, r.Err = client.Post(url, b, &r.Body, nil) // nolint
	return
}
