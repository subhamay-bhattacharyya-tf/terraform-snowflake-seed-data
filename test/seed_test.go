// File: test/seed_test.go
package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// TestSeedInlineSQL tests seeding data using inline SQL
func TestSeedInlineSQL(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())

	database := "TT_SEED_DB"
	schema := fmt.Sprintf("TT_SCHEMA_%s", unique)
	table := fmt.Sprintf("TT_TABLE_%s", unique)

	tfDir := "../examples/basic"

	// Pre-create the test table
	db := openSnowflake(t)
	createTestTable(t, db, database, schema, table)
	defer func() {
		dropTestSchema(t, db, database, schema)
		_ = db.Close()
	}()

	// Inline SQL to insert test data (single line to avoid Terraform multi-line string issues)
	sqlText := fmt.Sprintf("INSERT INTO %s.%s.%s (id, name) VALUES (1, 'Alice'), (2, 'Bob'), (3, 'Charlie');", database, schema, table)

	seedConfig := map[string]interface{}{
		"enabled":     true,
		"environment": "devl",
		"database":    database,
		"schema":      schema,
		"table":       table,
		"sql_text":    sqlText,
	}

	tfOptions := &terraform.Options{
		TerraformDir: tfDir,
		NoColor:      true,
		Vars: map[string]interface{}{
			"seed":                        seedConfig,
			"snowflake_organization_name": os.Getenv("SNOWFLAKE_ORGANIZATION_NAME"),
			"snowflake_account_name":      os.Getenv("SNOWFLAKE_ACCOUNT_NAME"),
			"snowflake_user":              os.Getenv("SNOWFLAKE_USER"),
			"snowflake_role":              os.Getenv("SNOWFLAKE_ROLE"),
			"snowflake_private_key":       os.Getenv("SNOWFLAKE_PRIVATE_KEY"),
		},
	}

	defer terraform.Destroy(t, tfOptions)
	terraform.InitAndApply(t, tfOptions)

	time.Sleep(retrySleep)

	// Verify data was inserted
	count := countRows(t, db, database, schema, table)
	require.Equal(t, 3, count, "Expected 3 rows in table after seeding")

	// Verify outputs
	seedExecuted := terraform.Output(t, tfOptions, "seed_executed")
	require.Equal(t, "true", seedExecuted)

	seedBlocked := terraform.Output(t, tfOptions, "seed_blocked")
	require.Equal(t, "false", seedBlocked)
}

// TestSeedBlockedEnvironment tests that seeding is blocked in production environments
func TestSeedBlockedEnvironment(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())

	database := "TT_SEED_DB"
	schema := fmt.Sprintf("TT_SCHEMA_%s", unique)
	table := fmt.Sprintf("TT_TABLE_%s", unique)

	tfDir := "../examples/basic"

	// Pre-create the test table
	db := openSnowflake(t)
	createTestTable(t, db, database, schema, table)
	defer func() {
		dropTestSchema(t, db, database, schema)
		_ = db.Close()
	}()

	// Inline SQL to insert test data (single line)
	sqlText := fmt.Sprintf("INSERT INTO %s.%s.%s (id, name) VALUES (1, 'Should Not Appear');", database, schema, table)

	// Set environment to "prod" which should be blocked
	seedConfig := map[string]interface{}{
		"enabled":     true,
		"environment": "prod",
		"database":    database,
		"schema":      schema,
		"table":       table,
		"sql_text":    sqlText,
	}

	tfOptions := &terraform.Options{
		TerraformDir: tfDir,
		NoColor:      true,
		Vars: map[string]interface{}{
			"seed":                        seedConfig,
			"snowflake_organization_name": os.Getenv("SNOWFLAKE_ORGANIZATION_NAME"),
			"snowflake_account_name":      os.Getenv("SNOWFLAKE_ACCOUNT_NAME"),
			"snowflake_user":              os.Getenv("SNOWFLAKE_USER"),
			"snowflake_role":              os.Getenv("SNOWFLAKE_ROLE"),
			"snowflake_private_key":       os.Getenv("SNOWFLAKE_PRIVATE_KEY"),
		},
	}

	defer terraform.Destroy(t, tfOptions)
	terraform.InitAndApply(t, tfOptions)

	time.Sleep(retrySleep)

	// Verify NO data was inserted (blocked)
	count := countRows(t, db, database, schema, table)
	require.Equal(t, 0, count, "Expected 0 rows - seeding should be blocked in prod")

	// Verify outputs
	seedExecuted := terraform.Output(t, tfOptions, "seed_executed")
	require.Equal(t, "false", seedExecuted)

	seedBlocked := terraform.Output(t, tfOptions, "seed_blocked")
	require.Equal(t, "true", seedBlocked)
}

// TestSeedDisabled tests that seeding doesn't run when disabled
func TestSeedDisabled(t *testing.T) {
	t.Parallel()

	retrySleep := 5 * time.Second
	unique := strings.ToUpper(random.UniqueId())

	database := "TT_SEED_DB"
	schema := fmt.Sprintf("TT_SCHEMA_%s", unique)
	table := fmt.Sprintf("TT_TABLE_%s", unique)

	tfDir := "../examples/basic"

	// Pre-create the test table
	db := openSnowflake(t)
	createTestTable(t, db, database, schema, table)
	defer func() {
		dropTestSchema(t, db, database, schema)
		_ = db.Close()
	}()

	// Inline SQL to insert test data (single line)
	sqlText := fmt.Sprintf("INSERT INTO %s.%s.%s (id, name) VALUES (1, 'Should Not Appear');", database, schema, table)

	// Seeding is disabled
	seedConfig := map[string]interface{}{
		"enabled":     false,
		"environment": "devl",
		"database":    database,
		"schema":      schema,
		"table":       table,
		"sql_text":    sqlText,
	}

	tfOptions := &terraform.Options{
		TerraformDir: tfDir,
		NoColor:      true,
		Vars: map[string]interface{}{
			"seed":                        seedConfig,
			"snowflake_organization_name": os.Getenv("SNOWFLAKE_ORGANIZATION_NAME"),
			"snowflake_account_name":      os.Getenv("SNOWFLAKE_ACCOUNT_NAME"),
			"snowflake_user":              os.Getenv("SNOWFLAKE_USER"),
			"snowflake_role":              os.Getenv("SNOWFLAKE_ROLE"),
			"snowflake_private_key":       os.Getenv("SNOWFLAKE_PRIVATE_KEY"),
		},
	}

	defer terraform.Destroy(t, tfOptions)
	terraform.InitAndApply(t, tfOptions)

	time.Sleep(retrySleep)

	// Verify NO data was inserted (disabled)
	count := countRows(t, db, database, schema, table)
	require.Equal(t, 0, count, "Expected 0 rows - seeding is disabled")

	// Verify outputs
	seedExecuted := terraform.Output(t, tfOptions, "seed_executed")
	require.Equal(t, "false", seedExecuted)

	seedBlocked := terraform.Output(t, tfOptions, "seed_blocked")
	require.Equal(t, "false", seedBlocked)
}
