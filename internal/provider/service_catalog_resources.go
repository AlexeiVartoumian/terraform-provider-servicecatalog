package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry"
	awstypes "github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//appregistry resource association defines the resource implementation

type AppRegistryResourceAssociation struct {
	client *servicecatalogappregistry.Client
}

// AppregistryResource Assocation model describes the resource data model
type AppRegistryResourceAssociationModel struct {
	ApplicationArn types.String `tfsdk:"application_arn"`
	ResourceType   types.String `tfsdk:"resource_type"`
	ResourceName   types.String `tfsdk:"resource_name"`
	ResourceArn    types.String `tfsdk:"resource_arn"`
	Options        types.List   `tfsdk:"options"`
	Id             types.String `tfsdk:"id"`
}

func (r *AppRegistryResourceAssociation) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_association"
}

func (r *AppRegistryResourceAssociation) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description: "Associates a resource with an AppRegistry application",
		Attributes: map[string]schema.Attribute{
			"application_arn": schema.StringAttribute{
				Description: "ARN of the Appregistry application",
				Required:    true,
			},
			"resource_type": schema.StringAttribute{
				Description: "Name of the reosurce to associate",
				Required:    true,
			},
			"resource_name": schema.StringAttribute{
				Description: "Name of the resource to associate",
				Required:    true,
			},
			"resource_arn": schema.StringAttribute{
				Description: "ARN of the associated resource",
				Computed:    true,
			},
			"options": schema.ListAttribute{
				Description: "Options for the association eg apply applcaition tag",
				ElementType: types.StringType,
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "Idenitfier of the resource association",
				Computed:    true,
			},
		},
	}
}

func (r *AppRegistryResourceAssociation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AppRegistryResourceAssociationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	applicationArn := plan.ApplicationArn.ValueString()
	resourceName := plan.ResourceName.ValueString()
	resourceType := plan.ResourceType.ValueString()
	//Associate resource
	var options []awstypes.AssociationOption

	if !plan.Options.IsNull() {
		var optionsList []string
		diags = plan.Options.ElementsAs(ctx, &optionsList, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, opt := range optionsList {
			options = append(options, awstypes.AssociationOption(opt))
		}
	} else {
		options = []awstypes.AssociationOption{awstypes.AssociationOption("APPLY_APPLICATION_TAG")}

	}
	input := &servicecatalogappregistry.AssociateResourceInput{
		Application:  &applicationArn,
		Resource:     &resourceName,
		ResourceType: awstypes.ResourceType(resourceType), //"AWS::CloudFormation::Stack", // Let's hardcode for now to test
		Options:      options,
	}

	output, err := r.client.AssociateResource(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Associating resource",
			fmt.Sprintf("Could not associate resource with AppRegistry: %s", err),
		)
		return
	}

	plan.ResourceArn = types.StringValue(*output.ResourceArn)
	plan.Id = types.StringValue(fmt.Sprintf("%s:%s", plan.ApplicationArn.ValueString(), *output.ResourceArn))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

}

func (r *AppRegistryResourceAssociation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AppRegistryResourceAssociationModel
	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	applicationArn := state.ApplicationArn.ValueString()
	// list associated resource to verigy assoication still exits
	input := &servicecatalogappregistry.ListAssociatedResourcesInput{
		Application: &applicationArn,
	}

	output, err := r.client.ListAssociatedResources(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource Assocation",
			fmt.Sprintf("Could not read Appregistry resource association: %s", err),
		)
		return
	}

	//check if resour is still associated
	resourceFound := false
	for _, res := range output.Resources {
		if *res.Arn == state.ResourceArn.ValueString() {
			resourceFound = true
			break
		}
	}

	if !resourceFound {
		resp.State.RemoveResource(ctx)
		return
	}

	//set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *AppRegistryResourceAssociation) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*servicecatalogappregistry.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Expected *servicatloagappregistry.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}
func (r *AppRegistryResourceAssociation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not supported for resource associations
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"AppRegistry resource associations cannot be updated. Delete and recreate the association instead.",
	)
}

func (r *AppRegistryResourceAssociation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state AppRegistryResourceAssociationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	applicationArn := state.ApplicationArn.ValueString()
	resourceName := state.ResourceName.ValueString()
	resourceType := state.ResourceType.ValueString()

	input := &servicecatalogappregistry.DisassociateResourceInput{
		Application:  &applicationArn,
		Resource:     &resourceName,
		ResourceType: awstypes.ResourceType(resourceType),
	}
	_, err := r.client.DisassociateResource(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Disassociating resource",
			fmt.Sprintf("could not disassociate resource from Appregistry: %s", err),
		)
		return
	}
}

func NewAppregistryResourceAssociation() resource.Resource {
	return &AppRegistryResourceAssociation{}
}
