// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"github.com/lrstanley/clix"
	"github.com/sethvargo/go-githubactions"
)

var (
	version = "master"
	commit  = "latest"
	date    = "-"
)

type Flags struct {
	OutputFile        string `long:"output-file"         env:"OUTPUT_FILE" default:"-" description:"output json file containing changeset (- for stdout)"`
	CheckDir          string `long:"check-dir"           env:"CHECK_DIR"   default:"." description:"directory to recursively check for ci-config.yaml files"`
	SupportPreRelease bool   `long:"support-pre-release" env:"SUPPORT_PRERELEASE"      description:"support pre-release versions when upgrading"`
	DryRun            bool   `long:"dry-run"             env:"DRY_RUN"                 description:"do not write update Chart.yaml files, only output changeset"`
}

func main() {
	config := &clix.CLI[Flags]{
		Links: clix.GithubLinks("github.com/lrstanley/helm-version-updater", "master", "https://liam.sh"),
		VersionInfo: &clix.VersionInfo[Flags]{
			Version: version,
			Commit:  commit,
			Date:    date,
		},
	}

	config.Parse(clix.OptDisableLogging | clix.OptDisableGlobalLogger)

	changeset, err := generateChangeset(config.Flags)
	if err != nil {
		githubactions.Fatalf("failed to generate changeset: %s", err)
	}

	if len(changeset.Changes) == 0 {
		githubactions.Noticef("no changes to chart versions detected")
	}

	err = changeset.Write()
	if err != nil {
		githubactions.Fatalf("failed to write changeset: %s", err)
	}
}
