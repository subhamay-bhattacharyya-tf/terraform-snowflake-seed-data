// File: test/helpers_test.go
package test

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

func openSnowflake(t *testing.T) *sql.DB {
	t.Helper()

	orgName := mustEnv(t, "SNOWFLAKE_ORGANIZATION_NAME")
	accountName := mustEnv(t, "SNOWFLAKE_ACCOUNT_NAME")
	user := mustEnv(t, "SNOWFLAKE_USER")
	privateKeyPEM := mustEnv(t, "SNOWFLAKE_PRIVATE_KEY")
	role := os.Getenv("SNOWFLAKE_ROLE")

	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	require.NotNil(t, block, "Failed to decode PEM block from private key")

	var privateKey *rsa.PrivateKey
	var err error

	// Try PKCS8 first, then PKCS1
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		require.NoError(t, err, "Failed to parse private key")
	} else {
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		require.True(t, ok, "Private key is not RSA")
	}

	// Build account identifier: orgname-accountname
	account := fmt.Sprintf("%s-%s", orgName, accountName)

	config := gosnowflake.Config{
		Account:       account,
		User:          user,
		Authenticator: gosnowflake.AuthTypeJwt,
		PrivateKey:    privateKey,
	}

	if role != "" {
		config.Role = role
	}

	dsn, err := gosnowflake.DSN(&config)
	require.NoError(t, err, "Failed to build DSN")

	db, err := sql.Open("snowflake", dsn)
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	return db
}

func tableExists(t *testing.T, db *sql.DB, database, schema, table string) bool {
	t.Helper()

	q := fmt.Sprintf("SHOW TABLES LIKE '%s' IN %s.%s;", escapeLike(table), database, schema)
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	return rows.Next()
}

func countRows(t *testing.T, db *sql.DB, database, schema, table string) int {
	t.Helper()

	q := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s.%s;", database, schema, table)
	var count int
	err := db.QueryRow(q).Scan(&count)
	require.NoError(t, err)
	return count
}

func createTestTable(t *testing.T, db *sql.DB, database, schema, table string) {
	t.Helper()

	// Create database if not exists
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", database))
	require.NoError(t, err)

	// Create schema if not exists
	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s.%s;", database, schema))
	require.NoError(t, err)

	// Create table
	q := fmt.Sprintf(`
		CREATE OR REPLACE TABLE %s.%s.%s (
			id INTEGER,
			name VARCHAR(100),
			created_at TIMESTAMP_NTZ DEFAULT CURRENT_TIMESTAMP()
		);
	`, database, schema, table)
	_, err = db.Exec(q)
	require.NoError(t, err)
}

func dropTestSchema(t *testing.T, db *sql.DB, database, schema string) {
	t.Helper()

	q := fmt.Sprintf("DROP SCHEMA IF EXISTS %s.%s CASCADE;", database, schema)
	_, err := db.Exec(q)
	require.NoError(t, err)
}

func mustEnv(t *testing.T, key string) string {
	t.Helper()
	v := strings.TrimSpace(os.Getenv(key))
	require.NotEmpty(t, v, "Missing required environment variable %s", key)
	return v
}

func escapeLike(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
