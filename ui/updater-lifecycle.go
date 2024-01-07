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
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/paths"
	"github.com/go-curses/cdk/lib/ptypes"
	"github.com/go-curses/cdk/log"
	"github.com/go-curses/ctk"
	"github.com/go-curses/ctk/lib/enums"

	update "github.com/go-curses/coreutils-go-mod-update"
)

func (u *CUpdater) startupInitChecks(ctx *cli.Context) (event cenums.EventFlag) {

	if ctx.Bool("direct") {
		u.goProxy = "direct"
	} else if ctx.IsSet("goproxy") {
		u.goProxy = ctx.String("goproxy")
	}

	u.tidy = ctx.Bool("tidy")

	if args := ctx.Args().Slice(); len(args) > 0 {
		for _, arg := range args {
			if paths.IsDir(arg) {
				if paths.IsFile(arg + "/go.mod") {
					if path, err := filepath.Abs(arg); err == nil {
						u.paths = append(u.paths, path)
					} else {
						u.LastError = fmt.Errorf("error: %q - %v", arg, err)
						log.Error(u.LastError)
						return cenums.EVENT_STOP
					}
				} else {
					u.LastError = fmt.Errorf("error: go.mod not found in %q", arg)
					log.Error(u.LastError)
					return cenums.EVENT_STOP
				}
			} else {
				u.LastError = fmt.Errorf("error: not a directory - %q", arg)
				log.Error(u.LastError)
				return cenums.EVENT_STOP
			}
		}
	}

	if len(u.paths) == 0 {
		if cwd, err := os.Getwd(); err == nil {
			if paths.IsFile("./go.mod") {
				u.paths = append(u.paths, cwd)
			} else {
				u.LastError = fmt.Errorf("error: go.mod not found in %q", cwd)
				log.Error(u.LastError)
				return cenums.EVENT_STOP
			}
		} else {
			u.LastError = err
			log.Error(err)
			return cenums.EVENT_STOP
		}
	}

	screenSize := ptypes.MakeRectangle(u.Display.Screen().Size())
	if screenSize.W < 60 || screenSize.H < 14 {
		u.LastError = fmt.Errorf("rpl requires a terminal with at least 80x24 dimensions")
		log.Error(u.LastError)
		return cenums.EVENT_STOP
	}
	return cenums.EVENT_PASS
}

func (u *CUpdater) startup(data []interface{}, argv ...interface{}) cenums.EventFlag {
	var ok bool
	if u.App, u.Display, _, _, _, ok = ctk.ArgvApplicationSignalStartup(argv...); ok {

		u.Display.CaptureCtrlC()
		defer u.App.NotifyStartupComplete()

		ctx := u.App.GetContext()
		if e := u.startupInitChecks(ctx); e != cenums.EVENT_PASS {
			return e
		}

		title := fmt.Sprintf("%s - %v", u.App.Title(), u.App.Version())

		ctk.GetAccelMap().LoadFromString(gAccelMap)

		u.Window = ctk.NewWindowWithTitle(title)
		u.Window.SetName("updater-window")
		// u.Window.Show()
		u.Window.SetTheme(WindowTheme)
		//u.Window.SetDecorated(false)
		// _ = u.Window.SetBoolProperty(ctk.PropertyDebug, true)
		if err := u.Window.ImportStylesFromString(gStyles); err != nil {
			u.Window.LogErr(err)
		}

		u.Window.Show()
		u.Window.AddAccelGroup(u.makeAccelmap())

		vbox := u.Window.GetVBox()
		vbox.SetSpacing(0)
		vbox.SetHomogeneous(false)

		u.ContentVBox = ctk.NewVBox(false, 0)
		u.ContentVBox.Show()
		vbox.PackStart(u.ContentVBox, true, true, 0)

		u.ContentScroll = ctk.NewScrolledViewport()
		u.ContentScroll.Show()
		u.ContentScroll.SetPolicy(enums.PolicyNever, enums.PolicyAutomatic)
		u.ContentScroll.SetSizeRequest(-1, -1)
		//_ = u.ContentScroll.SetBoolProperty(ctk.PropertyDebug, true)
		//_ = u.ContentScroll.SetBoolProperty(ctk.PropertyDebugChildren, true)
		u.ContentVBox.PackStart(u.ContentScroll, true, true, 0)

		u.ProjectList = ctk.NewVBox(false, 0)
		u.ProjectList.SetName("project-list")
		u.ProjectList.Show()
		u.ContentScroll.Add(u.ProjectList)
		//_ = u.ProjectList.SetBoolProperty(ctk.PropertyDebug, true)
		//_ = u.ProjectList.SetBoolProperty(ctk.PropertyDebugChildren, true)

		u.ErrorList = ctk.NewVBox(true, 0)
		u.ErrorList.SetName("error-list")
		u.ErrorList.Hide()
		//_ = u.ErrorList.SetBoolProperty(ctk.PropertyDebug, true)
		u.ContentVBox.PackStart(u.ErrorList, false, false, 0)

		vbox.PackEnd(u.makeActionButtonBox(), false, true, 0)

		u.refreshUI()
		u.Window.Resize()

		cdk.Go(u.requestDiscovery)

		u.Display.Connect(cdk.SignalEventResize, "display-resize-handler", func(data []interface{}, argv ...interface{}) cenums.EventFlag {
			if u.App.StartupCompleted() {
				u.refreshUI()
			}
			return cenums.EVENT_PASS
		})
		return cenums.EVENT_PASS
	}
	return cenums.EVENT_STOP
}

func (u *CUpdater) shutdown(_ []interface{}, _ ...interface{}) cenums.EventFlag {
	_ = update.StopDiscovery()
	if u.LastError != nil {
		fmt.Printf("%v\n", u.LastError)
		log.InfoF("> exiting (with error)")
	} else {
		log.InfoF("> exiting (without error)")
	}
	return cenums.EVENT_STOP
}
