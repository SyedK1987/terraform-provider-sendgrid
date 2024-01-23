package sendgrid

import (
	"context"
	"log"
	"os"

	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Provider is the main entrypoint for the provider.

var (
	_ provider.Provider = &sendgridProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &sendgridProvider{
			version: version,
		}
	}
}

type sendgridProvider struct {
	version string
}

func (p *sendgridProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sendgrid"
	resp.Version = p.version
}

func (p *sendgridProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Terraform provider to interact with SendGrid API...",
		Attributes: map[string]schema.Attribute{
			"apikey": schema.StringAttribute{
				Description: "Api key to authenticate SendGrid",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

type sendgridProviderModel struct {
	ApiKey types.String `tfsdk:"apikey"`
}

func (p *sendgridProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	log.Printf("[INFO] Configuring sendgrid provider: %s", req.Config)
	tflog.Info(ctx, "Configuring SendGrid client")

	var config sendgridProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown SendGrid API key",
			"The provider cannot authentication SendGrid API Client as there is an unknown configuration value for the SendGrid apikey. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SENDGRID_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	api_key := os.Getenv("SENDGRID_API_KEY")

	if !config.ApiKey.IsNull() {
		api_key = config.ApiKey.ValueString()
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing SendGrid API Key",
			"The provider cannot create the SendGrid API client as there is a missing or empty value for the Sendgrid API Key. "+
				"Set the apikey value in the configuration or use the SENDGRID_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "apikey")

	tflog.Debug(ctx, "Creating SendGrid API Client")

	client, err := sendgrid.NewClient(api_key)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create SendGrid API Client",
			"An unexpected error occurred when creating the SendGrid API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"SendGrid Client Error: "+err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Debug(ctx, "Finished configuring SendGrid provider", map[string]interface{}{"success": true})
}

func (p *sendgridProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewIpWhitelistResource,
		NewSingleSenderResource,
		NewApiKeyResource,
		NewTeammateResource,
		NewSubuserResource,
		NewDomainAuthResource,
		NewLinkbrandResource,
		NewLinkbrandValidateResource,
		NewDomainValidateResource,
		NewDomainSubuserResource,
		//NewValidateDomainResource,
		//	NewResendTmateResource,
	}
}

func (p *sendgridProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTeammateDataSource,
		NewipwhitelistDataSource,
		NewSubuserDataSource,
		NewdomainauthDataSource,
	}
}
