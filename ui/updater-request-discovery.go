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

package ui

import (
	"fmt"

	"github.com/go-corelibs/spinner"
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/log"

	update "github.com/go-curses/coreutils-go-mod-update"
)

func (u *CUI) requestDiscovery() {
	if !u.State().Idle() {
		log.DebugF("user requesting discovery when updater is not idle")
		return
	}

	u.modLock.Lock()
	defer u.modLock.Unlock()

	u.SetState(DiscoveryState)

	var idx int
	var s spinner.Spinner
	var project, previous *CProject
	projectCount := len(u.paths)

	u.Projects = make([]*CProject, 0)

	for _, child := range u.ProjectList.GetChildren() {
		u.ProjectList.Remove(child)
		child.Destroy()
	}

	s = spinner.New(spinner.DefaultSymbols, func(symbol string) {
		if project != nil {
			project.Frame.SetLabel(symbol + " " + project.Name)
			project.Frame.Resize()
		}
		if previous != nil && previous.Path != project.Path {
			previous.UpdateTitle()
			previous = project
		}
		if projectCount > 1 {
			u.StatusLabel.SetLabel(fmt.Sprintf("discovering... (%d of %d)", idx+1, projectCount))
		}
		u.Display.RequestDraw()
		u.Display.RequestShow()
	})
	cdk.Go(s.Start)

	for _, path := range u.paths {
		p := u.newProject(path)
		u.Projects = append(u.Projects, p)
		u.ProjectList.PackStart(p.Frame, false, false, 0)
	}
	u.ProjectList.Resize()
	u.Display.RequestDraw()
	u.Display.RequestShow()

	for idx, project = range u.Projects {
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
		previous = project
	}

	s.Stop()

	u.SetState(IdleState)
}
