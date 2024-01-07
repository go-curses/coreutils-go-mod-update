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

type State uint8

const (
	InitState State = iota
	IdleState
	DiscoveryState
	UpdatingState
	EndOfStates
)

func (s State) Valid() (valid bool) {
	valid = s < EndOfStates // State is unsigned, condition of >= 0 is implicit
	return
}

func (s State) Idle() (idle bool) {
	idle = s == InitState || s == IdleState
	return
}

func (u *CUpdater) State() (state State) {
	u.RLock()
	defer u.RUnlock()
	state = u.state
	return
}

func (u *CUpdater) SetState(state State) {
	if !state.Valid() {
		panic("invalid State given")
	}
	u.Lock()
	u.state = state
	u.Unlock()
	switch state {
	case InitState, IdleState:
		u.setIdleState()
	case DiscoveryState:
		u.setDiscoveryState()
	case UpdatingState:
		u.setUpdatingState()
	}
	u.refreshUI()
}

func (u *CUpdater) setIdleState() {
	u.refreshUpdateButton()
	u.refreshDiscoverButton()
	u.StateSpinner.StopSpinning()
	u.StateSpinner.Hide()
	u.StatusLabel.SetLabel("")
	u.StatusLabel.Hide()
	u.ActionHBox.Resize()
	u.Display.RequestDraw()
	u.Display.RequestShow()
}

func (u *CUpdater) setDiscoveryState() {
	u.refreshUpdateButton()
	u.refreshDiscoverButton()
	u.StateSpinner.Show()
	u.StateSpinner.StartSpinning()
	u.StatusLabel.Show()
	u.StatusLabel.SetLabel("discovering...")
	u.Display.RequestDraw()
	u.Display.RequestShow()
}

func (u *CUpdater) setUpdatingState() {
	u.refreshUpdateButton()
	u.refreshDiscoverButton()
	u.StateSpinner.Show()
	u.StateSpinner.StartSpinning()
	u.StatusLabel.Show()
	u.StatusLabel.SetLabel("updating...")
	u.Display.RequestDraw()
	u.Display.RequestShow()
}
