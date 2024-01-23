package sendgrid

import (
	"context"
	"fmt"
	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &teammateDataSource{}
	_ datasource.DataSourceWithConfigure = &teammateDataSource{}
)

func NewTeammateDataSource() datasource.DataSource {

	return &teammateDataSource{}

}

type teammateDataSource struct {
	client *sendgrid.Client
}

type teammateModel struct {
	Email    types.String `tfsdk:"email"`
	Username types.String `tfsdk:"username"`
}

func (d *teammateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teammate"
}

func (d *teammateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Allows to retrive the details for provided teammate",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Description: "Email of the teammate",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username of the teammate",
				Optional:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *teammateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state teammateModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	itemResponse, err := d.client.ReadUser(ctx, state.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting teammate",
			fmt.Sprintf("Error getting teammate: %s", err.Error()),
		)

		return
	}

	// if err := json.Unmarshal([]byte(itemResponse), &state); err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error unmarshalling teammate",
	// 		fmt.Sprintf("Error unmarshalling teammate: %s", err.Error()),
	// 	)

	// 	return
	// }

	// // Map response body to model
	state = teammateModel{
		Email:    types.StringValue(itemResponse.Email),
		Username: types.StringValue(itemResponse.Username),
	}

	// //set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, "Finished reading item data source", map[string]any{"success": true})

}

// Configure adds the provider configured client to the data source.
func (d *teammateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
