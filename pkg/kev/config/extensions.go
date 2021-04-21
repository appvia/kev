package config

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	composego "github.com/compose-spec/compose-go/types"
	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

const K8SExtensionKey = "x-k8s"

// ExtensionRoot represents the root of the docker-compose extensions
type ExtensionRoot struct {
	K8S K8SConfiguration `yaml:"x-k8s"`
}

// K8SConfiguration represents the root of the k8s specific fields supported by kev.
type K8SConfiguration struct {
	Enabled  bool     `yaml:"enabled,omitempty"`
	Workload Workload `yaml:"workload" validate:"required,dive"`
	Service  Service  `yaml:"service,omitempty"`
}

func (k K8SConfiguration) ToMap() (map[string]interface{}, error) {
	bs, err := yaml.Marshal(k)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = yaml.Unmarshal(bs, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (k K8SConfiguration) Merge(other K8SConfiguration) (K8SConfiguration, error) {
	k8s := k

	if err := mergo.Merge(&k8s, other, mergo.WithOverride); err != nil {
		return K8SConfiguration{}, err
	}

	return k8s, nil
}

func (k K8SConfiguration) Validate() error {
	err := validator.New().Struct(k.Workload)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			if e.Tag() == "required" {
				return fmt.Errorf("%s is required", e.StructNamespace())
			}
		}

		return errors.New(validationErrors[0].Error())
	}

	return nil
}

// DefaultK8SConfig returns a K8SServiceConfig with all the defaults set into it.
func DefaultK8SConfig() K8SConfiguration {
	return K8SConfiguration{
		Enabled: DefaultServiceEnabled,
		Workload: Workload{
			Type:           DefaultWorkload,
			LivenessProbe:  DefaultLivenessProbe(),
			ReadinessProbe: DefaultReadinessProbe(),
			Replicas:       1,
		},
		Service: Service{
			Type: "None",
		},
	}
}

type k8sConfigOptions struct {
	requireExtensions bool
	disableValidation bool
}

// K8SCfgOption will modify parsing behaviour of the x-k8s extension.
type K8SCfgOption func(*k8sConfigOptions)

func DisableValidation() K8SCfgOption {
	return func(kco *k8sConfigOptions) {
		kco.disableValidation = true
	}
}

// RequireExtensions will ensure that x-k8s is present and that it is validated.
func RequireExtensions() K8SCfgOption {
	return func(kco *k8sConfigOptions) {
		kco.requireExtensions = true
	}
}

func K8SCfgFromCompose(svc *composego.ServiceConfig) (K8SConfiguration, error) {
	var cfg K8SConfiguration

	cfg.Enabled = true
	cfg.Workload.Type = WorkloadTypeFromCompose(svc)
	cfg.Workload.Replicas = WorkloadReplicasFromCompose(svc)
	cfg.Workload.RestartPolicy = WorkloadRestartPolicyFromCompose(svc)
	cfg.Service.Type = NoService
	if svc.Ports != nil {
		cfg.Service.Type = ClusterIPService
	}

	cfg.Workload.LivenessProbe = LivenessProbeFromHealthcheck(svc.HealthCheck)
	cfg.Workload.ReadinessProbe = DefaultReadinessProbe()

	k8sext, err := ParseK8SCfgFromMap(svc.Extensions, DisableValidation())
	if err != nil {
		return K8SConfiguration{}, err
	}

	cfg, err = cfg.Merge(k8sext)
	if err != nil {
		return K8SConfiguration{}, err
	}

	if err := cfg.Validate(); err != nil {
		return K8SConfiguration{}, err
	}

	return cfg, nil
}

func WorkloadRestartPolicyFromCompose(svc *composego.ServiceConfig) string {
	if svc.Deploy == nil || svc.Deploy.RestartPolicy == nil {
		return ""
	}

	switch svc.Deploy.RestartPolicy.Condition {
	case "on-failure":
		return RestartPolicyOnFailure
	case "none":
		return RestartPolicyNever
	default:
		return RestartPolicyAlways
	}
}

func WorkloadReplicasFromCompose(svc *composego.ServiceConfig) int {
	if svc.Deploy == nil || svc.Deploy.Replicas == nil {
		return 1
	}

	return int(*svc.Deploy.Replicas)
}

// TODO: Turn these strings into enums
func WorkloadTypeFromCompose(svc *composego.ServiceConfig) string {
	if svc.Deploy != nil && svc.Deploy.Mode == "global" {
		return DaemonsetWorkload
	}

	if len(svc.Volumes) != 0 {
		return StatefulsetWorkload
	}

	return DeploymentWorkload
}

func LivenessProbeFromHealthcheck(healthcheck *composego.HealthCheckConfig) LivenessProbe {
	var res LivenessProbe

	if healthcheck == nil {
		return DefaultLivenessProbe()
	}

	if healthcheck.Disable {
		res.Type = ProbeTypeNone.String()
		return res
	}

	res.Type = ProbeTypeExec.String()

	test := healthcheck.Test
	if len(test) > 0 && strings.ToLower(test[0]) == "cmd" {
		test = test[1:]
	}
	res.Exec.Command = test

	if healthcheck.Timeout != nil {
		res.Timeout = time.Duration(*healthcheck.Timeout)
	}

	if healthcheck.Retries != nil {
		res.FailureThreashold = int(*healthcheck.Retries)
	}

	if healthcheck.StartPeriod != nil {
		res.InitialDelay = time.Duration(*healthcheck.StartPeriod)
	}

	if healthcheck.Interval != nil {
		res.Period = time.Duration(*healthcheck.Interval)
	}

	return res
}

// ParseK8SCfgFromMap handles the extraction of the k8s-specific extension values from the top level map.
func ParseK8SCfgFromMap(m map[string]interface{}, opts ...K8SCfgOption) (K8SConfiguration, error) {
	var options k8sConfigOptions
	for _, o := range opts {
		o(&options)
	}

	if _, ok := m[K8SExtensionKey]; !ok && !options.requireExtensions {
		return K8SConfiguration{}, nil
	}

	var extensions ExtensionRoot

	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf).Encode(m); err != nil {
		return K8SConfiguration{}, err
	}

	if err := yaml.NewDecoder(&buf).Decode(&extensions); err != nil {
		return K8SConfiguration{}, err
	}

	if options.disableValidation {
		return extensions.K8S, nil
	}

	extensions.K8S.Workload.Type = DefaultWorkload

	if err := extensions.K8S.Validate(); err != nil {
		return K8SConfiguration{}, err
	}

	return extensions.K8S, nil
}

// Workload holds all the workload-related k8s configurations.
type Workload struct {
	Type           string         `yaml:"type,omitempty" validate:"required,oneof=DaemonSet StatefulSet Deployment"`
	Replicas       int            `yaml:"replicas" validate:"required,gt=0"`
	LivenessProbe  LivenessProbe  `yaml:"livenessProbe" validate:"required"`
	ReadinessProbe ReadinessProbe `yaml:"readinessProbe,omitempty"`
	RestartPolicy  string         `yaml:"restartPolicy,omitempty"`
}

// Service will hold the service specific extensions in the future.
// TODO: expand with new properties.
type Service struct {
	Type string `yaml:"type" validate:"required,oneof=none, ClusterIP"`
}
