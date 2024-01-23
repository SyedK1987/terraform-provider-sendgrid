package sendgrid

import (
	"context"
	"fmt"
	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &domainsubuserResource{}
	_ resource.ResourceWithConfigure = &domainsubuserResource{}
)

func NewDomainSubuserResource() resource.Resource {
	return &domainsubuserResource{}
}

type domainsubuserResource struct {
	client *sendgrid.Client
}

type DomainsubuserResourceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	UserID   types.Int64  `tfsdk:"user_id"`
	//	Subusers types.List   `tfsdk:"subusers"`
}

func (r *domainsubuserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domainauth_add_subuser"
}

func (r *domainsubuserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage link branding",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The id of the domain to be associated with the subuser",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "The subuser username to be associate with the domain",
				Required:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "The subuser id to be associate with the domain",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *domainsubuserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var domainsubnewstate DomainsubuserResourceModel
	diags := req.Plan.Get(ctx, &domainsubnewstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Associating Subuser to Domain",
			fmt.Sprintf("Error associating subuser to domain: %s", diags.Errors()),
		)
		return
	}

	domainsubnewitem := sendgrid.DomainAuth{
		ID:       domainsubnewstate.ID.ValueInt64(),
		Username: domainsubnewstate.Username.ValueString(),
	}

	// Create the resource.
	domainsubuser, err := r.client.CreateDomainAuthSubuser(ctx, domainsubnewitem)
	if err != nil {
		resp.Diagnostics.AddError(
			"SK - Error Associating Subuser to Domain",
			fmt.Sprintf("Error associating subuser to domain: %s", err.Error()),
		)
		return
	}

	domainsubnewstate.ID = types.Int64Value(domainsubuser.ID)
	for i := 0; i < len(domainsubuser.Subusers); i++ {
		domainsubnewstate.UserID = types.Int64Value(domainsubuser.Subusers[i].UserID)
		domainsubnewstate.Username = types.StringValue(domainsubuser.Subusers[i].Username)
	}
	resp.State.Set(ctx, domainsubnewstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Error Associating Subuser to Domain",
			fmt.Sprintf("Error updating state for associating subuser to domain: %s", resp.Diagnostics.Errors()),
		)
		return
	}
}

// Read reads the Terraform state and returns the up-to-date configuration.
func (r *domainsubuserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var domainsubuserstate DomainsubuserResourceModel
	diags := req.State.Get(ctx, &domainsubuserstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Reading Domain Subuser",
			fmt.Sprintf("Error reading domain subuser: %s", diags.Errors()),
		)
		return
	}

	domainreadsubitem := sendgrid.DomainAuth{
		ID:       domainsubuserstate.ID.ValueInt64(),
		Username: domainsubuserstate.Username.ValueString(),
		UserId:   domainsubuserstate.UserID.ValueInt64(),
	}

	domainsubuser, err := r.client.GetDomainSubuser(ctx, domainreadsubitem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Domain Subuser",
			fmt.Sprintf("Error reading domain subuser: %s", err.Error()),
		)
		return
	}

	if domainsubuser.Subusers[0].Username != "" {
		domainsubuserstate = DomainsubuserResourceModel{
			ID:       types.Int64Value(domainsubuser.ID),
			Username: types.StringValue(domainsubuser.Subusers[0].Username),
			UserID:   types.Int64Value(domainsubuser.Subusers[0].UserID),
		}
	} else {
		domainsubuserstate = DomainsubuserResourceModel{
			ID:       domainsubuserstate.ID,
			Username: domainsubuserstate.Username,
			UserID:   domainsubuserstate.UserID,
		}
	}

	resp.State.Set(ctx, domainsubuserstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Error Reading Domain Subuser",
			fmt.Sprintf("Error updating state for reading domain subuser: %s", resp.Diagnostics.Errors()),
		)
		return
	}
}

// Update updates the resource and sets the Terraform state.
func (r *domainsubuserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource.
func (r *domainsubuserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var deldomainsubuserstate DomainsubuserResourceModel
	diags := req.State.Get(ctx, &deldomainsubuserstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error DisAssociating Domainauth Subuser",
			fmt.Sprintf("Error deleting domain subuser: %s", diags.Errors()),
		)
		return
	}

	deldomainsubuseritem := sendgrid.DomainAuth{
		ID:       deldomainsubuserstate.ID.ValueInt64(),
		Username: deldomainsubuserstate.Username.ValueString(),
	}

	_, err := r.client.DeleteDomainAuthSubuser(ctx, deldomainsubuseritem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error DisAssociating Domainauth Subuser",
			fmt.Sprintf("Error deleting domain subuser: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Deleted item resource", map[string]any{"success": true})
}

func (r *domainsubuserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
