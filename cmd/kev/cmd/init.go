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

package cmd

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/appvia/kev/pkg/kev"
	"github.com/spf13/cobra"
)

var initLongDesc = `Tracks compose sources & creates deployment environments.

Examples:

  ### Initialise kev.yaml with root docker-compose.yml and override file tracking. Adds the default dev deployment environment.
  $ kev init

  ### Use an alternate docker-compose.yml file.
  $ kev init -f docker-compose.dev.yaml

  ### Use multiple alternate docker-compose.yml files.
  $ kev init -f docker-compose.alternate.yaml -f docker-compose.other.yaml

  ### Use a specified environment.
  $ kev init -e staging

  ### Use multiple specified environments.
  $ kev init -e staging`

var initCmd = &cobra.Command{
	Use:      "init",
	Short:    "Tracks compose sources & creates deployment environments.",
	Long:     initLongDesc,
	RunE:     runInitCmd,
	PostRunE: runDetectSecretsCmd,
}

type skippableFile struct {
	FileName string
	Skipped  bool
	Updated  bool
}

func init() {
	flags := initCmd.Flags()
	flags.SortFlags = false

	flags.StringSliceP(
		"file",
		"f",
		[]string{},
		"Specify an alternate compose file\n(default: docker-compose.yml or docker-compose.yaml)",
	)

	flags.StringSliceP(
		"environment",
		"e",
		[]string{},
		"Specify a deployment environment\n(default: dev)",
	)

	flags.BoolP("skaffold", "s", false, "prepare the project for Skaffold")

	rootCmd.AddCommand(initCmd)
}

func runInitCmd(cmd *cobra.Command, _ []string) error {
	cmdName := "Init"

	files, _ := cmd.Flags().GetStringSlice("file")
	envs, _ := cmd.Flags().GetStringSlice("environment")
	skaffold, _ := cmd.Flags().GetBool("skaffold")

	workingDirAbs, err := os.Getwd()
	if err != nil {
		return displayError(cmdName, err)
	}

	manifestPath := path.Join(workingDirAbs, kev.ManifestName)
	if manifestExistsForPath(manifestPath) {
		err := fmt.Errorf("kev.yaml already exists at: %s", manifestPath)
		return displayError(cmdName, err)
	}

	manifest, err := kev.Init(files, envs, ".")
	if err != nil {
		return displayError(cmdName, err)
	}

	var results []skippableFile

	for _, environment := range manifest.Environments {
		envPath := path.Join(workingDirAbs, environment.File)

		if err := writeTo(envPath, environment); err != nil {
			return displayError(cmdName, err)
		}

		results = append(results, skippableFile{
			FileName: environment.File,
		})
	}

	if skaffold {
		skaffoldManifestPath := path.Join(workingDirAbs, kev.SkaffoldFileName)

		// set skaffold path in kev manifest
		manifest.Skaffold = kev.SkaffoldFileName

		if err := createOrUpdateSkaffoldManifest(skaffoldManifestPath, envs, &results); err != nil {
			return displayError(cmdName, err)
		}
	}

	if err := writeTo(manifestPath, manifest); err != nil {
		return displayError(cmdName, err)
	}

	results = append([]skippableFile{{
		FileName: kev.ManifestName,
	}}, results...)

	displayInitSuccess(os.Stdout, results)

	return nil
}

func manifestExistsForPath(manifestPath string) bool {
	_, err := os.Stat(manifestPath)
	return err == nil
}

func writeTo(filePath string, w io.WriterTo) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	if _, err := w.WriteTo(file); err != nil {
		return err
	}
	return file.Close()
}

func createOrUpdateSkaffoldManifest(path string, envs []string, results *[]skippableFile) error {
	cmdName := "Init"

	if manifestExistsForPath(path) {
		// skaffold manifest already present - add additional profiles to it!
		// Note: kev will skip profiles with names matching those of existing
		// profile names defined in skaffold to avoid profile "hijack".

		updatedSkaffold, err := kev.AddProfiles(path, envs, true)
		if err != nil {
			return displayError(cmdName, err)
		}
		if err := writeTo(path, updatedSkaffold); err != nil {
			return displayError(cmdName, err)
		}

		*results = append(*results, skippableFile{
			FileName: kev.SkaffoldFileName,
			Updated:  true,
		})

	} else {

		skaffoldManifest, err := kev.PrepareForSkaffold(envs)
		if err != nil {
			return displayError(cmdName, err)
		}

		if err := writeTo(kev.SkaffoldFileName, skaffoldManifest); err != nil {
			return displayError(cmdName, err)
		}

		*results = append(*results, skippableFile{
			FileName: kev.SkaffoldFileName,
		})
	}

	return nil
}
