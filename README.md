# Terraform Snowflake Module - Seed Data

![Release](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data/actions/workflows/ci.yaml/badge.svg)&nbsp;![Snowflake](https://img.shields.io/badge/Snowflake-29B5E8?logo=snowflake&logoColor=white)&nbsp;![Commit Activity](https://img.shields.io/github/commit-activity/t/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data)&nbsp;![Last Commit](https://img.shields.io/github/last-commit/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data)&nbsp;![Release Date](https://img.shields.io/github/release-date/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data)&nbsp;![File Count](https://img.shields.io/github/directory-file-count/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data)&nbsp;![Issues](https://img.shields.io/github/issues/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data)&nbsp;![Top Language](https://img.shields.io/github/languages/top/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data)&nbsp;![Custom Endpoint](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/bsubhamay/770827eb380b5d92a1f3e7b58e0e63ae/raw/terraform-snowflake-seed-data.json?)

A Terraform module for seeding data into Snowflake tables. Supports inline SQL or external script files with environment-based safety controls and re-run management.

## Features

- Seed data using inline SQL or external `.sql` files
- Environment-based blocking (e.g., prevent seeding in production)
- Version-controlled re-runs via `seed_version`
- Optional automatic re-run on script content changes
- Built-in input validation with descriptive error messages

## Usage

### Inline SQL

```hcl
module "seed" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data"

  seed = {
    enabled     = true
    environment = "dev"
    database    = "MY_DATABASE"
    schema      = "MY_SCHEMA"
    table       = "USERS"
    sql_text    = <<-EOT
      INSERT INTO MY_DATABASE.MY_SCHEMA.USERS (id, name, email)
      VALUES
        (1, 'Alice', 'alice@example.com'),
        (2, 'Bob', 'bob@example.com');
    EOT
  }
}
```

### External Script File

```hcl
module "seed" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data"

  seed = {
    enabled                = true
    environment            = "dev"
    database               = "MY_DATABASE"
    schema                 = "MY_SCHEMA"
    table                  = "USERS"
    script_path            = "${path.module}/seed.sql"
    rerun_on_script_change = true
  }
}
```

### Using JSON Configuration

```hcl
module "seed" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data"

  seed = jsondecode(file("seed.json"))
}
```

## Examples

- [Basic](examples/basic) - Seed data using inline SQL or script file

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.3.0 |
| snowflake | >= 0.87.0 |
| null | >= 3.0.0 |

## Providers

| Name | Source | Version |
|------|--------|---------|
| snowflake | snowflakedb/snowflake | >= 0.87.0 |
| null | hashicorp/null | >= 3.0.0 |

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| seed | Seed configuration object | `object` | yes |

### seed Object Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| enabled | bool | `false` | Enable/disable seeding |
| environment | string | `"dev"` | Current environment name |
| database | string | - | Target database name (required) |
| schema | string | - | Target schema name (required) |
| table | string | - | Target table name (required) |
| blocked_environments | list(string) | `["prod", "production"]` | Environments where seeding is blocked |
| script_path | string | `null` | Path to external `.sql` file |
| sql_text | string | `null` | Inline SQL text |
| seed_version | string | `"v1"` | Version string for re-run control |
| rerun_on_script_change | bool | `false` | Re-run when script content changes |

**Note:** Either `script_path` or `sql_text` must be provided. If both are provided, `script_path` takes precedence.

## Outputs

| Name | Description |
|------|-------------|
| seed_enabled | Whether seeding is enabled in configuration |
| seed_executed | Whether the seed was actually executed |
| seed_blocked | Whether seeding was blocked due to environment |
| seed_environment | The current environment |
| seed_target | Target table (database.schema.table) |
| seed_version | The seed version used |
| seed_trigger_key | Computed trigger key for debugging |
| seed_sql_source | Source of SQL (script_path or sql_text) |

## Re-run Controls

### Manual Re-run

Bump `seed_version` to force a re-run:

```hcl
seed = {
  # ...
  seed_version = "v2"  # Changed from "v1"
}
```

### Automatic Re-run on Script Changes

Enable `rerun_on_script_change` to automatically re-run when the SQL script content changes:

```hcl
seed = {
  # ...
  script_path            = "${path.module}/seed.sql"
  rerun_on_script_change = true
}
```

## Environment Safety

By default, seeding is blocked in `prod` and `production` environments (case-insensitive). Customize the blocked list:

```hcl
seed = {
  # ...
  blocked_environments = ["prod", "production", "staging", "uat"]
}
```

## Validation

The module validates inputs and provides descriptive error messages for:

- Missing `script_path` and `sql_text` (at least one required)
- Empty `database`, `schema`, or `table` names

## Testing

The module includes Terratest-based integration tests:

```bash
cd test
go mod tidy
go test -v -timeout 30m
```

Required environment variables for testing:
- `SNOWFLAKE_ORGANIZATION_NAME` - Snowflake organization name
- `SNOWFLAKE_ACCOUNT_NAME` - Snowflake account name
- `SNOWFLAKE_USER` - Snowflake username
- `SNOWFLAKE_ROLE` - Snowflake role (e.g., "SYSADMIN")
- `SNOWFLAKE_PRIVATE_KEY` - Snowflake private key for key-pair authentication

## License

MIT License - See [LICENSE](LICENSE) for details.
