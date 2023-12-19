// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package update

import (
	"os"
	"strings"

	"github.com/go-curses/corelibs/chdirs"

	"github.com/go-curses/cdk/log"
	"github.com/go-curses/corelibs/run"
)

func Update(module *Module, goProxy string) {
	var err error
	if err = chdirs.Push(module.Path); err != nil {
		cwd, _ := os.Getwd()
		log.ErrorF("run.Push error, cwd: %q", cwd)
		module.Err = err
		return
	}
	defer func() { _ = chdirs.Pop() }()

	var environ []string
	for _, line := range os.Environ() {
		if strings.HasPrefix(line, "GOWORK=") || strings.HasPrefix(line, "GOPROXY=") {
			continue
		}
		environ = append(environ, line)
	}
	environ = append(environ, "GOWORK=off", "GOPROXY="+goProxy)

	var status int
	if _, _, _, err = run.With(run.Options{
		Path:    module.Path,
		Name:    "go",
		Argv:    []string{"get", module.Name + "@v" + module.Next.String()},
		Environ: environ,
	}); err != nil {
		log.ErrorF(`"go get" exited with status %d: %v`, status, err)
		module.Err = err
	} else {
		module.Pick = false
		module.Done = true
	}

	return
}