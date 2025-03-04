package kaniko

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/devspace-cloud/devspace/pkg/devspace/config/configutil"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/generated"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"
	"github.com/devspace-cloud/devspace/pkg/devspace/docker"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/devspace-cloud/devspace/pkg/util/ptr"

	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const testNamespace = "test-kaniko-build"

func TestKanikoBuildWithEntrypointOverride(t *testing.T) {
	t.Skip("Package is untestable because of kubeClient stream usage")

	// 1. Write test dockerfile and context to a temp folder
	dir, err := ioutil.TempDir("", "testKaniko")
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

	// 5. Delete temp files
	defer os.Chdir(wdBackup)
	defer os.RemoveAll(dir)

	err = makeTestProject(dir)
	if err != nil {
		t.Fatalf("Error creating test project: %v", err)
	}

	// 2. Create kubectl client
	deployConfig := &latest.DeploymentConfig{
		Name: ptr.String("test-deployment"),
		Component: &latest.ComponentConfig{
			Containers: &[]*latest.ContainerConfig{
				{
					Image: ptr.String("nginx"),
				},
			},
			Service: &latest.ServiceConfig{
				Ports: &[]*latest.ServicePortConfig{
					{
						Port: ptr.Int(3000),
					},
				},
			},
		},
	}

	// Create fake devspace config
	testConfig := &latest.Config{
		Deployments: &[]*latest.DeploymentConfig{
			deployConfig,
		},
		// The images config will tell the deployment method to override the image name used in the component above with the tag defined in the generated config below
		Images: &map[string]*latest.ImageConfig{
			"default": &latest.ImageConfig{
				Image: ptr.String("nginx"),
			},
		},
	}
	configutil.SetFakeConfig(testConfig)

	// Create fake generated config
	generatedConfig := &generated.Config{
		ActiveConfig: "default",
		Configs: map[string]*generated.CacheConfig{
			"default": &generated.CacheConfig{
				Images: map[string]*generated.ImageCache{
					"default": &generated.ImageCache{
						Tag: "1.15", // This will be appended to nginx during deploy
					},
				},
			},
		},
	}
	generated.InitDevSpaceConfig(generatedConfig, "default")

	namespace := "test-kaniko-build"
	imageName := "testimage"
	buildArgs := make(map[string]*string)
	buildArgsNoPush := "true"
	buildArgs["--no-push"] = &buildArgsNoPush
	imageConfig := &latest.ImageConfig{
		Build: &latest.BuildConfig{
			Kaniko: &latest.KanikoConfig{
				Namespace: &namespace,
				Options: &latest.BuildOptions{
					BuildArgs: &buildArgs,
				},
			},
		},
		Image: &imageName,
	}

	// Create the fake client.
	kubeClient := fake.NewSimpleClientset()

	dockerClient, err := docker.NewClient(testConfig, true, log.GetInstance())
	if err != nil {
		t.Fatalf("Error creating docker client: %v", err)
	}

	builder, err := NewBuilder(testConfig, dockerClient, kubeClient, "", imageConfig, "v1", true, log.GetInstance())
	if err != nil {
		t.Fatalf("Error creating new kaniko builder: %v", err)
	}

	// 3. Create test namespace test-kaniko-build
	_, err = kubeClient.CoreV1().Namespaces().Create(&k8sv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	})
	if err != nil {
		t.Fatalf("Error creating namespace: %v", err)
	}
	//pod := k8sv1.Pod{}
	//kubeClient.CoreV1().Pods(namespace).Create(&pod)
	go func() {
		buildPod, err := kubeClient.CoreV1().Pods(namespace).Get("", metav1.GetOptions{})
		for err != nil {
			time.Sleep(1 * time.Millisecond)
			buildPod, err = kubeClient.CoreV1().Pods(namespace).Get("", metav1.GetOptions{})
		}
		buildPod.Status.InitContainerStatuses = make([]k8sv1.ContainerStatus, 1)
		buildPod.Status.InitContainerStatuses[0] = k8sv1.ContainerStatus{
			State: k8sv1.ContainerState{
				Running: &k8sv1.ContainerStateRunning{},
			},
		}
		kubeClient.CoreV1().Pods(namespace).Update(buildPod)
	}()

	// 4. Build image with kaniko, but don't push it (In kaniko options use "--no-push" as flag)
	entrypoint := make([]*string, 3)

	entrypoint0 := "go"
	entrypoint1 := "run"
	entrypoint2 := "main.go"
	entrypoint[0] = &entrypoint0
	entrypoint[1] = &entrypoint1
	entrypoint[2] = &entrypoint2
	err = builder.BuildImage(".", "Dockerfile", &entrypoint, log.GetInstance())
	if err != nil {
		t.Fatalf("Error building image: %v", err)
	}

	// 5. Delete test namespace
	err = kubeClient.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
	if err != nil {
		t.Fatalf("Error deleting namespace: %v", err)
	}
}

func makeTestProject(dir string) error {
	file, err := os.Create("package.json")
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(`{
  "name": "node-js-sample",
  "version": "0.0.1",
  "description": "A sample Node.js app using Express 4",
  "main": "index.js",
  "scripts": {
    "start": "nodemon index.js"
  },
  "dependencies": {
    "express": "^4.13.3",
    "nodemon": "^1.18.4",
    "request": "^2.88.0"
  },
  "keywords": [
    "node",
    "express"
  ],
  "license": "MIT"
}`))
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	file, err = os.Create("index.js")
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(`var express = require('express');
var request = require('request');
var app = express();

app.get('/', async (req, res) => {
  var body = await new Promise((resolve, reject) => {
    request('http://php/index.php', (err, res, body) => {
      if (err) { 
        reject(err);
        return;
      }

      resolve(body);
    });
  });

  res.send(body);
});

app.listen(3000, function () {
  console.log('Example app listening on port 3000!');
});`))
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	file, err = os.Create("Dockerfile")
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(`FROM node:8.11.4

RUN mkdir /app
WORKDIR /app

COPY package.json .
RUN npm install

COPY . .

CMD ["npm", "start"]`))
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	file, err = os.Create(".dockerignore")
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(`Dockerfile
.devspace/
chart/
node_modules/`))
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	fileInfo, err := os.Lstat(".")
	if err != nil {
		return err
	}
	err = os.Mkdir("kube", fileInfo.Mode())
	if err != nil {
		return err
	}

	file, err = os.Create("kube/deployment.yaml")
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(`apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: devspace
spec:
  replicas: 1
  template:
    metadata:
      labels:
        release: devspace-node
    spec:
      containers:
      - name: node
        image: node`))
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
