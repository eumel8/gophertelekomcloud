package networks

import (
  "github.com/mitchellh/mapstructure"
  "github.com/rackspace/gophercloud"
  "github.com/rackspace/gophercloud/pagination"
)

type commonResult struct {
  gophercloud.Result
}

// Extract is a function that accepts a result and extracts a network resource.
func (r commonResult) Extract() (*Network, error) {
  if r.Err != nil {
    return nil, r.Err
  }

  var res struct {
    Network *Network `json:"network"`
  }

  err := mapstructure.Decode(r.Body, &res)

  return res.Network, err
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
  commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
  commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult commonResult

// Network represents, well, a network.
type Network struct {
  // UUID for the network
  ID string `mapstructure:"id" json:"id"`

  // Human-readable name for the network. Might not be unique.
  Label string `mapstructure:"label" json:"label"`

  // Classless Inter-Domain Routing
  CIDR string `mapstructure:"cidr" json:"cidr"`
}

// NetworkPage is the page returned by a pager when traversing over a
// collection of networks.
type NetworkPage struct {
  pagination.SinglePageBase
}

// ExtractNetworks accepts a Page struct, specifically a NetworkPage struct,
// and extracts the elements into a slice of Network structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractNetworks(page pagination.Page) ([]Network, error) {
  var resp struct {
    Networks []Network `mapstructure:"networks" json:"networks"`
  }

  err := mapstructure.Decode(page.(NetworkPage).Body, &resp)

  return resp.Networks, err
}
