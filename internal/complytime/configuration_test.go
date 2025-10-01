// SPDX-License-Identifier: Apache-2.0

package complytime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oscal-compass/oscal-sdk-go/validation"
	"github.com/stretchr/testify/require"
)

func TestApplicationDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	appDir, err := newApplicationDirectory(tmpDir, false)
	require.NoError(t, err)

	expectedAppDir := filepath.Join(tmpDir, "complytime")
	expectedPluginDir := filepath.Join(tmpDir, "complytime", "plugins")
	expectedPluginManifestDir := filepath.Join(tmpDir, "complytime", "plugins")
	expectedBundleDir := filepath.Join(tmpDir, "complytime", "bundles")
	expectedControlDir := filepath.Join(tmpDir, "complytime", "controls")

	require.Equal(t, expectedAppDir, appDir.AppDir())
	require.Equal(t, expectedPluginDir, appDir.PluginDir())
	require.Equal(t, expectedPluginManifestDir, appDir.PluginManifestDir())
	require.Equal(t, expectedBundleDir, appDir.BundleDir())
	require.Equal(t, expectedControlDir, appDir.ControlDir())
	require.Equal(t, []string{expectedAppDir, expectedPluginDir, expectedPluginManifestDir, expectedBundleDir, expectedControlDir}, appDir.Dirs())

	appDir, err = newApplicationDirectory(tmpDir, true)
	require.NoError(t, err)
	_, err = os.Stat(appDir.AppDir())
	require.NoError(t, err)
	_, err = os.Stat(appDir.PluginDir())
	require.NoError(t, err)
	_, err = os.Stat(appDir.BundleDir())
	require.NoError(t, err)
	_, err = os.Stat(appDir.ControlDir())
	require.NoError(t, err)
}

func TestFindComponentDefinitions(t *testing.T) {
	compDefs, err := FindComponentDefinitions("testdata/complytime/bundles", validation.NoopValidator{})
	require.NoError(t, err)
	require.Len(t, compDefs, 1)

	_, err = FindComponentDefinitions("testdata/", validation.NoopValidator{})
	require.ErrorIs(t, err, ErrNoComponentDefinitionsFound)

}

func TestEnsureUserWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	testPlanPath := filepath.Join(tmpDir, "test_workspace")

	// Ensure user workspace directory is properly created
	err := EnsureUserWorkspace(testPlanPath)
	require.NoError(t, err)

	info, err := os.Stat(testPlanPath)
	require.NoError(t, err)
	require.NotNil(t, info)
	require.True(t, info.IsDir(), "expected a directory, got something else")

	// Now lets remove the created directory and set the parent dir as read-only
	err = os.RemoveAll(testPlanPath)
	require.NoError(t, err, "failed to remove test_workspace directory")

	err = os.Chmod(tmpDir, 0500) // read + execute only, no write
	require.NoError(t, err, "failed to chmod parent dir")
	defer os.Chmod(tmpDir, 0755) // restore permissions so TempDir cleanup works

	// Try to create the user workspace again but now it should fail
	err = EnsureUserWorkspace(testPlanPath)
	require.Error(t, err, "expected error when trying to create dir in read-only parent")
}
