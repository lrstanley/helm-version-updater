// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"encoding/json"
	"io"
	"os"

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
