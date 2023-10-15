// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lrstanley/clix"
	"github.com/sethvargo/go-githubactions"
)

var (
	version = "master"
	commit  = "latest"
	date    = "-"

	cli = &clix.CLI[Flags]{
		Links: clix.GithubLinks("github.com/lrstanley/helm-version-updater", "master", "https://liam.sh"),
		VersionInfo: &clix.VersionInfo[Flags]{
			Version: version,
			Commit:  commit,
			Date:    date,
		},
	}
)

type Flags struct {
	OutputFile        string `long:"output-file"         env:"OUTPUT_FILE" default:"-" description:"output json file containing changeset (- for stdout)"`
	CheckDir          string `long:"check-dir"           env:"CHECK_DIR"   default:"." description:"directory to recursively check for ci-config.yaml files"`
	SupportPreRelease bool   `long:"support-pre-release" env:"SUPPORT_PRERELEASE"      description:"support pre-release versions when upgrading"`
}

func main() {
	cli.Parse(clix.OptDisableLogging | clix.OptDisableGlobalLogger)

	var chartDirs []string

	// Recursively find ci-config.yaml files, and parse them.
	err := filepath.Walk(cli.Flags.CheckDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			githubactions.Errorf("failed to walk path %s: %s", path, err)
			return nil
		}

		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}

			return nil
		}

		if info.Name() == "ci-config.yaml" {
			chartDirs = append(chartDirs, filepath.Dir(path))
		}

		return nil
	})
	if len(chartDirs) == 0 && err != nil {
		githubactions.Fatalf("failed to find ci-config.yaml files: %s", err)
	} else if len(chartDirs) == 0 {
		githubactions.Fatalf("failed to find any ci-config.yaml files, validate check-dir")
	}

	var chart *Chart
	var change *Change

	changeset := &ChangeSet{Path: cli.Flags.OutputFile}

	for _, dir := range chartDirs {
		githubactions.Group(fmt.Sprintf("chart-dir:%q", dir))

		// Read the chart.yaml file, and parse it.
		chart, err = ParseChart(filepath.Join(dir, "Chart.yaml"))
		if err != nil {
			githubactions.Fatalf("failed to parse chart.yaml file: %s", err)
		}

		githubactions.Noticef(
			"chart metadata name:%q description:%q version:%q appVersion: %q",
			chart.Name,
			chart.Description,
			chart.OriginalVersion,
			chart.OriginalAppVersion,
		)

		change, err = CheckImageUpdates(chart)
		if err != nil {
			githubactions.Fatalf("failed to check image updates: %s", err)
		}

		if change != nil {
			changeset.Changes = append(changeset.Changes, change)
		}

		githubactions.EndGroup()
	}

	if len(changeset.Changes) == 0 {
		githubactions.Noticef("no changes to chart versions detected")
	}

	err = changeset.Write()
	if err != nil {
		githubactions.Fatalf("failed to write changeset: %s", err)
	}
}
