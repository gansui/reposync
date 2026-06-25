# RepoSync

this repos is forked and improved from https://github.com/cidverse/reposync,the main improvement
- 1.when the repo is an absolute url,like https://github.com/xxxx/xxxx,it will clone others public repo
- 2.add sync command,the sync = clone + pull
- 3.add --log-file dir switch，the log could be stored in file at dir. if no dir,the log file will be stored in file where the  config.yaml is saved.


> You can configure one or more platforms as source to mirror the structure from the remote server onto your local system (repositories will be moved if necessary).
> You can use complex rules to include or exclude projects based on the namespace or other properties.

## Configuration

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/cidverse/reposync/main/configschema/v1.json
servers:
  - url: https://github.com
    type: github
    auth:
      username: YourAccount
      password: <readOnlyPersonalAccessToken>
      password-file: /path/to/file
      password-command: pass show personal/github.com/read-only-pat
    mirror:
      dir: /tmp/github
      default-action: exclude
      naming-style: slug
      rules:
        - rule: group == "my-org"
          action: include
      rules:
        - rule: group == "https://github.com/xxxx/xxxx"
          action: include
```
when you configure the group to https://xxx. the reposync will clone others public repo.
The configuration is read from `~/.config/reposync/config.yaml` by default, but you can also specify a custom path by setting the `REPOSYNC_CONFIG` environment variable.

Supported platforms:

- `github`
- `gitlab`

> The `git` commands will use your local git installation, so you can use ssh keys or other authentication methods.
> The personal access tokens are only used to query the repositories you have access to and not to clone them.

## Usage

### sync
`reposync sync` will clone and pull all repositories you have access to into the target directory
### Clone

`reposync clone` will clone all repositories you have access to into the target directory, keeping the remote structure as close as possible.

**Notes:**

- tracked repositories might be moved to track the remote structure changes (`<namespace>/<projectName>`).

### Pull

`reposync pull` will pull the latest changes for all tracked repositories.

- `git fetch`
- `git pull --ff-only`

### List

`reposync list` will list all tracked repositories.

Additional options:

- `--format` (`-f`) -> specify the output format (table, json, csv, ...)
- `--columns` (`-c`) -> filter output by columns

### HouseKeeping

`reposync housekeeping` (`reposync hk`) will run the following tasks for all repositories:

- `repack` -> `git repack -a -d --write-bitmap-index`
- `prune` -> `git prune --expire now`
- `gc` -> `git gc --auto`
- `fsck` -> `git fsck --full --unreachable --strict`
- `commit-graph` -> `git commit-graph write --reachable`

You can choose to only run a subset of the tasks by specifying them as arguments. For example, `reposync hk --jobs repack prune`.

### Index

> This is work-in-progress and not yet implemented.

Before your first run, you can use `reposync index /old-project-dir`, to add your local projects to the known repositories.
When running `reposync clone`, it will then move the projects to the new location, instead of cloning them.

### Rules

Rules support the following variables to match against:

| Variable   | Example           | Description                             |
|------------|-------------------|-----------------------------------------|
| `uniqueId` | github-com/123456 | The unique id of the project            |
| `id`       | 123456            | The id of the project                   |
| `group`    | my-org            | The group / namespace the project is in |
| `name`     | my-project        | The name of the project                 |
| `path`     | my-org/my-project | The full path of the project            |
| `is_fork`  | false             | Whether the project is a fork           |

Rules follow the [Common Expression Language](https://github.com/google/cel-spec) syntax.

## License

Released under the [MIT license](./LICENSE).
