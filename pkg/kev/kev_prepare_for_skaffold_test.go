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

package kev_test

import (
	"github.com/appvia/kev/pkg/kev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PrepareForSkaffold", func() {
	var (
		skaffoldManifest *kev.SkaffoldManifest
		err              error
	)

	JustBeforeEach(func() {
		skaffoldManifest, err = kev.PrepareForSkaffold([]string{})
	})

	It("generates skaffold config for the project", func() {
		Expect(skaffoldManifest).ToNot(BeNil())
		Expect(err).ToNot(HaveOccurred())
	})

})
