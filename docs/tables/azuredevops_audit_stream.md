---
title: "Steampipe Table: azuredevops_audit_stream - Query Azure DevOps audit stream using SQL"
description: "Allows users to query the audit streams that forward an Azure DevOps organization's audit log to a SIEM or log store."
---

# Table: azuredevops_audit_stream - Query Azure DevOps audit stream using SQL

Audit streams send Azure DevOps audit events to Splunk, Azure Monitor Logs or Azure Event Grid. A stream can only exist while auditing is on.

## Table Usage Guide

The `azuredevops_audit_stream` table lists the organization's audit streams. As a security engineer, use it to check that audit events leave Azure DevOps for retention and alerting.

**Important Notes**
- Lists the whole organization; there are no required qualifiers.
- The stream's consumer inputs, which hold the destination's credentials, are never exposed.
- Needs the PAT scope `vso.auditlog` (Audit Log: Read), and the token's owner needs the View audit log permission. Without them the query fails with an error that starts `azuredevops: permission denied`.

## Examples

### Enabled audit streams

```sql+postgres
select
  display_name,
  consumer_type,
  updated_time
from
  azuredevops_audit_stream
where
  status = 'enabled';
```

```sql+sqlite
select
  display_name,
  consumer_type,
  updated_time
from
  azuredevops_audit_stream
where
  status = 'enabled';
```
