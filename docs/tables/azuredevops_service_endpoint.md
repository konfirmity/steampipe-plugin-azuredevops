---
title: "Steampipe Table: azuredevops_service_endpoint - Query Azure DevOps service endpoint using SQL"
description: "Allows users to query the service connections of an Azure DevOps project and which pipelines may use them."
---

# Table: azuredevops_service_endpoint - Query Azure DevOps service endpoint using SQL

Service connections hold the credentials pipelines use to reach Azure, GitHub, registries and other services. A connection authorized for all pipelines can be used by any pipeline in the project without approval.

## Table Usage Guide

The `azuredevops_service_endpoint` table lists the service connections in a project with their pipeline authorization. As a security engineer, use it to find connections open to every pipeline.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. The `project` column echoes the value you give.
- `all_pipelines_authorized` and `authorized_pipeline_ids` come from the Pipeline Permissions API, one extra request per connection, made only when you select them.
- Authorization parameters are never exposed; `authorization_scheme` gives the scheme only.
- The older `azuredevops_serviceendpoint` table lists the same connections through the azure-devops-go-api client, without pipeline permissions.
- Needs the PAT scopes `vso.serviceendpoint` (Service Connections: Read) and, for the permission columns, `vso.build` (Build: Read).

## Examples

### Service connections any pipeline can use

```sql+postgres
select
  name,
  type,
  is_shared
from
  azuredevops_service_endpoint
where
  project = 'my-project'
  and all_pipelines_authorized;
```

```sql+sqlite
select
  name,
  type,
  is_shared
from
  azuredevops_service_endpoint
where
  project = 'my-project'
  and all_pipelines_authorized;
```
