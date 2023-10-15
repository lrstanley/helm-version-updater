// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sethvargo/go-githubactions"
)

type ChangeSet struct {
	Path string `json:"-"`

	Changes []*Change `json:"changes"`
}

type Change struct {
	Chart         *Chart `json:"chart"`
	OldVersion    string `json:"old_version"`
	NewVersion    string `json:"new_version"`
	OldAppVersion string `json:"old_app_version"`
	NewAppVersion string `json:"new_app_version"`
}

func (c *ChangeSet) Write() error {
	var f io.WriteCloser
	var err error

	if c.Path == "-" {
		f = os.Stdout
	} else {
		f, err = os.Create(c.Path)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")

	err = enc.Encode(c)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	githubactions.SetOutput("changeset", string(data))

	return nil
}

func generateChangeset(config *Flags) (*ChangeSet, error) {
	var chartDirs []string

	// Recursively find ci-config.yaml files, and parse them.
	err := filepath.Walk(config.CheckDir, func(path string, info os.FileInfo, err error) error {
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
		return nil, fmt.Errorf("failed to find ci-config.yaml files: %s", err)
	} else if len(chartDirs) == 0 {
		return nil, fmt.Errorf("failed to find any ci-config.yaml files, validate check-dir")
	}

	var chart *Chart
	var change *Change

	changeset := &ChangeSet{Path: config.OutputFile}

	for _, dir := range chartDirs {
		githubactions.Group(fmt.Sprintf("chart-dir:%q", dir))

		// Read the chart.yaml file, and parse it.
		chart, err = ParseChart(filepath.Join(dir, "Chart.yaml"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse chart.yaml file: %s", err)
		}

		githubactions.Infof(
			"chart metadata name:%q description:%q version:%q appVersion: %q",
			chart.Name,
			chart.Description,
			chart.OriginalVersion,
			chart.OriginalAppVersion,
		)

		change, err = CheckImageUpdates(config, chart)
		if err != nil {
			return nil, fmt.Errorf("failed to check image updates: %s", err)
		}

		if change != nil {
			changeset.Changes = append(changeset.Changes, change)
		}

		githubactions.EndGroup()
	}

	return changeset, nil
}
