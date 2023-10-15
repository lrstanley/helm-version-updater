// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
)

var (
	reVersionMatcher    = regexp.MustCompile(`(?m)^version:\s+?(.*)$`)
	reAppVersionMatcher = regexp.MustCompile(`(?m)^appVersion:\s+?(.*)$`)
)

// ChartConfig contains the CI configuration for a chart, which is stored in the
// <chart>/ci-config.yaml file.
type ChartConfig struct {
	Source struct {
		Image string `yaml:"image" json:"image"`
	} `yaml:"source" json:"source"`
}

func ParseChartConfig(path string) (*ChartConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)

	chart := &ChartConfig{}

	err = dec.Decode(chart)
	if err != nil {
		return nil, err
	}

	return chart, nil
}

// Chart contains partial fields from Chart.yaml -- goal is to only use this
// for reading, not writing.
type Chart struct {
	Path       string          `json:"path"   yaml:"-"`
	Config     *ChartConfig    `json:"config" yaml:"-"`
	Version    *semver.Version `json:"-"      yaml:"-"`
	AppVersion *semver.Version `json:"-"      yaml:"-"`

	Name               string `yaml:"name"        json:"name"`
	Description        string `yaml:"description" json:"description"`
	OriginalVersion    string `yaml:"version"     json:"version"`
	OriginalAppVersion string `yaml:"appVersion"  json:"appVersion"`
}

func (c *Chart) SetAppVersion(version *semver.Version) {
	if version.Major() > c.AppVersion.Major() {
		*c.Version = c.Version.IncMajor()
	} else if version.Minor() > c.AppVersion.Minor() {
		*c.Version = c.Version.IncMinor()
	} else {
		*c.Version = c.Version.IncPatch()
	}

	c.AppVersion = version
}

// WriteUpdatedVersions writes back to file via string replacement (can't use encode,
// as the order may not be the same).
func (c *Chart) WriteUpdatedVersions() error {
	body, err := os.ReadFile(c.Path)
	if err != nil {
		return err
	}

	v := strings.TrimPrefix(c.Version.String(), "v")
	appv := strings.TrimPrefix(c.AppVersion.String(), "v")

	if strings.HasPrefix(c.OriginalVersion, "v") {
		v = "v" + v
	}

	if strings.HasPrefix(c.OriginalAppVersion, "v") {
		appv = "v" + appv
	}

	body = reVersionMatcher.ReplaceAll(body, []byte("version: "+v))
	body = reAppVersionMatcher.ReplaceAll(body, []byte("appVersion: "+appv))

	return os.WriteFile(c.Path, body, 0o600)
}

func ParseChart(path string) (*Chart, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)

	chart := &Chart{Path: path}

	err = dec.Decode(chart)
	if err != nil {
		return nil, err
	}

	chart.Config, err = ParseChartConfig(filepath.Join(filepath.Dir(path), "ci-config.yaml"))
	if err != nil {
		return nil, err
	}

	chart.Version, err = semver.StrictNewVersion(chart.OriginalVersion)
	if err != nil {
		return nil, err
	}

	chart.AppVersion, err = semver.StrictNewVersion(chart.OriginalAppVersion)
	if err != nil {
		return nil, err
	}

	return chart, nil
}
