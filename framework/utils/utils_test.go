package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	// arrange
	json := `{
				"id": "525b5fd9-700d-4feb-89c0-415a1e6e148c",
				"file_path": "convite.mp4",
				"status": "pending"
			}`

	// act & assert
	err := IsJson(json)
	require.Nil(t, err)

	json = `sidarta`
	err = IsJson(json)
	require.NotNil(t, err)
}
