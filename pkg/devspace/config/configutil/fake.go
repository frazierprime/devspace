package configutil

import (
	"github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"
	"github.com/devspace-cloud/devspace/pkg/util/ptr"
)

// TestNamespace is the test namespace to use
const TestNamespace = "test-namespace"

// SetFakeConfig initializes the config objects
func SetFakeConfig(fakeConfig *latest.Config) {
	getConfigOnce.Do(func() {})

	if fakeConfig.Deployments == nil {
		fakeConfig.Deployments = &[]*latest.DeploymentConfig{}
	}
	if fakeConfig.Images == nil {
		fakeConfig.Images = &map[string]*latest.ImageConfig{}
	}
	if fakeConfig.Dev == nil {
		fakeConfig.Dev = &latest.DevConfig{}
	}

	fakeConfig.Cluster = &latest.Cluster{
		Namespace: ptr.String(TestNamespace),
	}

	config = fakeConfig
}
