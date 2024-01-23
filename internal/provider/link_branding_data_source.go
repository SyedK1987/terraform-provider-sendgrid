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
	_ datasource.DataSource              = &linkbrandDataSource{}
	_ datasource.DataSourceWithConfigure = &linkbrandDataSource{}
)

func NewlinkbrandDataSource() datasource.DataSource {

	return &linkbrandDataSource{}

}

type linkbrandDataSource struct {
	client *sendgrid.Client
}

func (d *linkbrandDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_linkbrand"
}

func (d *linkbrandDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Allows to retrive link brand details",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of the link brand",
				Required:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "The ID of the user",
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The domain of the link brand",
				Computed:    true,
			},
			"subdomain": schema.StringAttribute{
				Description: "The subdomain of the link brand",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username of the link brand",
				Computed:    true,
			},
			"valid": schema.BoolAttribute{
				Description: "The valid of the link brand",
				Computed:    true,
			},
			"default": schema.BoolAttribute{
				Description: "The default of the link brand",
				Computed:    true,
			},
			"legacy": schema.BoolAttribute{
				Description: "The legacy of the link brand",
				Computed:    true,
			},
			"domain_cname": schema.SingleNestedAttribute{
				Computed: true,
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
			},

			"owner_cname": schema.SingleNestedAttribute{
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
		},
	}
}

func (d *linkbrandDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *linkbrandDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var datalinkbrand LinkbrandResourceModel

	diags := req.Config.Get(ctx, &datalinkbrand)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link brand",
			fmt.Sprintf("Error reading link brand: %s", resp.Diagnostics.Errors()),
		)
		return
	}

	getlinkbranditem, err := d.client.Getlinkbrand(ctx, sendgrid.LinkAuth{ID: datalinkbrand.ID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading link brand",
			fmt.Sprintf("Error retrieving link brand details: %s", err),
		)
		return
	}

	recievedcnameValue := map[string]attr.Value{
		"valid": types.BoolValue(getlinkbranditem.DNSDetails.DCNAME.Valid),
		"types": types.StringValue(getlinkbranditem.DNSDetails.DCNAME.Type),
		"host":  types.StringValue(getlinkbranditem.DNSDetails.DCNAME.Host),
		"data":  types.StringValue(getlinkbranditem.DNSDetails.DCNAME.Data),
	}

	receivedcnameType := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	receivedcnameMapVlaue, diags := types.ObjectValue(receivedcnameType, recievedcnameValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	receivedocnameValue := map[string]attr.Value{
		"valid": types.BoolValue(getlinkbranditem.DNSDetails.OCNAME.Valid),
		"types": types.StringValue(getlinkbranditem.DNSDetails.OCNAME.Type),
		"host":  types.StringValue(getlinkbranditem.DNSDetails.OCNAME.Host),
		"data":  types.StringValue(getlinkbranditem.DNSDetails.OCNAME.Data),
	}

	receivedocnameType := map[string]attr.Type{
		"valid": types.BoolType,
		"types": types.StringType,
		"host":  types.StringType,
		"data":  types.StringType,
	}

	receivedocnameMapVlaue, diags := types.ObjectValue(receivedocnameType, receivedocnameValue)

	//dkimMapVlaue, diags := types.SetValue(types.StringType, elementsValue)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error reading link branding",
			fmt.Sprintf("Unable to convert dat into Objectvalue: %s", diags.Errors()),
		)
		return
	}

	datalinkbrand = LinkbrandResourceModel{
		ID:            types.Int64Value(getlinkbranditem.ID),
		UserId:        types.Int64Value(getlinkbranditem.UserId),
		Domain:        types.StringValue(getlinkbranditem.Domain),
		Subdomain:     types.StringValue(getlinkbranditem.Subdomain),
		Username:      types.StringValue(getlinkbranditem.Username),
		Defaultdomain: types.BoolValue(getlinkbranditem.Defaultdomain),
		Legacy:        types.BoolValue(getlinkbranditem.Legacy),
		Valid:         types.BoolValue(getlinkbranditem.Valid),
		DCNAME:        receivedcnameMapVlaue,
		OCNAME:        receivedocnameMapVlaue,
	}

	diags = resp.State.Set(ctx, datalinkbrand)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
