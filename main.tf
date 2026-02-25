# -----------------------------------------------------------------------------
# Terraform Snowflake Module - Seed Data Main
# -----------------------------------------------------------------------------
# This module seeds data into a Snowflake table using either a SQL script file
# or inline SQL text. Includes safety controls for environment blocking.
# -----------------------------------------------------------------------------

locals {
  # Normalize environment names for comparison (case-insensitive)
  current_env_lower  = lower(var.seed.environment)
  blocked_envs_lower = [for env in var.seed.blocked_environments : lower(env)]

  # Check if current environment is blocked
  is_blocked = contains(local.blocked_envs_lower, local.current_env_lower)

  # Determine if seeding should run
  should_seed = var.seed.enabled && !local.is_blocked

  # Read SQL from file if script_path is provided, otherwise use sql_text
  sql_from_file = var.seed.script_path != null ? (
    fileexists(var.seed.script_path) ? file(var.seed.script_path) : null
  ) : null

  # script_path wins if both are provided
  effective_sql = coalesce(local.sql_from_file, var.seed.sql_text, "")

  # Generate hash for change detection
  script_hash = local.should_seed ? sha256(local.effective_sql) : ""

  # Trigger key combines version and optionally script hash
  trigger_key = local.should_seed ? (
    var.seed.rerun_on_script_change
    ? "${var.seed.seed_version}-${local.script_hash}"
    : var.seed.seed_version
  ) : ""
}

# Execute the seed SQL using snowflake_execute
# This resource will re-run when trigger_key changes
resource "snowflake_execute" "seed" {
  count = local.should_seed && local.effective_sql != "" ? 1 : 0

  execute = local.effective_sql
  revert  = "SELECT 1" # No-op revert - seeding is typically additive

  lifecycle {
    # Force replacement when trigger_key changes
    replace_triggered_by = [null_resource.seed_trigger[0]]
  }
}

# Null resource to track trigger changes
resource "null_resource" "seed_trigger" {
  count = local.should_seed ? 1 : 0

  triggers = {
    trigger_key = local.trigger_key
    database    = var.seed.database
    schema      = var.seed.schema
    table       = var.seed.table
  }
}
