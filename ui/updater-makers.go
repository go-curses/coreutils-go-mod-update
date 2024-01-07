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
	"strings"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paint"
	"github.com/go-curses/ctk"
)

func (u *CUpdater) makeAccelmap() (ag ctk.AccelGroup) {
	ag = ctk.NewAccelGroup()
	ag.ConnectByPath(
		"<updater-window>/File/Update",
		"update-accel",
		func(argv ...interface{}) (handled bool) {
			if u.State().Idle() {
				cdk.Go(u.requestUpdates)
				ag.LogDebug("update-accel called")
			} else {
				ag.LogDebug("update-accel ignored")
			}
			return
		},
	)
	ag.ConnectByPath(
		"<updater-window>/File/Rediscover",
		"rediscover-accel",
		func(argv ...interface{}) (handled bool) {
			if u.State().Idle() {
				cdk.Go(u.requestDiscovery)
				ag.LogDebug("rediscover-accel called")
			} else {
				ag.LogDebug("rediscover-accel ignored")
			}
			return
		},
	)
	ag.ConnectByPath(
		"<updater-window>/File/Quit",
		"quit-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("quit-accel called")
			u.requestQuit()
			return
		},
	)
	return
}

func (u *CUpdater) makeActionButtonBox() ctk.HButtonBox {
	u.ActionHBox = ctk.NewHButtonBox(false, 1)
	u.ActionHBox.Show()
	u.ActionHBox.SetSizeRequest(-1, 1)

	u.StateSpinner = ctk.NewSpinner()
	u.StateSpinner.Hide()
	u.StateSpinner.SetSizeRequest(1, 1)
	u.StateSpinner.SetSpinnerRunes(paint.OrbitDotSpinnerRuneSet...)
	u.ActionHBox.PackStart(u.StateSpinner, false, false, 1)

	u.StatusLabel = u.makeLabel("status-label", "")
	u.StatusLabel.SetCanFocus(false)
	u.StatusLabel.SetCanDefault(false)
	u.ActionHBox.PackStart(u.StatusLabel, true, true, 0)

	sep := ctk.NewSeparator()
	sep.Show()
	u.ActionHBox.PackEnd(sep, true, true, 0)

	u.UpdateButton = u.makeButton("update", "F8", "_Update", "", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		cdk.Go(u.requestUpdates)
		return enums.EVENT_PASS
	})
	u.ActionHBox.PackEnd(u.UpdateButton, false, false, 0)

	u.DiscoverButton = u.makeButton("discover", "F5", "_Reload", "", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		cdk.Go(u.requestDiscovery)
		return enums.EVENT_PASS
	})
	u.ActionHBox.PackEnd(u.DiscoverButton, false, false, 0)

	u.QuitButton = u.makeButton("quit", "F10", "_Quit", "", func(data []interface{}, argv ...interface{}) enums.EventFlag {
		u.requestQuit()
		return enums.EVENT_PASS
	})
	u.QuitButton.SetSensitive(true)
	u.ActionHBox.PackEnd(u.QuitButton, false, false, 0)

	return u.ActionHBox
}

func (u *CUpdater) makeButton(key, fKey, label, tooltip string, actFn cdk.SignalListenerFn) (button ctk.Button) {
	fKey = strings.ToUpper(fKey)
	shortTitle := fmt.Sprintf("%s <%s>", label, fKey)
	button = ctk.NewButtonWithMnemonic(shortTitle)
	button.Show()
	button.SetSizeRequest(-1, 1)
	button.SetSensitive(false)
	if tooltip != "" {
		button.SetHasTooltip(true)
		button.SetTooltipText(tooltip)
	}
	if actFn != nil {
		button.Connect(ctk.SignalActivate, key+"-handler", actFn)
	}
	return
}

func (u *CUpdater) makeLabel(name, text string) (label ctk.Label) {
	label = ctk.NewLabel(text)
	label.Show()
	label.SetName(name)
	label.SetEllipsize(true)
	label.SetSizeRequest(-1, 1)
	label.SetLineWrap(false)
	label.SetLineWrapMode(enums.WRAP_NONE)
	label.SetSingleLineMode(true)
	return
}

func (u *CUpdater) makeError(err error) (label ctk.Label) {
	text := err.Error()
	label = ctk.NewLabel(text)
	label.SetName("module-error")
	label.Show()
	label.SetEllipsize(true)
	label.SetLineWrap(false)
	label.SetLineWrapMode(enums.WRAP_NONE)
	label.SetSingleLineMode(true)
	w, _ := u.Display.Screen().Size()
	label.SetSizeRequest(w-4, 1)
	return
}
