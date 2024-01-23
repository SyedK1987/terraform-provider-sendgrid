package sendgrid

import (
	"context"
	"fmt"
	"strconv"

	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &linkbrandvalidateResource{}
	_ resource.ResourceWithConfigure = &linkbrandvalidateResource{}
)

func NewLinkbrandValidateResource() resource.Resource {
	return &linkbrandvalidateResource{}
}

type linkbrandvalidateResource struct {
	client *sendgrid.Client
}

type LinkbrandvalidResourceModel struct {
	ID    types.Int64 `tfsdk:"id"`
	Valid types.Bool  `tfsdk:"valid"`
}

func (r *linkbrandvalidateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_linkbrand_validate"
}

func (r *linkbrandvalidateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage link branding",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of the link branding",
				Required:    true,
			},
			"valid": schema.BoolAttribute{
				Description: "The valid domain",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *linkbrandvalidateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var linkvalnewstate LinkbrandvalidResourceModel
	diags := req.Plan.Get(ctx, &linkvalnewstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Validating Link Brand",
			fmt.Sprintf("Error creating link branding: %s", diags.Errors()),
		)
		return
	}

	linkitemState := sendgrid.LinkAuth{
		ID: linkvalnewstate.ID.ValueInt64(),
	}

	newlinkItem, err := r.client.Validatelinkbrand(ctx, linkitemState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating link brand",
			fmt.Sprintf("Error validating link brand: %s", err),
		)
		return
	}

	if newlinkItem.Valid {
		//convertToMap := structToMap(newItem.DNSDetails)
		tflog.Debug(ctx, "RetriveData:", map[string]any{"item": newlinkItem})
		linkvalnewstate = LinkbrandvalidResourceModel{
			ID:    types.Int64Value(newlinkItem.ID),
			Valid: types.BoolValue(newlinkItem.Valid),
		}

		//getelemements := make(map[string]DomainAuthRecord)

		diags = resp.State.Set(ctx, linkvalnewstate)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Error validating link brand",
			fmt.Sprintf("Error validating link brand: %s", err),
		)
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *linkbrandvalidateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var linkreadstate LinkbrandvalidResourceModel
	diags := req.State.Get(ctx, &linkreadstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Error reading link branding: %s", diags.Errors()),
		)
		return
	}

	linkreaditem, err := r.client.Getlinkbrand(ctx, sendgrid.LinkAuth{ID: linkreadstate.ID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Error reading link branding: %s", err),
		)
		return
	}
	linkreadstate = LinkbrandvalidResourceModel{
		ID:    types.Int64Value(linkreaditem.ID),
		Valid: types.BoolValue(linkreaditem.Valid),
	}

	//getelemements := make(map[string]DomainAuthRecord)

	diags = resp.State.Set(ctx, linkreadstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *linkbrandvalidateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *linkbrandvalidateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// Configure adds the provider configured client to the resource.
func (r *linkbrandvalidateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *linkbrandvalidateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing Link Branding",
			fmt.Sprintf("Error importing Link Branding: %s", err.Error()),
		)

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
