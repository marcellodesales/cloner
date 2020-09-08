# cloner

This is a CLI that makes git clone simpler with an optional config to clone a given URL 
in a location for all of your git projects based on the host.

# Why cloner?

* You don't need to change to the directory where your github repos are located.
* The base git repo will be based on the host, so they all are on the same place.

```
$ tree -L 4 ~/dev/
/Users/marcellodesales/dev/
├── github.com
│   ├── comsysto
│   │   └── redis-locks-with-grafana
│   └── marcellodesales
│       └── alpine-git-hub-docker-image
│   ├── intuit
│   │   ├── intuit-spring-cloud-config-inspector
│   │   ├── intuit-spring-cloud-config-validator
│   │   └── unmazedboot
├── github.google.com
│   ├── docker
│   ├── kubernetes
├── gitlab.com
│   └── supercash
│       └── services
│           └── sms-gateway-service
```

# Config

* `~/.cloner.yaml`
  **git.cloneBaseDir**: as the location for base git host dirs

```yaml
version: 1.0
git:
  cloneBaseDir: ~/dev
```

When the CLI runs, it will create the dirs `git.cloneBaseDir/git.host/git.org/git.repo`

> NOTE: Gitlab's and other hosts may be located in deeper folders.

# Running

```go
$ cloner git --repo https://github.com/comsysto/redis-locks-with-grafana
INFO[0000] Loading the config object 'git' from '/Users/marcellodesales/.cloner.yaml'
INFO[2020-09-08T12:28:11-03:00] Cloning into '/Users/marcellodesales/dev/github.com/comsysto/redis-locks-with-grafana'
Enumerating objects: 233, done.
Total 233 (delta 0), reused 0 (delta 0), pack-reused 233
INFO[2020-09-08T12:28:18-03:00] Done...

$ cloner git --repo https://github.com/comsysto/redis-locks-with-grafana
INFO[0000] Loading the config object 'git' from '/Users/marcellodesales/.cloner.yaml'
ERRO[2020-09-08T12:29:58-03:00] Can't clone repo: clone location '/Users/marcellodesales/dev/github.com/comsysto/redis-locks-with-grafana' exists and it's not empty
ERRO[2020-09-08T12:29:58-03:00] You can specify --force or -f to delete the existing dir and clone again. Make sure there are no panding changes!

$ cloner git --repo https://github.com/comsysto/redis-locks-with-grafana -f
INFO[0000] Loading the config object 'git' from '/Users/marcellodesales/.cloner.yaml'
INFO[2020-09-08T12:30:42-03:00] Forcing clone...
INFO[2020-09-08T12:30:42-03:00] Deleted dir '/Users/marcellodesales/dev/github.com/comsysto/redis-locks-with-grafana'
INFO[2020-09-08T12:30:42-03:00] Cloning into '/Users/marcellodesales/dev/github.com/comsysto/redis-locks-with-grafana'
Enumerating objects: 233, done.
Total 233 (delta 0), reused 0 (delta 0), pack-reused 233
INFO[2020-09-08T12:28:18-03:00] Done...
```

# Development

* Here's how we are doing it!

## Design

The design now is as follows:

```
            +-----+       +-----+      +------+
main +----->+ CMD +------>+ API +----->+ UTIL |
            +--+--+       +-----+      +------+
               ^
               |      +--------+
               +------+ CONFIG |
                      +--------+
```

* `CMD`: CLI implementation using `Cobra`
* `CONFIG`: Yaml configuration using `Viper`
* `API`: Abstraction of services that implement capabilities
* `UTIL`: Utility functions serving the API services

## Binaries

You can use Golang locally to build a local executable as follows (MacOS). The Github Actions CI [build-develop](https://github.com/marcellodesales/cloner/actions?query=workflow%3Abuild-develop) implements the steps below, publishing binary artifacts at every commit under the `develop` branch.

```console
$ make local
```

* The local build will be available at `dist/cloner-darwin-local`

Other binary versions are available at the `Artifacts` section of the `build-develop` Github Actions. If you want to build locally, here's what you need to do:

> Requirement: install the following:
> * `docker`
> * `build tools`

Run the following commands:

```console
make dist
```

You should get the CLI binaries in different platforms-architectures

The CLI will print the help

> **ATTENTION**: Make sure to select the proper binary for your current host!

```console
$ ./dist/cloner-darwin-amd64
```

## Bug Reports & Feature Requests

Please use the [issue tracker](https://github.com/marcellodesales/cloner/issues) to report any bugs or file feature requests.

## Social Coding and Contributing

1. Create an issue to discuss about your idea
2. [Fork it] (https://github.com/marcellodesales/cloner/fork)
3. Create your feature branch (`git checkout -b feature/my-new-feature`)
4. Commit your changes (`git commit -am 'feature x: Add some capability'`)
5. Push to the branch (`git push origin feature/my-new-feature`)
6. Create a new Pull Request
7. Profit! :white_check_mark:

Our Gitflow is as follows:

* `develop`: CI that runs all builds and generate binaries.
  * Target your PR here.
* `master`: CD that creates a new release based on a TAG
  * Performed from PRs from `develop`
  * Publishes latest binary version from `develop`

## Releases

Requires `make` and `docker` for a dockerized build and cross-compilation for Golang. We can run the same exact script on local host or in a CI environemtn. See our Github Action `release-master`, which runs at every PR merge to the `master` branch.  

> ATTENTION: Make sure to have a Github token with write permissions to the repo
>  * https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line

* Initial releases MUST have the code fully pushed to the repo
* Git with versions are created in the format of `yy.mm.#`
  * Where, # is incremented at each release automatically.

```console
$ make release PUBLISH_GITHUB_TOKEN=6100**********dbeb326
fatal: No names found, cannot describe anything.
rm -rf dist
mkdir dist
#docker build -t cloner --build-arg BIN_NAME=cloner --build-arg BIN_VERSION=19.11.1 .
Building next version 19.11.1
BIN_VERSION=19.11.1 docker-compose build --build-arg BIN_VERSION=19.11.1 cli
Building cli
Step 1/24 : FROM golang:1.13.0-stretch as builder
 ---> d68f79d0e22c
Step 2/24 : WORKDIR /build
 ---> Using cache
 ---> 22957f5469a2
Step 3/24 : RUN apt-get update && apt-get install -y file
 ---> Using cache
 ---> 1ec0f47f5d8c
Step 4/24 : COPY go.mod /build/go.mod
 ---> Using cache
 ---> 1799a4b09df6
Step 5/24 : COPY go.sum /build/go.sum
 ---> Using cache
 ---> 358b83dc11ce
Step 6/24 : ENV GO111MODULE=on
 ---> Using cache
 ---> b2eb62f39971
Step 7/24 : RUN go mod download
 ---> Using cache
 ---> d5704e9c7950
Step 8/24 : COPY main.go /build/main.go
 ---> Using cache
 ---> 2e40e04f3424
Step 9/24 : COPY api /build/api
 ---> Using cache
 ---> c31b316588cb
Step 10/24 : COPY cmd /build/cmd
 ---> Using cache
 ---> 7556207f1bf7
Step 11/24 : COPY config /build/config
 ---> Using cache
 ---> 2e44679ae6b9
Step 12/24 : COPY util /build/util
 ---> Using cache
 ---> 0d4c826dc7da
Step 13/24 : ARG BIN_NAME
 ---> Using cache
 ---> 7da993355ee9
Step 14/24 : ARG BIN_VERSION
 ---> Using cache
 ---> 2e380a0cc0e7
Step 15/24 : ENV BIN_NAME ${BIN_NAME:-unknown}
 ---> Using cache
 ---> 2b9e4a4772b8
Step 16/24 : ENV BUILD_VERSION ${BIN_VERSION:-0.1.0}
 ---> Using cache
 ---> a86be0c52bcb
Step 17/24 : ENV PLATFORMS "darwin linux windows"
 ---> Using cache
 ---> 46dcf836a303
Step 18/24 : ENV ARCHS "amd64"
 ---> Using cache
 ---> e53be57a615f
Step 19/24 : COPY .git/ /build/.git/
 ---> f00b7d01fd82
Step 20/24 : RUN export GO_MODULE_FULL_NAME=$(git -C /build/.git remote -v | grep fetch | awk '{print $2}' | sed -En "s/git@//p" | sed -En "s/.git//p" | sed -En "s/:/\//p") &&         export BUILD_GIT_SHA=$(git rev-parse --short HEAD) &&         for GOOS in ${PLATFORMS}; do for GOARCH in ${ARCHS}; do BINARY="${BIN_NAME}-$GOOS-$GOARCH"; if [ $GOOS = "windows" ]; then BINARY="${BINARY}.exe"; fi; export BUILD_TIME="$(date -u +"%Y-%m-%d_%H:%M:%S_GMT")"; echo "Cross-compiling $GO_MODULE_FULL_NAME@$BUILD_GIT_SHA as ${BINARY}@$BUILD_VERSION at $BUILD_TIME"; GO_MODULE_FULL_NAME=$GO_MODULE_FULL_NAME BUILD_GIT_SHA=$BUILD_GIT_SHA BUILD_VERSION=$BUILD_VERSION BUILD_TIME=$BUILD_TIME GO111MODULE=$GO111MODULE CGO_ENABLED=0 GOARCH=$GOARCH GOOS=$GOOS GOPRIVATE=$GOPRIVATE go build -o ${BINARY} -a -ldflags "-X ${GO_MODULE_FULL_NAME}/config.BuildModule=${GO_MODULE_FULL_NAME} -X ${GO_MODULE_FULL_NAME}/config.BuildVersion=${BUILD_VERSION} -X ${GO_MODULE_FULL_NAME}/config.BuildGitSha=${BUILD_GIT_SHA} -X ${GO_MODULE_FULL_NAME}/config.BuildTime=${BUILD_TIME}"; ls -la "/build/${BINARY}"; file "/build/${BINARY}"; chmod +x "/build/${BINARY}"; if [ $GOOS = "linux" ]; then sh -c "/build/${BINARY}"; fi; done; done
 ---> Running in 2453c3ae25a8
Cross-compiling github.com/marcellodesales/cloner@65d3043 as cloner-darwin-amd64@19.11.1 at 2019-11-25_01:18:25_GMT
-rwxr-xr-x 1 root root 12332944 Nov 25 01:18 /build/cloner-darwin-amd64
/build/cloner-darwin-amd64: Mach-O 64-bit x86_64 executable
Cross-compiling github.com/marcellodesales/cloner@65d3043 as cloner-linux-amd64@19.11.1 at 2019-11-25_01:18:43_GMT
-rwxr-xr-x 1 root root 12468409 Nov 25 01:19 /build/cloner-linux-amd64
/build/cloner-linux-amd64: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), statically linked, not stripped
cloner knows how to clone version-control urls by simply making sure
it can place software in specific location designed.

Usage:
  cloner [command]

Available Commands:
  git         Clones a given git repo
  help        Help about any command

Flags:
      --cloner string      Spec file (default is $HOME/cloner.yaml) (default "cloner")
  -h, --help               help for cloner
  -v, --verbosity string   Log level (debug, info, warn, error, fatal, panic (default "info")

Use "cloner [command] --help" for more information about a command.
Cross-compiling github.com/marcellodesales/cloner@65d3043 as cloner-windows-amd64.exe@19.11.1 at 2019-11-25_01:19:00_GMT
-rwxr-xr-x 1 root root 12219392 Nov 25 01:19 /build/cloner-windows-amd64.exe
/build/cloner-windows-amd64.exe: PE32+ executable (console) x86-64 (stripped to external PDB), for MS Windows
Removing intermediate container 2453c3ae25a8
 ---> d0e20c40137d
Step 21/24 : FROM alpine:latest
 ---> 961769676411
Step 22/24 : WORKDIR /root/
 ---> Using cache
 ---> 17b62865a5d4
Step 23/24 : COPY --from=builder /build/${BIN_NAME}* /bin/
 ---> 412bd633d2da
Step 24/24 : ENTRYPOINT /bin/${ARG BIN_NAME}-linux-amd64
 ---> Running in 44b3f771668a
Removing intermediate container 44b3f771668a
 ---> 1715a2e68e3a
Successfully built 1715a2e68e3a
Successfully tagged marcellodesales/cloner:19.11.1
Distribution libraries for version 19.11.1
docker run --rm -ti --entrypoint sh -v /Users/marcellodesales/dev/github.com/marcellodesales/cloner/dist:/bins marcellodesales/cloner:19.11.1 -c "cp /bin/cloner-darwin-amd64 /bins"
docker run --rm -ti --entrypoint sh -v /Users/marcellodesales/dev/github.com/marcellodesales/cloner/dist:/bins marcellodesales/cloner:19.11.1 -c "cp /bin/cloner-linux-amd64 /bins"
docker run --rm -ti --entrypoint sh -v /Users/marcellodesales/dev/github.com/marcellodesales/cloner/dist:/bins marcellodesales/cloner:19.11.1 -c "cp /bin/cloner-windows-amd64.exe /bins"
ls -la /Users/marcellodesales/dev/github.com/marcellodesales/cloner/dist
total 73880
drwxr-xr-x   5 marcellodesales  CORP\Domain Users       160 Nov 24 17:19 .
drwxr-xr-x  19 marcellodesales  CORP\Domain Users       608 Nov 24 17:18 ..
-rwxr-xr-x   1 marcellodesales  CORP\Domain Users  12332944 Nov 24 17:19 cloner-darwin-amd64
-rwxr-xr-x   1 marcellodesales  CORP\Domain Users  12468409 Nov 24 17:19 cloner-linux-amd64
-rwxr-xr-x   1 marcellodesales  CORP\Domain Users  12219392 Nov 24 17:19 cloner-windows-amd64.exe
echo "Releasing next version 19.11.1"
Releasing next version 19.11.1
git tag v19.11.1 || true
git push origin v19.11.1 || true
Enumerating objects: 25, done.
Counting objects: 100% (25/25), done.
Delta compression using up to 8 threads
Compressing objects: 100% (13/13), done.
Writing objects: 100% (14/14), 2.04 KiB | 2.04 MiB/s, done.
Total 14 (delta 8), reused 0 (delta 0)
remote: Resolving deltas: 100% (8/8), completed with 7 local objects.
To github.com:marcellodesales/cloner.git
 * [new tag]         v19.11.1 -> v19.11.1
docker run -ti -e GITHUB_HOST=github.com -e GITHUB_USER=marcellodesales  -e GITHUB_TOKEN=6100758b68072e0570ce0c250a6e398cadbeb326 -e GITHUB_REPOSITORY=marcellodesales/cloner -e HUB_PROTOCOL=https -v /Users/marcellodesales/dev/github.com/marcellodesales/cloner:/git marcellodesales/github-hub release create --prerelease --attach dist/cloner-darwin-amd64 --attach dist/cloner-linux-amd64 --attach dist/cloner-windows-amd64.exe -m "cloner 19.11.1 release" v19.11.1
Attaching release asset `dist/cloner-darwin-amd64'...
Attaching release asset `dist/cloner-linux-amd64'...
Attaching release asset `dist/cloner-windows-amd64.exe'...
https://github.com/marcellodesales/cloner/releases/tag/v19.11.1
```

## License

```
The MIT License (MIT)

Copyright (c) 2019-2020 Marcello de Sales

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```