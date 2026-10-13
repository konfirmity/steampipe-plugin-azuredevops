---
title: "Steampipe Table: azuredevops_user_entitlement - Query Azure DevOps user entitlement using SQL"
description: "Allows users to query the users of an Azure DevOps organization with their access level, status and last access date."
---

# Table: azuredevops_user_entitlement - Query Azure DevOps user entitlement using SQL

User entitlements record each user's license (access level), status and when they last accessed the organization.

## Table Usage Guide

The `azuredevops_user_entitlement` table lists every user in the organization. As a security engineer, use it to find inactive users who still hold access.

**Important Notes**
- Lists the whole organization; there are no required qualifiers.
- `last_accessed_date` is null for a user who has never accessed the organization.
- Needs the PAT scope `vso.memberentitlementmanagement` (Member Entitlement Management: Read). Without it the query fails with an error that starts `azuredevops: permission denied`.

## Examples

### Active users who have not accessed the organization for 90 days

```sql+postgres
select
  principal_name,
  access_level,
  last_accessed_date
from
  azuredevops_user_entitlement
where
  status = 'active'
  and (last_accessed_date is null or last_accessed_date < now() - interval '90 days');
```

```sql+sqlite
select
  principal_name,
  access_level,
  last_accessed_date
from
  azuredevops_user_entitlement
where
  status = 'active'
  and (last_accessed_date is null or last_accessed_date < datetime('now', '-90 days'));
```
