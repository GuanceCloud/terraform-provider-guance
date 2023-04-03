package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/avast/retry-go"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-retryablehttp"

	"github.com/GuanceCloud/terraform-provider-guance/internal/sdk/types"
)

type Client interface {
	CreateResource(ctx context.Context, req *CreateResourceRequest, opts ...Option) (rsp *CreateResourceResponse, err error)
	DeleteResource(ctx context.Context, req *DeleteResourceRequest, opts ...Option) (rsp *DeleteResourceResponse, err error)
	GetResource(ctx context.Context, req *GetResourceRequest, opts ...Option) (rsp *GetResourceResponse, err error)
	ListResources(ctx context.Context, req *ListResourcesRequest, opts ...Option) (rsp *ListResourcesResponse, err error)
	UpdateResource(ctx context.Context, req *UpdateResourceRequest, opts ...Option) (rsp *UpdateResourceResponse, err error)
}

type Option func(Client) error

func NewClient(options ...Option) (Client, error) {
	c := &client{}
	var mErr error
	for _, option := range options {
		err := option(c)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}
	return c, mErr
}

type client struct {
	endpoint             string
	wait                 bool
	maxRetries           int
	nonIdempotentRetried bool
}

// WithEndpoint is an option to set the endpoint of the service.
func WithEndpoint(endpoint string) Option {
	return func(c Client) error {
		c.(*client).endpoint = endpoint
		return nil
	}
}

// WithWait is an option to set if need to wait the request is completed.
func WithWait(wait bool) Option {
	return func(c Client) error {
		c.(*client).wait = wait
		return nil
	}
}

// WithMaxRetries is an option to set the max retries of the request.
func WithMaxRetries(maxRetries int) Option {
	return func(c Client) error {
		c.(*client).maxRetries = maxRetries
		return nil
	}
}

// WithNonIdempotentRetried is an option to set if need to retry the non idempotent request.
func WithNonIdempotentRetried(nonIdempotentRetried bool) Option {
	return func(c Client) error {
		c.(*client).nonIdempotentRetried = nonIdempotentRetried
		return nil
	}
}

func (c client) clone() client {
	return client{
		endpoint:   c.endpoint,
		wait:       c.wait,
		maxRetries: c.maxRetries,
	}
}

func (c client) withOptions(opts ...Option) (*client, error) {
	var mErr error
	cc := c.clone()
	for _, opt := range opts {
		if err := opt(&cc); err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}
	return &cc, mErr
}

func (c client) CreateResource(ctx context.Context, req *CreateResourceRequest, opts ...Option) (rsp *CreateResourceResponse, err error) {
	if !c.nonIdempotentRetried {
		opts = append(opts, WithMaxRetries(0))
	}

	cc, err := c.withOptions(opts...)
	if err != nil {
		return nil, err
	}

	// request to create resource
	rsp = &CreateResourceResponse{}
	err = cc.invoke(ctx, "/cloud-control/create-resource", req, rsp)
	if err != nil {
		return nil, err
	}

	if !cc.wait {
		return rsp, nil
	}

	waiter := RequestWaiter{RequestId: rsp.ProgressEvent.RequestToken}
	if err := waiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("failed when retring to create resource: %w", err)
	}
	return rsp, nil
}

func (c client) DeleteResource(ctx context.Context, req *DeleteResourceRequest, opts ...Option) (rsp *DeleteResourceResponse, err error) {
	if !c.nonIdempotentRetried {
		opts = append(opts, WithMaxRetries(0))
	}

	cc, err := c.withOptions(opts...)
	if err != nil {
		return nil, err
	}

	// request to create resource
	rsp = &DeleteResourceResponse{}
	err = cc.invoke(ctx, "/cloud-control/delete-resource", req, rsp)
	if err != nil {
		return nil, err
	}

	if !cc.wait {
		return rsp, nil
	}

	waiter := RequestWaiter{RequestId: rsp.ProgressEvent.RequestToken}
	if err := waiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("failed when retring to create resource: %w", err)
	}
	return rsp, nil
}

func (c client) GetResource(ctx context.Context, req *GetResourceRequest, opts ...Option) (rsp *GetResourceResponse, err error) {
	cc, err := c.withOptions(opts...)
	if err != nil {
		return nil, err
	}

	// request to create resource
	rsp = &GetResourceResponse{}
	err = cc.invoke(ctx, "/cloud-control/get-resource", req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c client) ListResources(ctx context.Context, req *ListResourcesRequest, opts ...Option) (rsp *ListResourcesResponse, err error) {
	cc, err := c.withOptions(opts...)
	if err != nil {
		return nil, err
	}

	// request to create resource
	rsp = &ListResourcesResponse{}
	err = cc.invoke(ctx, "/cloud-control/list-resource", req, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (c client) UpdateResource(ctx context.Context, req *UpdateResourceRequest, opts ...Option) (rsp *UpdateResourceResponse, err error) {
	if !c.nonIdempotentRetried {
		opts = append(opts, WithMaxRetries(0))
	}

	cc, err := c.withOptions(opts...)
	if err != nil {
		return nil, err
	}

	// request to create resource
	rsp = &UpdateResourceResponse{}
	err = cc.invoke(ctx, "/cloud-control/update-resource", req, rsp)
	if err != nil {
		return nil, err
	}

	if !cc.wait {
		return rsp, nil
	}

	waiter := RequestWaiter{RequestId: rsp.ProgressEvent.RequestToken}
	if err := waiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("failed when retring to create resource: %w", err)
	}
	return rsp, nil
}

func (c client) invoke(ctx context.Context, path string, in interface{}, out interface{}) error {
	reqBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}

	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryMax = c.maxRetries

	// Send request
	resp, err := retryableClient.Post(
		fmt.Sprintf("http://127.0.0.1:8000%s", path),
		"application/json",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return err
	}

	// Got errors
	if resp.StatusCode >= 400 {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body failed with status code %d", resp.StatusCode)
		}
		outErr := &errors.Error{}
		if err := json.Unmarshal(respBytes, outErr); err != nil {
			return fmt.Errorf("unmarshal failed with status code %d: %s", resp.StatusCode, string(respBytes))
		}
		return outErr
	}

	// Got response
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(respBytes, out); err != nil {
		return err
	}
	return nil
}

func (c client) getResourceRequestStatus(ctx context.Context, g *GetResourceRequestStatusRequest) (*GetResourceRequestStatusResponse, error) {
	rsp := &GetResourceRequestStatusResponse{}
	err := c.invoke(ctx, "/cloud-control/get-resource-request-status", g, rsp)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

type RequestWaiter struct {
	RequestId string
	client    *client
}

func (w RequestWaiter) Wait(ctx context.Context) error {
	return retry.Do(
		func() error {
			// wait for the request is completed
			r, err := w.client.getResourceRequestStatus(ctx, &GetResourceRequestStatusRequest{
				RequestToken: w.RequestId,
			})
			if err != nil {
				return err
			}

			switch r.ProgressEvent.OperationStatus {
			case types.RequestStatusPending,
				types.RequestStatusInProgress,
				types.RequestStatusCancelInProgress:
				return fmt.Errorf("request is not completed")
			case types.RequestStatusSuccess,
				types.RequestStatusCancelComplete:
				return nil
			case types.RequestStatusFailed:
				return fmt.Errorf("failed to create resource: %s", r.ProgressEvent.StatusMessage)
			default:
				return fmt.Errorf("unknown status %s", r.ProgressEvent.OperationStatus)
			}
		},
	)
}
