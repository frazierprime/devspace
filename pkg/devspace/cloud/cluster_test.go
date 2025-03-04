package cloud

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/devspace-cloud/devspace/pkg/devspace/cloud/config/versions/latest"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/constants"
	"github.com/devspace-cloud/devspace/pkg/util/survey"

	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"gotest.tools/assert"
)

func TestConnectCluster(t *testing.T) {
	provider := &Provider{}
	options := &ConnectClusterOptions{
		ClusterName: "#",
	}
	err := provider.ConnectCluster(options)
	assert.Error(t, err, "Cluster name # can only contain letters, numbers and dashes (-)", "Wrong or no error when connecting cluster with wrong clustername")

	options.ClusterName = ""
	survey.SetNextAnswer("aaa")
	options.KubeContext = "invalidContext"

	err = provider.ConnectCluster(options)
	assert.Error(t, err, "new kubectl client: invalid configuration: [context was not found for specified context: invalidContext, cluster has no server defined]", "Wrong or no error when connecting cluster with invalid context")
}

func TestDefaultClusterSpaceDomain(t *testing.T) {
	// @Florian make test faster (currently around 10 seconds)
	t.Skip("Takes too long")

	kubeClient := fake.NewSimpleClientset()
	err := defaultClusterSpaceDomain(&Provider{}, kubeClient, true, 0, "")
	assert.Error(t, err, "Couldn't find a node in cluster", "Wrong or no error when trying to get the spacedomain of the default cluster from empty setting")

	kubeClient.CoreV1().Nodes().Create(&k8sv1.Node{})
	err = defaultClusterSpaceDomain(&Provider{}, kubeClient, true, 0, "")
	assert.Error(t, err, "Couldn't find a node with a valid external ip in cluster, make sure your nodes are accessable from the outside", "Wrong or no error when trying to get the spacedomain of the default cluster without any ip")

	kubeClient.CoreV1().Nodes().Update(&k8sv1.Node{
		Status: k8sv1.NodeStatus{
			Addresses: []k8sv1.NodeAddress{
				k8sv1.NodeAddress{
					Type:    k8sv1.NodeExternalIP,
					Address: "someAddress",
				},
			},
		},
	})
	err = defaultClusterSpaceDomain(&Provider{}, kubeClient, true, 0, "")
	assert.Error(t, err, "get token: Provider has no key specified", "Wrong or no error when trying to get the spacedomain of the default cluster without a token")

	waitTimeout = time.Second * 8
	err = defaultClusterSpaceDomain(&Provider{}, kubeClient, false, 0, "")
	assert.Error(t, err, "Loadbalancer didn't receive a valid ip in time. Skipping configuration of default cluster space url", "Wrong or no error when trying to get the spacedomain of the default cluster without services")

	kubeClient.CoreV1().Services(constants.DevSpaceCloudNamespace).Create(&k8sv1.Service{
		Spec: k8sv1.ServiceSpec{
			Type: k8sv1.ServiceTypeLoadBalancer,
		},
		Status: k8sv1.ServiceStatus{
			LoadBalancer: k8sv1.LoadBalancerStatus{
				Ingress: []k8sv1.LoadBalancerIngress{
					k8sv1.LoadBalancerIngress{
						IP:       "SomeIp",
						Hostname: "SomeHost",
					},
				},
			},
		},
	})
	err = defaultClusterSpaceDomain(&Provider{}, kubeClient, false, 0, "")
	assert.Error(t, err, "get token: Provider has no key specified", "Wrong or no error when trying to get the spacedomain of the default cluster without a token")
}

func TestDeleteClusterUnexported(t *testing.T) {
	provider := &Provider{}
	err := deleteCluster(provider, 0, "")
	assert.Error(t, err, "get token: Provider has no key specified", "Wrong or no error when trying to delete a cluster without a token")
}

func TestSpecifyDomain(t *testing.T) {
	provider := &Provider{}
	survey.SetNextAnswer("some.Domain")
	err := provider.specifyDomain(0, &ConnectClusterOptions{})
	assert.Error(t, err, "update cluster domain: get token: Provider has no key specified", "Wrong or no error when trying to delete a space without a token")
}

func TestInitCore(t *testing.T) {
	provider := &Provider{}
	err := provider.initCore(0, "", true)
	assert.Error(t, err, "get token: Provider has no key specified", "Wrong or no error when trying to init the core without a token")
}

func TestGetServiceAccountCredentials(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()
	kubeClient.CoreV1().ServiceAccounts(DevSpaceCloudNamespace).Create(&k8sv1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: DevSpaceServiceAccount,
		},
		Secrets: []k8sv1.ObjectReference{
			k8sv1.ObjectReference{
				Name: "secret",
			},
		},
	})

	_, _, err := getServiceAccountCredentials(kubeClient)
	assert.Error(t, err, "secrets \"secret\" not found", "Wrong or no error when getting non-existent service account credentials")

	flag := []byte("flag")
	kubeClient.CoreV1().Secrets(DevSpaceCloudNamespace).Create(&k8sv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "secret",
		},
		Data: map[string][]byte{
			"token":  flag,
			"ca.crt": flag,
		},
	})
	returnedToken, cert, err := getServiceAccountCredentials(kubeClient)
	assert.NilError(t, err, "Error getting service account credentials")
	assert.Equal(t, string(flag), string(returnedToken), "Wrong token returned")
	decodedCert, err := base64.StdEncoding.DecodeString(cert)
	assert.NilError(t, err, "Error decoding returned cert")
	assert.Equal(t, string(flag), string(decodedCert), "Wrong cert returned")
}

func TestGetKey(t *testing.T) {
	provider := &Provider{
		latest.Provider{
			ClusterKey: map[int]string{
				5: "onlyKey",
			},
		},
	}
	returnedKey, err := getKey(provider, false)
	assert.NilError(t, err, "Error getting the only clusteKey without forcing a question")
	assert.Equal(t, returnedKey, "onlyKey", "Wrong key returned from Clusterkey with only one item")
}

type checkResourcesTestCase struct {
	name         string
	createdNodes []*k8sv1.Node
	expectedErr  string
}

func TestCheckResources(t *testing.T) {
	testCases := []checkResourcesTestCase{
		checkResourcesTestCase{
			name:         "Test without nodes",
			createdNodes: []*k8sv1.Node{},
			expectedErr:  "The cluster specified has no nodes, please choose a cluster where at least one node is up and running",
		},
		checkResourcesTestCase{
			name:         "Test without group versions",
			createdNodes: []*k8sv1.Node{&k8sv1.Node{}},
			expectedErr:  "Group version rbac.authorization.k8s.io/v1beta1 does not exist in cluster, but is required. Is RBAC enabled?",
		},
	}

	for _, testCase := range testCases {
		kubeClient := fake.NewSimpleClientset()
		for _, node := range testCase.createdNodes {
			kubeClient.CoreV1().Nodes().Create(node)
		}

		_, err := checkResources(kubeClient)
		if testCase.expectedErr == "" {
			assert.NilError(t, err, "Error checking resources in testCase %s", testCase.name)
		} else {
			assert.Error(t, err, testCase.expectedErr, "Wrong or no error from checking resources in testCase %s", testCase.name)
		}
	}
}
