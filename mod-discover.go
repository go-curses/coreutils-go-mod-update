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
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/go-curses/corelibs/chdirs"

	"github.com/go-curses/cdk/log"
	"github.com/go-curses/corelibs/run"
)

const (
	gNilPid = -255
)

var (
	gDiscoveryPid = gNilPid
)

func StopDiscovery() (err error) {
	if gDiscoveryPid > 0 {
		if err = syscall.Kill(gDiscoveryPid, 0); err == nil {
			err = syscall.Kill(gDiscoveryPid, syscall.SIGINT)
		}
		gDiscoveryPid = gNilPid
	}
	return
}

func StartDiscovery(path string, goProxy string) (modules Modules, err error) {

	if gDiscoveryPid > 0 {
		err = fmt.Errorf("existing discover process found: %d", gDiscoveryPid)
		return
	}

	if err = chdirs.Push(path); err != nil {
		cwd, _ := os.Getwd()
		log.ErrorF("run.Push error: cwd=%q; %v", cwd, err)
		return
	}
	defer func() {
		if ee := chdirs.Pop(); ee != nil {
			cwd, _ := os.Getwd()
			log.ErrorF("run.Pop error: cwd=%q; %v", cwd, ee)
		}
	}()

	var environ []string
	for _, line := range os.Environ() {
		if strings.HasPrefix(line, "GOWORK=") || strings.HasPrefix(line, "GOPROXY=") {
			continue
		}
		environ = append(environ, line)
	}
	environ = append(environ, "GOWORK=off", "GOPROXY="+goProxy)

	var await chan bool
	if gDiscoveryPid, await, err = run.Callback(run.Options{
		Path: path,
		Name: "go",
		Argv: []string{
			"list",
			"-u",
			"-mod=readonly",
			"-f",
			"{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}\t{{.Version}}\t{{.Update.Version}}{{end}}",
			"-m",
			"all",
		},
		Environ: environ,
	}, func(line string) {
		if line = strings.TrimSpace(line); line == "" {
			return
		}
		if name, versions, found := strings.Cut(line, "\t"); found {
			versions = strings.TrimSpace(versions)
			if this, next, ok := strings.Cut(versions, "\t"); ok {
				module := NewModule(path, name, this, next)
				modules = append(modules, module)
			} else {
				log.ErrorF("discover error, line missing second tab: %q", line)
			}
		} else {
			log.ErrorF("discover error, line missing first tab: %q", line)
		}
	}, func(line string) {
		log.ErrorF("discover stderr: %v", line)
	}); err != nil {
		return
	}

	<-await
	return
}