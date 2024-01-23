package sendgrid

import (
	"context"
	"fmt"
	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ipwhitelistDataSource{}
	_ datasource.DataSourceWithConfigure = &ipwhitelistDataSource{}
)

func NewipwhitelistDataSource() datasource.DataSource {

	return &ipwhitelistDataSource{}

}

type ipwhitelistDataSource struct {
	client *sendgrid.Client
}

type IpwhitelistDataModel struct {
	Ip types.String `tfsdk:"ip"`
	Id types.Int64  `tfsdk:"id"`
}

func (d *ipwhitelistDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipwhitelist"
}

func (d *ipwhitelistDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Allows to retrive IP Access management details",
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description: "IP address",
				Computed:    true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.RequiresReplace(),
				// },
			},
			"id": schema.Int64Attribute{
				Description: "ID of the IP address",
				Required:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ipwhitelistDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var refstate IpwhitelistModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &refstate)...)

	itemResponse, err := d.client.GetIPMgmt(ctx, fmt.Sprintf("%d", refstate.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get IP Management",
			fmt.Sprintf("Failed to get IP Management: %s", err.Error()),
		)

		return
	}

	refstate = IpwhitelistModel{
		Ip: types.StringValue(itemResponse.Result.IP),
		Id: types.Int64Value(itemResponse.Result.ID),
	}

	diags := resp.State.Set(ctx, refstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *ipwhitelistDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sendgrid.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Sendgrid.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
