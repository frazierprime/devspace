package generated

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/devspace-cloud/devspace/pkg/util/fsutil"

	"gotest.tools/assert"
)

func TestLoadConfigFromPath(t *testing.T) {
	//Create tempDir and go into it
	dir, err := ioutil.TempDir("", "testDir")
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

	// Delete temp folder after test
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

	ConfigPath = "NotExistent"
	loadedConfigOnce = sync.Once{}
	returnedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config from non existent path: %v", err)
	}
	assert.Equal(t, DefaultConfigName, returnedConfig.ActiveConfig, "Wrong initial config returned")
	assert.Equal(t, 1, len(returnedConfig.Configs), "Wrong initial config returned")
	assert.Equal(t, false, returnedConfig.GetActive() == nil, "Active config not initialized")

	ConfigPath = "generated.yaml"
	loadedConfigOnce = sync.Once{}
	fsutil.WriteToFile([]byte(""), "generated.yaml")
	returnedConfig, err = LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config from existent path with empty content: %v", err)
	}
	assert.Equal(t, DefaultConfigName, returnedConfig.ActiveConfig, "Wrong initial config returned")
	assert.Equal(t, 1, len(returnedConfig.Configs), "Wrong initial config returned")
	assert.Equal(t, false, returnedConfig.GetActive() == nil, "Active config not initialized")
}

func TestSaveConfig(t *testing.T) {
	//Create tempDir and go into it
	dir, err := ioutil.TempDir("", "testDir")
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

	// Delete temp folder after test
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

	testDontSaveConfig = false
	err = SaveConfig(&Config{
		ActiveConfig: "SavedActiveConfig",
		Configs: map[string]*CacheConfig{
			"SavedActiveConfig": &CacheConfig{},
		},
		CloudSpace: &CloudSpaceConfig{
			Name: "SavedCloudSpaceConfig",
		},
	})
	if err != nil {
		t.Fatalf("Error saving config: %v", err)
	}

	returnedConfig, err := LoadConfigFromPath(ConfigPath)
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}
	assert.Equal(t, "SavedActiveConfig", returnedConfig.ActiveConfig, "Wrong config saved or returned")
	assert.Equal(t, "SavedCloudSpaceConfig", returnedConfig.CloudSpace.Name, "Wrong config saved or returned")

	//Now with testDontSaveConfig set true. Loaded config shouldn't change
	SetTestConfig(&Config{})
	err = SaveConfig(&Config{
		ActiveConfig: "NewActiveConfig",
		Configs: map[string]*CacheConfig{
			"NewActiveConfig": &CacheConfig{},
		},
		CloudSpace: &CloudSpaceConfig{
			Name: "NewCloudSpaceConfig",
		},
	})
	if err != nil {
		t.Fatalf("Error saving config: %v", err)
	}

	returnedConfig, err = LoadConfigFromPath(ConfigPath)
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}
	assert.Equal(t, "SavedActiveConfig", returnedConfig.ActiveConfig, "Wrong config saved or returned")
	assert.Equal(t, "SavedCloudSpaceConfig", returnedConfig.CloudSpace.Name, "Wrong config saved or returned")
}

func TestGetCaches(t *testing.T) {
	dsConfig := &Config{
		Configs: map[string]*CacheConfig{
			"SomeConfig": &CacheConfig{},
		},
	}
	InitDevSpaceConfig(dsConfig, "SomeConfig")
	cacheConfig := dsConfig.Configs["SomeConfig"]
	assert.Equal(t, 0, len(cacheConfig.Deployments), "Deployments wrong initialized")
	assert.Equal(t, 0, len(cacheConfig.Images), "Images wrong initialized")
	assert.Equal(t, 0, len(cacheConfig.Dependencies), "Dependencies wrong initialized")
	assert.Equal(t, 0, len(cacheConfig.Vars), "Vars wrong initialized")

	imageCache := cacheConfig.GetImageCache("NewImageCache")
	assert.Equal(t, 1, len(cacheConfig.Images), "New imageCache not added to cache")
	assert.Equal(t, "", imageCache.ImageConfigHash, "ImageCache wrong initialized")
	assert.Equal(t, "", imageCache.DockerfileHash, "ImageCache wrong initialized")
	assert.Equal(t, "", imageCache.ContextHash, "ImageCache wrong initialized")
	assert.Equal(t, "", imageCache.EntrypointHash, "ImageCache wrong initialized")
	assert.Equal(t, "", imageCache.CustomFilesHash, "ImageCache wrong initialized")
	assert.Equal(t, "", imageCache.ImageName, "ImageCache wrong initialized")
	assert.Equal(t, "", imageCache.Tag, "ImageCache wrong initialized")

	deploymentCache := cacheConfig.GetDeploymentCache("NewDeploymentCache")
	assert.Equal(t, 1, len(cacheConfig.Deployments), "New deploymentCache not added to cache")
	assert.Equal(t, "", deploymentCache.DeploymentConfigHash, "DeploymentCache wrong initialized")
	assert.Equal(t, "", deploymentCache.HelmOverridesHash, "DeploymentCache wrong initialized")
	assert.Equal(t, "", deploymentCache.HelmChartHash, "DeploymentCache wrong initialized")
	assert.Equal(t, "", deploymentCache.KubectlManifestsHash, "DeploymentCache wrong initialized")
}
