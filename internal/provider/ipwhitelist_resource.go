package sendgrid

import (
	"context"
	"fmt"

	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &ipwhitelistResource{}
	_ resource.ResourceWithConfigure   = &ipwhitelistResource{}
	_ resource.ResourceWithImportState = &ipwhitelistResource{}
)

func NewIpWhitelistResource() resource.Resource {
	return &ipwhitelistResource{}
}

type ipwhitelistResource struct {
	client *sendgrid.Client
}

type IpwhitelistModel struct {
	Ip types.String `tfsdk:"ip"`
	Id types.Int64  `tfsdk:"id"`
}

func (r *ipwhitelistResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipwhitelist"
}

func (r *ipwhitelistResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Helps to add IP address to IP Access Management",
		Attributes: map[string]schema.Attribute{
			"ip": schema.StringAttribute{
				Description: "IP address",
				Required:    true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.RequiresReplace(),
				// },
			},
			"id": schema.Int64Attribute{
				Description: "ID of the IP address",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ipwhitelistResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var state IpwhitelistModel
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ip := state.Ip.ValueString()

	item := sendgrid.Ipmgmt{
		IP: ip,
	}

	ipwlressponse, err := r.client.CreateIPMgmt(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ip whitelist",
			fmt.Sprintf("Error creating ip whitelist: %s", err.Error()),
		)

		return
	}

	state = IpwhitelistModel{
		Ip: types.StringValue(ipwlressponse.IP),
		Id: types.Int64Value(ipwlressponse.ID),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *ipwhitelistResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var refstate IpwhitelistModel
	var inputItem string
	diags := req.State.Get(ctx, &refstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//tflog.Debug(ctx, "Item to be read", map[string]any{"ID: %+v": refstate.Ip.ValueString()})
	if refstate.Ip.ValueString() != "" && refstate.Id.ValueInt64() == 0 {
		inputItem = refstate.Ip.ValueString()
	}

	if refstate.Id.ValueInt64() != 0 && refstate.Ip.ValueString() == "" {
		inputItem = fmt.Sprintf("%d", refstate.Id.ValueInt64())
	}

	if refstate.Id.ValueInt64() != 0 && refstate.Ip.ValueString() != "" {
		inputItem = fmt.Sprintf("%d", refstate.Id.ValueInt64())
	}

	tflog.Debug(ctx, "Callme to read", map[string]any{"ID: %s": inputItem})

	ipwlgetlist, err := r.client.GetIPMgmt(ctx, inputItem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ip whitelist",
			fmt.Sprintf("Error reading ip whitelist: %s", err.Error()),
		)

		return
	}

	refstate = IpwhitelistModel{
		Ip: types.StringValue(ipwlgetlist.Result.IP),
		Id: types.Int64Value(ipwlgetlist.Result.ID),
	}

	diags = resp.State.Set(ctx, refstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ipwhitelistResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ipwhitelistResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	tflog.Debug(ctx, "Preparing to delete item resource")
	var delema IpwhitelistModel
	diags := req.State.Get(ctx, &delema)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Item to be deleted", map[string]any{"ID: %+v": delema})

	_, err := r.client.DeleteIPMgmt(ctx, fmt.Sprintf("%d", delema.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ip whitelist",
			fmt.Sprintf("Error deleting ip whitelist: %s", err.Error()),
		)

		return
	}

	tflog.Debug(ctx, "Deleted item resource", map[string]any{"success": true})
}

// Configure adds the provider configured client to the resource.
func (r *ipwhitelistResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sendgrid.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sendgrid.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ipwhitelistResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// id, err := strconv.ParseInt(req.ID, 10, 64)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error importing ip whitelist",
	// 		fmt.Sprintf("Error importing ip whitelist: %s", err.Error()),
	// 	)

	// 	return
	// }
	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)

	resource.ImportStatePassthroughID(ctx, path.Root("ip"), req, resp)

}
