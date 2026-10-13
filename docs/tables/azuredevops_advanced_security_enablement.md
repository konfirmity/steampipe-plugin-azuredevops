---
title: "Steampipe Table: azuredevops_advanced_security_enablement - Query Azure DevOps advanced security enablement using SQL"
description: "Allows users to query whether GitHub Advanced Security for Azure DevOps (Secret Protection and Code Security) is enabled on each repository of a project."
---

# Table: azuredevops_advanced_security_enablement - Query Azure DevOps advanced security enablement using SQL

GitHub Advanced Security for Azure DevOps adds secret scanning with push protection (Secret Protection) and code and dependency scanning (Code Security) to repositories.

## Table Usage Guide

The `azuredevops_advanced_security_enablement` table returns one row per repository. As a security engineer, use it to find repositories without secret scanning or push protection.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. `repository_id` is optional; leave it out to read every repository in the project.
- An organization without Advanced Security reports every flag false rather than an error.
- Disabled repositories are not read and report false, with `is_disabled` true.
- Needs the PAT scopes `vso.code` (Code: Read) to list repositories and `vso.advsec` (Advanced Security: Read). Without `vso.advsec` the query fails with an error that starts `azuredevops: permission denied`.

## Examples

### Repositories without push protection

```sql+postgres
select
  repository_name,
  secret_protection_enabled,
  block_pushes
from
  azuredevops_advanced_security_enablement
where
  project = 'my-project'
  and not is_disabled
  and not (secret_protection_enabled and block_pushes);
```

```sql+sqlite
select
  repository_name,
  secret_protection_enabled,
  block_pushes
from
  azuredevops_advanced_security_enablement
where
  project = 'my-project'
  and not is_disabled
  and not (secret_protection_enabled and block_pushes);
```
