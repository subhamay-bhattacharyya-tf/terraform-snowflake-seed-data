# -----------------------------------------------------------------------------
# Terraform Snowflake Module - Seed Data Basic Example Variables
# -----------------------------------------------------------------------------
# This file defines the input variables for the basic example.
# -----------------------------------------------------------------------------

variable "seed" {
  description = "Seed configuration object"
  type = object({
    enabled                = optional(bool, false)
    environment            = optional(string, "dev")
    database               = string
    schema                 = string
    table                  = string
    blocked_environments   = optional(list(string), ["prod", "production"])
    script_path            = optional(string)
    sql_text               = optional(string)
    seed_version           = optional(string, "v1")
    rerun_on_script_change = optional(bool, false)
  })
}

# Snowflake authentication variables
variable "snowflake_organization_name" {
  description = "Snowflake organization name"
  type        = string
  default     = null
}

variable "snowflake_account_name" {
  description = "Snowflake account name"
  type        = string
  default     = null
}

variable "snowflake_user" {
  description = "Snowflake username"
  type        = string
  default     = null
}

variable "snowflake_role" {
  description = "Snowflake role"
  type        = string
  default     = null
}

variable "snowflake_private_key" {
  description = "Snowflake private key for key-pair authentication"
  type        = string
  sensitive   = true
  default     = null
}
