package qovery

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/qovery/qovery-client-go"
)

type applicationDataSourceType struct{}

func (t applicationDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Use this data source to retrieve information about an existing application.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Id of the application.",
				Type:        types.StringType,
				Required:    true,
			},
			"environment_id": {
				Description: "Id of the environment.",
				Type:        types.StringType,
				Computed:    true,
			},
			"name": {
				Description: "Name of the application.",
				Type:        types.StringType,
				Computed:    true,
			},
			"description": {
				Description: "Description of the application.",
				Type:        types.StringType,
				Computed:    true,
			},
			"git_repository": {
				Description: "Git repository of the application.",
				Computed:    true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"url": {
						Description: "URL of the git repository.",
						Type:        types.StringType,
						Computed:    true,
					},
					"branch": {
						Description: "Branch of the git repository.",
						Type:        types.StringType,
						Computed:    true,
					},
					"root_path": {
						Description: "Root path of the application.",
						Type:        types.StringType,
						Computed:    true,
					},
				}),
			},
			"build_mode": {
				Description: "Build Mode of the application.",
				Type:        types.StringType,
				Computed:    true,
			},
			"dockerfile_path": {
				Description: "Dockerfile Path of the application.",
				Type:        types.StringType,
				Computed:    true,
			},
			"buildpack_language": {
				Description: "Buildpack Language framework.",
				Type:        types.StringType,
				Computed:    true,
			},
			"cpu": {
				Description: "CPU of the application in millicores (m) [1000m = 1 CPU].",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"memory": {
				Description: "RAM of the application in MB [1024MB = 1GB].",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"min_running_instances": {
				Description: "Minimum number of instances running for the application.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"max_running_instances": {
				Description: "Maximum number of instances running for the application.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"auto_preview": {
				Description: "Specify if the environment preview option is activated or not for this application.",
				Type:        types.BoolType,
				Computed:    true,
			},
			"storage": {
				Description: "List of storages linked to this application.",
				Computed:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Description: "Id of the storage.",
						Type:        types.StringType,
						Computed:    true,
					},
					"type": {
						Description: "Type of the storage for the application.",
						Type:        types.StringType,
						Computed:    true,
					},
					"size": {
						Description: "Size of the storage for the application in GB [1024MB = 1GB].",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"mount_point": {
						Description: "Mount point of the storage for the application.",
						Type:        types.StringType,
						Computed:    true,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
			"ports": {
				Description: "List of storages linked to this application.",
				Computed:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Description: "Id of the port.",
						Type:        types.StringType,
						Computed:    true,
					},
					"name": {
						Description: "Name of the port.",
						Type:        types.StringType,
						Computed:    true,
					},
					"internal_port": {
						Description: "Internal port of the application.",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"external_port": {
						Description: "External port of the application.",
						Type:        types.Int64Type,
						Computed:    true,
					},
					"publicly_accessible": {
						Description: "Specify if the port is exposed to the world or not for this application.",
						Type:        types.BoolType,
						Computed:    true,
					},
					"protocol": {
						Description: "Protocol used for the port of the application.",
						Type:        types.StringType,
						Computed:    true,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

func (t applicationDataSourceType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return applicationDataSource{
		client: p.(*provider).GetClient(),
	}, nil
}

type applicationDataSource struct {
	client *qovery.APIClient
}

// Read qovery application data source
func (d applicationDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	// Get current state
	var data Application
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get application from API
	application, res, err := d.client.ApplicationMainCallsApi.
		GetApplication(ctx, data.Id.Value).
		Execute()
	if err != nil || res.StatusCode >= 400 {
		apiErr := applicationReadAPIError(data.Id.Value, res, err)
		resp.Diagnostics.AddError(apiErr.Summary(), apiErr.Detail())
		return
	}

	state := convertResponseToApplication(application)
	tflog.Trace(ctx, "read application", "application_id", state.Id.Value)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}