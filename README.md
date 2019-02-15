# slim

[![Build-Status](https://api.travis-ci.org/openacid/slim.svg?branch=master)](https://travis-ci.org/openacid/slim)

A collection of unbelievably memory efficient data structure in Golang

# For Users

```
go get github.com/openacid/slim
```


# For Developers

## Smoke-test

-   If you already have your go environment setup, just:

    ```
    go get github.com/openacid/slim

    go test github.com/openacid/slim/array
    ```

-   If you do not, or you have a customized development environment,
    the following is a quick copy-paste:

    ```
    git clone git@github.com:openacid/slim.git
    cd build/pseudo-gopath
    export GOPATH=slim/build/pseudo-gopath

    go test github.com/openacid/slim/array
    ```

    Another reason to have a `pseudo-gopath` in it is that some tool have their
    own way conducting source code tree.
    E.g. [git-worktree](https://git-scm.com/docs/git-worktree)
    checkouts source code into another dir other than the GOPATH work space.


## Update Dependency

Dependencies are tracked by [dep](https://github.com/golang/dep).
All dependencies are kept in `vender/` dir thus you do not need to do anything
to run it.

You need to update dependency only when you bring in new feature with other dependency.

-   Install `dep`

    ```
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
    ```

-   Download dependency

    ```
    dep ensure
    ```

    > dep uses Gopkg.toml Gopkg.lock to track dependency info.
    > 
    > Gopkg.toml Gopkg.lock is created with `dep init`.
    >
    > dep creates a `vender` dir to have all dependency package there.

See more: [dep-install](https://github.com/golang/dep#installation)


## Directory Layout

Reference: [golang-standards-project-layout](https://github.com/golang-standards/project-layout)
