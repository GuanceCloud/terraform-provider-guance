package sdk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tidwall/gjson"

	"github.com/GuanceCloud/terraform-provider-guance/internal/helpers/tfcodec"
	ccv1 "github.com/GuanceCloud/terraform-provider-guance/internal/sdk/api/cloudcontrol/v1"
)

// Resource is the interface that all Guance Cloud resources must implement.
type Resource interface {
	// GetResourceType returns the type of the resource.
	GetResourceType() string

	// GetId returns the ID of the resource.
	GetId() string

	// SetId sets the ID of the resource.
	SetId(s string)

	// SetCreatedAt sets the creation time of the resource.
	SetCreatedAt(time string)
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
	desiredState, err := json.Marshal(tfcodec.Encode(plan))
	if err != nil {
		return err
	}
	tflog.Info(ctx, fmt.Sprintf("[DESIRED STATE]: %+v", string(desiredState)))
	resp, err := c.Client.CreateResource(ctx, &ccv1.CreateResourceRequest{
		TypeName:     plan.GetResourceType(),
		DesiredState: string(desiredState),
	})
	tflog.Info(ctx, fmt.Sprintf("[RESP]: %+v; [ERROR] %+v;", resp, err))
	if err != nil {
		return err
	}
	plan.SetId(resp.ProgressEvent.Identifier)
	return nil
}

// Delete deletes a resource.
func (c *Client[T]) Delete(ctx context.Context, plan T) error {
	_, err := c.Client.DeleteResource(ctx, &ccv1.DeleteResourceRequest{
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
	id := out.GetId()
	rsResp, err := c.Client.GetResource(ctx, &ccv1.GetResourceRequest{
		Identifier: id,
	})
	if err != nil {
		return fmt.Errorf("failed to get resource: %w", err)
	}

	// Map response body to schema and populate Computed attribute values
	if err := tfcodec.DecodeJSON([]byte(rsResp.ResourceDescription.Properties), out); err != nil {
		return fmt.Errorf("failed to decode properties: %w", err)
	}
	out.SetId(rsResp.ResourceDescription.Identifier)
	out.SetCreatedAt(rsResp.ResourceDescription.CreatedAt)
	return nil
}

type Filter struct {
	Name   types.String   `tfsdk:"name"`
	Values []types.String `tfsdk:"values"`
}

func (f *Filter) IsOK(state string) bool {
	for _, value := range f.Values {
		if gjson.Get(state, f.Name.ValueString()).String() != value.ValueString() {
			return false
		}
	}
	return true
}

func FilterAllSuccess(state string, filters ...*Filter) bool {
	for _, filter := range filters {
		if !filter.IsOK(state) {
			return false
		}
	}
	return true
}

// List lists resources.
func (c *Client[T]) List(ctx context.Context, options *ListOptions) (*ccv1.ListResourcesResponse, error) {
	return c.Client.ListResources(ctx, &ccv1.ListResourcesRequest{
		TypeName:   options.TypeName,
		MaxResults: options.MaxResults,
	})
}
