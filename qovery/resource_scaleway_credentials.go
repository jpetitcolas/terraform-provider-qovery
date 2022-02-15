package qovery

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/qovery/qovery-client-go"

	"terraform-provider-qovery/qovery/apierror"
)

const scalewayCredentialsAPIResource = "scaleway credentials"

type scalewayCredentialsResourceData struct {
	Id                types.String `tfsdk:"id"`
	OrganizationId    types.String `tfsdk:"organization_id"`
	Name              types.String `tfsdk:"name"`
	ScalewayAccessKey types.String `tfsdk:"scaleway_access_key"`
	ScalewaySecretKey types.String `tfsdk:"scaleway_secret_key"`
	ScalewayProjectId types.String `tfsdk:"scaleway_project_id"`
}

type scalewayCredentialsResourceType struct{}

func (r scalewayCredentialsResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Provides a Qovery SCALEWAY credentials resource. This can be used to create and manage Qovery SCALEWAY credentials.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Id of the SCALEWAY credentials.",
				Type:        types.StringType,
				Computed:    true,
			},
			"organization_id": {
				Description: "Id of the organization.",
				Type:        types.StringType,
				Required:    true,
			},
			"name": {
				Description: "Name of the scaleway credentials.",
				Type:        types.StringType,
				Required:    true,
			},
			"scaleway_access_key": {
				Description: "Your SCALEWAY access key id.",
				Type:        types.StringType,
				Required:    true,
				Sensitive:   true,
			},
			"scaleway_secret_key": {
				Description: "Your SCALEWAY secret key.",
				Type:        types.StringType,
				Required:    true,
				Sensitive:   true,
			},
			"scaleway_project_id": {
				Description: "Your SCALEWAY project ID.",
				Type:        types.StringType,
				Required:    true,
				Sensitive:   true,
			},
		},
	}, nil
}

func (r scalewayCredentialsResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return scalewayCredentialsResource{
		client: p.(*provider).GetClient(),
	}, nil
}

type scalewayCredentialsResource struct {
	client *qovery.APIClient
}

// Create qovery scaleway credentials resource
func (r scalewayCredentialsResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// Retrieve values from plan
	var plan scalewayCredentialsResourceData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new credentials
	credentials, res, err := r.client.CloudProviderCredentialsApi.
		CreateScalewayCredentials(ctx, plan.OrganizationId.Value).
		ScalewayCredentialsRequest(qovery.ScalewayCredentialsRequest{
			Name:              plan.Name.Value,
			ScalewayAccessKey: &plan.ScalewayAccessKey.Value,
			ScalewaySecretKey: &plan.ScalewaySecretKey.Value,
			ScalewayProjectId: &plan.ScalewayProjectId.Value,
		}).
		Execute()
	if err != nil || res.StatusCode >= 400 {
		apiErr := scalewayCredentialsCreateAPIError(plan.Name.Value, res, err)
		resp.Diagnostics.AddError(apiErr.Summary(), apiErr.Detail())
		return
	}

	// Initialize state values
	state := scalewayCredentialsResourceData{
		Id:                types.String{Value: *credentials.Id},
		Name:              types.String{Value: *credentials.Name},
		OrganizationId:    plan.OrganizationId,
		ScalewayAccessKey: plan.ScalewayAccessKey,
		ScalewaySecretKey: plan.ScalewaySecretKey,
		ScalewayProjectId: plan.ScalewayProjectId,
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read qovery scaleway credentials resource
func (r scalewayCredentialsResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Get current state
	var state scalewayCredentialsResourceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get credentials from API
	credentials, res, err := r.client.CloudProviderCredentialsApi.
		ListScalewayCredentials(ctx, state.OrganizationId.Value).
		Execute()
	if err != nil || res.StatusCode >= 400 {
		apiErr := scalewayCredentialsReadAPIError(state.Id.Value, res, err)
		resp.Diagnostics.AddError(apiErr.Summary(), apiErr.Detail())
		return
	}

	var toRefresh *scalewayCredentialsResourceData
	for _, creds := range credentials.GetResults() {
		if state.Id.Value == *creds.Id {
			toRefresh = &scalewayCredentialsResourceData{
				Name: types.String{Value: *creds.Name},
			}
			break
		}
	}

	// If credential id is not in list
	// Returning Not Found error
	if toRefresh == nil {
		res.StatusCode = 404
		apiErr := scalewayCredentialsReadAPIError(state.Id.Value, res, nil)
		resp.Diagnostics.AddError(apiErr.Summary(), apiErr.Detail())
		return
	}

	// Refresh state values
	state.Name = toRefresh.Name

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update qovery scaleway credentials resource
func (r scalewayCredentialsResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan and current state
	var plan, state scalewayCredentialsResourceData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update credentials in the backend
	credentials, res, err := r.client.CloudProviderCredentialsApi.
		EditScalewayCredentials(ctx, state.OrganizationId.Value, state.Id.Value).
		ScalewayCredentialsRequest(qovery.ScalewayCredentialsRequest{
			Name:              plan.Name.Value,
			ScalewayAccessKey: &plan.ScalewayAccessKey.Value,
			ScalewaySecretKey: &plan.ScalewaySecretKey.Value,
			ScalewayProjectId: &plan.ScalewayProjectId.Value,
		}).
		Execute()
	if err != nil || res.StatusCode >= 400 {
		apiErr := scalewayCredentialsUpdateAPIError(state.Id.Value, res, err)
		resp.Diagnostics.AddError(apiErr.Summary(), apiErr.Detail())
		return
	}

	toUpdate := scalewayCredentialsResourceData{
		Name:              types.String{Value: *credentials.Name},
		ScalewayAccessKey: plan.ScalewayAccessKey,
		ScalewaySecretKey: plan.ScalewaySecretKey,
		ScalewayProjectId: plan.ScalewayProjectId,
	}

	// Update state values
	state.Name = toUpdate.Name
	state.ScalewayAccessKey = toUpdate.ScalewayAccessKey
	state.ScalewaySecretKey = toUpdate.ScalewaySecretKey
	state.ScalewayProjectId = toUpdate.ScalewayProjectId

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete qovery scaleway credentials resource
func (r scalewayCredentialsResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	// Get current state
	var state scalewayCredentialsResourceData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete credentials in the backend
	res, err := r.client.CloudProviderCredentialsApi.
		DeleteScalewayCredentials(ctx, state.OrganizationId.Value, state.Id.Value).
		Execute()
	if err != nil || res.StatusCode >= 400 {
		apiErr := scalewayCredentialsDeleteAPIError(state.Id.Value, res, err)
		resp.Diagnostics.AddError(apiErr.Summary(), apiErr.Detail())
		return
	}

	// Remove credentials from state
	resp.State.RemoveResource(ctx)
}

// ImportState imports a qovery scaleway credentials resource using its id
func (r scalewayCredentialsResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: scaleway_credentials_id,organization_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("organization_id"), idParts[1])...)
}

func scalewayCredentialsCreateAPIError(credentialsName string, res *http.Response, err error) *apierror.APIError {
	return apierror.New(scalewayCredentialsAPIResource, credentialsName, apierror.Create, res, err)
}

func scalewayCredentialsReadAPIError(credentialsID string, res *http.Response, err error) *apierror.APIError {
	return apierror.New(scalewayCredentialsAPIResource, credentialsID, apierror.Read, res, err)
}

func scalewayCredentialsUpdateAPIError(credentialsID string, res *http.Response, err error) *apierror.APIError {
	return apierror.New(scalewayCredentialsAPIResource, credentialsID, apierror.Update, res, err)
}

func scalewayCredentialsDeleteAPIError(credentialsID string, res *http.Response, err error) *apierror.APIError {
	return apierror.New(scalewayCredentialsAPIResource, credentialsID, apierror.Delete, res, err)
}