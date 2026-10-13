---
title: "Steampipe Table: azuredevops_advanced_security_alert - Query Azure DevOps advanced security alert using SQL"
description: "Allows users to query the Advanced Security alerts (dependency, secret and code scanning) on the repositories of an Azure DevOps project."
---

# Table: azuredevops_advanced_security_alert - Query Azure DevOps advanced security alert using SQL

GitHub Advanced Security for Azure DevOps raises alerts for vulnerable dependencies, committed secrets and code weaknesses.

## Table Usage Guide

The `azuredevops_advanced_security_alert` table lists alerts per repository. As a security engineer, use it to count open critical alerts.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. `repository_id`, `state` and `alert_type` are optional and are sent to Azure DevOps as filters.
- Azure DevOps returns alerts on the default branch only, and high-confidence secret alerts only.
- A repository without Advanced Security contributes no rows. Join `azuredevops_advanced_security_enablement` to tell "no alerts" from "not scanned".
- Secret material (`truncatedSecret`, validation fingerprints) is never exposed.
- Needs the PAT scopes `vso.code` (Code: Read) and `vso.advsec` (Advanced Security: Read).

## Examples

### Open alerts by repository and severity

```sql+postgres
select
  repository_id,
  alert_type,
  severity,
  count(*)
from
  azuredevops_advanced_security_alert
where
  project = 'my-project'
  and state = 'active'
group by
  repository_id,
  alert_type,
  severity;
```

```sql+sqlite
select
  repository_id,
  alert_type,
  severity,
  count(*)
from
  azuredevops_advanced_security_alert
where
  project = 'my-project'
  and state = 'active'
group by
  repository_id,
  alert_type,
  severity;
```
