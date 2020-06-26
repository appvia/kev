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

package bootstrap

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/appvia/kube-devx/pkg/kev/app"
	"github.com/appvia/kube-devx/pkg/kev/config"
	"github.com/appvia/kube-devx/pkg/kev/transform"
	"github.com/compose-spec/compose-go/loader"
	compose "github.com/compose-spec/compose-go/types"
	"github.com/goccy/go-yaml"
)

// NewApp creates a new Definition using
// provided docker compose files and app root
func NewApp(root string, composeFiles, envs []string) (*app.Definition, error) {
	baseCompose, err := loadAndParse(composeFiles)
	if err != nil {
		return nil, err
	}

	composeData, err := yaml.Marshal(baseCompose)
	if err != nil {
		return nil, err
	}

	composeData, err = transform.AugmentOrAddDeploy(composeData)
	if err != nil {
		return nil, err
	}

	composeData, err = transform.HealthCheckBase(composeData)
	if err != nil {
		return nil, err
	}

	composeData, err = transform.ExternaliseSecrets(composeData)
	if err != nil {
		return nil, err
	}

	composeData, err = transform.ExternaliseConfigs(composeData)
	if err != nil {
		return nil, err
	}

	inferred, err := config.Infer(composeData)
	if err != nil {
		return nil, err
	}

	return app.Init(root, inferred.ComposeWithPlaceholders, inferred.AppConfig, envs)
}

func loadAndParse(paths []string) (*compose.Config, error) {
	var configFiles []compose.ConfigFile
	envVars := map[string]string{}

	for _, p := range paths {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}

		parsed, err := loader.ParseYAML(b)
		if err != nil {
			return nil, err
		}
		mineHostEnvVars(b, envVars)
		configFiles = append(configFiles, compose.ConfigFile{Filename: p, Config: parsed})
	}

	return loader.Load(compose.ConfigDetails{
		WorkingDir:  path.Dir(paths[0]),
		ConfigFiles: configFiles,
		Environment: envVars,
	})
}

func mineHostEnvVars(data []byte, target map[string]string) {
	pattern := regexp.MustCompile(`\$\{(.*)\}`)
	found := pattern.FindAllSubmatch(data, -1)
	for _, f := range found {
		envVar := string(f[1])
		val := os.Getenv(envVar)
		// ensures an env var like PORT: ${PORT} is not ignored post parse and load
		if len(val) < 1 {
			val = " "
		}
		target[envVar] = val
	}
}
