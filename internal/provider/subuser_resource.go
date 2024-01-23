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
	_ resource.Resource                = &subuserResource{}
	_ resource.ResourceWithConfigure   = &subuserResource{}
	_ resource.ResourceWithImportState = &subuserResource{}
)

func NewSubuserResource() resource.Resource {
	return &subuserResource{}
}

type subuserResource struct {
	client *sendgrid.Client
}

type SubuserModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Email    types.String `tfsdk:"email"`
	Ips      []string     `tfsdk:"ips"`
	Disabled types.Bool   `tfsdk:"disabled"`
	ID       types.Int64  `tfsdk:"id"`
}

func (r *subuserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *subuserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subuser"
}

func (r *subuserResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Allows to create and manage subuser for your SendGrid account",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Description: "Email address of the subuser",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username of the subuser",
				Computed:    true,
				Optional:    true,
			},
			"id": schema.Int64Attribute{
				Description: "ID of the subuser",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Is read only of the subuser",
				Required:    true,
				//	Computed:    true,
				//	Sensitive:   true,
			},
			"ips": schema.ListAttribute{
				Description: "Scopes of the subuser",
				Required:    true,
				ElementType: types.StringType,
			},
			"disabled": schema.BoolAttribute{
				Description: "Token of the Pending subuser",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func (r *subuserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var newstate SubuserModel
	//var customiplist []string
	diags := req.Plan.Get(ctx, &newstate)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	if newstate.Disabled.ValueBool() {
		resp.Diagnostics.AddError(
			"Incorrect value for disabled",
			fmt.Sprintf("Unable to create subuser: %s", "Disabled cannot be true while creating subuser"),
		)
		return
	}

	item := sendgrid.Subuser{
		Username: newstate.Username.ValueString(),
		Password: newstate.Password.ValueString(),
		Email:    newstate.Email.ValueString(),
		Ips:      newstate.Ips,
		Disabled: newstate.Disabled.ValueBool(),
	}

	subuserrespBody, err := r.client.CreateSubuser(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create subuser",
			fmt.Sprintf("Unable to create subuser: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "ReadingResource:", map[string]any{"item": subuserrespBody})

	newstate = SubuserModel{
		Username: types.StringValue(subuserrespBody.Username),
		Email:    types.StringValue(subuserrespBody.Email),
		Disabled: types.BoolValue(subuserrespBody.Disabled),
		Password: newstate.Password,
		Ips:      newstate.Ips,
		ID:       types.Int64Value(subuserrespBody.ID),
	}

	diags = resp.State.Set(ctx, newstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *subuserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var readstate SubuserModel
	diags := req.State.Get(ctx, &readstate)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	getitem := sendgrid.Subuser{
		Username: readstate.Username.ValueString(),
		Email:    readstate.Email.ValueString(),
		Disabled: readstate.Disabled.ValueBool(),
		Password: readstate.Password.ValueString(),
		Ips:      readstate.Ips,
	}

	subuserrespBody, err := r.client.GetSubuser(ctx, getitem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read subuser",
			fmt.Sprintf("Unable to read subuser: %s", err),
		)
		return
	}

	readstate = SubuserModel{
		Username: types.StringValue(subuserrespBody.Username),
		Email:    types.StringValue(subuserrespBody.Email),
		Disabled: types.BoolValue(subuserrespBody.Disabled),
		Password: readstate.Password,
		Ips:      readstate.Ips,
		ID:       types.Int64Value(subuserrespBody.ID),
	}

	diags = resp.State.Set(ctx, readstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *subuserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var updatestate SubuserModel
	diags := req.Plan.Get(ctx, &updatestate)
	if diags.HasError() {
		return
	}

	item := sendgrid.Subuser{
		Username: updatestate.Username.ValueString(),
		Disabled: updatestate.Disabled.ValueBool(),
		Password: updatestate.Password.ValueString(),
		Ips:      updatestate.Ips,
		Email:    updatestate.Email.ValueString(),
	}

	subuserrespBody, err := r.client.UpdateSubuser(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update subuser",
			fmt.Sprintf("Unable to update subuser: %s", err),
		)
		return
	}

	updatestate = SubuserModel{
		Username: types.StringValue(subuserrespBody.Username),
		Email:    types.StringValue(subuserrespBody.Email),
		Disabled: types.BoolValue(subuserrespBody.Disabled),
		Password: updatestate.Password,
		Ips:      updatestate.Ips,
		ID:       types.Int64Value(subuserrespBody.ID),
	}

	diags = resp.State.Set(ctx, updatestate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *subuserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Preparing to delete item resource")
	var deletestate SubuserModel
	diags := req.State.Get(ctx, &deletestate)
	if diags.HasError() {
		return
	}

	_, err := r.client.DeleteSubuser(ctx, deletestate.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete subuser",
			fmt.Sprintf("Unable to delete subuser: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "DeletedResource:", map[string]any{"success": true})

}

func (r *subuserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("username"), req, resp)
}
