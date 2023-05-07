package guance

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/GuanceCloud/terraform-provider-guance/internal/sdk"
	ccv1 "github.com/GuanceCloud/terraform-provider-guance/internal/sdk/api/cloudcontrol/v1"
)

//go:embed README.md
var doc string

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
	Region      types.String `tfsdk:"region"`
	AccessToken types.String `tfsdk:"access_token"`
}

// Metadata returns the provider type name.
func (p *guanceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "guance"
}

// Schema defines the provider-level schema for configuration data.
func (p *guanceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Interact with Guance Cloud.",
		MarkdownDescription: doc,
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Description:         "Region for Guance Cloud API. May also be provided via GUANCE_REGION environment variable. See https://github.com/GuanceCloud/terraform-provider-guance for a list of available regions.",
				MarkdownDescription: "Region for Guance Cloud API. May also be provided via GUANCE_REGION environment variable. See [GitHub](https://github.com/GuanceCloud/terraform-provider-guance) for a list of available regions.",
				Optional:            true,
			},
			"access_token": schema.StringAttribute{
				Description:         "Access token for Guance Cloud API. May also be provided via GUANCE_ACCESS_TOKEN environment variable. Get an Key ID from https://console.guance.com/workspace/apiManage as access token.",
				MarkdownDescription: "Access token for Guance Cloud API. May also be provided via GUANCE_ACCESS_TOKEN environment variable. Get an Key ID from [Guance Cloud](https://console.guance.com/workspace/apiManage) as access token.",
				Optional:            true,
				Sensitive:           true,
			},
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

	region := getConfigField("region", config.Region, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	accessToken := getConfigField("access_token", config.AccessToken, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "guance_region", region)
	ctx = tflog.SetField(ctx, "guance_access_token", accessToken)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "guance_access_token")

	tflog.Debug(ctx, "Creating Guance Cloud client")

	// Create a new Guance Cloud client using the configuration values
	client, err := ccv1.NewClient(
		ccv1.WithRegion(region),
		ccv1.WithWait(true),
		ccv1.WithAccessToken(accessToken),
	)
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

func getConfigField(name string, value types.String, resp *provider.ConfigureResponse) string {
	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	envName := fmt.Sprintf("GUANCE_%s", strings.ToUpper(name))

	if value.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root(name),
			fmt.Sprintf("Unknown Guance Cloud API %s", strings.ToTitle(name)),
			"The provider cannot create the Guance Cloud API client as there is an unknown configuration value for the Guance Cloud API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the "+envName+" environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return ""
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	valueString := os.Getenv(envName)
	if !value.IsNull() {
		valueString = value.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if valueString == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root(name),
			fmt.Sprintf("Missing Guance Cloud API %s", strings.ToTitle(name)),
			"The provider cannot create the Guance Cloud API client as there is a missing or empty value for the Guance Cloud API "+name+". "+
				"Set the host value in the configuration or use the "+envName+" environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	return valueString
}
