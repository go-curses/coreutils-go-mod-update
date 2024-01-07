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

	"github.com/go-corelibs/chdirs"
	"github.com/go-corelibs/run"
	"github.com/go-curses/cdk/log"
)

func Update(module *Module, goProxy string) {
	var err error
	if err = chdirs.Push(module.Path); err != nil {
		cwd, _ := os.Getwd()
		log.ErrorF("run.Push error, cwd: %q", cwd)
		module.Err = err
		return
	}
	defer func() {
		if ee := chdirs.Pop(); ee != nil {
			cwd, _ := os.Getwd()
			log.ErrorF("run.Pop error: cwd=%q; %v", cwd, ee)
		}
	}()

	var status int
	if _, _, _, err = run.With(run.Options{
		Path:    module.Path,
		Name:    "go",
		Argv:    []string{"get", module.Name + "@v" + module.Next.String()},
		Environ: goWorkProxyEnviron(goProxy),
	}); err != nil {
		log.ErrorF(`"go get" exited with status %d: %v`, status, err)
		module.Err = err
	} else {
		module.Pick = false
		module.Done = true
	}

	return
}
