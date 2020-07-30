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

package kubernetes

import (
	"fmt"
	"os"
	"path"

	"github.com/appvia/kube-devx/pkg/kev/app"
	compose "github.com/appvia/kube-devx/pkg/kev/compose"
)

const (
	// Name of the converter
	Name                  = "kubernetes"
	singleFileDefaultName = "k8s.yaml"
	multiFileSubDir       = ".k8s"
)

// K8s is a native kubernetes manifests converter
type K8s struct{}

// New return a native Kubernetes converter
func New() *K8s {
	return &K8s{}
}

// Render generates outcome
func (c *K8s) Render(singleFile bool, dir string, appDef *app.Definition) error {
	rendered := map[string]app.FileConfig{}
	for env, bc := range appDef.GetMappedBuildInfo() {
		fmt.Printf("\n🖨️  Rendering %s environment\n", env)

		// @step Override output directory if specified
		outDirPath := ""
		if dir != "" {
			// adding env name suffix to the custom directory to differentiate
			outDirPath = path.Join(dir, env)
		} else {
			outDirPath = path.Join(appDef.RootDir(), multiFileSubDir, env)
		}

		// @step Create output directory
		// To generate outcome as a set of separate manifests first must create out directory
		// as Kompose logic checks for this and only will do that for existing directories,
		// otherwise will treat OutFile as regular file and output all manifests to that single file.
		if err := os.MkdirAll(outDirPath, os.ModePerm); err != nil {
			return err
		}

		// @step Generate multiple / single file
		outFilePath := ""
		if singleFile {
			outFilePath = path.Join(outDirPath, singleFileDefaultName)
		} else {
			outFilePath = outDirPath
		}

		// @step Kubernetes manifests output options
		convertOpts := ConvertOptions{
			InputFiles:   []string{bc.Compose.File},
			OutFile:      outFilePath,
			Provider:     Name,
			YAMLIndent:   2,
			GenerateYaml: true,
		}

		// @step Load compose file(s) as compose project
		project, err := compose.LoadProject(convertOpts.InputFiles)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		// @step Get Kubernete transformer that maps compose project to Kubernetes primitives
		k := &Kubernetes{Opt: convertOpts, Project: project}

		// @step Do the transformation
		objects, err := k.Transform()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		// @step Produce objects
		err = PrintList(objects, convertOpts, rendered)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	for _, fileConfig := range rendered {
		appDef.Rendered = append(appDef.Rendered, fileConfig)
	}
	return nil
}
