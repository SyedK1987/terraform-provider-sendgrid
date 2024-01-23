package sendgrid

import (
	"context"
	"fmt"
	sendgrid "terraform-provider-sendgrid/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &domainAuthDataSource{}
	_ datasource.DataSourceWithConfigure = &domainAuthDataSource{}
)

func NewdomainauthDataSource() datasource.DataSource {

	return &domainAuthDataSource{}

}

type domainAuthDataSource struct {
	client *sendgrid.Client
}

type ReadDomainauthResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	UserId        types.Int64  `tfsdk:"user_id"`
	Domain        types.String `tfsdk:"domain"`
	Subdomain     types.String `tfsdk:"subdomain"`
	CustomDKIM    types.String `tfsdk:"custom_dkim"`
	Username      types.String `tfsdk:"username"`
	Ips           types.List   `tfsdk:"ips"`
	CusomSPF      types.Bool   `tfsdk:"custom_spf"`
	Defaultdomain types.Bool   `tfsdk:"default"`
	Legacy        types.Bool   `tfsdk:"legacy"`
	AutoSecurity  types.Bool   `tfsdk:"auto_security"`
	Valid         types.Bool   `tfsdk:"valid"`
	DKIM1         types.Object `tfsdk:"dkim1"`
	DKIM2         types.Object `tfsdk:"dkim2"`
	MCNAME        types.Object `tfsdk:"mail_cname"`
	Subusers      types.List   `tfsdk:"subusers"`
}

func (d *domainAuthDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_authentication"
}

func (d *domainAuthDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Allows to retrive authenticataed domain details",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of the domain authentication",
				Required:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "The ID of the user",
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The domain name",
				Computed:    true,
			},
			"custom_dkim": schema.StringAttribute{
				Description: "The custom DKIM",
				Computed:    true,
			},
			"subdomain": schema.StringAttribute{
				Description: "The subdomain name",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username",
				Computed:    true,
			},
			"ips": schema.ListAttribute{
				Description: "The list of IP addresses",
				Computed:    true,
				ElementType: types.StringType,
			},
			"custom_spf": schema.BoolAttribute{
				Description: "The custom SPF",
				Computed:    true,
			},
			"default": schema.BoolAttribute{
				Description: "The default domain",
				Computed:    true,
			},
			"legacy": schema.BoolAttribute{
				Description: "The legacy domain",
				Computed:    true,
			},
			"auto_security": schema.BoolAttribute{
				Description: "The auto security",
				Computed:    true,
				Optional:    true,
			},
			"valid": schema.BoolAttribute{
				Description: "The valid domain",
				Computed:    true,
				Optional:    true,
			},

			"dkim1": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"valid": schema.BoolAttribute{
						Description: "The valid domain",
						Computed:    true,
					},
					"types": schema.StringAttribute{
						Description: "The type of domain",
						Computed:    true,
					},
					"host": schema.StringAttribute{
						Description: "The host of domain",
						Computed:    true,
					},
					"data": schema.StringAttribute{
						Description: "The data of domain",
						Computed:    true,
					},
				},
				Computed: true,
			},

			"dkim2": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"valid": schema.BoolAttribute{
						Description: "The valid domain",
						Computed:    true,
					},
					"types": schema.StringAttribute{
						Description: "The type of domain",
						Computed:    true,
					},
					"host": schema.StringAttribute{
						Description: "The host of domain",
						Computed:    true,
					},
					"data": schema.StringAttribute{
						Description: "The data of domain",
						Computed:    true,
					},
				},
				Computed: true,
			},

			"mail_cname": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"valid": schema.BoolAttribute{
						Description: "The valid domain",
						Computed:    true,
					},
					"types": schema.StringAttribute{
						Description: "The type of domain",
						Computed:    true,
					},
					"host": schema.StringAttribute{
						Description: "The host of domain",
						Computed:    true,
					},
					"data": schema.StringAttribute{
						Description: "The data of domain",
						Computed:    true,
					},
				},
				Computed: true,
			},
			"subusers": schema.ListNestedAttribute{
				Computed: true,
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

func (d *domainAuthDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var refstate ReadDomainauthResourceModel

	diags := req.Config.Get(ctx, &refstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link brand",
			fmt.Sprintf("Error reading link brand: %s", resp.Diagnostics.Errors()),
		)
		return
	}
	refitem, err := d.client.GetDomainAuth(ctx, sendgrid.DomainAuth{ID: refstate.ID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Error reading domain authentication: %s", err),
		)
		return
	}

	refkim1Value := map[string]attr.Value{
		"valid": types.BoolValue(refitem.DNSDetails.DKIM1.Valid),
		"types": types.StringValue(refitem.DNSDetails.DKIM1.Type),
		"host":  types.StringValue(refitem.DNSDetails.DKIM1.Host),
		"data":  types.StringValue(refitem.DNSDetails.DKIM1.Data),
	}

	refdkim1Type := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	refdkimMapVlaue, diags := types.ObjectValue(refdkim1Type, refkim1Value)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Try Harded SK: %s", diags.Errors()),
		)
		return
	}

	refdkim2Value := map[string]attr.Value{
		"valid": types.BoolValue(refitem.DNSDetails.DKIM2.Valid),
		"types": types.StringValue(refitem.DNSDetails.DKIM2.Type),
		"host":  types.StringValue(refitem.DNSDetails.DKIM2.Host),
		"data":  types.StringValue(refitem.DNSDetails.DKIM2.Data),
	}

	refdkim2MapVlaue, diags := types.ObjectValue(refdkim1Type, refdkim2Value)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Try Harded SK: %s", diags.Errors()),
		)
		return
	}

	refmcnameValue := map[string]attr.Value{
		"valid": types.BoolValue(refitem.DNSDetails.MailCNAME.Valid),
		"types": types.StringValue(refitem.DNSDetails.MailCNAME.Type),
		"host":  types.StringValue(refitem.DNSDetails.MailCNAME.Host),
		"data":  types.StringValue(refitem.DNSDetails.MailCNAME.Data),
	}

	refmcnameMapVlaue, diags := types.ObjectValue(refdkim1Type, refmcnameValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Try Harded SK: %s", diags.Errors()),
		)
		return
	}

	getIplist, diags := types.ListValueFrom(ctx, types.StringType, refitem.Ips)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading domain authentication",
			fmt.Sprintf("Unable to convert IPList from API: %s", diags.Errors()),
		)
		return
	}

	refstate = ReadDomainauthResourceModel{
		ID:            types.Int64Value(refitem.ID),
		UserId:        types.Int64Value(refitem.UserId),
		Domain:        types.StringValue(refitem.Domain),
		Subdomain:     types.StringValue(refitem.Subdomain),
		Username:      types.StringValue(refitem.Username),
		Ips:           getIplist,
		CusomSPF:      types.BoolValue(refitem.CustomSPF),
		Defaultdomain: types.BoolValue(refitem.Defaultdomain),
		Legacy:        types.BoolValue(refitem.Legacy),
		Valid:         types.BoolValue(refitem.Valid),
		DKIM1:         refdkimMapVlaue,
		DKIM2:         refdkim2MapVlaue,
		MCNAME:        refmcnameMapVlaue,
	}

	if len(refitem.Subusers) > 0 {
		elements := []attr.Value{}
		for _, subuser := range refitem.Subusers {
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

		refstate.Subusers = types.ListValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username": types.StringType,
					"user_id":  types.Int64Type,
				},
			}, elements,
		)
	} else {
		refstate.Subusers = types.ListValueMust(
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

	diags = resp.State.Set(ctx, refstate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *domainAuthDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sendgrid.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Sendgrid.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
