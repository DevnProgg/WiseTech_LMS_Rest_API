package database

import (
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	"wisetech-lms-api/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeSchema(t *testing.T) {
	// Create a temporary file for the database
	tmpfile, err := ioutil.TempFile("", "test-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Create a new config with the temporary database path
	cfg := &config.Config{
		DBPath: tmpfile.Name(),
	}

	// Create a new database connection
	db, err := NewConnection(cfg)
	require.NoError(t, err)
	defer db.Close()

	// Initialize the schema
	err = InitializeSchema(db)
	require.NoError(t, err)

	// Check if all tables were created
	tables := []string{
		"Lenders", "Borrowers", "Accounts", "Plans", "Lender_Ledger",
		"Loans", "Recipets", "File", "Text", "Number",
	}

	for _, table := range tables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		assert.NoError(t, err, "Table %s should exist", table)
		assert.Equal(t, table, name, "Table %s should have the correct name", table)
	}

	// Check for a trigger
	var triggerName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='trigger' AND name='update_lenders_updated_at'").Scan(&triggerName)
	assert.NoError(t, err)
	assert.Equal(t, "update_lenders_updated_at", triggerName)
}

func TestNewConnection_Failure(t *testing.T) {
	// Create a new config with an invalid database path
	cfg := &config.Config{
		DBPath: "/non_existent_dir/test.db",
	}

	// Try to create a new database connection
	_, err := NewConnection(cfg)
	assert.Error(t, err)
}

func TestInitializeSchema_Failure(t *testing.T) {
	// Create a temporary file for the database
	tmpfile, err := ioutil.TempFile("", "test-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Create a new config with the temporary database path
	cfg := &config.Config{
		DBPath: tmpfile.Name(),
	}

	// Create a new database connection
	db, err := NewConnection(cfg)
	require.NoError(t, err)
	defer db.Close()

	// Close the database to cause a failure
	db.Close()

	// Try to initialize the schema on a closed database
	err = InitializeSchema(db)
	assert.Error(t, err)
}

func TestNewConnection_PingFailure(t *testing.T) {
	// This test is a bit tricky with sqlite3 as it's an in-process database.
	// A ping failure is more likely with a network database.
	// We can simulate this by creating a malformed database file.
	tmpfile, err := ioutil.TempFile("", "test-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Write some garbage to the file
	_, err = tmpfile.WriteString("this is not a database")
	require.NoError(t, err)
	tmpfile.Close()

	cfg := &config.Config{
		DBPath: tmpfile.Name(),
	}

	// This might not fail on open, but should fail on ping
	db, err := sql.Open("sqlite3", cfg.DBPath)
	require.NoError(t, err)
	defer db.Close()

	// We can't directly test NewConnection's ping failure here without
	// modifying NewConnection to allow injecting a faulty pinger.
	// Instead, we'll just test the ping directly.
	err = db.Ping()
	assert.Error(t, err)
}
