# -----------------------------------------------------------------------------
# Terraform Snowflake Module - Seed Data Basic Example Outputs
# -----------------------------------------------------------------------------
# This file defines the output values for the basic example.
# -----------------------------------------------------------------------------

output "seed_enabled" {
  description = "Whether seeding is enabled"
  value       = module.seed.seed_enabled
}

output "seed_executed" {
  description = "Whether the seed was actually executed"
  value       = module.seed.seed_executed
}

output "seed_blocked" {
  description = "Whether seeding was blocked due to environment"
  value       = module.seed.seed_blocked
}

output "seed_target" {
  description = "The target table for seeding"
  value       = module.seed.seed_target
}
