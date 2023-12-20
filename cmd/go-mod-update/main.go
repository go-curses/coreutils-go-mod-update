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

	"github.com/go-curses/cdk"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"

	"github.com/go-curses/coreutils-go-mod-update/ui/updater"
)

// Build Configuration Flags
// setting these will enable command line flags and their corresponding features
// use `go build -v -ldflags="-X 'main.IncludeLogFullPaths=false'"`
var (
	IncludeProfiling          = "false"
	IncludeLogFile            = "false"
	IncludeLogFormat          = "false"
	IncludeLogFullPaths       = "false"
	IncludeLogLevel           = "false"
	IncludeLogLevels          = "false"
	IncludeLogTimestamps      = "false"
	IncludeLogTimestampFormat = "false"
	IncludeLogOutput          = "false"
)

var (
	BuildVersion = "v0.0.0"
	BuildRelease = "trunk"
)

func init() {
	cdk.Build.Profiling = cstrings.IsTrue(IncludeProfiling)
	cdk.Build.LogFile = cstrings.IsTrue(IncludeLogFile)
	cdk.Build.LogFormat = cstrings.IsTrue(IncludeLogFormat)
	cdk.Build.LogFullPaths = cstrings.IsTrue(IncludeLogFullPaths)
	cdk.Build.LogLevel = cstrings.IsTrue(IncludeLogLevel)
	cdk.Build.LogLevels = cstrings.IsTrue(IncludeLogLevels)
	cdk.Build.LogTimestamps = cstrings.IsTrue(IncludeLogTimestamps)
	cdk.Build.LogTimestampFormat = cstrings.IsTrue(IncludeLogTimestampFormat)
	cdk.Build.LogOutput = cstrings.IsTrue(IncludeLogOutput)
}

func main() {
	updater := updater.NewUpdater(
		"go-mod-update",
		"go.mod updater",
		"command line utility for maintaining golang dependencies",
		BuildVersion+" ("+BuildRelease+")",
		"go-mod-update",
		"go.mod updater",
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