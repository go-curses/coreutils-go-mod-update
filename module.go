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

	"github.com/Masterminds/semver/v3"

	"github.com/go-curses/cdk/log"
)

type Module struct {
	Path string
	Name string
	This *semver.Version
	Next *semver.Version
	Pick bool
	Done bool
	Err  error
}

func NewModule(path, name, this, next string) (m *Module) {
	var err error
	var t, n *semver.Version
	if t, err = semver.NewVersion(this); err != nil {
		log.ErrorF("this is not a semver: %q - %v", this, err)
	}
	if n, err = semver.NewVersion(next); err != nil {
		log.ErrorF("next is not a semver: %q - %v", next, err)
	}
	m = &Module{
		Path: path,
		Name: name,
		This: t,
		Next: n,
	}
	return
}

func (m Module) String() (line string) {
	line = fmt.Sprintf("%v: %v -> %v", m.Name, m.This, m.Next)
	return
}

func (m Module) ThisVer() (version string) {
	version = Version(m.This)
	return
}

func (m Module) NextVer() (version string) {
	version = Version(m.Next)
	return
}
