package sendgrid

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sendgrid "terraform-provider-sendgrid/client"

	backoff "github.com/cenkalti/backoff"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &domainauthResource{}
	_ resource.ResourceWithConfigure   = &domainauthResource{}
	_ resource.ResourceWithImportState = &domainauthResource{}
)

func NewDomainAuthResource() resource.Resource {
	return &domainauthResource{}
}

type domainauthResource struct {
	client *sendgrid.Client
}

type DomainauthResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	UserId        types.Int64  `tfsdk:"user_id"`
	Domain        types.String `tfsdk:"domain"`
	Subdomain     types.String `tfsdk:"subdomain"`
	Environment   types.String `tfsdk:"environment"`
	CustomDKIM    types.String `tfsdk:"custom_dkim_selector"`
	Username      types.String `tfsdk:"username"`
	Ips           []string     `tfsdk:"ips"`
	CusomSPF      types.Bool   `tfsdk:"custom_spf"`
	Defaultdomain types.Bool   `tfsdk:"default"`
	Legacy        types.Bool   `tfsdk:"legacy"`
	//	AutoSecurity  types.Bool   `tfsdk:"automatic_security"`
	Valid    types.Bool   `tfsdk:"valid"`
	DKIM1    types.Object `tfsdk:"dkim1"`
	DKIM2    types.Object `tfsdk:"dkim2"`
	MCNAME   types.Object `tfsdk:"mail_cname"`
	Subusers types.List   `tfsdk:"subusers"`
}

type DomainAuthRecord struct {
	Valid types.Bool   `tfsdk:"valid"`
	Types types.String `tfsdk:"types"`
	Host  types.String `tfsdk:"host"`
	Data  types.String `tfsdk:"data"`
}

func (r *domainauthResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_authentication"
}

func (r *domainauthResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *domainauthResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Resource to manage domain authentication",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of the domain authentication",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"environment": schema.StringAttribute{
				Description: "The environment of the sendgrid account",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"prod", "nonprod"}...),
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_dkim_selector": schema.StringAttribute{
				Description: "The custom DKIM",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				//	Required:    true,
				// Validators: []validator.String{
				// 	stringvalidator.LengthBetween(3, 3),
				// 	stringvalidator.RegexMatches(
				// 		regexp.MustCompile(`^[a-z0-9]+$`),
				// 		"must be lowercase alphanumeric characters only",
				// 	),
				// },
			},
			"username": schema.StringAttribute{
				Description: "The username",
				//Optional:    true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ips": schema.ListAttribute{
				Description: "The list of IP addresses",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_spf": schema.BoolAttribute{
				Description: "The custom SPF",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default": schema.BoolAttribute{
				Description: "The default domain",
				Computed:    true,
				Optional:    true,
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

			"dkim1": schema.SingleNestedAttribute{
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

			"dkim2": schema.SingleNestedAttribute{
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

			"mail_cname": schema.SingleNestedAttribute{
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
			"subusers": schema.ListNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							Description: "The subuser username to be associate with the domain",
							Computed:    true,
						},
						"user_id": schema.Int64Attribute{
							Description: "The subuser id to be associate with the domain",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *domainauthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var newstate DomainauthResourceModel
	var autodecidedDKIM string
	diags := req.Plan.Get(ctx, &newstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Creating domain authentication",
			fmt.Sprintf("Error creating domain authentication: %s", diags.Errors()),
		)
		return
	}

	if newstate.Environment.ValueString() == "prod" {
		autodecidedDKIM = "sp1"
	} else {
		autodecidedDKIM = "sn1"
	}

	itemState := sendgrid.DomainAuth{
		Domain:        newstate.Domain.ValueString(),
		Subdomain:     newstate.Subdomain.ValueString(),
		CustomDKIM:    autodecidedDKIM,
		Ips:           newstate.Ips,
		CustomSPF:     newstate.CusomSPF.ValueBool(),
		Defaultdomain: newstate.Defaultdomain.ValueBool(),
		//Username:      newstate.Username.ValueString(),
		//	AutoSecurity:  newstate.AutoSecurity.ValueBool(),
	}

	if newstate.Valid.ValueBool() {
		resp.Diagnostics.AddError(
			"Error creating domain authentication:",
			"Valid parameter can not be true while creating domain.",
		)
		return
	}

	//newItem, err := r.client.CreateDomainAuth(ctx, itemState)

	// create retrycontext with backoff
	retryctx := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)

	var newItem *sendgrid.DomainAuth
	var err error
	err = backoff.Retry(func() error {
		newItem, err = r.client.CreateDomainAuth(ctx, itemState)
		return err
	}, retryctx)

	//newItem, err := newstate.Timeout.Create(ctx, func(ctx context.Context) (interface{}, error) {

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating domain authentication",
			fmt.Sprintf("Error creating domain authentication: %s", err),
		)
		return
	}

	//convertToMap := structToMap(newItem.DNSDetails)
	tflog.Debug(ctx, "RetriveData:", map[string]any{"item": newItem})
	elementsValue := map[string]attr.Value{
		"valid": types.BoolValue(newItem.DNSDetails.DKIM1.Valid),
		"types": types.StringValue(newItem.DNSDetails.DKIM1.Type),
		"host":  types.StringValue(newItem.DNSDetails.DKIM1.Host),
		"data":  types.StringValue(newItem.DNSDetails.DKIM1.Data),
	}

	elementsType := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	dkimMapVlaue, diags := types.ObjectValue(elementsType, elementsValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	dkim2elementsValue := map[string]attr.Value{
		"valid": types.BoolValue(newItem.DNSDetails.DKIM2.Valid),
		"types": types.StringValue(newItem.DNSDetails.DKIM2.Type),
		"host":  types.StringValue(newItem.DNSDetails.DKIM2.Host),
		"data":  types.StringValue(newItem.DNSDetails.DKIM2.Data),
	}

	dkim2elementsType := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	dkim2MapVlaue, diags := types.ObjectValue(dkim2elementsType, dkim2elementsValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	mcnameelementsValue := map[string]attr.Value{
		"valid": types.BoolValue(newItem.DNSDetails.MailCNAME.Valid),
		"types": types.StringValue(newItem.DNSDetails.MailCNAME.Type),
		"host":  types.StringValue(newItem.DNSDetails.MailCNAME.Host),
		"data":  types.StringValue(newItem.DNSDetails.MailCNAME.Data),
	}

	mcnameelementsType := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	mcnameMapVlaue, diags := types.ObjectValue(mcnameelementsType, mcnameelementsValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	// var receiveIps []string
	// if len(newItem.Ips) > 0 {
	// 	receiveIps = append(receiveIps, newItem.Ips...)

	// } else {
	// 	receiveIps = append(receiveIps, "")

	// }

	newstate = DomainauthResourceModel{
		ID:            types.Int64Value(newItem.ID),
		UserId:        types.Int64Value(newItem.UserId),
		Domain:        types.StringValue(newItem.Domain),
		Subdomain:     types.StringValue(newItem.Subdomain),
		Environment:   newstate.Environment,
		CustomDKIM:    types.StringValue(autodecidedDKIM),
		Username:      types.StringValue(newItem.Username),
		Ips:           newItem.Ips,
		CusomSPF:      newstate.CusomSPF,
		Defaultdomain: types.BoolValue(newItem.Defaultdomain),
		Legacy:        types.BoolValue(newItem.Legacy),
		//	AutoSecurity:  newstate.AutoSecurity,
		Valid:  types.BoolValue(newItem.Valid),
		DKIM1:  dkimMapVlaue,
		DKIM2:  dkim2MapVlaue,
		MCNAME: mcnameMapVlaue,
	}

	if len(newItem.Subusers) > 0 {
		elements := []attr.Value{}
		for _, subuser := range newItem.Subusers {
			elements = append(elements, types.ObjectValueMust(
				map[string]attr.Type{
					"username": types.StringType,
					"user_id":  types.Int64Type,
				},
				map[string]attr.Value{
					"username": types.StringValue(subuser.Username),
					"user_id":  types.Int64Value(subuser.UserID),
				},
			))
		}

		newstate.Subusers = types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username": types.StringType,
					"user_id":  types.Int64Type,
				},
			}, elements,
		)
	} else {
		newstate.Subusers = types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username": types.StringType,
					"user_id":  types.Int64Type,
				},
			}, []attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"username": types.StringType,
						"user_id":  types.Int64Type,
					},
					map[string]attr.Value{
						"username": types.StringValue(""),
						"user_id":  types.Int64Value(0),
					},
				),
			},
		)
	}

	diags = resp.State.Set(ctx, newstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *domainauthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var readstate DomainauthResourceModel
	diags := req.State.Get(ctx, &readstate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Error reading domain authentication: %s", diags.Errors()),
		)
		return
	}

	readitem, err := r.client.GetDomainAuth(ctx, sendgrid.DomainAuth{ID: readstate.ID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Error reading domain authentication: %s", err),
		)
		return
	}

	readdkim1Value := map[string]attr.Value{
		"valid": types.BoolValue(readitem.DNSDetails.DKIM1.Valid),
		"types": types.StringValue(readitem.DNSDetails.DKIM1.Type),
		"host":  types.StringValue(readitem.DNSDetails.DKIM1.Host),
		"data":  types.StringValue(readitem.DNSDetails.DKIM1.Data),
	}

	readdkim1Type := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	readkimMapVlaue, diags := types.ObjectValue(readdkim1Type, readdkim1Value)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	readdkim2Value := map[string]attr.Value{
		"valid": types.BoolValue(readitem.DNSDetails.DKIM2.Valid),
		"types": types.StringValue(readitem.DNSDetails.DKIM2.Type),
		"host":  types.StringValue(readitem.DNSDetails.DKIM2.Host),
		"data":  types.StringValue(readitem.DNSDetails.DKIM2.Data),
	}

	readdkim2Type := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	readdkim2MapVlaue, diags := types.ObjectValue(readdkim2Type, readdkim2Value)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	readmcnameValue := map[string]attr.Value{
		"valid": types.BoolValue(readitem.DNSDetails.MailCNAME.Valid),
		"types": types.StringValue(readitem.DNSDetails.MailCNAME.Type),
		"host":  types.StringValue(readitem.DNSDetails.MailCNAME.Host),
		"data":  types.StringValue(readitem.DNSDetails.MailCNAME.Data),
	}

	readmcnameType := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	readmcnameMapVlaue, diags := types.ObjectValue(readmcnameType, readmcnameValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	getcdkim := strings.Split(readitem.DNSDetails.DKIM1.Host, ".")[0]

	// var receiveIps []string
	// if len(readitem.Ips) > 0 {
	// 	receiveIps = append(receiveIps, readitem.Ips...)
	// } else {
	// 	receiveIps = append(receiveIps, "")
	// }

	if (readitem.Ips == nil) || (len(readitem.Ips) == 0) {
		readitem.Ips = []string{}
	}

	readstate = DomainauthResourceModel{
		ID:            types.Int64Value(readitem.ID),
		UserId:        types.Int64Value(readitem.UserId),
		Domain:        types.StringValue(readitem.Domain),
		Subdomain:     types.StringValue(readitem.Subdomain),
		Environment:   readstate.Environment,
		CustomDKIM:    types.StringValue(getcdkim),
		Username:      types.StringValue(readitem.Username),
		Ips:           readitem.Ips,
		CusomSPF:      types.BoolValue(readitem.CustomSPF),
		Defaultdomain: types.BoolValue(readitem.Defaultdomain),
		Legacy:        types.BoolValue(readitem.Legacy),
		Valid:         types.BoolValue(readitem.Valid),
		DKIM1:         readkimMapVlaue,
		DKIM2:         readdkim2MapVlaue,
		MCNAME:        readmcnameMapVlaue,
	}

	if len(readitem.Subusers) > 0 {
		elements := []attr.Value{}
		for _, subuser := range readitem.Subusers {
			elements = append(elements, types.ObjectValueMust(
				map[string]attr.Type{
					"username": types.StringType,
					"user_id":  types.Int64Type,
				},
				map[string]attr.Value{
					"username": types.StringValue(subuser.Username),
					"user_id":  types.Int64Value(subuser.UserID),
				},
			))
		}

		readstate.Subusers = types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username": types.StringType,
					"user_id":  types.Int64Type,
				},
			}, elements,
		)
	} else {
		readstate.Subusers = types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username": types.StringType,
					"user_id":  types.Int64Type,
				},
			}, []attr.Value{
				types.ObjectValueMust(
					map[string]attr.Type{
						"username": types.StringType,
						"user_id":  types.Int64Type,
					},
					map[string]attr.Value{
						"username": types.StringValue(""),
						"user_id":  types.Int64Value(0),
					},
				),
			},
		)
	}

	diags = resp.State.Set(ctx, readstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *domainauthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var updateplan, updatestate DomainauthResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &updateplan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &updatestate)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Error updating domain authentication",
			fmt.Sprintf("Error updating domain authentication: %s", resp.Diagnostics.Errors()),
		)
		return
	}

	itemState := sendgrid.DomainAuth{
		ID:            updatestate.ID.ValueInt64(),
		CustomSPF:     updatestate.CusomSPF.ValueBool(),
		Defaultdomain: updatestate.Defaultdomain.ValueBool(),
		Valid:         updatestate.Valid.ValueBool(),
	}

	if updatestate.CusomSPF != updateplan.CusomSPF || updatestate.Defaultdomain != updateplan.Defaultdomain {

		tflog.Debug(ctx, "Preparing to update item resource:", map[string]any{"id:": fmt.Sprintf("%+v", updatestate)})

		updaterespBody, err := r.client.UpdateDomainAuth(ctx, itemState)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating domain authentication",
				fmt.Sprintf("Error updating domain authentication: %s", err),
			)
			return
		}

		if updatestate.CusomSPF.ValueBool() && !updaterespBody.CustomSPF {
			resp.Diagnostics.AddError(
				"Error updating domain authentication",
				fmt.Sprintf("Error updating domain authentication: %s", "Failed to set custom_spf true."),
			)
			return
		} else if updatestate.Defaultdomain.ValueBool() && !updaterespBody.Defaultdomain {
			resp.Diagnostics.AddError(
				"Error updating domain authentication",
				fmt.Sprintf("Error updating domain authentication: %s", "Failed to set default true."),
			)
			return
		} else {
			tflog.Debug(ctx, "RetriveData:", map[string]any{"item": updaterespBody})
			keypairforupdateset := map[string]interface{}{
				"custom_spf": updaterespBody.CustomSPF,
				"default":    updaterespBody.Defaultdomain,
			}

			for k, v := range keypairforupdateset {
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(k), v)...)
				if resp.Diagnostics.HasError() {
					return
				}
			}
		}
	} else {
		return
	}

}

func (r *domainauthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Preparing to delete item resource")
	var deletestate DomainauthResourceModel
	diags := req.State.Get(ctx, &deletestate)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Error reading domain authentication: %s", diags.Errors()),
		)
		return
	}

	_, err := r.client.DeleteDomainAuth(ctx, fmt.Sprintf("%d", deletestate.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting domain authentication",
			fmt.Sprintf("Error deleting domain authentication: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Deleted item resource", map[string]any{"success": true})
}

func (r *domainauthResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing Domain Authentication",
			fmt.Sprintf("Error importing domain Authentication: %s", err.Error()),
		)

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
