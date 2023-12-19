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
	"github.com/go-curses/ctk"
)

func (u *CUpdater) makeAccelmap() (ag ctk.AccelGroup) {
	ag = ctk.NewAccelGroup()
	ag.ConnectByPath(
		"<updater-window>/File/Update",
		"update-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("update-accel called")
			u.requestUpdates()
			return
		},
	)
	ag.ConnectByPath(
		"<updater-window>/File/Rediscover",
		"rediscover-accel",
		func(argv ...interface{}) (handled bool) {
			ag.LogDebug("rediscover-accel called")
			cdk.Go(u.requestDiscovery)
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