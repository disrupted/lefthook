package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInstallCmdExecutor(t *testing.T) {
	// if branch
	fs := afero.NewMemMapFs()

	InstallCmdExecutor([]string{}, fs)

	expectedFile := "hookah.yml"

	_, err := fs.Stat(expectedFile)
	assert.Equal(t, os.IsNotExist(err), false, "hookah.yml not exists after install command")

	// else branch
	fs = afero.NewMemMapFs()
	presetConfig(fs)

	InstallCmdExecutor([]string{}, fs)

	expectedFiles := []string{
		"commit-msg",
		"pre-commit",
	}

	files, err := afero.ReadDir(fs, filepath.Join(getRootPath(), ".git/hooks"))
	assert.NoError(t, err)

	actualFiles := []string{}
	for _, f := range files {
		actualFiles = append(actualFiles, f.Name())
	}
	assert.Equal(t, expectedFiles, actualFiles, "Expected files not exists")
}

func TestAddCmdExecutor(t *testing.T) {
	fs := afero.NewMemMapFs()
	presetConfig(fs)

	addCmdExecutor([]string{"pre-push"}, fs)

	expectedFiles := []string{
		"pre-push",
	}

	expectedDirs := []string{
		"commit-msg",
		"pre-commit",
		"pre-push",
	}

	files, _ := afero.ReadDir(fs, filepath.Join(getRootPath(), ".git/hooks"))
	actualFiles := []string{}
	for _, f := range files {
		actualFiles = append(actualFiles, f.Name())
	}

	dirs, _ := afero.ReadDir(fs, filepath.Join(getRootPath(), ".hookah"))
	actualDirs := []string{}
	for _, f := range dirs {
		actualDirs = append(actualDirs, f.Name())
	}

	assert.Equal(t, expectedFiles, actualFiles, "Expected files not exists")
	assert.Equal(t, expectedDirs, actualDirs, "Expected dirs not exists")

	addCmdExecutor(expectedFiles, fs)

	expectedFiles = []string{
		"pre-push",
		"pre-push.old",
	}

	files, _ = afero.ReadDir(fs, filepath.Join(getRootPath(), ".git/hooks"))
	actualFiles = []string{}
	for _, f := range files {
		actualFiles = append(actualFiles, f.Name())
	}

	assert.Equal(t, expectedDirs, actualDirs, "Haven`t renamed file with .old extension")
}

// TODO: little tricky to call exec.Command with virtual file system
// func TestRunCmdExecutor(t *testing.T) {
// 	fs := afero.NewMemMapFs()
// 	presetConfig(fs)
// 	presetExecutable("fail_script", "pre-commit", "1", fs)

// 	err := RunCmdExecutor([]string{"pre-commit"}, fs)
// 	assert.Error(t, err)
// }

func presetConfig(fs afero.Fs) {
	viper.SetDefault(configSourceDirKey, ".hookah")

	AddConfigYaml(fs)

	fs.Mkdir(filepath.Join(getRootPath(), ".hookah/commit-msg"), defaultFilePermission)
	fs.Mkdir(filepath.Join(getRootPath(), ".hookah/pre-commit"), defaultFilePermission)

	fs.Mkdir(filepath.Join(getRootPath(), ".git/hooks"), defaultFilePermission)
}

func presetExecutable(hookName string, hookGroup string, exitCode string, fs afero.Fs) {
	template := "#!/bin/sh\nexit " + exitCode + "\n"
	pathToFile := filepath.Join(".hookah", hookGroup, hookName)
	err := afero.WriteFile(fs, pathToFile, []byte(template), defaultFilePermission)
	check(err)
}