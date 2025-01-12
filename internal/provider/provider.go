// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	//"github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &AppRegistryProvider{}

//var _ provider.ProviderWithFunctions = &AppRegistryProvider{}

type AppRegistryProvider struct {
	version string
}

type AppRegistryProviderModel struct {
	Region types.String `tfsdk:"region"`
}

func (p *AppRegistryProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "servicecatalog" // name of the provider
}

func (p *AppRegistryProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with aws service catalog app registry",
		Attributes: map[string]schema.Attribute{
			// "access_key" : schema.StringAttribute{
			// 	Description : "AWS Access Key",
			// 	Optional:	true,
			// },
			// "secret_key": schema.StringAttribute{
			// 	Description: "Aws secret Key",
			// 	Optional: true,
			// },
			"region": schema.StringAttribute{
				Description: "AWS region",
				Required:    true,
			},
		},
	}
}

func (p *AppRegistryProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var providerConfig AppRegistryProviderModel

	diags := req.Config.Get(ctx, &providerConfig)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(
			providerConfig.Region.ValueString()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create AWS Config",
			err.Error(),
		)
		return
	}

	client := servicecatalogappregistry.NewFromConfig(cfg)
	resp.DataSourceData = client
	resp.ResourceData = client

}

func (p *AppRegistryProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAppregistryApplicationsDataSource,
	}
}

func (p *AppRegistryProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAppregistryResourceAssociation,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AppRegistryProvider{
			version: version,
		}
	}
}

// todo , use the scaffolding framework in combination with the aws go sdk v2
// to convert two of the cli commands on service catalog app registry into a data source
// and then into a resource

//#aws servicecatalog-appregistry list-applications
//aws servicecatalog-appregistry associate-resource --application application_ARN
//--resource-type type --resource name --option "APPLY_APPLICATION_TAG"

/*
aws servicecatalog-appregistry list-applications
{
    "applications": [
        {
            "id": "08eyt0oo157qjgw5x6ieigqsgw",
            "arn": "arn:aws:servicecatalog:eu-west-2:390746273208:/applications/08eyt0oo157qjgw5x6ieigqsgw",
            "name": "instance-scheduler-on-aws-eu-west-2-390746273208-instance-scheduler",
            "description": "Service Catalog application to track and manage all your resources for the solution instance-scheduler-on-aws",
            "creationTime": "2025-01-11T16:24:28.680000+00:00",
            "lastUpdateTime": "2025-01-11T16:24:28.680000+00:00"
        }
    ]
}
*/

// func (p *ScaffoldingProvider) Resources(ctx context.Context) []func() resource.Resource {
// 	return []func() resource.Resource{
// 		NewExampleResource,
// 	}
// }

// func (p *ScaffoldingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
// 	return []func() datasource.DataSource{
// 		NewExampleDataSource,
// 	}
// }

// func (p *ScaffoldingProvider) Functions(ctx context.Context) []func() function.Function {
// 	return []func() function.Function{
// 		NewExampleFunction,
// 	}
// }

// func New(version string) func() provider.Provider {
// 	return func() provider.Provider {
// 		return &ScaffoldingProvider{
// 			version: version,
// 		}
// 	}
// }
