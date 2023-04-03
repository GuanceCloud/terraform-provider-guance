package guance

import (
	"context"
	"os"

	"github.com/GuanceCloud/terraform-provider-guance/internal/sdk"

	ccv1 "github.com/GuanceCloud/openapi/api/cloudcontrol/v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &guanceProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &guanceProvider{}
}

// guanceProvider is the provider implementation.
type guanceProvider struct{}

// guanceProviderModel maps provider schema data to a Go type.
type guanceProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

// Metadata returns the provider type name.
func (p *guanceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "guance"
}

// Schema defines the provider-level schema for configuration data.
func (p *guanceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Guance Cloud.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "URI for Guance Cloud API. May also be provided via GUANCE_ENDPOINT environment variable.",
				Optional:    true,
			},
			// TODO: add credential provider
		},
	}
}

// Configure prepares a Guance Cloud API client for data sources and resources.
func (p *guanceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Guance Cloud client")

	// Retrieve provider data from configuration
	var config guanceProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Guance Cloud API endpoint",
			"The provider cannot create the Guance Cloud API client as there is an unknown configuration value for the Guance Cloud API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the GUANCE_ENDPOINT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	endpoint := os.Getenv("GUANCE_ENDPOINT")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Guance Cloud API Host",
			"The provider cannot create the Guance Cloud API client as there is a missing or empty value for the Guance Cloud API host. "+
				"Set the host value in the configuration or use the GUANCE_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "guance_endpoint", endpoint)

	tflog.Debug(ctx, "Creating Guance Cloud client")

	// Create a new Guance Cloud client using the configuration values
	client, err := ccv1.NewClient(ccv1.WithEndpoint(endpoint))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Guance Cloud API Client",
			"An unexpected error occurred when creating the Guance Cloud API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Guance Cloud Client Error: "+err.Error(),
		)
		return
	}

	// Make the Guance Cloud client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = &sdk.Client[sdk.Resource]{Client: client}
	resp.ResourceData = &sdk.Client[sdk.Resource]{Client: client}

	tflog.Info(ctx, "Configured Guance Cloud client", map[string]any{"success": true})
}
