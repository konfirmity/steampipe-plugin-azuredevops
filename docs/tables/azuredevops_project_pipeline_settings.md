---
title: "Steampipe Table: azuredevops_project_pipeline_settings - Query Azure DevOps project pipeline settings using SQL"
description: "Allows users to query the pipeline general settings of an Azure DevOps project, such as job authorization scope, settable variables and fork protection."
---

# Table: azuredevops_project_pipeline_settings - Query Azure DevOps project pipeline settings using SQL

Each Azure DevOps project has pipeline general settings that limit what a pipeline's job token can reach and how pull requests from forks are built.

## Table Usage Guide

The `azuredevops_project_pipeline_settings` table returns one row per project. As a security engineer, use it to check that job tokens are scoped to the project and that builds of forks cannot read secrets.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. The `project` column echoes the value you give.
- A setting the API does not report is null rather than false. `settings` holds every setting exactly as Azure DevOps returns it.
- Organization-level pipeline settings have no public REST API. The values here are the project's, as the project API reports them.
- To check fork settings per pipeline, read `triggers` on `azuredevops_build_definition`: a `pullRequest` trigger carries `forks.enabled` and `forks.allowSecrets`.
- Needs the PAT scope `vso.project` (Project and Team: Read).

## Examples

### Projects whose job tokens are not limited to the project

```sql+postgres
select
  project,
  enforce_job_auth_scope,
  enforce_job_auth_scope_for_releases,
  enforce_referenced_repo_scoped_token
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project'
  and not (enforce_job_auth_scope and enforce_job_auth_scope_for_releases);
```

```sql+sqlite
select
  project,
  enforce_job_auth_scope,
  enforce_job_auth_scope_for_releases,
  enforce_referenced_repo_scoped_token
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project'
  and not (enforce_job_auth_scope and enforce_job_auth_scope_for_releases);
```

### Fork protection

```sql+postgres
select
  project,
  builds_enabled_for_forks,
  enforce_no_access_to_secrets_from_forks,
  enforce_job_auth_scope_for_forks
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project';
```

```sql+sqlite
select
  project,
  builds_enabled_for_forks,
  enforce_no_access_to_secrets_from_forks,
  enforce_job_auth_scope_for_forks
from
  azuredevops_project_pipeline_settings
where
  project = 'my-project';
```
