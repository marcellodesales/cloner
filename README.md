# cloner

Clones a given github URL in a collection of repos.

# Why cloner?

* You don't need to change to the directory where your github repos are located.
* The base git repo will be based on the host, so they all are on the same place.

# What is cloner?

This is a CLI that makes git clone simpler.

# Design

* **cloner.yaml**: The metadata of the default settings to be overridden.

## Git

* [x] **clone**: Clones a given URL

# cloner.yaml

```yaml
git:
  dockerImage: alpine/git
```

# Development

## Design

The design now is as follows:

```
main -> cmd -> api/module/service -> config/module
```

This 4 layers enables the implementation to be entirely in the api/service layer that depends on the conversion from yaml to the struct.

## Build

You can use Golang locally to build a local executable as follows (MacOS)

```console
GO111MODULE=on CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o dist/cloner-darwin-local main.go
```

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
