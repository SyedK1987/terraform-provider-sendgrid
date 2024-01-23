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
	_ resource.Resource                = &teammateResource{}
	_ resource.ResourceWithConfigure   = &teammateResource{}
	_ resource.ResourceWithImportState = &teammateResource{}
)

func NewTeammateResource() resource.Resource {
	return &teammateResource{}
}

type teammateResource struct {
	client *sendgrid.Client
}

type TeammateModel struct {
	Username  types.String `tfsdk:"username"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	// Address   types.String `tfsdk:"address"`
	// Address2  types.String `tfsdk:"address2"`
	// City      types.String `tfsdk:"city"`
	// State     types.String `tfsdk:"state"`
	// Zip       types.String `tfsdk:"zip"`
	// Coutnry   types.String `tfsdk:"country"`
	// Company   types.String `tfsdk:"company"`
	// Phone     types.String `tfsdk:"phone"`
	IsAdmin        types.Bool  `tfsdk:"is_admin"`
	IsReadOnly     types.Bool  `tfsdk:"is_read_only"`
	ExpirationDate types.Int64 `tfsdk:"expiration_date"`
	//	IsSSO          types.Bool   `tfsdk:"is_sso"`
	UserType types.String `tfsdk:"user_type"`
	Scopes   types.List   `tfsdk:"scopes"`
	Token    types.String `tfsdk:"token"`
}

func (r *teammateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teammate"
}

func (r *teammateResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Helps to create and manage teammate to your account",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Description: "Email address of the teammate",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username of the teammate",
				Computed:    true,
			},
			"first_name": schema.StringAttribute{
				Description: "First name of the teammate",
				Computed:    true,
			},
			"last_name": schema.StringAttribute{
				Description: "Last name of the teammate",
				Computed:    true,
			},
			// "address": schema.StringAttribute{
			// 	Description: "Address of the teammate",
			// 	Computed:    true,
			// },
			// "address2": schema.StringAttribute{
			// 	Description: "Address2 of the teammate",
			// 	Computed:    true,
			// },
			// "city": schema.StringAttribute{
			// 	Description: "City of the teammate",
			// 	Computed:    true,
			// },
			// "state": schema.StringAttribute{
			// 	Description: "State of the teammate",
			// 	Computed:    true,
			// },
			// "zip": schema.StringAttribute{
			// 	Description: "Zip of the teammate",
			// 	Computed:    true,
			// },
			// "country": schema.StringAttribute{
			// 	Description: "Country of the teammate",
			// 	Computed:    true,
			// },
			// "company": schema.StringAttribute{
			// 	Description: "Company of the teammate",
			// 	Computed:    true,
			// },
			// "phone": schema.StringAttribute{
			// 	Description: "Phone of the teammate",
			// 	Computed:    true,
			// },
			"is_read_only": schema.BoolAttribute{
				Description: "Is read only of the teammate",
				Computed:    true,
			},
			"expiration_date": schema.Int64Attribute{
				Description: "Expiration date of the teammate invite",
				Computed:    true,
			},
			"is_admin": schema.BoolAttribute{
				Description: "Is admin of the teammate",
				Required:    true,
			},
			// "is_sso": schema.BoolAttribute{
			// 	Description: "Is sso of the teammate",
			// 	Computed:    true,
			// },
			"user_type": schema.StringAttribute{
				Description: "User type of the teammate",
				Computed:    true,
			},
			"scopes": schema.ListAttribute{
				Description: "Scopes of the teammate",
				Computed:    true,
				ElementType: types.StringType,
				// Validators: []validator.List{
				// 	// Validate this block has a value (not null).
				// 	listvalidator.IsRequired(),
				// },
			},
			"token": schema.StringAttribute{
				Description: "Token of the Pending teammate",
				Computed:    true,
			},
		},
	}
}

func (r *teammateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *teammateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var newstate TeammateModel
	var customTeammateScopes []string
	var err error

	diags := req.Plan.Get(ctx, &newstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !newstate.IsAdmin.ValueBool() {
		customTeammateScopes, err = DeveloperPermission()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected error",
				fmt.Sprintf("Unable to get DeveloperPermission: %s", err),
			)
			return
		}
	}

	teammateitem := sendgrid.User{
		Email:   newstate.Email.ValueString(),
		Scopes:  customTeammateScopes,
		IsAdmin: newstate.IsAdmin.ValueBool(),
	}

	teammaterespBody, err := r.client.CreateTeammate(ctx, teammateitem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("Unable to create teammate: %s", err),
		)
		return
	}

	tmatescopelist, _ := types.ListValueFrom(ctx, types.StringType, teammaterespBody.Scopes)
	newstate = TeammateModel{
		Username:  types.StringValue(teammaterespBody.Username),
		Email:     types.StringValue(teammaterespBody.Email),
		FirstName: types.StringValue(teammaterespBody.FirstName),
		LastName:  types.StringValue(teammaterespBody.LastName),
		// Address:   types.StringValue(teammaterespBody.Address),
		// Address2:  types.StringValue(teammaterespBody.Address2),
		// City:      types.StringValue(teammaterespBody.City),
		// State:     types.StringValue(teammaterespBody.State),
		// Zip:       types.StringValue(teammaterespBody.Zip),
		// Coutnry:   types.StringValue(teammaterespBody.Country),
		// Company:   types.StringValue(teammaterespBody.Company),
		// Phone:     types.StringValue(teammaterespBody.Phone),
		IsAdmin: types.BoolValue(teammaterespBody.IsAdmin),
		//IsSSO:    types.BoolValue(teammaterespBody.IsSSO),
		UserType:       types.StringValue(teammaterespBody.UserType),
		Scopes:         tmatescopelist,
		IsReadOnly:     types.BoolValue(teammaterespBody.IsReadOnly),
		ExpirationDate: types.Int64Value(teammaterespBody.ExpirationDate),
		Token:          types.StringValue(teammaterespBody.Token),
	}

	diags = resp.State.Set(ctx, &newstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *teammateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var readstate TeammateModel
	var err error

	diags := req.State.Get(ctx, &readstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if readstate.Token.ValueString() != "" {
	// 	resp.Diagnostics.AddError(
	// 		"Pending Invite",
	// 		fmt.Sprintf("Unable to update teammate: %s", "Invite is still pending. Please accept the invite first."),
	// 	)
	// 	return
	// }

	teammaterespBody, err := r.client.RefreshTeammate(ctx, readstate.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("Unable to refresh teammate: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "ReadingResource:", map[string]any{"item": teammaterespBody.Token})

	refreshtmatescopelist, _ := types.ListValueFrom(ctx, types.StringType, teammaterespBody.Scopes)
	readstate = TeammateModel{
		Username:       types.StringValue(teammaterespBody.Username),
		FirstName:      types.StringValue(teammaterespBody.FirstName),
		LastName:       types.StringValue(teammaterespBody.LastName),
		Email:          types.StringValue(teammaterespBody.Email),
		IsAdmin:        types.BoolValue(teammaterespBody.IsAdmin),
		UserType:       types.StringValue(teammaterespBody.UserType),
		Scopes:         refreshtmatescopelist,
		IsReadOnly:     types.BoolValue(teammaterespBody.IsReadOnly),
		ExpirationDate: types.Int64Value(teammaterespBody.ExpirationDate),
		Token:          types.StringValue(teammaterespBody.Token),
	}

	diags = resp.State.Set(ctx, &readstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *teammateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var updatestate TeammateModel
	var customTeammateScopes []string
	var err error

	diags := req.Plan.Get(ctx, &updatestate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if updatestate.Token.ValueString() != "" {
		resp.Diagnostics.AddError(
			"Pending Invite",
			fmt.Sprintf("Unable to update teammate: %s", "Invite is still pending. Please accept the invite first."),
		)
		return
	}

	if !updatestate.IsAdmin.ValueBool() {
		customTeammateScopes, err = DeveloperPermission()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected error",
				fmt.Sprintf("Unable to get DeveloperPermission: %s", err),
			)
			return
		}
	}

	teammateitem := sendgrid.User{
		Email:   updatestate.Email.ValueString(),
		Scopes:  customTeammateScopes,
		IsAdmin: updatestate.IsAdmin.ValueBool(),
	}

	updatetmaterespBody, err := r.client.UpdateTeammate(ctx, teammateitem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("Unable to update teammate: %s", err),
		)
		return
	}

	updatetmatescopelist, _ := types.ListValueFrom(ctx, types.StringType, updatetmaterespBody.Scopes)
	updatestate = TeammateModel{
		Username:       types.StringValue(updatetmaterespBody.Username),
		FirstName:      types.StringValue(updatetmaterespBody.FirstName),
		LastName:       types.StringValue(updatetmaterespBody.LastName),
		Email:          types.StringValue(updatetmaterespBody.Email),
		IsAdmin:        types.BoolValue(updatetmaterespBody.IsAdmin),
		UserType:       types.StringValue(updatetmaterespBody.UserType),
		Scopes:         updatetmatescopelist,
		IsReadOnly:     types.BoolValue(updatetmaterespBody.IsReadOnly),
		ExpirationDate: types.Int64Value(updatetmaterespBody.ExpirationDate),
		Token:          types.StringValue(updatetmaterespBody.Token),
	}

	diags = resp.State.Set(ctx, &updatestate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *teammateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	tflog.Debug(ctx, "Preparing to delete item resource")
	var deletestate TeammateModel
	var err error

	diags := req.State.Get(ctx, &deletestate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err = r.client.DeleteTeammate(ctx, deletestate.Email.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("Unable to delete teammate: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Deleted item resource", map[string]any{"success": true})
}

func (r *teammateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("email"), req, resp)
}
