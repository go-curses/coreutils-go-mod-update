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
	"time"

	"github.com/go-curses/corelibs/spinner"
	update "github.com/go-curses/coreutils-go-mod-update"
)

func (u *CUpdater) requestUpdates() {
	u.modLock.Lock()
	defer u.modLock.Unlock()

	u.SetState(UpdatingState)

	var s spinner.Spinner
	var idx, moduleCount int
	var overrideMessage string
	var project, previous *CProject

	for _, p := range u.Projects {
		for _, pkg := range p.Packages {
			if pkg.Module.Pick && !pkg.Module.Done {
				moduleCount += 1
			}
		}
	}

	statusUpdater := func(symbol string) {
		if project != nil {
			project.Frame.SetLabel(symbol + " " + project.Name)
			project.Frame.Resize()
		}
		if previous != nil && previous.Path != project.Path {
			previous.UpdateTitle()
			previous = project
		}
		if overrideMessage == "" {
			if moduleCount > 1 {
				u.StatusLabel.SetLabel(fmt.Sprintf("updating... (%d of %d)", idx+1, moduleCount))
			} else {
				u.StatusLabel.SetLabel("updating...")
			}
		} else {
			u.StatusLabel.SetLabel(overrideMessage)
		}
		u.Display.RequestDraw()
		u.Display.RequestShow()
		return
	}

	s = spinner.NewSpinner(spinner.DefaultSymbols, statusUpdater)

	for _, project = range u.Projects {
		var updated bool
		for _, pkg := range project.Packages {
			if updated = pkg.Module.Pick && !pkg.Module.Done; updated {
				statusUpdater(s.String())
				pkg.GoModUpdate()
				pkg.Module.Pick = false
				statusUpdater(s.String())
				idx += 1
			}
		}
		if u.tidy && updated {
			overrideMessage = "go mod tidy: " + project.Name
			statusUpdater(s.String())
			time.Sleep(time.Millisecond * 500)
			if err := update.Tidy(project.Path, u.goProxy); err != nil {
				err = fmt.Errorf("%q error: %v", project.Name, err)
				u.ErrorList.PackStart(u.makeError(err), true, true, 0)
			}
			time.Sleep(time.Millisecond * 500)
			overrideMessage = ""
			statusUpdater(s.String())
		}
		previous = project
	}

	s.Stop()

	u.SetState(IdleState)
}