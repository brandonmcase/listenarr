package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	// Create a temporary database file
	testDBPath := filepath.Join(os.TempDir(), "test_listenarr.db")
	defer os.Remove(testDBPath)

	db, err := Initialize(testDBPath)
	require.NoError(t, err)
	assert.NotNil(t, db)

	// Test that we can perform a simple query
	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestInitialize_InvalidPath(t *testing.T) {
	// Try to initialize with an invalid path (directory instead of file)
	invalidPath := os.TempDir()

	db, err := Initialize(invalidPath)
	// This might succeed or fail depending on SQLite behavior
	// We just want to ensure it doesn't panic
	if err != nil {
		assert.Nil(t, db)
	}
}

func TestInitialize_CreatesFile(t *testing.T) {
	testDBPath := filepath.Join(os.TempDir(), "test_creates.db")
	defer os.Remove(testDBPath)

	// Ensure file doesn't exist
	os.Remove(testDBPath)

	db, err := Initialize(testDBPath)
	require.NoError(t, err)
	assert.NotNil(t, db)

	// Verify file was created
	_, err = os.Stat(testDBPath)
	assert.NoError(t, err, "Database file should be created")
}

