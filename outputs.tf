# -----------------------------------------------------------------------------
# Terraform Snowflake Module - Seed Data Outputs
# -----------------------------------------------------------------------------
# This file defines the output values for the seed data module.
# -----------------------------------------------------------------------------

output "seed_enabled" {
  description = "Whether seeding is enabled in configuration."
  value       = var.seed.enabled
}

output "seed_executed" {
  description = "Whether the seed was actually executed (enabled and not blocked)."
  value       = local.should_seed
}

output "seed_blocked" {
  description = "Whether seeding was blocked due to environment restrictions."
  value       = local.is_blocked
}

output "seed_environment" {
  description = "The current environment for seeding."
  value       = var.seed.environment
}

output "seed_target" {
  description = "The target table for seeding (database.schema.table)."
  value       = "${var.seed.database}.${var.seed.schema}.${var.seed.table}"
}

output "seed_version" {
  description = "The seed version used for re-run control."
  value       = var.seed.seed_version
}

output "seed_trigger_key" {
  description = "The computed trigger key (for debugging re-run behavior)."
  value       = local.trigger_key
}

output "seed_sql_source" {
  description = "The source of the SQL (script_path or sql_text)."
  value       = var.seed.script_path != null && local.sql_from_file != null ? "script_path" : "sql_text"
}
