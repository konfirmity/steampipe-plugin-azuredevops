---
title: "Steampipe Table: azuredevops_git_item - Query Azure DevOps Git repository files using SQL"
description: "Allows users to list the files and folders of an Azure DevOps Git repository at a branch, tag or commit, and read the content of a single file."
---

# Table: azuredevops_git_item - Query Azure DevOps Git repository files using SQL

The Azure DevOps Git Items API returns a repository's tree at any version in a single request, and the content of individual files.

## Table Usage Guide

The `azuredevops_git_item` table lists the files in a repository and reads individual files without cloning it. As a security engineer, use it to find dependency manifests and lockfiles across repositories, or to check that a required file (a `CODEOWNERS`, a pipeline definition) exists and says what it should.

**Important Notes**
- You must specify `project` and `repository_id` in the `where` or `join` clause to query this table. Both accept a name or an ID.
- Give `path` to read one file and its `content`, or leave it out to list the whole subtree under `scope_path` (the repository root by default) in one request. `path` must name a file, and the two cannot be combined.
- `version` takes a branch, tag or commit and defaults to the repository's default branch. `version_type` (`branch`, `tag` or `commit`) defaults to `commit` for a 40-character hash and `branch` otherwise.
- Paths start with `/`. The listing does not report file sizes; `size` is set when a text file is read with `path`.
- `content` and `size` are null for binary files.
- Requests honour Azure DevOps rate limiting: a throttled request is retried after the `Retry-After` it names, and a response that reports the limit is near delays the next request to the organization.

## Examples

### List every file in a repository

```sql+postgres
select
  path,
  commit_id
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and not is_folder;
```

```sql+sqlite
select
  path,
  commit_id
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and not is_folder;
```

### Find dependency lockfiles on a branch

```sql+postgres
select
  path
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and version = 'develop'
  and git_object_type = 'blob'
  and (path like '%/package-lock.json' or path like '%/packages.lock.json' or path like '%/go.mod');
```

```sql+sqlite
select
  path
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and version = 'develop'
  and git_object_type = 'blob'
  and (path like '%/package-lock.json' or path like '%/packages.lock.json' or path like '%/go.mod');
```

### Read one file at a commit

```sql+postgres
select
  size,
  content
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and version = '23d0bc5b128a10056dc68afece360d8a0fabb014'
  and path = '/package.json';
```

```sql+sqlite
select
  size,
  content
from
  azuredevops_git_item
where
  project = 'my-project'
  and repository_id = 'my-repo'
  and version = '23d0bc5b128a10056dc68afece360d8a0fabb014'
  and path = '/package.json';
```
