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
	_ datasource.DataSource              = &subuserDataSource{}
	_ datasource.DataSourceWithConfigure = &subuserDataSource{}
)

func NewSubuserDataSource() datasource.DataSource {
	return &subuserDataSource{}
}

type subuserDataSource struct {
	client *sendgrid.Client
}

type DataSubUserModel struct {
	Username types.String `tfsdk:"username"`
	Email    types.String `tfsdk:"email"`
	Disabled types.Bool   `tfsdk:"disabled"`
	ID       types.Int64  `tfsdk:"id"`
}

func (d *subuserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subuser"
}

func (d *subuserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *subuserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Allows to retrive the details for provided subuser",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Description: "Email address of the subuser",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username of the subuser",
				Required:    true,
			},
			"id": schema.Int64Attribute{
				Description: "ID of the subuser",
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: "Token of the Pending subuser",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *subuserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var datastate DataSubUserModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &datastate)...)

	itemResponse, err := d.client.ReadSubuser(ctx, datastate.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting subuser",
			"Error getting subuser: "+err.Error(),
		)

		return
	}

	datastate = DataSubUserModel{
		Username: types.StringValue(itemResponse.Username),
		Email:    types.StringValue(itemResponse.Email),
		Disabled: types.BoolValue(itemResponse.Disabled),
		ID:       types.Int64Value(itemResponse.ID),
	}

	diags := resp.State.Set(ctx, datastate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
