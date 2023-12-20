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
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/corelibs/spinner"
	"github.com/go-curses/ctk"

	"github.com/go-curses/coreutils-go-mod-update"
)

var (
	ModUnselected = `[ ]`
	ModSelected   = `[🗸]`
	ModDone       = ` 🗸 `
	ModError      = ` 🗴 `
)

type CPackage struct {
	u *CUpdater

	Project *CProject

	Module *update.Module

	HBox    ctk.HBox
	Button  ctk.Button
	Label   ctk.Label
	ThisVer ctk.Label
	NextVer ctk.Label

	Error ctk.Label
}

func (u *CUpdater) newPackage(project *CProject, module *update.Module) (p *CPackage) {
	p = &CPackage{
		u:       u,
		Project: project,
		Module:  module,
	}
	p.HBox = ctk.NewHBox(false, 1)
	p.HBox.Show()
	p.HBox.SetName("module-entry")
	p.HBox.SetSizeRequest(-1, 1)
	//_ = p.HBox.SetBoolProperty(ctk.PropertyDebug, true)

	p.Button = ctk.NewButtonWithLabel(ModUnselected)
	p.Button.Show()
	p.Button.SetName("module-pick")
	p.Button.SetSizeRequest(5, 1)
	if !module.Done && module.Err == nil {
		p.Button.Connect(ctk.SignalActivate, "entry-handler", func(data []interface{}, argv ...interface{}) enums.EventFlag {
			if module.Pick {
				module.Pick = false
				p.Button.SetLabel(ModUnselected)
			} else {
				module.Pick = true
				p.Button.SetLabel(ModSelected)
			}
			u.refreshUpdateButton()
			return enums.EVENT_PASS
		})
	}
	p.UpdateButton()
	p.HBox.PackStart(p.Button, false, false, 0)

	p.Label = u.makeLabel("module-name", module.Name)
	p.Label.SetJustify(enums.JUSTIFY_CENTER)
	p.HBox.PackStart(p.Label, true, true, 0)

	p.ThisVer = u.makeLabel("module-this", module.ThisVer())
	p.ThisVer.SetJustify(enums.JUSTIFY_CENTER)
	p.ThisVer.SetSizeRequest(12, 1)
	p.HBox.PackStart(p.ThisVer, false, false, 0)

	p.NextVer = u.makeLabel("module-next", module.NextVer())
	p.NextVer.SetJustify(enums.JUSTIFY_CENTER)
	p.NextVer.SetSizeRequest(12, 1)
	p.HBox.PackStart(p.NextVer, false, false, 0)

	if module.Err != nil {
		p.Error = u.makeLabel("module-err", " ╰─ "+module.Err.Error())
	}
	return
}

func (p *CPackage) UpdateButton() {
	if !p.u.State().Idle() || p.Module.Err != nil || p.Module.Done {
		p.Button.SetSensitive(false)
	} else {
		p.Button.SetSensitive(true)
	}
	if p.Module.Err != nil {
		p.Button.SetName("module-pick--err")
		p.Button.SetLabel(ModError)
	} else if p.Module.Done {
		p.Button.SetName("module-pick--done")
		p.Button.SetLabel(ModDone)
	} else if p.Module.Pick {
		p.Button.SetName("module-pick")
		p.Button.SetLabel(ModSelected)
	} else {
		p.Button.SetName("module-pick")
		p.Button.SetLabel(ModUnselected)
	}
}

func (p *CPackage) Resize() {
	w, _ := p.u.Display.Screen().Size()
	p.Button.SetSizeRequest(5, 1)
	p.Label.SetSizeRequest(-1, 1)
	p.ThisVer.SetSizeRequest(w/5, 1)
	p.NextVer.SetSizeRequest(w/5, 1)
	p.HBox.Resize()
}

func (p *CPackage) GoModUpdate() {
	p.Button.SetSensitive(false)

	var s spinner.Spinner
	s = spinner.NewSpinner(spinner.DefaultSymbols, func(symbol string) {
		p.Button.SetLabel(" " + symbol + " ")
		p.u.Display.RequestDraw()
		p.u.Display.RequestShow()
	})
	cdk.Go(s.Start)

	update.Update(p.Module, p.u.goProxy)

	s.Stop()
	p.UpdateButton()
	return
}