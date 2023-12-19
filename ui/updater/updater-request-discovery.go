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

package updater

import (
	"fmt"

	"github.com/go-curses/corelibs/spinner"

	"github.com/go-curses/cdk"
	"github.com/go-curses/coreutils-go-mod-update"
)

func (u *CUpdater) requestDiscovery() {

	u.modLock.Lock()
	defer u.modLock.Unlock()
	d := u.Display

	u.startUiChanges("")

	u.Projects = make([]*CProject, 0)
	for _, child := range u.ProjectList.GetChildren() {
		u.ProjectList.Remove(child)
		child.Destroy()
	}

	pathCount := len(u.paths)

	var idx int
	var note string
	var s spinner.Spinner

	s = spinner.NewSpinner(spinner.DefaultSymbols, func(symbol string) {
		var message string
		if pathCount > 1 {
			message = fmt.Sprintf("%s discovering (%d of %d): %s", symbol, idx+1, pathCount, note)
		} else {
			message = fmt.Sprintf("%s discovering: %s", symbol, note)
		}
		u.setStatus(message)
		d.RequestDraw()
		d.RequestShow()
	})
	cdk.Go(s.Start)

	for _, path := range u.paths {
		project := u.newProject(path)
		u.Projects = append(u.Projects, project)
		u.ProjectList.PackStart(project.Frame, false, false, 0)
	}
	u.ProjectList.Resize()
	d.RequestDraw()
	d.RequestShow()

	var project *CProject
	for idx, project = range u.Projects {
		note = "../" + project.Name
		if found, err := update.StartDiscovery(project.Path, u.goProxy); err == nil {
			if count := len(found); count > 0 {
				project.Add(found...)
				u.resizeUI()
			} else {
				project.Resize()
			}
		} else {
			err = fmt.Errorf("%q error: %v", project.Name, err)
			u.ErrorList.PackStart(u.makeError(err), true, true, 0)
		}
		_ = update.StopDiscovery()
	}

	s.Stop()

	u.finishUiChanges()
}