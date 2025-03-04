package latest

import (
	"github.com/devspace-cloud/devspace/pkg/devspace/config/versions/config"
	"github.com/devspace-cloud/devspace/pkg/util/ptr"
)

// Version is the current api version
const Version string = "v1beta2"

// GetVersion returns the version
func (c *Config) GetVersion() string {
	return Version
}

// New creates a new config object
func New() config.Config {
	return NewRaw()
}

// NewRaw creates a new config object
func NewRaw() *Config {
	return &Config{
		Version: ptr.String(Version),
		Cluster: &Cluster{},
		Dev:     &DevConfig{},
		Images:  &map[string]*ImageConfig{},
	}
}

// Config defines the configuration
type Config struct {
	Version      *string                  `yaml:"version"`
	Images       *map[string]*ImageConfig `yaml:"images,omitempty"`
	Deployments  *[]*DeploymentConfig     `yaml:"deployments,omitempty"`
	Dev          *DevConfig               `yaml:"dev,omitempty"`
	Dependencies *[]*DependencyConfig     `yaml:"dependencies,omitempty"`
	Hooks        *[]*HookConfig           `yaml:"hooks,omitempty"`
	Cluster      *Cluster                 `yaml:"cluster,omitempty"`
}

// ImageConfig defines the image specification
type ImageConfig struct {
	Image            *string      `yaml:"image"`
	Tag              *string      `yaml:"tag,omitempty"`
	Dockerfile       *string      `yaml:"dockerfile,omitempty"`
	Context          *string      `yaml:"context,omitempty"`
	CreatePullSecret *bool        `yaml:"createPullSecret,omitempty"`
	Build            *BuildConfig `yaml:"build,omitempty"`
}

// BuildConfig defines the build process for an image
type BuildConfig struct {
	Disabled *bool         `yaml:"disabled,omitempty"`
	Docker   *DockerConfig `yaml:"docker,omitempty"`
	Kaniko   *KanikoConfig `yaml:"kaniko,omitempty"`
	Custom   *CustomConfig `yaml:"custom,omitempty"`
}

// DockerConfig tells the DevSpace CLI to build with Docker on Minikube or on localhost
type DockerConfig struct {
	PreferMinikube  *bool         `yaml:"preferMinikube,omitempty"`
	SkipPush        *bool         `yaml:"skipPush,omitempty"`
	DisableFallback *bool         `yaml:"disableFallback,omitempty"`
	Options         *BuildOptions `yaml:"options,omitempty"`
}

// KanikoConfig tells the DevSpace CLI to build with Docker on Minikube or on localhost
type KanikoConfig struct {
	Cache        *bool         `yaml:"cache,omitempty"`
	SnapshotMode *string       `yaml:"snapshotMode,omitempty"`
	Flags        *[]*string    `yaml:"flags,omitempty"`
	Namespace    *string       `yaml:"namespace,omitempty"`
	Insecure     *bool         `yaml:"insecure,omitempty"`
	PullSecret   *string       `yaml:"pullSecret,omitempty"`
	Options      *BuildOptions `yaml:"options,omitempty"`
}

// CustomConfig tells the DevSpace CLI to build with a custom build script
type CustomConfig struct {
	Command   *string    `yaml:"command,omitempty"`
	Args      *[]*string `yaml:"flags,omitempty"`
	ImageFlag *string    `yaml:"imageFlag,omitempty"`
	OnChange  *[]*string `yaml:"onChange,omitempty"`
}

// BuildOptions defines options for building Docker images
type BuildOptions struct {
	Target    *string             `yaml:"target,omitempty"`
	Network   *string             `yaml:"network,omitempty"`
	BuildArgs *map[string]*string `yaml:"buildArgs,omitempty"`
}

// DeploymentConfig defines the configuration how the devspace should be deployed
type DeploymentConfig struct {
	Name      *string          `yaml:"name"`
	Namespace *string          `yaml:"namespace,omitempty"`
	Component *ComponentConfig `yaml:"component,omitempty"`
	Helm      *HelmConfig      `yaml:"helm,omitempty"`
	Kubectl   *KubectlConfig   `yaml:"kubectl,omitempty"`
}

// ComponentConfig holds the component information
type ComponentConfig struct {
	Containers          *[]*ContainerConfig     `yaml:"containers,omitempty"`
	Replicas            *int                    `yaml:"replicas,omitempty"`
	Autoscaling         *AutoScalingConfig      `yaml:"autoScaling,omitempty"`
	RollingUpdate       *RollingUpdateConfig    `yaml:"rollingUpdate,omitempty"`
	Labels              *map[string]*string     `yaml:"labels,omitempty"`
	Annotations         *map[string]*string     `yaml:"annotations,omitempty"`
	Volumes             *[]*VolumeConfig        `yaml:"volumes,omitempty"`
	Service             *ServiceConfig          `yaml:"service,omitempty"`
	ServiceName         *string                 `yaml:"serviceName,omitempty"`
	Ingress             *IngressConfig          `yaml:"ingress,omitempty"`
	PodManagementPolicy *string                 `yaml:"podManagementPolicy,omitempty"`
	PullSecrets         *[]*string              `yaml:"pullSecrets,omitempty"`
	Options             *ComponentConfigOptions `yaml:"options,omitempty"`
}

// ContainerConfig holds the configurations of a container
type ContainerConfig struct {
	Name           *string                         `yaml:"name,omitempty"`
	Image          *string                         `yaml:"image,omitempty"`
	Command        *[]*string                      `yaml:"command,omitempty"`
	Args           *[]*string                      `yaml:"args,omitempty"`
	Env            *[]*map[interface{}]interface{} `yaml:"env,omitempty"`
	VolumeMounts   *[]*VolumeMountConfig           `yaml:"volumeMounts,omitempty"`
	Resources      *map[interface{}]interface{}    `yaml:"resources,omitempty"`
	LivenessProbe  *map[interface{}]interface{}    `yaml:"livenessProbe,omitempty"`
	ReadinessProbe *map[interface{}]interface{}    `yaml:"readinessProbe,omitempty"`
}

// VolumeMountConfig holds the configuration for a specific mount path
type VolumeMountConfig struct {
	ContainerPath *string                  `yaml:"containerPath,omitempty"`
	Volume        *VolumeMountVolumeConfig `yaml:"volume,omitempty"`
}

// VolumeMountVolumeConfig holds the configuration for a specfic mount path volume
type VolumeMountVolumeConfig struct {
	Name     *string `yaml:"name,omitempty"`
	SubPath  *string `yaml:"subPath,omitempty"`
	ReadOnly *bool   `yaml:"readOnly,omitempty"`
}

// AutoScalingConfig holds the autoscaling config of a component
type AutoScalingConfig struct {
	Horizontal *AutoScalingHorizontalConfig `yaml:"horizontal,omitempty"`
}

// AutoScalingHorizontalConfig holds the horizontal autoscaling config of a component
type AutoScalingHorizontalConfig struct {
	MaxReplicas   *int    `yaml:"maxReplicas,omitempty"`
	AverageCPU    *string `yaml:"averageCPU,omitempty"`
	AverageMemory *string `yaml:"averageMemory,omitempty"`
}

// RollingUpdateConfig holds the configuration for rolling updates
type RollingUpdateConfig struct {
	Enabled        *bool   `yaml:"enabled,omitempty"`
	MaxSurge       *string `yaml:"maxSurge,omitempty"`
	MaxUnavailable *string `yaml:"maxUnavailable,omitempty"`
	Partition      *int    `yaml:"partition,omitempty"`
}

// VolumeConfig holds the configuration for a specific volume
type VolumeConfig struct {
	Name        *string                      `yaml:"name,omitempty"`
	Size        *string                      `yaml:"size,omitempty"`
	ConfigMap   *map[interface{}]interface{} `yaml:"configMap,omitempty"`
	Secret      *map[interface{}]interface{} `yaml:"secret,omitempty"`
	Labels      *map[string]*string          `yaml:"labels,omitempty"`
	Annotations *map[string]*string          `yaml:"annotations,omitempty"`
}

// ServiceConfig holds the configuration of a component service
type ServiceConfig struct {
	Name        *string               `yaml:"name,omitempty"`
	Type        *string               `yaml:"type,omitempty"`
	Ports       *[]*ServicePortConfig `yaml:"ports,omitempty"`
	ExternalIPs *[]*string            `yaml:"externalIPs,omitempty"`
	Labels      *map[string]*string   `yaml:"labels,omitempty"`
	Annotations *map[string]*string   `yaml:"annotations,omitempty"`
}

// ServicePortConfig holds the port configuration of a component service
type ServicePortConfig struct {
	Port          *int    `yaml:"port,omitempty"`
	ContainerPort *int    `yaml:"containerPort,omitempty"`
	Protocol      *string `yaml:"protocol,omitempty"`
}

// IngressConfig holds the configuration of a component ingress
type IngressConfig struct {
	Name        *string               `yaml:"name,omitempty"`
	TLS         *string               `yaml:"tls,omitempty"`
	Labels      *map[string]*string   `yaml:"labels,omitempty"`
	Annotations *map[string]*string   `yaml:"annotations,omitempty"`
	Rules       *[]*IngressRuleConfig `yaml:"rules,omitempty"`
}

// IngressRuleConfig holds the port configuration of a component service
type IngressRuleConfig struct {
	Host        *string `yaml:"host,omitempty"`
	ServicePort *int    `yaml:"servicePort,omitempty"`
	Path        *string `yaml:"path,omitempty"`
	TLS         *string `yaml:"tls,omitempty"`
}

// ComponentConfigOptions defines the specific helm options used during deployment of a component
type ComponentConfigOptions struct {
	Wait            *bool   `yaml:"wait,omitempty"`
	Rollback        *bool   `yaml:"rollback,omitempty"`
	Force           *bool   `yaml:"force,omitempty"`
	Timeout         *int64  `yaml:"timeout,omitempty"`
	TillerNamespace *string `yaml:"tillerNamespace,omitempty"`
}

// HelmConfig defines the specific helm options used during deployment
type HelmConfig struct {
	Chart           *ChartConfig                 `yaml:"chart,omitempty"`
	Wait            *bool                        `yaml:"wait,omitempty"`
	Rollback        *bool                        `yaml:"rollback,omitempty"`
	Force           *bool                        `yaml:"force,omitempty"`
	Timeout         *int64                       `yaml:"timeout,omitempty"`
	TillerNamespace *string                      `yaml:"tillerNamespace,omitempty"`
	DevSpaceValues  *bool                        `yaml:"devSpaceValues,omitempty"`
	ValuesFiles     *[]*string                   `yaml:"valuesFiles,omitempty"`
	Values          *map[interface{}]interface{} `yaml:"values,omitempty"`
}

// ChartConfig defines the helm chart options
type ChartConfig struct {
	Name     *string `yaml:"name,omitempty"`
	Version  *string `yaml:"version,omitempty"`
	RepoURL  *string `yaml:"repo,omitempty"`
	Username *string `yaml:"username,omitempty"`
	Password *string `yaml:"password,omitempty"`
}

// KubectlConfig defines the specific kubectl options used during deployment
type KubectlConfig struct {
	CmdPath   *string    `yaml:"cmdPath,omitempty"`
	Manifests *[]*string `yaml:"manifests,omitempty"`
	Kustomize *bool      `yaml:"kustomize,omitempty"`
	Flags     *[]*string `yaml:"flags,omitempty"`
}

// DevConfig defines the devspace deployment
type DevConfig struct {
	OverrideImages *[]*ImageOverrideConfig  `yaml:"overrideImages,omitempty"`
	Terminal       *Terminal                `yaml:"terminal,omitempty"`
	Ports          *[]*PortForwardingConfig `yaml:"ports,omitempty"`
	Sync           *[]*SyncConfig           `yaml:"sync,omitempty"`
	AutoReload     *AutoReloadConfig        `yaml:"autoReload,omitempty"`
	Selectors      *[]*SelectorConfig       `yaml:"selectors,omitempty"`
}

// ImageOverrideConfig holds information about what parts of the image config are overwritten during devspace dev
type ImageOverrideConfig struct {
	Name       *string    `yaml:"name"`
	Entrypoint *[]*string `yaml:"entrypoint,omitempty"`
	Dockerfile *string    `yaml:"dockerfile,omitempty"`
	Context    *string    `yaml:"context,omitempty"`
}

// Terminal describes the terminal options
type Terminal struct {
	Disabled      *bool               `yaml:"disabled,omitempty"`
	Selector      *string             `yaml:"selector,omitempty"`
	LabelSelector *map[string]*string `yaml:"labelSelector,omitempty"`
	Namespace     *string             `yaml:"namespace,omitempty"`
	ContainerName *string             `yaml:"containerName,omitempty"`
	Command       *[]*string          `yaml:"command,omitempty"`
}

// PortForwardingConfig defines the ports for a port forwarding to a DevSpace
type PortForwardingConfig struct {
	Selector      *string             `yaml:"selector,omitempty"`
	Namespace     *string             `yaml:"namespace,omitempty"`
	LabelSelector *map[string]*string `yaml:"labelSelector,omitempty"`
	PortMappings  *[]*PortMapping     `yaml:"forward"`
}

// PortMapping defines the ports for a PortMapping
type PortMapping struct {
	LocalPort   *int    `yaml:"port"`
	RemotePort  *int    `yaml:"remotePort,omitempty"`
	BindAddress *string `yaml:"bindAddress,omitempty"`
}

// SyncConfig defines the paths for a SyncFolder
type SyncConfig struct {
	Selector             *string             `yaml:"selector,omitempty"`
	Namespace            *string             `yaml:"namespace,omitempty"`
	LabelSelector        *map[string]*string `yaml:"labelSelector,omitempty"`
	ContainerName        *string             `yaml:"containerName,omitempty"`
	LocalSubPath         *string             `yaml:"localSubPath,omitempty"`
	ContainerPath        *string             `yaml:"containerPath,omitempty"`
	WaitInitialSync      *bool               `yaml:"waitInitialSync,omitempty"`
	ExcludePaths         *[]string           `yaml:"excludePaths,omitempty"`
	DownloadExcludePaths *[]string           `yaml:"downloadExcludePaths,omitempty"`
	UploadExcludePaths   *[]string           `yaml:"uploadExcludePaths,omitempty"`
	BandwidthLimits      *BandwidthLimits    `yaml:"bandwidthLimits,omitempty"`
}

// BandwidthLimits defines the struct for specifying the sync bandwidth limits
type BandwidthLimits struct {
	Download *int64 `yaml:"download,omitempty"`
	Upload   *int64 `yaml:"upload,omitempty"`
}

// AutoReloadConfig defines the struct for auto reloading devspace with additional paths
type AutoReloadConfig struct {
	Paths       *[]*string `yaml:"paths,omitempty"`
	Deployments *[]*string `yaml:"deployments,omitempty"`
	Images      *[]*string `yaml:"images,omitempty"`
}

// SelectorConfig defines the selectors that belong to the devspace
type SelectorConfig struct {
	Name          *string             `yaml:"name,omitempty"`
	Namespace     *string             `yaml:"namespace,omitempty"`
	LabelSelector *map[string]*string `yaml:"labelSelector"`
	ContainerName *string             `yaml:"containerName,omitempty"`
}

// DependencyConfig defines the devspace dependency
type DependencyConfig struct {
	Source             *SourceConfig `yaml:"source"`
	Config             *string       `yaml:"config"`
	SkipBuild          *bool         `yaml:"skipBuild,omitempty"`
	IgnoreDependencies *bool         `yaml:"ignoreDependencies,omitempty"`
	Namespace          *string       `yaml:"namespace,omitempty"`
}

// SourceConfig defines the dependency source
type SourceConfig struct {
	Git      *string `yaml:"git,omitempty"`
	Branch   *string `yaml:"branch,omitempty"`
	Tag      *string `yaml:"tag,omitempty"`
	Revision *string `yaml:"revision,omitempty"`

	Path *string `yaml:"path,omitempty"`
}

// HookConfig defines a hook
type HookConfig struct {
	Command *string    `yaml:"command"`
	Args    *[]*string `yaml:"args,omitempty"`

	When *HookWhenConfig `yaml:"when,omitempty"`
}

// HookWhenConfig defines when the hook should be executed
type HookWhenConfig struct {
	Before *HookWhenAtConfig `yaml:"before,omitempty"`
	After  *HookWhenAtConfig `yaml:"after,omitempty"`
}

// HookWhenAtConfig defines at which stage the hook should be executed
type HookWhenAtConfig struct {
	Images      *string `yaml:"images,omitempty"`
	Deployments *string `yaml:"deployments,omitempty"`
}

// Cluster is a struct that contains data for a Kubernetes-Cluster
type Cluster struct {
	KubeContext *string `yaml:"kubeContext,omitempty"`
	Namespace   *string `yaml:"namespace,omitempty"`
}
