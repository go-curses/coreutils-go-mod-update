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
	_ "embed"
	"sort"

	"github.com/urfave/cli/v2"

	clcli "github.com/go-corelibs/cli"
	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/env"
	"github.com/go-curses/cdk/lib/sync"
	"github.com/go-curses/ctk"
)

//go:embed go-mod-update.accelmap
var gAccelMap string

//go:embed go-mod-update.styles
var gStyles string

type CUI struct {
	App ctk.Application

	ContentVBox ctk.VBox
	ActionHBox  ctk.HButtonBox
	Display     cdk.Display
	Window      ctk.Window

	QuitButton     ctk.Button
	UpdateButton   ctk.Button
	DiscoverButton ctk.Button

	ErrorList     ctk.VBox
	ProjectList   ctk.VBox
	StatusLabel   ctk.Label
	StateSpinner  ctk.Spinner
	ContentScroll ctk.ScrolledViewport

	LastError error

	Projects Projects

	paths   []string
	goProxy string
	tidy    bool
	modLock *sync.RWMutex

	state State

	sync.RWMutex
}

func NewUI(name string, usage string, description string, version string, tag string, title string, ttyPath string) (u *CUI) {
	u = &CUI{
		App:     ctk.NewApplication(name, usage, description, version, tag, title, ttyPath),
		modLock: &sync.RWMutex{},
	}
	c := u.App.CLI()
	c.Flags = append(c.Flags,
		&cli.BoolFlag{
			Name:    "direct",
			Usage:   `specify the GOPROXY setting of "direct" (overrides --goproxy)`,
			EnvVars: []string{"GO_MOD_UPDATE_GOPROXY_DIRECT"},
			Aliases: []string{"d"},
		},
		&cli.StringFlag{
			Name:    "goproxy",
			Usage:   "specify the GOPROXY setting to use",
			EnvVars: []string{"GO_MOD_UPDATE_GOPROXY"},
			Aliases: []string{"p"},
			Value:   env.Get("GOPROXY", "https://proxy.golang.org,direct"),
		},
		&cli.BoolFlag{
			Name:    "tidy",
			Usage:   `run "go mod tidy" after updates`,
			EnvVars: []string{"GO_MOD_UPDATE_TIDY"},
			Aliases: []string{"t"},
		},
	)
	clcli.ClearEmptyCategories(c.Flags)
	sort.Sort(cli.FlagsByName(c.Flags))
	u.App.Connect(cdk.SignalStartup, "updater-startup-handler", u.startup)
	u.App.Connect(cdk.SignalShutdown, "updater-shutdown-handler", u.shutdown)
	return
}

func (u *CUI) Run(argv []string) (err error) {
	err = u.App.Run(argv)
	return
}
