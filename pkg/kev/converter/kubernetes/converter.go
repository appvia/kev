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
	"os"
	"path"

	"github.com/appvia/kube-devx/pkg/kev/log"
	composego "github.com/compose-spec/compose-go/types"
)

const (
	// Name of the converter
	Name                  = "kubernetes"
	singleFileDefaultName = "k8s.yaml"

	// MultiFileSubDir is default output directory name for kubernetes manifests
	MultiFileSubDir = "k8s"
)

// K8s is a native kubernetes manifests converter
type K8s struct{}

// New return a native Kubernetes converter
func New() *K8s {
	return &K8s{}
}

// Render generates outcome
func (c *K8s) Render(singleFile bool, dir, workDir string, projects map[string]*composego.Project, files map[string][]string, rendered map[string][]byte) error {

	for env, project := range projects {
		log.Infof("🖨️  Rendering %s environment", env)

		// @step override output directory if specified
		outDirPath := ""
		if dir != "" {
			// adding env name suffix to the custom directory to differentiate
			outDirPath = path.Join(dir, env)
		} else {
			outDirPath = path.Join(workDir, MultiFileSubDir, env)
		}

		// @step create output directory
		// To generate outcome as a set of separate manifests first must create out directory
		// as Kompose logic checks for this and only will do that for existing directories,
		// otherwise will treat OutFile as regular file and output all manifests to that single file.
		if err := os.MkdirAll(outDirPath, os.ModePerm); err != nil {
			return err
		}

		// @step generate multiple / single file
		outFilePath := ""
		if singleFile {
			outFilePath = path.Join(outDirPath, singleFileDefaultName)
		} else {
			outFilePath = outDirPath
		}

		// @step kubernetes manifests output options
		convertOpts := ConvertOptions{
			InputFiles: files[env],
			OutFile:    outFilePath,
		}

		// @step Get Kubernete transformer that maps compose project to Kubernetes primitives
		k := &Kubernetes{Opt: convertOpts, Project: project}

		// @step Do the transformation
		objects, err := k.Transform()
		if err != nil {
			return err
		}

		// @step Produce objects
		err = PrintList(objects, convertOpts, rendered)
		if err != nil {
			return err
		}
	}

	return nil
}
