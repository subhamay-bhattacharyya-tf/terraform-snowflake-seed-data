# -----------------------------------------------------------------------------
# Terraform Snowflake Module - Seed Data Basic Example
# -----------------------------------------------------------------------------
# This example demonstrates how to seed data into a Snowflake table.
# -----------------------------------------------------------------------------

module "seed" {
  source = "../.."

  seed = var.seed
}
