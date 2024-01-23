package sendgrid

import (
	"context"
	"fmt"

	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &domainvalidateResource{}
	_ resource.ResourceWithConfigure = &domainvalidateResource{}
)

func NewDomainValidateResource() resource.Resource {
	return &domainvalidateResource{}
}

type domainvalidateResource struct {
	client *sendgrid.Client
}

type DomainvalidResourceModel struct {
	ID    types.Int64 `tfsdk:"id"`
	Valid types.Bool  `tfsdk:"valid"`
}

func (r *domainvalidateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_validate_domain"
}

func (r *domainvalidateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *domainvalidateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var domainvalnewstate DomainvalidResourceModel
	diags := req.Plan.Get(ctx, &domainvalnewstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Validating Domain",
			fmt.Sprintf("Error validating domain: %s", diags.Errors()),
		)
		return
	}

	// Create the resource.
	domainvalitem := sendgrid.DomainAuth{
		ID:    domainvalnewstate.ID.ValueInt64(),
		Valid: domainvalnewstate.Valid.ValueBool(),
	}

	domainvalresponse, err := r.client.ValidateDomainAuth(ctx, domainvalitem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Validating Domain",
			fmt.Sprintf("Error validating domain: %s", err.Error()),
		)
		return
	}

	// Set the Terraform state.
	if domainvalresponse.Valid {
		tflog.Debug(ctx, "RetriveData:", map[string]any{"item": domainvalresponse})
		domainvalnewstate = DomainvalidResourceModel{
			ID:    types.Int64Value(domainvalresponse.ID),
			Valid: types.BoolValue(domainvalresponse.Valid),
		}

		//getelemements := make(map[string]DomainAuthRecord)

		diags = resp.State.Set(ctx, domainvalnewstate)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Domain Validation Failed",
			fmt.Sprintf("Error validating domain: %s", err),
		)
		return
	}
}

// Read reads the Terraform state and returns the up-to-date configuration.
func (r *domainvalidateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var domainvalreadstate DomainvalidResourceModel
	diags := req.State.Get(ctx, &domainvalreadstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain",
			fmt.Sprintf("Error reading domain: %s", diags.Errors()),
		)
		return
	}

	domainvalreaditem, err := r.client.GetDomainAuth(ctx, sendgrid.DomainAuth{ID: domainvalreadstate.ID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading domain",
			fmt.Sprintf("Error reading domain: %s", err),
		)
		return
	}
	domainvalreadstate = DomainvalidResourceModel{
		ID:    types.Int64Value(domainvalreaditem.ID),
		Valid: types.BoolValue(domainvalreaditem.Valid),
	}

	//getelemements := make(map[string]DomainAuthRecord)

	diags = resp.State.Set(ctx, domainvalreadstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates an existing resource with new values.
func (r *domainvalidateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes an existing resource.
func (r *domainvalidateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// Configure configures the resource, which is called on both create and update.
func (r *domainvalidateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
