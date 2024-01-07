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

func (u *CUpdater) getContentWidth() (w int) {
	if screen := u.Display.Screen(); screen != nil {
		w, _ = screen.Size()
		w -= 2
		if sMin, sMax := u.ContentScroll.GetVScrollbar().GetRange(); sMin != sMax {
			w -= 1
		}
	}
	return
}

func (u *CUpdater) resizeUI() {
	u.ProjectList.Freeze()

	contentWidth := u.getContentWidth()
	listHeight := u.Projects.Height()
	u.ProjectList.SetSizeRequest(contentWidth, listHeight)
	u.ProjectList.Resize()

	errors := u.ErrorList.GetChildren()
	if count := len(errors); count > 0 {
		u.ErrorList.Show()
		u.ErrorList.SetSizeRequest(contentWidth, count)
		u.ErrorList.Resize()
	} else {
		u.ErrorList.Hide()
	}

	u.ProjectList.Thaw()
	u.Window.Resize()
	u.Display.RequestDraw()
	u.Display.RequestShow()
}

func (u *CUpdater) refreshUI() {
	u.Window.Freeze()
	u.Projects.Refresh()
	u.refreshUpdateButton()
	u.Window.Thaw()
	u.resizeUI()
}

func (u *CUpdater) refreshDiscoverButton() {
	if !u.State().Idle() {
		u.DiscoverButton.SetSensitive(false)
		u.DiscoverButton.Hide()
		return
	}
	u.DiscoverButton.SetSensitive(true)
	u.DiscoverButton.Show()
}

func (u *CUpdater) refreshUpdateButton() {
	if u.State().Idle() {
		for _, project := range u.Projects {
			for _, pkg := range project.Packages {
				if hasPicked := pkg.Module.Pick && !pkg.Module.Done; hasPicked {
					u.UpdateButton.Show()
					u.UpdateButton.SetSensitive(hasPicked)
					return
				}
			}
		}
	}
	u.UpdateButton.SetSensitive(false)
	u.UpdateButton.Hide()
}
