package azuredevops

import (
	"context"
	"net/url"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const (
	serviceEndpointsAPIVersion    = "7.1"
	pipelinePermissionsAPIVersion = "7.1-preview.1"
)

// https://learn.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints/get-service-endpoints
// https://learn.microsoft.com/en-us/rest/api/azure/devops/approvalsandchecks/pipeline-permissions/get
func tableAzureDevOpsServiceEndpointPermissions(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azuredevops_service_endpoint",
		Description: "Service connections in an Azure DevOps project and which pipelines are authorized to use them.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "project", Require: plugin.Required},
			},
			Hydrate: listServiceEndpointsWithPermissions,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "project", Description: "Name or ID of the project, as given in the qualifier.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("project")},
			{Name: "id", Description: "ID of the service connection.", Type: proto.ColumnType_STRING},
			{Name: "name", Description: "Name of the service connection.", Type: proto.ColumnType_STRING},
			{Name: "type", Description: "Type of the service connection, e.g. azurerm, github, dockerregistry.", Type: proto.ColumnType_STRING},
			{Name: "owner", Description: "Owner of the service connection: Library or AgentCloud.", Type: proto.ColumnType_STRING},
			{Name: "url", Description: "URL the service connection targets.", Type: proto.ColumnType_STRING},
			{Name: "description", Description: "Description of the service connection.", Type: proto.ColumnType_STRING},
			{Name: "is_shared", Description: "True if the service connection is shared with other projects.", Type: proto.ColumnType_BOOL},
			{Name: "is_ready", Description: "True if the service connection is ready to use.", Type: proto.ColumnType_BOOL},
			{Name: "authorization_scheme", Description: "Authentication scheme, e.g. ServicePrincipal, WorkloadIdentityFederation, UsernamePassword. Parameters are not exposed.", Type: proto.ColumnType_STRING},
			{Name: "all_pipelines_authorized", Description: "True if every pipeline in the project may use the service connection without approval.", Type: proto.ColumnType_BOOL, Hydrate: getServiceEndpointPipelinePermissions, Transform: transform.FromField("AllPipelinesAuthorized")},
			{Name: "authorized_pipeline_ids", Description: "IDs of the pipeline definitions individually authorized to use the service connection.", Type: proto.ColumnType_JSON, Hydrate: getServiceEndpointPipelinePermissions, Transform: transform.FromField("AuthorizedPipelineIds")},
			{Name: "title", Description: "Title of the resource.", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name")},
		}),
	}
}

type serviceEndpointResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Owner         string `json:"owner"`
	Url           string `json:"url"`
	Description   string `json:"description"`
	IsShared      bool   `json:"isShared"`
	IsReady       bool   `json:"isReady"`
	Authorization struct {
		Scheme string `json:"scheme"`
	} `json:"authorization"`
}

type ServiceEndpointRow struct {
	Id                  string
	Name                string
	Type                string
	Owner               string
	Url                 string
	Description         string
	IsShared            bool
	IsReady             bool
	AuthorizationScheme string
}

type pipelinePermissions struct {
	AllPipelinesAuthorized bool
	AuthorizedPipelineIds  []int
}

func listServiceEndpointsWithPermissions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	project := d.EqualsQualString("project")
	if project == "" {
		return nil, nil
	}
	client, err := restClientFor(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_service_endpoint.listServiceEndpointsWithPermissions", "connection_error", err)
		return nil, err
	}
	err = readServiceEndpoints(ctx, client, project, func(row ServiceEndpointRow) bool {
		d.StreamListItem(ctx, row)
		return d.RowsRemaining(ctx) != 0
	})
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_service_endpoint.listServiceEndpointsWithPermissions", "api_error", err)
		return nil, err
	}
	return nil, nil
}

func readServiceEndpoints(ctx context.Context, client *restClient, project string, emit func(ServiceEndpointRow) bool) error {
	params := apiVersion(serviceEndpointsAPIVersion)
	params.Set("includeFailed", "true")
	return getList(ctx, client, projectURL(client.baseURL, project, "serviceendpoint/endpoints"), params, scopeServiceEndpoint, func(e serviceEndpointResponse) bool {
		return emit(ServiceEndpointRow{
			Id: e.Id, Name: e.Name, Type: e.Type, Owner: e.Owner, Url: e.Url, Description: e.Description,
			IsShared: e.IsShared, IsReady: e.IsReady, AuthorizationScheme: e.Authorization.Scheme,
		})
	})
}

func getServiceEndpointPipelinePermissions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	row := h.Item.(ServiceEndpointRow)
	client, err := restClientFor(ctx, d)
	if err != nil {
		return nil, err
	}
	perms, err := readPipelinePermissions(ctx, client, d.EqualsQualString("project"), "endpoint", row.Id)
	if err != nil {
		plugin.Logger(ctx).Error("azuredevops_service_endpoint.getServiceEndpointPipelinePermissions", "api_error", err)
		return nil, err
	}
	return perms, nil
}

func readPipelinePermissions(ctx context.Context, client *restClient, project, resourceType, resourceID string) (*pipelinePermissions, error) {
	endpoint := projectURL(client.baseURL, project, "pipelines/pipelinepermissions/"+url.PathEscape(resourceType)+"/"+url.PathEscape(resourceID)) +
		"?" + apiVersion(pipelinePermissionsAPIVersion).Encode()
	var resp struct {
		AllPipelines *struct {
			Authorized bool `json:"authorized"`
		} `json:"allPipelines"`
		Pipelines []struct {
			Id         int  `json:"id"`
			Authorized bool `json:"authorized"`
		} `json:"pipelines"`
	}
	if err := client.getJSON(ctx, endpoint, scopeBuild, &resp); err != nil {
		return nil, err
	}
	perms := &pipelinePermissions{AuthorizedPipelineIds: []int{}}
	perms.AllPipelinesAuthorized = resp.AllPipelines != nil && resp.AllPipelines.Authorized
	for _, p := range resp.Pipelines {
		if p.Authorized {
			perms.AuthorizedPipelineIds = append(perms.AuthorizedPipelineIds, p.Id)
		}
	}
	return perms, nil
}
