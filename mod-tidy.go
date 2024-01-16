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

func Tidy(path string, goProxy string) (err error) {
	if err = chdirs.Push(path); err != nil {
		cwd, _ := os.Getwd()
		log.ErrorF("run.Push error, cwd: %q", cwd)
		return
	}
	defer func() { _ = chdirs.Pop() }()

	var status int
	if _, _, _, err = run.With(run.Options{
		Path:    path,
		Name:    "go",
		Argv:    []string{"mod", "tidy"},
		Environ: goWorkProxyEnviron(goProxy),
	}); err != nil {
		log.ErrorF(`"go mod tidy" exited with status %d: %v`, status, err)
	}

	return
}
