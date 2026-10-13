---
title: "Steampipe Table: azuredevops_variable_group - Query Azure DevOps variable group using SQL"
description: "Allows users to query the variable groups of an Azure DevOps project, counting secret variables and flagging plain variables whose names look like credentials."
---

# Table: azuredevops_variable_group - Query Azure DevOps variable group using SQL

Variable groups store values shared across pipelines. Values not marked secret are stored and shown in plain text.

## Table Usage Guide

The `azuredevops_variable_group` table lists the variable groups in a project. As a security engineer, use it to find credentials kept as plain variables instead of secrets or a key vault.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. The `project` column echoes the value you give.
- Variable values are never read into the table: only names and the secret flag are decoded.
- `plaintext_variable_names_suspected` lists variables not marked secret whose names look like credentials (password, secret, token, API key, connection string, PAT, names ending in key). It is a name heuristic and can include false positives. It is always empty for key vault groups, whose variables come from the vault.
- Needs the PAT scope `vso.variablegroups_read` (Variable Groups: Read).

## Examples

### Variable groups with credential-like plain variables

```sql+postgres
select
  name,
  secret_variable_count,
  plaintext_variable_names_suspected
from
  azuredevops_variable_group
where
  project = 'my-project'
  and jsonb_array_length(plaintext_variable_names_suspected) > 0;
```

```sql+sqlite
select
  name,
  secret_variable_count,
  plaintext_variable_names_suspected
from
  azuredevops_variable_group
where
  project = 'my-project'
  and json_array_length(plaintext_variable_names_suspected) > 0;
```
