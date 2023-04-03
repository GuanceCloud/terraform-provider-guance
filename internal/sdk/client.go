package sdk

import (
	"context"
	"encoding/json"

	ccv1 "github.com/GuanceCloud/terraform-provider-guance/internal/sdk/api/v1"
	"github.com/hashicorp/go-multierror"
)

// Resource is the interface that all Guance Cloud resources must implement.
type Resource interface {
	// GetResourceType returns the type of the resource.
	GetResourceType() string

	// GetId returns the ID of the resource.
	GetId() string

	// SetId sets the ID of the resource.
	SetId(s string)
}

// ListOptions are the options for listing resources.
type ListOptions struct {
	// Count of reserved records.
	MaxResults int64

	// The name of the resource type.
	TypeName string
}

// Client is a Client for the Guance Cloud API.
type Client[T Resource] struct {
	Client ccv1.Client
}

// Create creates a resource.
func (c *Client[T]) Create(ctx context.Context, plan T) error {
	desiredState, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	_, err = c.Client.CreateResource(ctx, &ccv1.CreateResourceRequest{
		TypeName:     plan.GetResourceType(),
		DesiredState: string(desiredState),
	})
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a resource.
func (c *Client[T]) Delete(ctx context.Context, plan T) error {
	_, err := c.Client.DeleteResource(ctx, &ccv1.DeleteResourceRequest{
		TypeName:   plan.GetResourceType(),
		Identifier: plan.GetId(),
	})
	if err != nil {
		return err
	}
	return nil
}

// Update updates a resource.
func (c *Client[T]) Update(ctx context.Context, plan T) error {
	return nil
}

// Read reads a resource.
func (c *Client[T]) Read(ctx context.Context, out T) error {
	rsResp, err := c.Client.GetResource(ctx, &ccv1.GetResourceRequest{
		TypeName:   out.GetResourceType(),
		Identifier: out.GetId(),
	})
	if err != nil {
		return err
	}

	// Map response body to schema and populate Computed attribute values
	if err := json.Unmarshal([]byte(rsResp.ResourceDescription.Properties), out); err != nil {
		return err
	}
	return nil
}

// List lists resources.
func (c *Client[T]) List(ctx context.Context, options *ListOptions) ([]T, error) {
	resp, err := c.Client.ListResources(ctx, &ccv1.ListResourcesRequest{
		TypeName:   options.TypeName,
		MaxResults: options.MaxResults,
	})
	if err != nil {
		return nil, err
	}

	var mErr error
	var results []T
	for _, descriptor := range resp.ResourceDescriptions {
		var item T
		if err := json.Unmarshal([]byte(descriptor.Properties), &item); err != nil {
			mErr = multierror.Append(mErr, err)
			continue
		}
		results = append(results, item)
	}
	return results, mErr
}
