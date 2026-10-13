---
title: "Steampipe Table: azuredevops_git_branch_policy - Query Azure DevOps git branch policy using SQL"
description: "Allows users to query the branch and repository policies of an Azure DevOps project, such as required reviewers, build validation, comment resolution and linked work items."
---

# Table: azuredevops_git_branch_policy - Query Azure DevOps git branch policy using SQL

Azure DevOps branch policies protect branches by requiring pull requests to meet conditions before they merge: a minimum number of reviewers, a passing build, resolved comments or linked work items.

## Table Usage Guide

The `azuredevops_git_branch_policy` table lists every policy configured in a project, one row per policy scope. As a security engineer, use it to check that default branches require reviewed, built pull requests.

**Important Notes**
- You must specify `project` (name or ID) in the `where` clause. The `project` column echoes the value you give.
- A policy can cover several scopes, so there is one row per scope and `id` can repeat. `repository_id`, `ref_name` and `match_kind` come from that scope.
- A scope that covers every repository in the project has a null `repository_id`. A default-branch scope has `match_kind = 'DefaultBranch'` and a null `ref_name`.
- `repository_id` and `ref_name` are filtered by Postgres, not sent to Azure DevOps. Filtering on `repository_id = '...'` therefore drops project-wide policies that also apply to that repository, so use `repository_id is null or repository_id = '...'`.
- `settings` holds the policy settings exactly as Azure DevOps returns them, e.g. `minimumApproverCount`, `creatorVoteCounts`, `resetOnSourcePush`, `allowDownvotes` for reviewers, or `buildDefinitionId` and `validDuration` for build validation.
- Needs the PAT scope `vso.code` (Code: Read).

## Examples

### Enabled, blocking minimum-reviewer policies and their settings

```sql+postgres
select
  id,
  repository_id,
  ref_name,
  match_kind,
  settings ->> 'minimumApproverCount' as minimum_approvers,
  settings ->> 'creatorVoteCounts' as creator_vote_counts,
  settings ->> 'resetOnSourcePush' as reset_on_push
from
  azuredevops_git_branch_policy
where
  project = 'my-project'
  and type_display_name = 'Minimum number of reviewers'
  and is_enabled
  and is_blocking;
```

```sql+sqlite
select
  id,
  repository_id,
  ref_name,
  match_kind,
  json_extract(settings, '$.minimumApproverCount') as minimum_approvers,
  json_extract(settings, '$.creatorVoteCounts') as creator_vote_counts,
  json_extract(settings, '$.resetOnSourcePush') as reset_on_push
from
  azuredevops_git_branch_policy
where
  project = 'my-project'
  and type_display_name = 'Minimum number of reviewers'
  and is_enabled
  and is_blocking;
```

### Policies that apply to one repository, including project-wide ones

```sql+postgres
select
  id,
  type_display_name,
  ref_name,
  match_kind
from
  azuredevops_git_branch_policy
where
  project = 'my-project'
  and (repository_id is null or repository_id = '49c1d4d0-be28-4f20-9e0a-4ecfafc39257');
```

```sql+sqlite
select
  id,
  type_display_name,
  ref_name,
  match_kind
from
  azuredevops_git_branch_policy
where
  project = 'my-project'
  and (repository_id is null or repository_id = '49c1d4d0-be28-4f20-9e0a-4ecfafc39257');
```
