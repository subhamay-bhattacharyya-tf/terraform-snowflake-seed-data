# Basic Seed Example

This example demonstrates how to seed data into a Snowflake table using inline SQL.

## Usage

```hcl
module "seed" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data"

  seed = {
    enabled     = true
    environment = "dev"
    database    = "MY_DATABASE"
    schema      = "MY_SCHEMA"
    table       = "MY_TABLE"
    sql_text    = <<-EOT
      INSERT INTO MY_DATABASE.MY_SCHEMA.MY_TABLE (id, name)
      VALUES (1, 'Test'), (2, 'Sample');
    EOT
  }
}
```

## Using a JSON file

```hcl
module "seed" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-seed-data"

  seed = jsondecode(file("seed.json"))
}
```

Example `seed.json`:
```json
{
  "enabled": true,
  "environment": "dev",
  "database": "MY_DATABASE",
  "schema": "MY_SCHEMA",
  "table": "MY_TABLE",
  "sql_text": "INSERT INTO MY_DATABASE.MY_SCHEMA.MY_TABLE (id, name) VALUES (1, 'Test');"
}
```
