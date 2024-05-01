package sendgrid

import (
	"context"
	"fmt"
	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &apikeyResource{}
	_ resource.ResourceWithConfigure   = &apikeyResource{}
	_ resource.ResourceWithImportState = &apikeyResource{}
)

func NewApiKeyResource() resource.Resource {
	return &apikeyResource{}
}

type ApiKey struct {
	Name   types.String `tfsdk:"name"`
	Scopes []string     `tfsdk:"scopes"`
	ID     types.String `tfsdk:"api_key_id"`
	Apikey types.String `tfsdk:"api_key"`
	//	Permission  types.String `tfsdk:"permission"`
	//	Environment types.String `tfsdk:"environment"`
}

type apikeyResource struct {
	client *sendgrid.Client
}

func (r *apikeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *apikeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (r *apikeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the API key",
				Required:    true,
			},
			"scopes": schema.ListAttribute{
				Description: "List of scopes for the API key",
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			// "permission": schema.StringAttribute{
			// 	Description: "Role of the Api Key",
			// 	Optional:    true,
			// 	Computed:    true,
			// 	Default:     stringdefault.StaticString("custom"),
			// 	Validators: []validator.String{
			// 		stringvalidator.OneOf([]string{"full", "custom"}...),
			// 	},
			// },
			"api_key_id": schema.StringAttribute{
				Description: "ID of the API key",
				Computed:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "API key",
				Computed:    true,
			},
		},
	}
}

func (r *apikeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var newstate ApiKey
	//var customDefinedScopes []string
	//stringArray := make([]string, 0)
	var err error

	diags := req.Plan.Get(ctx, &newstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(newstate.Scopes) == 0 {
		resp.Diagnostics.AddError(
			"Unable to create API key",
			fmt.Sprintf("Unable to create API key: %s", "Scopes cannot be empty"),
		)
		return
	}

	if newstate.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Unable to create API key",
			fmt.Sprintf("Unable to create API key: %s", "Name cannot be empty"),
		)
		return
	}

	itemapikey := sendgrid.ChildApiKey{
		Name:   newstate.Name.ValueString(),
		Scopes: newstate.Scopes,
	}

	apikeyrespBody, err := r.client.CreateApiKey(ctx, itemapikey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create API key",
			fmt.Sprintf("Unable to create API key: %s", err.Error()),
		)

		return
	}

	//listmyvalue, _ := types.ListValueFrom(ctx, types.StringType, apikeyrespBody.Scopes)
	var listmyvalue []string
	if len(apikeyrespBody.Scopes) > 0 {
		listmyvalue = append(listmyvalue, apikeyrespBody.Scopes...)
	} else {
		listmyvalue = append(listmyvalue, "")
	}
	newstate = ApiKey{
		ID:     types.StringValue(apikeyrespBody.ID),
		Name:   types.StringValue(apikeyrespBody.Name),
		Scopes: listmyvalue,
		Apikey: types.StringValue(apikeyrespBody.Apikey),
	}

	diags = resp.State.Set(ctx, &newstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *apikeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var readstate ApiKey
	var err error

	diags := req.State.Get(ctx, &readstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readapikeyresponse, err := r.client.ReadApiKey(ctx, readstate.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read API key",
			fmt.Sprintf("Unable to read API key: %s", err),
		)

		return
	}

	//readlistmyvalue, _ := types.ListValueFrom(ctx, types.StringType, readapikeyresponse.Scopes)
	var readlistmyvalue []string
	if len(readapikeyresponse.Scopes) > 0 {
		readlistmyvalue = append(readlistmyvalue, readapikeyresponse.Scopes...)
	} else {
		readlistmyvalue = append(readlistmyvalue, "")
	}
	readstate = ApiKey{
		ID:     types.StringValue(readapikeyresponse.ID),
		Name:   types.StringValue(readapikeyresponse.Name),
		Scopes: readlistmyvalue,
		Apikey: types.StringValue(readapikeyresponse.Apikey),
	}

	diags = resp.State.Set(ctx, readstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *apikeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var updateplan, updatestate ApiKey
	//	var customDefinedScopes []string
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &updatestate)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &updateplan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Error updating ApiKey",
			fmt.Sprintf("Error updating apikey: %s", resp.Diagnostics.Errors()),
		)
		return
	}

	updateitemapikey := sendgrid.ChildApiKey{
		Name:   updateplan.Name.ValueString(),
		Scopes: updateplan.Scopes,
		ID:     updatestate.ID.ValueString(),
	}

	updateapikeyrespBody, err := r.client.UpdateApiKey(ctx, updateitemapikey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update API key Permission",
			fmt.Sprintf("Unable to update API key: %s", err),
		)

		return
	}

	//updatelistmyvalue, _ := types.ListValueFrom(ctx, types.StringType, updateapikeyrespBody.Scopes)
	var updatelistmyvalue []string
	if len(updateapikeyrespBody.Scopes) > 0 {
		updatelistmyvalue = append(updatelistmyvalue, updateapikeyrespBody.Scopes...)
	} else {
		updatelistmyvalue = append(updatelistmyvalue, "")
	}
	updatestate = ApiKey{
		ID:     types.StringValue(updateapikeyrespBody.ID),
		Name:   types.StringValue(updateapikeyrespBody.Name),
		Scopes: updatelistmyvalue,
	}

	//resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scop"), updateplan.Permission)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, updatestate)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Unable to update API key Permission",
			fmt.Sprintf("Unable to update API key: %s", err),
		)

		return
	}
}

func (r *apikeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	tflog.Debug(ctx, "Preparing to delete item resource")
	var deletestate ApiKey
	var err error

	diags := req.State.Get(ctx, &deletestate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err = r.client.DeleteApiKey(ctx, deletestate.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete API key",
			fmt.Sprintf("Unable to delete API key: %s", err),
		)

		return
	}
	tflog.Debug(ctx, "Deleted item resource", map[string]any{"success": true})
}

func (r *apikeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("api_key_id"), req, resp)
}
