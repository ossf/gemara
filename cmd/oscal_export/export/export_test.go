package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuidance(t *testing.T) {
	tempDir := t.TempDir()
	mockYAML := `
metadata:
  id: Test
  title: Test
  description: ""
categories:
  - id: TEST
    title: Test
    description: Test
`
	inputFilePath := filepath.Join(tempDir, "guidance.yaml")
	require.NoError(t, os.WriteFile(inputFilePath, []byte(mockYAML), 0600))

	t.Run("Success/Defaults", func(t *testing.T) {
		catalogFilePath := filepath.Join(tempDir, "guidance.json")
		profileFilePath := filepath.Join(tempDir, "profile.json")

		args := []string{"--catalog-output", catalogFilePath, "--profile-output", profileFilePath}
		err := Guidance(inputFilePath, args)
		require.NoError(t, err)

		if _, err := os.Stat(catalogFilePath); os.IsNotExist(err) {
			t.Fatalf("Catalog output file not created: %s", catalogFilePath)
		}

		var catalogModel oscal.OscalModels
		catalogData, _ := os.ReadFile(catalogFilePath)
		require.NoError(t, json.Unmarshal(catalogData, &catalogModel))
		assert.NotNil(t, catalogModel.Catalog)

		if _, err := os.Stat(profileFilePath); os.IsNotExist(err) {
			t.Fatalf("Profile output file not created: %s", profileFilePath)
		}

		var profileModel oscal.OscalModels
		profileData, _ := os.ReadFile(profileFilePath)
		require.NoError(t, json.Unmarshal(profileData, &profileModel))
		assert.NotNil(t, profileModel.Profile)
	})

	t.Run("Failure/NotExists", func(t *testing.T) {
		err := Guidance("non-existent-file.yaml", []string{})
		require.Error(t, err)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("Failure/InvalidInput", func(t *testing.T) {
		failYAMLPath := filepath.Join(t.TempDir(), "fail-profile.yaml")
		require.NoError(t, os.WriteFile(failYAMLPath, []byte("fail-profile"), 0600))
		err := Guidance(failYAMLPath, []string{})
		require.ErrorContains(t, err, "string was used where mapping is expected")
	})
}

func TestCatalog(t *testing.T) {
	tempDir := t.TempDir()

	mockYAML := `
metadata:
  id: Test
  title: Test
  description: ""
control-families:
  - id: TEST
    title: Test
`
	inputFilePath := filepath.Join(tempDir, "catalog.yaml")
	require.NoError(t, os.WriteFile(inputFilePath, []byte(mockYAML), 0600))

	t.Run("Success/Defaults", func(t *testing.T) {
		catalogFilePath := filepath.Join(tempDir, "catalog.json")
		args := []string{"--output", catalogFilePath}
		err := Catalog(inputFilePath, args)
		require.NoError(t, err)

		if _, err := os.Stat(catalogFilePath); os.IsNotExist(err) {
			t.Fatalf("Catalog output file not created: %s", catalogFilePath)
		}

		var catalogModel oscal.OscalModels
		catalogData, _ := os.ReadFile(catalogFilePath)
		require.NoError(t, json.Unmarshal(catalogData, &catalogModel))
		assert.NotNil(t, catalogModel.Catalog)
	})

	t.Run("Failure/NotExists", func(t *testing.T) {
		err := Catalog("non-existent-file.yaml", []string{})
		require.Error(t, err)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("Failure/InvalidInput", func(t *testing.T) {
		failYAMLPath := filepath.Join(t.TempDir(), "fail.yaml")
		require.NoError(t, os.WriteFile(failYAMLPath, []byte("fail"), 0600))
		err := Catalog(failYAMLPath, []string{})
		require.ErrorContains(t, err, "string was used where mapping is expected")
	})
}
