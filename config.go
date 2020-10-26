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
	"time"
)

const (
	// BaseURLGovCloud is the base URL for AWS GovCloud
	BaseURLGovCloud = "amazonaws-us-gov.com"
	// BaseURLDefault is the normal AWS base URL
	BaseURLDefault = "aws.amazon.com"
)

// Config holds all the config
type Config struct {
	AWSBaseURL string
	Path       string
	Role       string
	// BufferSeconds is the number of seconds to renew before expiration
	BufferSeconds int64
}

// NewConfig returns a config object
func NewConfig(partition, role, path string, buffer int64) *Config {
	ret := &Config{
		Role:          role,
		Path:          path,
		BufferSeconds: buffer,
	}
	switch partition {
	case "aws":
		ret.AWSBaseURL = BaseURLDefault
	case "gov":
		ret.AWSBaseURL = BaseURLGovCloud
	}

	return ret
}

// expired checks to see if a set of credentials are expired
func expired(buffer, duration int64, created time.Time) bool {
	elapsed := time.Since(created)

	// Take 30 seconds off the duration of the valid time to account for use
	d := time.Second * time.Duration(duration-buffer)
	return elapsed > d
}
