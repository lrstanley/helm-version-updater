// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sethvargo/go-githubactions"
)

func CheckImageUpdates(config *Flags, chart *Chart) (change *Change, err error) {
	if chart.Config.Source.Image == "" {
		githubactions.Warningf("chart %q does not have a source image defined", chart.Name)
		return nil, nil //nolint:nilnil // no error, just no change.
	}

	ref, err := name.ParseReference(chart.Config.Source.Image)
	if err != nil {
		return nil, err
	}

	tags, err := remote.List(ref.Context(), remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, err
	}

	versions := semver.Collection{}
	var v *semver.Version

	for _, tag := range tags {
		v, err = semver.StrictNewVersion(tag)
		if err != nil {
			continue
		}

		if v.Prerelease() != "" && !config.SupportPreRelease {
			continue
		}

		versions = append(versions, v)
	}

	sort.Sort(sort.Reverse(versions))

	if len(versions) == 0 {
		githubactions.Warningf("chart %q has no valid versions matching our constraints", chart.Name)
		return nil, nil //nolint:nilnil // no error, just no change.
	}

	latest := versions[0]

	if chart.AppVersion.LessThan(latest) {
		chart.SetAppVersion(latest)
		githubactions.Noticef("updating chart %q app version from %q to %q", chart.Name, chart.OriginalAppVersion, chart.AppVersion.String())
		githubactions.Noticef("updating chart %q main version from %q to %q", chart.Name, chart.OriginalVersion, chart.Version.String())

		if !config.DryRun {
			err = chart.WriteUpdatedVersions()
			if err != nil {
				return nil, err
			}
		}

		return &Change{
			Chart:         chart,
			Path:          filepath.Dir(chart.Path),
			OldVersion:    chart.OriginalVersion,
			NewVersion:    chart.Version.String(),
			OldAppVersion: chart.OriginalAppVersion,
			NewAppVersion: chart.AppVersion.String(),
		}, nil
	}

	githubactions.Noticef("chart %q is already up-to-date", chart.Name)
	return nil, nil //nolint:nilnil // no error, just no change.
}
