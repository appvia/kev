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

package main

import (
	"os"

	"github.com/appvia/kev/cmd/kev/cmd"
	"github.com/appvia/kev/pkg/kev/log"
	"github.com/spf13/cobra/doc"
)

func main() {
	outputDir := os.Args[1]

	kev := cmd.NewRootCmd()
	err := doc.GenMarkdownTree(kev, outputDir)
	if err != nil {
		log.Fatal(err)
	}
}
