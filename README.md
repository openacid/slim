<!--
based on the a great readme template
https://gist.github.com/PurpleBooth/109311bb0361f32d87a2
-->

# Slim - space efficient data structures in Golang

[![Build-Status](https://api.travis-ci.org/openacid/slim.svg?branch=master)](https://travis-ci.org/openacid/slim)

Slim provides with in-memory data-structures and coresponding serialization
APIs for persisting them on-disk or tranport.

<!-- TODO list data types -->
<!-- TODO slim.Trie -->
<!-- TODO other data types -->

<!-- TODO toc -->

## Status

-   `1.0.0`(Current):

## Synoposis

<!-- TODO refine this synoposis -->
```go
package main

import "fmt"
import "github.com/openacid/slim/bit"

func main() {
	fmt.Println("Hello", bit.PopCnt64Before(3, 2))
}
```

<!-- ## Performance and benchmark -->

<!-- TODO  -->

<!-- ## FAQ -->

## Getting started

### Install

```sh
go get github.com/openacid/slim
```

All dependency packages for running it are already included in `vendor/` dir.


<!-- TODO add FAQ -->
<!-- TODO add serialization explaination, on-disk data structure etc. -->

### Prerequisites

-   **For users** (who'd like to build cool stuff based on `slim`):

    **You do NOT need to install anything else before using it**.

-   **For contributors** (who'd like to make `slim` better):

    -   `dep`:
        for dependency management.
    -   `protobuf`:
        for re-compiling `*.proto` file if on-disk data structure chnages.

    Max OS X:
    ```sh
    brew install dep
    brew install protobuf
    ```

    On other platforms you can read more:
    [dep-install][],
    [protoc-install][].


## Who are using slim

<span style="display: inline-block; text-align: center;"> <span> ![][baishancloud-favicon] </span> <br/> <span> [baishancloud][] </span> </span>

## Slim internal

### Built With

- [protobuf][] - Define on-disk data-structure and serialization engine.
- [semver][] - For versioning data-structure.
- [dep][] - Dependency Management.

<!-- ### Directory Layout -->

<!-- We follow the: [golang-standards-project-layout][]. -->

<!-- [> TODO read the doc and add more standards <] -->

<!-- -   `vendor/`: dependency packages. -->
<!-- -   `prototype/`: on-disk data-structure. -->
<!-- -   `docs/`: documents about design, trade-off, etc -->
<!-- -   `tools/`: documents about design, trade-off, etc -->
<!-- -   `expamples/`: documents about design, trade-off, etc -->

<!-- Other directories are sub-package. -->


<!-- ### Versioning -->

<!-- We use [SemVer](http://semver.org/) for versioning. -->

<!-- For the versions available, see the [tags on this repository](https://github.com/your/project/tags).  -->

<!-- ### Data structure explained -->
<!-- [> TODO  <] -->

<!-- ## Limitation -->
<!-- [> TODO  <] -->

## Roadmap

-   [x] slimtrie
-   [ ] bitrie: 1 byte-per-key implementation.
-   [ ] balanced bitrie: which gives better worst-case performance.
-   [ ] generalized API as a drop-in replacement for map etc.


## Feedback

**Feedback and Contributions are greatly appreciated**.

At this stage, the maintainers are most interested in feedback centered on:

-   Do you have a real life scenario that the tool supports well, or doesn't support at all?
-   Do any of the APIs fulfill your needs well?

Let us know by filing an issue, describing what you did or wanted to do, what
you expected to happen, and what actually happened:

-   [bug-report][]
-   [improve-document][]
-   [feature-request][]

Or other type of [issue][new-issue].

<!-- ## Contributing -->
<!-- The maintainers actively manage the issues list, and try to highlight issues -->
<!-- suitable for newcomers. -->

<!-- [> TODO dep CONTRIBUTING <] -->
<!-- The project follows the typical GitHub pull request model. See CONTRIBUTING.md for more details. -->

<!-- Before starting any work, please either comment on an existing issue, -->
<!-- or file a new one. -->

<!-- [> TODO  <] -->
<!-- Please read [CONTRIBUTING.md][] -->
<!-- for details on our code of conduct, and the process for submitting pull requests to us. -->
<!-- https://gist.github.com/PurpleBooth/b24679402957c63ec426 -->


<!-- ### Code style -->

<!-- ### Tool chain -->

<!-- ### Customized install -->

<!-- Alternatively, if you have a customized go develop environment, you could also -->
<!-- clone it: -->

<!-- ```sh -->
<!-- git clone git@github.com:openacid/slim.git -->
<!-- ``` -->

<!-- As a final step you'd like have a test to see if everything goes well: -->

<!-- ```sh -->
<!-- cd path/to/slim/build/pseudo-gopath -->
<!-- export GOPATH=$(pwd) -->
<!-- go test github.com/openacid/slim/array -->
<!-- ``` -->

<!-- Another reason to have a `pseudo-gopath` in it is that some tool have their -->
<!-- own way conducting source code tree. -->
<!-- E.g. [git-worktree](https://git-scm.com/docs/git-worktree) -->
<!-- checkouts source code into another dir other than the GOPATH work space. -->

<!-- ## Update dependency -->

<!-- Dependencies are tracked by [dep](https://github.com/golang/dep). -->
<!-- All dependencies are kept in `vendor/` dir thus you do not need to do anything -->
<!-- to run it. -->

<!-- You need to update dependency only when you bring in new feature with other dependency. -->

<!-- -   Install `dep` -->

<!--     ``` -->
<!--     curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh -->
<!--     ``` -->

<!-- -   Download dependency -->

<!--     ``` -->
<!--     dep ensure -->
<!--     ``` -->

<!--     > dep uses Gopkg.toml Gopkg.lock to track dependency info. -->
<!--     >  -->
<!--     > Gopkg.toml Gopkg.lock is created with `dep init`. -->
<!--     > -->
<!--     > dep creates a `vendor` dir to have all dependency package there. -->

<!-- See more: [dep-install][] -->


## Authors

<!-- ordered by unicode of author's name -->
<!-- leave 3 to 5 major jobs you have done in this project -->

- ![][刘保海-img-sml] **[刘保海][]** *what you do*
- ![][吴义谱-img-sml] **[吴义谱][]** *compacted-array,*
- ![][张炎泼-img-sml] **[张炎泼][]** *slimtrie design,*
- ![][李文博-img-sml] **[李文博][]** *trie-compressing,*
- ![][李树龙-img-sml] **[李树龙][]** *serialization,*


See also the list of [contributors][] who participated in this project.


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

<!-- ## Acknowledgments -->

<!-- [> TODO  <] -->
<!-- - Hat tip to anyone whose code was used -->

<!-- - Inspiration -->
<!--     patricial tree -->
<!--     fusion tree -->
<!--     critic trie -->
<!-- - etc -->

<!-- links -->

<!-- Bio -->

[刘保海]: https://github.com/liubaohai
[吴义谱]: https://github.com/pengsven
[张炎泼]: https://github.com/drmingdrmer
[李文博]: https://github.com/wenbobuaa
[李树龙]: https://github.com/lishulong

<!-- avatar -->

[刘保海-img-sml]: https://avatars1.githubusercontent.com/u/26271283?s=36&v=4
[吴义谱-img-sml]: https://avatars3.githubusercontent.com/u/6927668?s=36&v=4
[张炎泼-img-sml]: https://avatars3.githubusercontent.com/u/44069?s=36&v=4
[李文博-img-sml]: https://avatars1.githubusercontent.com/u/11748387?s=36&v=4
[李树龙-img-sml]: https://avatars2.githubusercontent.com/u/13903162?s=36&v=4

[contributors]: https://github.com/openacid/slim/contributors

[dep]: https://github.com/golang/dep
[protobuf]: https://github.com/protocolbuffers/protobuf
[semver]: http://semver.org/

[protoc-install]: http://google.github.io/proto-lens/installing-protoc.html
[dep-install]: https://github.com/golang/dep#installation

[CONTRIBUTING.md]: CONTRIBUTING.md

[baishancloud]: http://www.baishancdnx.com
[baishancloud-favicon]: http://www.baishancdnx.com/public/favicon.ico
[golang-standards-project-layout]: https://github.com/golang-standards/project-layout

<!-- issue links -->

[bug-report]:       https://github.com/openacid/slim/issues/new?labels=bug&template=bug_report.md
[improve-document]: https://github.com/openacid/slim/issues/new?labels=doc&template=doc_improve.md
[feature-request]:  https://github.com/openacid/slim/issues/new?labels=feature&template=feature_request.md

[new-issue]: https://github.com/openacid/slim/issues/new/choose
