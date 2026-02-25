# -----------------------------------------------------------------------------
# Terraform Snowflake Module - Seed Data Variables
# -----------------------------------------------------------------------------
# This file defines the input variables for the seed data module.
# -----------------------------------------------------------------------------

variable "seed" {
  description = <<EOT
Seed configuration object. Designed to be passed via:
- jsondecode(file("seed.json"))
- or as a Terraform object directly.

Notes:
- Provide either script_path OR sql_text.
- If both provided, script_path wins.
EOT

  type = object({
    enabled     = optional(bool, false)
    environment = optional(string, "devl")

    # Target table information
    database = string
    schema   = string
    table    = string

    # Safety: block seeding in these envs (case-insensitive match)
    blocked_environments = optional(list(string), ["prod", "production"])

    # One of these two must be provided:
    script_path = optional(string) # path to .sql file
    sql_text    = optional(string) # inline SQL

    # Re-run controls
    seed_version           = optional(string, "v1") # bump intentionally to re-run
    rerun_on_script_change = optional(bool, false)  # re-run when script file changes
  })

  validation {
    condition = (
      try(length(var.seed.script_path) > 0, false) ||
      try(length(var.seed.sql_text) > 0, false)
    )
    error_message = "seed.script_path or seed.sql_text must be provided (non-empty)."
  }

  validation {
    condition     = length(var.seed.database) > 0
    error_message = "seed.database must be provided (non-empty)."
  }

  validation {
    condition     = length(var.seed.schema) > 0
    error_message = "seed.schema must be provided (non-empty)."
  }

  validation {
    condition     = length(var.seed.table) > 0
    error_message = "seed.table must be provided (non-empty)."
  }
}
