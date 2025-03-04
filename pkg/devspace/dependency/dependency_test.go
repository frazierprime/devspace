package dependency

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/devspace-cloud/devspace/pkg/devspace/config/generated"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"
	"github.com/devspace-cloud/devspace/pkg/util/fsutil"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/devspace-cloud/devspace/pkg/util/ptr"

	"gotest.tools/assert"
)

var logOutput string

type testLogger struct {
	log.DiscardLogger
}

func (t testLogger) Done(args ...interface{}) {
	logOutput = logOutput + "\nDone " + fmt.Sprint(args...)
}
func (t testLogger) Donef(format string, args ...interface{}) {
	logOutput = logOutput + "\nDone " + fmt.Sprintf(format, args...)
}

type updateAllTestCase struct {
	name string

	files            map[string]string
	dependencyTasks  []*latest.DependencyConfig
	activeConfig     *generated.CacheConfig
	allowCyclicParam bool

	expectedErr string
	expectedLog string
}

func TestUpdateAll(t *testing.T) {
	testCases := []updateAllTestCase{
		updateAllTestCase{
			name: "No Dependencies to update",
		},
		updateAllTestCase{
			name: "Update one dependency",
			files: map[string]string{
				"devspace.yaml":         "",
				"someDir/devspace.yaml": "",
			},
			dependencyTasks: []*latest.DependencyConfig{
				&latest.DependencyConfig{
					Source: &latest.SourceConfig{
						Path: ptr.String("someDir"),
					},
					Config: ptr.String("someDir/devspace.yaml"),
				},
			},
			activeConfig: &generated.CacheConfig{
				Images: map[string]*generated.ImageCache{
					"default": &generated.ImageCache{
						Tag: "1.15", // This will be appended to nginx during deploy
					},
				},
				Dependencies: map[string]string{},
			},
			allowCyclicParam: true,
			expectedLog: `
Done Resolved 1 dependencies`,
		},
	}

	dir, err := ioutil.TempDir("", "testFolder")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}

	wdBackup, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Error changing working directory: %v", err)
	}

	// Delete temp folder
	defer func() {
		err = os.Chdir(wdBackup)
		if err != nil {
			t.Fatalf("Error changing dir back: %v", err)
		}
		err = os.RemoveAll(dir)
		if err != nil {
			t.Fatalf("Error removing dir: %v", err)
		}
	}()

	for _, testCase := range testCases {
		for path, content := range testCase.files {
			err = fsutil.WriteToFile([]byte(content), path)
			assert.NilError(t, err, "Error writing file in testCase %s", testCase.name)
		}

		logOutput = ""

		testConfig := &latest.Config{
			Dependencies: &testCase.dependencyTasks,
		}
		generatedConfig := &generated.Config{
			ActiveConfig: "default",
			Configs: map[string]*generated.CacheConfig{
				"default": testCase.activeConfig,
			},
		}

		err = UpdateAll(testConfig, generatedConfig, testCase.allowCyclicParam, &testLogger{})

		if testCase.expectedErr == "" {
			assert.NilError(t, err, "Error updating all in testCase %s", testCase.name)
		} else {
			assert.Error(t, err, testCase.expectedErr, "Wrong or no error from UpdateALl in testCase %s", testCase.name)
		}
		assert.Equal(t, logOutput, testCase.expectedLog, "Unexpected log output in testCase %s", testCase.name)

		for path := range testCase.files {
			err = os.Remove(path)
			assert.NilError(t, err, "Error removing file in testCase %s", testCase.name)
		}
	}
}

type deployAllTestCase struct {
	name string

	files                        map[string]string
	dependencyTasks              []*latest.DependencyConfig
	activeConfig                 *generated.CacheConfig
	allowCyclicParam             bool
	updateDependenciesParam      bool
	skipPushParam                bool
	forceBuildParam              bool
	forceDeployDependenciesParam bool
	forceDeployParam             bool

	expectedErr string
	expectedLog string
}

func TestDeployAll(t *testing.T) {
	dir, err := ioutil.TempDir("", "testFolder")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}

	wdBackup, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Error changing working directory: %v", err)
	}

	// Delete temp folder
	defer func() {
		err = os.Chdir(wdBackup)
		if err != nil {
			t.Fatalf("Error changing dir back: %v", err)
		}
		err = os.RemoveAll(dir)
		if err != nil {
			t.Fatalf("Error removing dir: %v", err)
		}
	}()

	testCases := []deployAllTestCase{
		deployAllTestCase{
			name: "No Dependencies to deploy",
		},
		deployAllTestCase{
			name: "Deploy one dependency",
			files: map[string]string{
				"devspace.yaml":         "",
				"someDir/devspace.yaml": "",
			},
			dependencyTasks: []*latest.DependencyConfig{
				&latest.DependencyConfig{
					Source: &latest.SourceConfig{
						Path: ptr.String("someDir"),
					},
					Config: ptr.String("someDir/devspace.yaml"),
				},
			},
			activeConfig: &generated.CacheConfig{
				Images: map[string]*generated.ImageCache{
					"default": &generated.ImageCache{
						Tag: "1.15", // This will be appended to nginx during deploy
					},
				},
				Dependencies: map[string]string{},
			},
			allowCyclicParam: true,
			expectedLog: `
Done Resolved 1 dependencies`,
			expectedErr: fmt.Sprintf("Error deploying dependency %s:  Unable to create new kubectl client: invalid configuration: no configuration has been provided", dir+string(os.PathSeparator)+"someDir"),
		},
	}

	for _, testCase := range testCases {
		for path, content := range testCase.files {
			err = fsutil.WriteToFile([]byte(content), path)
			assert.NilError(t, err, "Error writing file in testCase %s", testCase.name)
		}

		logOutput = ""

		testConfig := &latest.Config{
			Dependencies: &testCase.dependencyTasks,
		}
		generatedConfig := &generated.Config{
			ActiveConfig: "default",
			Configs: map[string]*generated.CacheConfig{
				"default": testCase.activeConfig,
			},
		}

		err = DeployAll(testConfig, generatedConfig, testCase.allowCyclicParam, testCase.updateDependenciesParam, testCase.skipPushParam, testCase.forceDeployDependenciesParam, false, testCase.forceBuildParam, testCase.forceDeployParam, &testLogger{})

		if testCase.expectedErr == "" {
			assert.NilError(t, err, "Error deploying all in testCase %s", testCase.name)
		} else {
			assert.Error(t, err, testCase.expectedErr, "Wrong or no error from DeployALl in testCase %s", testCase.name)
		}
		assert.Equal(t, logOutput, testCase.expectedLog, "Unexpected log output in testCase %s", testCase.name)

		for path := range testCase.files {
			err = os.Remove(path)
			assert.NilError(t, err, "Error removing file in testCase %s", testCase.name)
		}
	}
}

type purgeAllTestCase struct {
	name string

	files            map[string]string
	dependencyTasks  []*latest.DependencyConfig
	activeConfig     *generated.CacheConfig
	allowCyclicParam bool

	expectedErr string
	expectedLog string
}

func TestPurgeAll(t *testing.T) {
	dir, err := ioutil.TempDir("", "testFolder")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}

	wdBackup, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Error changing working directory: %v", err)
	}

	// Delete temp folder
	defer func() {
		err = os.Chdir(wdBackup)
		if err != nil {
			t.Fatalf("Error changing dir back: %v", err)
		}
		err = os.RemoveAll(dir)
		if err != nil {
			t.Fatalf("Error removing dir: %v", err)
		}
	}()

	testCases := []purgeAllTestCase{
		purgeAllTestCase{
			name: "No Dependencies to update",
		},
		purgeAllTestCase{
			name: "Update one dependency",
			files: map[string]string{
				"devspace.yaml":         "",
				"someDir/devspace.yaml": "",
			},
			dependencyTasks: []*latest.DependencyConfig{
				&latest.DependencyConfig{
					Source: &latest.SourceConfig{
						Path: ptr.String("someDir"),
					},
					Config: ptr.String("someDir/devspace.yaml"),
				},
			},
			activeConfig: &generated.CacheConfig{
				Images: map[string]*generated.ImageCache{
					"default": &generated.ImageCache{
						Tag: "1.15", // This will be appended to nginx during deploy
					},
				},
				Dependencies: map[string]string{},
			},
			allowCyclicParam: true,
			expectedLog: `
Done Resolved 1 dependencies`,
			expectedErr: fmt.Sprintf("Error deploying dependency %s:  Unable to create new kubectl client: invalid configuration: no configuration has been provided", dir+string(os.PathSeparator)+"someDir"),
		},
	}

	for _, testCase := range testCases {
		for path, content := range testCase.files {
			err = fsutil.WriteToFile([]byte(content), path)
			assert.NilError(t, err, "Error writing file in testCase %s", testCase.name)
		}

		logOutput = ""

		testConfig := &latest.Config{
			Dependencies: &testCase.dependencyTasks,
		}
		generatedConfig := &generated.Config{
			ActiveConfig: "default",
			Configs: map[string]*generated.CacheConfig{
				"default": testCase.activeConfig,
			},
		}

		err = PurgeAll(testConfig, generatedConfig, testCase.allowCyclicParam, &testLogger{})

		if testCase.expectedErr == "" {
			assert.NilError(t, err, "Error purging all in testCase %s", testCase.name)
		} else {
			assert.Error(t, err, testCase.expectedErr, "Wrong or no error from PurgeALl in testCase %s", testCase.name)
		}
		assert.Equal(t, logOutput, testCase.expectedLog, "Unexpected log output in testCase %s", testCase.name)

		for path := range testCase.files {
			err = os.Remove(path)
			assert.NilError(t, err, "Error removing file in testCase %s", testCase.name)
		}
	}
}
