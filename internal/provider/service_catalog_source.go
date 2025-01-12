package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AppRegistryApplicationsDataSource struct {
	client *servicecatalogappregistry.Client
}

type AppRegistryApplicationsDataSourceModel struct {
	Applications []ApplicationModel `tfsdk:"applications"`
	Id           types.String       `tfsdk:"id"`
}

type ApplicationModel struct {
	ID          types.String `tfsdk:"id"`
	Arn         types.String `tfsdk:"arn"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (d *AppRegistryApplicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_applications"
}

func (d *AppRegistryApplicationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Lists Appregistry applications",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "identifier for this data source",
				Computed:    true,
			},
			"applications": schema.ListNestedAttribute{
				Description: "list of Appregistryapplications",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the application",
							Computed:    true,
						},
						"arn": schema.StringAttribute{
							Description: "Name of the application",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "name of the application",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "name of the application",
							Computed:    true,
						},
					},
				},
			},
		},
	}

}

// read refreshes the terraform state with latest data
func (d *AppRegistryApplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state AppRegistryApplicationsDataSourceModel

	//list applications from appregistry
	output, err := d.client.ListApplications(ctx, &servicecatalogappregistry.ListApplicationsInput{})

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list appregistry Applications",
			err.Error(),
		)
		return
	}

	//now we map the response to the model
	for _, app := range output.Applications {
		application := ApplicationModel{
			ID:          types.StringValue(*app.Id),
			Arn:         types.StringValue(*app.Arn),
			Name:        types.StringValue(*app.Name),
			Description: types.StringValue(*app.Description),
		}
		state.Applications = append(state.Applications, application)
	}
	state.Id = types.StringValue("appregistry-applications")

	//set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func NewAppregistryApplicationsDataSource() datasource.DataSource {
	return &AppRegistryApplicationsDataSource{}
}

func (d *AppRegistryApplicationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*servicecatalogappregistry.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected data source confiure type",
			fmt.Sprintf("Expected *servicecatalogappregistry.Client, got %T", req.ProviderData),
		)
		return
	}
	d.client = client
}
