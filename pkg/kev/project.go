/**
 * Copyright 2021 Appvia Ltd <info@appvia.io>
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

package kev

import (
	"fmt"

	"github.com/appvia/kev/pkg/kev/terminal"
)

func (p *Project) Init(opts ...Options) {
	var cfg runConfig
	for _, o := range opts {
		o(p, &cfg)
	}
	p.config = cfg
	if p.UI == nil {
		p.UI = terminal.ConsoleUI()
	}
}

func WithComposeSources(c []string) Options {
	return func(project *Project, cfg *runConfig) {
		cfg.composeSources = c
	}
}

func WithEnvs(c []string) Options {
	return func(project *Project, cfg *runConfig) {
		cfg.envs = c
	}
}

func WithSkaffold(c bool) Options {
	return func(project *Project, cfg *runConfig) {
		cfg.skaffold = c
		fmt.Printf("%v", cfg.skaffold)
	}
}

func WithWorkingDir(wd string) Options {
	return func(project *Project, cfg *runConfig) {
		project.workingDir = wd
	}
}

func WithUI(ui terminal.UI) Options {
	return func(project *Project, cfg *runConfig) {
		project.UI = ui
	}
}