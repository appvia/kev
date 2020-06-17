/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

// Contains structs and logic around config.yaml content

// Workload defines app default Kubernetes workload parameters.
type Workload struct {
	// image-pull-policy, Default: IfNotPresent. Possible options: IfNotPresent / Always.
	ImagePullPolicy string `yaml:"image-pull-policy,omitempty" json:"image_pull_policy,omitempty" default:"IfNotPresent"`
	// image-pull-secret, Default: nil - don't use private registry pull secret.
	ImagePullSecret string `yaml:"image-pull-secret,omitempty" json:"image_pull_secret,omitempty" default:""`
	// restart, Default: Always. Possible options: Always / OnFailure / Never.
	Restart string `yaml:",omitempty" json:"restart,omitempty" default:"Always"`
	// service-account-name, Default: default.
	ServiceAccountName string `yaml:"service-account-name,omitempty" json:"service_account_name,omitempty" default:""`
	// security-context-run-as-user, Default: nil - don't set.
	SecurityContextRunAsUser string `yaml:"security-context-run-as-user,omitempty" json:"security_context_run_as_user,omitempty" default:""`
	// security-context-run-as-group, Default: nil - don't set.
	SecurityContextRunAsGroup string `yaml:"security-context-run-as-group,omitempty" json:"security_context_run_as_group,omitempty" default:""`
	// security-context-fs-group, Default: nil - don't set.
	SecurityContextFsGroup string `yaml:"security-context-fs-group,omitempty" json:"security_context_fs_group,omitempty" default:""`
	// type, Default: deployment. Possible option: pod | deployment | statefulset | daemonset | job.
	Type string `yaml:",omitempty" json:"type,omitempty" default:"deployment"`
	// replicas, Default: 1. Number of replicas per workload.
	Replicas *uint64 `yaml:",omitempty" json:",omitempty" default:"1"`
	// rolling-update-max-surge, Default: 1. Maximum number of containers to be updated at a time.
	RollingUpdateMaxSurge *uint64 `yaml:"rolling-update-max-surge,omitempty" json:"rolling_update_max_surge,omitempty" default:"1"`
	// cpu, Default: 0.1. CPU request per workload.
	CPU string `yaml:",omitempty" json:"cpu,omitempty" default:"0.1"`
	// memory, Default: 50M. Memory request per workload.
	Memory string `yaml:",omitempty" json:"memory,omitempty" default:"50M"`
	// max-cpu, Default: 0.2. CPU limit per workload.
	MaxCPU string `yaml:"max-cpu,omitempty" json:"max_cpu,omitempty" default:"0.2"`
	// max-memory, Default: 100M. Memory limit per workload.
	MaxMemory string `yaml:"max-memory,omitempty" json:"max_memory,omitempty" default:"100M"`
}

// Service defines app default component K8s service parameters.
type Service struct {
	// type, Default: none (no service). Possible options: none | headless | clusterip | nodeport | loadbalancer.
	Type string `yaml:",omitempty" json:"type,omitempty" default:"none"`
	// nodeport, Default: nil. Only taken into account when working with service.type: nodeport
	Nodeport uint32 `yaml:",omitempty" json:"node_port,omitempty" default:"0"`
	// expose, Default: false (no ingress). Possible options: false | true | domain.com,otherdomain.com (comma separated domain names). When true / domain(s) - it'll set ingress object.
	Expose interface{} `yaml:",omitempty" json:"expose,omitempty" default:"false"`
	// tls-secret, Default: nil (no tls). Secret name where certs will be loaded from.
	TLSSecret string `yaml:"tls-secret,omitempty" json:"tls_secret,omitempty" default:""`
}

// Volume defines individual volume configuration
type Volume struct {
	// volume name
	Name string `yaml:",omitempty" json:"name,omitempty"`
	// volume class
	Class string `yaml:",omitempty" json:"class,omitempty" default:"default"`
	// volume size
	Size string `yaml:",omitempty" json:"size,omitempty" default:"1Gi"`
}

// Component defines configuration for specific compose service
type Component struct {
	// Compose service name
	Name     string   `yaml:",omitempty" json:"name,omitempty"`
	Workload Workload `yaml:",omitempty" json:"workload,omitempty"`
	Service  Service  `yaml:",omitempty" json:"service,omitempty"`
	// Environment variable value formats:
	// - secret.{secret-name}.{secret-key} # Refer to the value stored in a secret key
	// - config.{config-name}.{config-key} # Refer to the value stored in a configmap key
	// - literal-value                     # Use literal value for Env variable
	Environment map[string]string `yaml:",omitempty" json:"environment,omitempty"`
}

// Config definition
type Config struct {
	// App name
	Name string `yaml:",omitempty" json:"name,omitempty"`
	// App description
	Description string `yaml:",omitempty" json:"description,omitempty"`
	// Defines app default Kubernetes workload parameters.
	Workload Workload `yaml:",omitempty" json:"workload,omitempty"`
	// Defines app default component K8s service parameters.
	Service Service `yaml:",omitempty" json:"service,omitempty"`
	// Control volumes defined in compose file by specifing storage class and size.
	Volumes map[string]Volume `yaml:",omitempty" json:"volumes,omitempty"`
	// Map of defined compose services
	Components map[string]Component `yaml:",omitempty,inline" json:"components,omitempty,inline"`
}