package sendgrid

import (
	"context"
	"fmt"
	"strconv"

	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Provider is the main entrypoint for the provider.
var (
	_ resource.Resource                = &singlesenderResource{}
	_ resource.ResourceWithConfigure   = &singlesenderResource{}
	_ resource.ResourceWithImportState = &singlesenderResource{}
)

func NewSingleSenderResource() resource.Resource {
	return &singlesenderResource{}
}

type singlesenderResource struct {
	client *sendgrid.Client
}

type SingleSenderModel struct {
	Nickname    types.String `tfsdk:"nickname"`
	FromEmail   types.String `tfsdk:"from_email"`
	FromName    types.String `tfsdk:"from_name"`
	ReplyTo     types.String `tfsdk:"reply_to"`
	ReplyToName types.String `tfsdk:"reply_to_name"`
	Address     types.String `tfsdk:"address"`
	Address2    types.String `tfsdk:"address2"`
	State       types.String `tfsdk:"state"`
	City        types.String `tfsdk:"city"`
	Country     types.String `tfsdk:"country"`
	Zip         types.String `tfsdk:"zip"`
	ID          types.Int64  `tfsdk:"id"`
	Verified    types.Bool   `tfsdk:"verified"`
	Locked      types.Bool   `tfsdk:"locked"`
}

func (r *singlesenderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_single_sender"
}

func (r *singlesenderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Helps to add single sender authentication to your sendgrid account",
		Attributes: map[string]schema.Attribute{
			"nickname": schema.StringAttribute{
				Description: "Nickname of the sender",
				Required:    true,
			},
			"from_email": schema.StringAttribute{
				Description: "Email address of the sender",
				Required:    true,
			},
			"from_name": schema.StringAttribute{
				Description: "Name of the sender",
				Required:    true,
			},
			"reply_to": schema.StringAttribute{
				Description: "Reply to email address of the sender",
				Required:    true,
			},
			"reply_to_name": schema.StringAttribute{
				Description: "Reply to name of the sender",
				Required:    true,
			},
			"address": schema.StringAttribute{
				Description: "Address of the sender",
				Required:    true,
			},
			"address2": schema.StringAttribute{
				Description: "Address2 of the sender",
				Required:    true,
			},
			"state": schema.StringAttribute{
				Description: "State of the sender",
				Required:    true,
			},
			"city": schema.StringAttribute{
				Description: "City of the sender",
				Required:    true,
			},
			"country": schema.StringAttribute{
				Description: "Country of the sender",
				Required:    true,
			},
			"zip": schema.StringAttribute{
				Description: "Zip of the sender",
				Required:    true,
			},
			"id": schema.Int64Attribute{
				Description: "ID of the sender",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"verified": schema.BoolAttribute{
				Description: "Verified of the sender",
				Computed:    true,
			},
			"locked": schema.BoolAttribute{
				Description: "Locked of the sender",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *singlesenderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var newstate SingleSenderModel
	diags := req.Plan.Get(ctx, &newstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := sendgrid.Singlesender{
		Nickname:    newstate.Nickname.ValueString(),
		FromEmail:   newstate.FromEmail.ValueString(),
		FromName:    newstate.FromName.ValueString(),
		ReplyTo:     newstate.ReplyTo.ValueString(),
		ReplyToName: newstate.ReplyToName.ValueString(),
		Address:     newstate.Address.ValueString(),
		Address2:    newstate.Address2.ValueString(),
		State:       newstate.State.ValueString(),
		City:        newstate.City.ValueString(),
		Country:     newstate.Country.ValueString(),
		Zip:         newstate.Zip.ValueString(),
	}

	singlesenderresponse, err := r.client.CreateSingleSender(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating single sender",
			fmt.Sprintf("Error creating single sender: %s", err.Error()),
		)

		return
	}

	newstate = SingleSenderModel{
		Nickname:    types.StringValue(singlesenderresponse.Nickname),
		FromEmail:   types.StringValue(singlesenderresponse.FromEmail),
		FromName:    types.StringValue(singlesenderresponse.FromName),
		ReplyTo:     types.StringValue(singlesenderresponse.ReplyTo),
		ReplyToName: types.StringValue(singlesenderresponse.ReplyToName),
		Address:     types.StringValue(singlesenderresponse.Address),
		Address2:    types.StringValue(singlesenderresponse.Address2),
		State:       types.StringValue(singlesenderresponse.State),
		City:        types.StringValue(singlesenderresponse.City),
		Country:     types.StringValue(singlesenderresponse.Country),
		Zip:         types.StringValue(singlesenderresponse.Zip),
		ID:          types.Int64Value(singlesenderresponse.ID),
		Verified:    types.BoolValue(singlesenderresponse.Verified),
		Locked:      types.BoolValue(singlesenderresponse.Locked),
	}

	diags = resp.State.Set(ctx, newstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *singlesenderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var readstate SingleSenderModel
	diags := req.State.Get(ctx, &readstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//tflog.Debug(ctx, "ReadingResource:", map[string]any{"id": readstate.ID.ValueInt64()})
	readsinglesenderresponse, err := r.client.ReadSingleSender(ctx, fmt.Sprintf("%d", readstate.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating single sender",
			fmt.Sprintf("Error creating single sender: %s", err.Error()),
		)

		return
	}

	readstate = SingleSenderModel{
		Nickname:    types.StringValue(readsinglesenderresponse.Nickname),
		FromEmail:   types.StringValue(readsinglesenderresponse.FromEmail),
		FromName:    types.StringValue(readsinglesenderresponse.FromName),
		ReplyTo:     types.StringValue(readsinglesenderresponse.ReplyTo),
		ReplyToName: types.StringValue(readsinglesenderresponse.ReplyToName),
		Address:     types.StringValue(readsinglesenderresponse.Address),
		Address2:    types.StringValue(readsinglesenderresponse.Address2),
		State:       types.StringValue(readsinglesenderresponse.State),
		City:        types.StringValue(readsinglesenderresponse.City),
		Country:     types.StringValue(readsinglesenderresponse.Country),
		Zip:         types.StringValue(readsinglesenderresponse.Zip),
		ID:          types.Int64Value(readsinglesenderresponse.ID),
		Verified:    types.BoolValue(readsinglesenderresponse.Verified),
		Locked:      types.BoolValue(readsinglesenderresponse.Locked),
	}

	diags = resp.State.Set(ctx, readstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *singlesenderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var updatestate SingleSenderModel
	diags := req.Plan.Get(ctx, &updatestate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Preparing to update item resource:", map[string]any{"id:": fmt.Sprintf("%d", updatestate.ID.ValueInt64())})

	updateitem := sendgrid.Singlesender{
		Nickname:    updatestate.Nickname.ValueString(),
		FromEmail:   updatestate.FromEmail.ValueString(),
		FromName:    updatestate.FromName.ValueString(),
		ReplyTo:     updatestate.ReplyTo.ValueString(),
		ReplyToName: updatestate.ReplyToName.ValueString(),
		Address:     updatestate.Address.ValueString(),
		Address2:    updatestate.Address2.ValueString(),
		State:       updatestate.State.ValueString(),
		City:        updatestate.City.ValueString(),
		Country:     updatestate.Country.ValueString(),
		Zip:         updatestate.Zip.ValueString(),
		ID:          updatestate.ID.ValueInt64(),
	}

	updatesinglesenderresponse, err := r.client.UpdateSingleSender(ctx, updateitem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating single sender",
			fmt.Sprintf("Error updating single sender: %s", err.Error()),
		)

		return
	}

	updatestate = SingleSenderModel{
		Nickname:    types.StringValue(updatesinglesenderresponse.Nickname),
		FromEmail:   types.StringValue(updatesinglesenderresponse.FromEmail),
		FromName:    types.StringValue(updatesinglesenderresponse.FromName),
		ReplyTo:     types.StringValue(updatesinglesenderresponse.ReplyTo),
		ReplyToName: types.StringValue(updatesinglesenderresponse.ReplyToName),
		Address:     types.StringValue(updatesinglesenderresponse.Address),
		Address2:    types.StringValue(updatesinglesenderresponse.Address2),
		State:       types.StringValue(updatesinglesenderresponse.State),
		City:        types.StringValue(updatesinglesenderresponse.City),
		Country:     types.StringValue(updatesinglesenderresponse.Country),
		Zip:         types.StringValue(updatesinglesenderresponse.Zip),
		ID:          types.Int64Value(updatesinglesenderresponse.ID),
		Verified:    types.BoolValue(updatesinglesenderresponse.Verified),
		Locked:      types.BoolValue(updatesinglesenderresponse.Locked),
	}

	diags = resp.State.Set(ctx, updatestate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *singlesenderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Preparing to delete item resource")
	var delsender SingleSenderModel
	diags := req.State.Get(ctx, &delsender)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteSingleSender(ctx, fmt.Sprintf("%d", delsender.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Single Sender",
			fmt.Sprintf("Error deleting Single Sender: %s", err.Error()),
		)

		return
	}

	tflog.Debug(ctx, "Deleted item resource", map[string]any{"success": true})
}

// Configure adds the provider configured client to the resource.
func (r *singlesenderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ImportState imports the resource state from the Terraform state.
func (r *singlesenderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing Single Sender Authentication",
			fmt.Sprintf("Error importing Single Sender Authentication: %s", err.Error()),
		)

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
