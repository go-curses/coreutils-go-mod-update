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

package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/go-curses/cdk/log"

	clcli "github.com/go-corelibs/cli"

	"github.com/go-curses/coreutils-go-mod-update/ui"
)

const (
	APP_TAG         = "go-mod-update"
	APP_NAME        = "go-mod-update"
	APP_TITLE       = "go.mod update"
	APP_USAGE       = "go.mod update"
	APP_DESCRIPTION = "command line utility for updating golang dependencies"
)

var (
	BuildVersion = "v0.2.3"
	BuildRelease = "trunk"
)

func init() {
	cli.FlagStringer = clcli.NewFlagStringer().
		PruneDefaultBools(true).
		DetailsOnNewLines(true).
		Make()
}

func main() {
	updater := ui.NewUI(
		APP_NAME,
		APP_USAGE,
		APP_DESCRIPTION,
		BuildVersion+" ("+BuildRelease+")",
		APP_TAG,
		APP_TITLE,
		"/dev/tty",
	)
	appCLI := updater.App.CLI()
	appCLI.UsageText = "go-mod-update [options] [/source/paths ...]"
	appCLI.HideHelpCommand = true
	appCLI.EnableBashCompletion = true
	appCLI.UseShortOptionHandling = true

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Usage:   "display the version",
		Aliases: []string{"v"},
	}
	if err := updater.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
