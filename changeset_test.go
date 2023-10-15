// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"testing"
)

func Test_generateChangeset(t *testing.T) {
	tests := []struct {
		name    string
		want    []*Change
		wantErr bool
		config  *Flags
	}{
		{
			name: "test-1",
			want: []*Change{
				{
					OldVersion:    "1.0.0",
					NewVersion:    "1.1.0",
					OldAppVersion: "1.1.0",
					NewAppVersion: "1.2.0",
				},
			},
			wantErr: false,
			config:  &Flags{CheckDir: "tests/test-1/", DryRun: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateChangeset(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateChangeset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil && tt.want != nil || got != nil && tt.want == nil {
				t.Errorf("generateChangeset() = %v, want %v", got, tt.want)
			}

			if got == nil {
				return
			}

			for i, v := range got.Changes {
				if v.OldVersion != tt.want[i].OldVersion {
					t.Errorf("OldVersion = %v, want %v", v.OldVersion, tt.want[i].OldVersion)
				}

				if v.NewVersion != tt.want[i].NewVersion {
					t.Errorf("NewVersion = %v, want %v", v.NewVersion, tt.want[i].NewVersion)
				}

				if v.OldAppVersion != tt.want[i].OldAppVersion {
					t.Errorf("OldAppVersion = %v, want %v", v.OldAppVersion, tt.want[i].OldAppVersion)
				}

				if v.NewAppVersion != tt.want[i].NewAppVersion {
					t.Errorf("NewAppVersion = %v, want %v", v.NewAppVersion, tt.want[i].NewAppVersion)
				}
			}
		})
	}
}
