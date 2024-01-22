# go-mod-update

Utility for updating Golang `go.mod` projects.

## INSTALLATION

``` shell
> go install github.com/go-curses/coreutils-go-mod-update/cmd/go-mod-update@latest
```

## DOCUMENTATION

``` shell
> go-mod-update --help
NAME:
   go-mod-update - go.mod update

USAGE:
   go-mod-update [options] [/source/paths ...]

VERSION:
   v0.2.4 (trunk)

DESCRIPTION:
   command line utility for updating golang dependencies

GLOBAL OPTIONS:
   --direct, -d               specify the GOPROXY setting of "direct" (overrides --goproxy) (default: false) [$GO_MOD_UPDATE_GOPROXY_DIRECT]
   --goproxy value, -p value  specify the GOPROXY setting to use (default: "https://proxy.golang.org,direct") [$GO_MOD_UPDATE_GOPROXY]
   --tidy, -t                 run "go mod tidy" after updates (default: false) [$GO_MOD_UPDATE_TIDY]
   --help, -h, --usage        display command-line usage information (default: false)
   --version, -v              display the version (default: false)
```


## LICENSE

```
Copyright 2023  The Go-Curses Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
