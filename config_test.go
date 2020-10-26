// Copyright 2020 Fairwinds
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vaultutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpired(t *testing.T) {
	tests := []struct {
		name     string
		duration int64
		created  time.Time
		want     bool
	}{
		{
			name:     "true",
			duration: 100,
			created:  time.Now().Add(time.Second * -200),
			want:     true,
		},
		{
			name:     "false",
			duration: 100,
			created:  time.Now(),
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expired(int64(30), tt.duration, tt.created)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewConfig(t *testing.T) {
	type args struct {
		partition string
		role      string
		path      string
		buffer    int64
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "basic gov",
			args: args{
				partition: "gov",
				path:      "someone",
				role:      "arole",
				buffer:    30,
			},
			want: &Config{
				AWSBaseURL:    "amazonaws-us-gov.com",
				Path:          "someone",
				Role:          "arole",
				BufferSeconds: 30,
			},
		},
		{
			name: "basic commercial",
			args: args{
				partition: "aws",
				path:      "someone",
				role:      "arole",
				buffer:    30,
			},
			want: &Config{
				AWSBaseURL:    "aws.amazon.com",
				Path:          "someone",
				Role:          "arole",
				BufferSeconds: 30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConfig(tt.args.partition, tt.args.role, tt.args.path, tt.args.buffer)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
