<!-- template:define:options
{
  "nodescription": true
}
-->
![logo](https://liam.sh/-/gh/svg/lrstanley/helm-version-updater?icon=simple-icons%3Ahelm&icon.height=80&bg=topography&bgcolor=rgba(2%2C+0%2C+26%2C+1)&layout=left)

<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<p align="center">
  <a href="https://github.com/lrstanley/helm-version-updater/releases">
    <img title="Release Downloads" src="https://img.shields.io/github/downloads/lrstanley/helm-version-updater/total?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/helm-version-updater/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/helm-version-updater?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/helm-version-updater/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/helm-version-updater?style=flat-square">
  </a>



  <a href="https://github.com/lrstanley/helm-version-updater/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/actions/workflow/status/lrstanley/helm-version-updater/test.yml?branch=master&label=test&style=flat-square">
  </a>


  <a href="https://codecov.io/gh/lrstanley/helm-version-updater">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/lrstanley/helm-version-updater/master?style=flat-square">
  </a>

  <a href="https://pkg.go.dev/github.com/lrstanley/helm-version-updater">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/lrstanley/helm-version-updater?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/helm-version-updater">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/lrstanley/helm-version-updater?style=flat-square">
  </a>
</p>
<p align="center">
  <a href="https://github.com/lrstanley/helm-version-updater/issues?q=is:open+is:issue+label:bug">
    <img title="Bug reports" src="https://img.shields.io/github/issues/lrstanley/helm-version-updater/bug?label=issues&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/helm-version-updater/issues?q=is:open+is:issue+label:enhancement">
    <img title="Feature requests" src="https://img.shields.io/github/issues/lrstanley/helm-version-updater/enhancement?label=feature%20requests&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/helm-version-updater/pulls">
    <img title="Open Pull Requests" src="https://img.shields.io/github/issues-pr/lrstanley/helm-version-updater?label=prs&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/helm-version-updater/releases">
    <img title="Latest Semver Release" src="https://img.shields.io/github/v/release/lrstanley/helm-version-updater?style=flat-square">
    <img title="Latest Release Date" src="https://img.shields.io/github/release-date/lrstanley/helm-version-updater?label=date&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/helm-version-updater/discussions/new?category=q-a">
    <img title="Ask a Question" src="https://img.shields.io/badge/support-ask_a_question!-blue?style=flat-square">
  </a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/discord-bytecord-blue.svg?style=flat-square" title="Discord Chat"></a>
</p>
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [Why](#grey_question-why)
  - [Usage](#gear-usage)
    - [Build from Source](#build-from-source)
  - [Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [Contributing](#handshake-contributing)
  - [License](#balance_scale-license)
<!-- template:end:toc -->

## :grey_question: Why

Multiple of my helm charts have their `appVersion` synced with docker images. With this
GitHub action, I can easily automatically update the chart versions, even through a PR
workflow.

## :gear: Usage

```yaml
name: helm-version-updater

on:
  schedule:
    - cron: "0 13 * * *" # run once a day.
  workflow_dispatch: {} # be able to trigger manually.

jobs:
  helm-version-updater:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: lrstanley/helm-version-updater@latest
        with:
          check-dir: charts/
```

-------------------------

Full list of supported options specified below.

| input name         | required | default            | description                                                                         |
|--------------------|----------|--------------------|-------------------------------------------------------------------------------------|
| version            | false    | `<latest-version>` | Version of helm-version-updater to use (defaults to the same version as the action) |
| output-file        | false    | `-`                | Output json file containing change set (- for stdout)                               |
| check-dir          | false    | `.`                | Directory to recursively check for ci-config.yaml files                             |
| support-prerelease | false    | `false`            | Support pre-release tags as versions                                                |

-------------------------

Check out the [releases](https://github.com/lrstanley/helm-version-updater/releases)
page for prebuilt versions of the binary.

### Build from Source

Note that you must have [Go](https://golang.org/doc/install) installed (latest is usually best).

    git clone https://github.com/lrstanley/helm-version-updater.git && cd helm-version-updater
    make
    ./helm-version-updater --help

<!-- template:begin:support -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :raising_hand_man: Support & Assistance

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for
     guidelines on ensuring everyone has the best experience interacting with
     the community.
* :raising_hand_man: Take a look at the [support](.github/SUPPORT.md) document on
     guidelines for tips on how to ask the right questions.
* :lady_beetle: For all features/bugs/issues/questions/etc, [head over here](https://github.com/lrstanley/helm-version-updater/issues/new/choose).
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :handshake: Contributing

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for guidelines
     on ensuring everyone has the best experience interacting with the
    community.
* :clipboard: Please review the [contributing](.github/CONTRIBUTING.md) doc for submitting
     issues/a guide on submitting pull requests and helping out.
* :old_key: For anything security related, please review this repositories [security policy](https://github.com/lrstanley/helm-version-updater/security/policy).
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :balance_scale: License

```
MIT License

Copyright (c) 2023 Liam Stanley <me@liamstanley.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

_Also located [here](LICENSE)_
<!-- template:end:license -->
