package sendgrid

import (
	"context"
	"fmt"
	"strconv"

	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &linkbrandResource{}
	_ resource.ResourceWithConfigure   = &linkbrandResource{}
	_ resource.ResourceWithImportState = &linkbrandResource{}
)

func NewLinkbrandResource() resource.Resource {
	return &linkbrandResource{}
}

type linkbrandResource struct {
	client *sendgrid.Client
}

type LinkbrandResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Domain        types.String `tfsdk:"domain"`
	Subdomain     types.String `tfsdk:"subdomain"`
	Username      types.String `tfsdk:"username"`
	UserId        types.Int64  `tfsdk:"user_id"`
	Defaultdomain types.Bool   `tfsdk:"default"`
	Valid         types.Bool   `tfsdk:"valid"`
	Legacy        types.Bool   `tfsdk:"legacy"`
	DCNAME        types.Object `tfsdk:"domain_cname"`
	OCNAME        types.Object `tfsdk:"owner_cname"`
}

func (r *linkbrandResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_linkbrand"
}

func (r *linkbrandResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage link branding",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of the link branding",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "The ID of the user",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Description: "The domain name",
				Required:    true,
			},
			"subdomain": schema.StringAttribute{
				Description: "The subdomain name",
				Computed:    true,
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default": schema.BoolAttribute{
				Description: "The default domain",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"legacy": schema.BoolAttribute{
				Description: "The legacy domain",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"valid": schema.BoolAttribute{
				Description: "The valid domain",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},

			"domain_cname": schema.SingleNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"valid": schema.BoolAttribute{
						Description: "The valid domain",
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"types": schema.StringAttribute{
						Description: "The type of domain",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"host": schema.StringAttribute{
						Description: "The host of domain",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"data": schema.StringAttribute{
						Description: "The data of domain",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},

			"owner_cname": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"valid": schema.BoolAttribute{
						Description: "The valid domain",
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"types": schema.StringAttribute{
						Description: "The type of domain",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"host": schema.StringAttribute{
						Description: "The host of domain",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"data": schema.StringAttribute{
						Description: "The data of domain",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *linkbrandResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var newstate LinkbrandResourceModel
	diags := req.Plan.Get(ctx, &newstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Creating Link Branding",
			fmt.Sprintf("Error creating link branding: %s", diags.Errors()),
		)
		return
	}

	itemState := sendgrid.LinkAuth{
		Domain:        newstate.Domain.ValueString(),
		Subdomain:     newstate.Subdomain.ValueString(),
		Defaultdomain: newstate.Defaultdomain.ValueBool(),
	}

	if newstate.Valid.ValueBool() {
		resp.Diagnostics.AddError(
			"Error creating link branding:",
			"Valid parameter can not be true while creating link brand.",
		)
		return
	}

	newItem, err := r.client.CreateLinkBrand(ctx, itemState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating link branding",
			fmt.Sprintf("Error creating link branding: %s", err),
		)
		return
	}

	//convertToMap := structToMap(newItem.DNSDetails)
	tflog.Debug(ctx, "RetriveData:", map[string]any{"item": newItem})
	elementsValue := map[string]attr.Value{
		"valid": types.BoolValue(newItem.DNSDetails.DCNAME.Valid),
		"types": types.StringValue(newItem.DNSDetails.DCNAME.Type),
		"host":  types.StringValue(newItem.DNSDetails.DCNAME.Host),
		"data":  types.StringValue(newItem.DNSDetails.DCNAME.Data),
	}

	elementsType := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	dcnameMapVlaue, diags := types.ObjectValue(elementsType, elementsValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	ocnmaelementsValue := map[string]attr.Value{
		"valid": types.BoolValue(newItem.DNSDetails.OCNAME.Valid),
		"types": types.StringValue(newItem.DNSDetails.OCNAME.Type),
		"host":  types.StringValue(newItem.DNSDetails.OCNAME.Host),
		"data":  types.StringValue(newItem.DNSDetails.OCNAME.Data),
	}

	ocnamelementsType := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	ocnameMapVlaue, diags := types.ObjectValue(ocnamelementsType, ocnmaelementsValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	newstate = LinkbrandResourceModel{
		ID:            types.Int64Value(newItem.ID),
		UserId:        types.Int64Value(newItem.UserId),
		Domain:        types.StringValue(newItem.Domain),
		Subdomain:     types.StringValue(newItem.Subdomain),
		Username:      types.StringValue(newItem.Username),
		Defaultdomain: types.BoolValue(newItem.Defaultdomain),
		Legacy:        types.BoolValue(newItem.Legacy),
		Valid:         types.BoolValue(newItem.Valid),
		DCNAME:        dcnameMapVlaue,
		OCNAME:        ocnameMapVlaue,
	}

	//getelemements := make(map[string]DomainAuthRecord)

	diags = resp.State.Set(ctx, newstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *linkbrandResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var readstate LinkbrandResourceModel
	diags := req.State.Get(ctx, &readstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Error reading link branding: %s", diags.Errors()),
		)
		return
	}

	readitem, err := r.client.Getlinkbrand(ctx, sendgrid.LinkAuth{ID: readstate.ID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Error reading link branding: %s", err),
		)
		return
	}

	readcnameValue := map[string]attr.Value{
		"valid": types.BoolValue(readitem.DNSDetails.DCNAME.Valid),
		"types": types.StringValue(readitem.DNSDetails.DCNAME.Type),
		"host":  types.StringValue(readitem.DNSDetails.DCNAME.Host),
		"data":  types.StringValue(readitem.DNSDetails.DCNAME.Data),
	}

	readcnameType := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	readcnameMapVlaue, diags := types.ObjectValue(readcnameType, readcnameValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	readocnameValue := map[string]attr.Value{
		"valid": types.BoolValue(readitem.DNSDetails.OCNAME.Valid),
		"types": types.StringValue(readitem.DNSDetails.OCNAME.Type),
		"host":  types.StringValue(readitem.DNSDetails.OCNAME.Host),
		"data":  types.StringValue(readitem.DNSDetails.OCNAME.Data),
	}

	readocnameType := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	readocnameMapVlaue, diags := types.ObjectValue(readocnameType, readocnameValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	readstate = LinkbrandResourceModel{
		ID:            types.Int64Value(readitem.ID),
		UserId:        types.Int64Value(readitem.UserId),
		Domain:        types.StringValue(readitem.Domain),
		Subdomain:     types.StringValue(readitem.Subdomain),
		Username:      types.StringValue(readitem.Username),
		Defaultdomain: types.BoolValue(readitem.Defaultdomain),
		Legacy:        types.BoolValue(readitem.Legacy),
		Valid:         types.BoolValue(readitem.Valid),
		DCNAME:        readcnameMapVlaue,
		OCNAME:        readocnameMapVlaue,
	}

	diags = resp.State.Set(ctx, readstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *linkbrandResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var updateplan, updatestate LinkbrandResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &updateplan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &updatestate)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Error updating link branding",
			fmt.Sprintf("Error updaing link branding: %s", resp.Diagnostics.Errors()),
		)
		return
	}

	itemState := sendgrid.LinkAuth{
		ID:            updateplan.ID.ValueInt64(),
		Defaultdomain: updateplan.Defaultdomain.ValueBool(),
		Valid:         updateplan.Valid.ValueBool(),
	}

	if updateplan.Defaultdomain != updatestate.Defaultdomain {
		updaterespBody, err := r.client.Updatelinkbrand(ctx, itemState)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating link branding",
				fmt.Sprintf("Error validating link brand: %s", err),
			)
			return
		}

		if updaterespBody.Defaultdomain {
			valdcnameValue := map[string]attr.Value{
				"valid": types.BoolValue(updaterespBody.DNSDetails.DCNAME.Valid),
				"types": types.StringValue(updaterespBody.DNSDetails.DCNAME.Type),
				"host":  types.StringValue(updaterespBody.DNSDetails.DCNAME.Host),
				"data":  types.StringValue(updaterespBody.DNSDetails.DCNAME.Data),
			}

			valdcnameType := map[string]attr.Type{
				"valid": types.BoolType,
				"types": types.StringType,
				"host":  types.StringType,
				"data":  types.StringType,
			}

			valdcnameMapVlaue, diags := types.ObjectValue(valdcnameType, valdcnameValue)

			//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
			if diags.HasError() {
				resp.Diagnostics.AddError(
					"Error reading link branding",
					fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
				)
				return
			}

			valocnameValue := map[string]attr.Value{
				"valid": types.BoolValue(updaterespBody.DNSDetails.OCNAME.Valid),
				"types": types.StringValue(updaterespBody.DNSDetails.OCNAME.Type),
				"host":  types.StringValue(updaterespBody.DNSDetails.OCNAME.Host),
				"data":  types.StringValue(updaterespBody.DNSDetails.OCNAME.Data),
			}

			valocnameType := map[string]attr.Type{
				"valid": types.BoolType,
				"types": types.StringType,
				"host":  types.StringType,
				"data":  types.StringType,
			}

			valocnameMapVlaue, diags := types.ObjectValue(valocnameType, valocnameValue)

			//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
			if diags.HasError() {
				resp.Diagnostics.AddError(
					"Error reading link branding",
					fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
				)
				return
			}

			updateplan = LinkbrandResourceModel{
				ID:            types.Int64Value(updaterespBody.ID),
				UserId:        types.Int64Value(updaterespBody.UserId),
				Domain:        types.StringValue(updaterespBody.Domain),
				Subdomain:     types.StringValue(updaterespBody.Subdomain),
				Username:      types.StringValue(updaterespBody.Username),
				Defaultdomain: types.BoolValue(updaterespBody.Defaultdomain),
				Legacy:        types.BoolValue(updaterespBody.Legacy),
				Valid:         types.BoolValue(updaterespBody.Valid),
				DCNAME:        valdcnameMapVlaue,
				OCNAME:        valocnameMapVlaue,
			}

			resp.Diagnostics.Append(resp.State.Set(ctx, updateplan)...)
			if resp.Diagnostics.HasError() {
				return
			}
		} else {
			resp.Diagnostics.AddError(
				"Error updating link branding",
				fmt.Sprintf("Error updating link brand: %s", "updating link brand failed. Please check."),
			)
			return
		}
	} else {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *linkbrandResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Preparing to delete item resource")
	var deletestate LinkbrandResourceModel
	diags := req.State.Get(ctx, &deletestate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Error reading link branding: %s", diags.Errors()),
		)
		return
	}

	_, err := r.client.Deletelinkbrand(ctx, fmt.Sprintf("%d", deletestate.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting link branding",
			fmt.Sprintf("Error deleting link branding: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Deleted item resource", map[string]any{"success": true})
}

// Configure adds the provider configured client to the resource.
func (r *linkbrandResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *linkbrandResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
