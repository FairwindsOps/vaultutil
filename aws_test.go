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

func TestAWSCredentials_buildEnv(t *testing.T) {
	testTime := time.Unix(1, 0)
	type fields struct {
		AccessKeyID     string
		SecretAccessKey string
		SessionToken    string
		Created         time.Time
		Duration        int64
		LeaseID         string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr bool
	}{
		{
			name: "simple no error",
			fields: fields{
				AccessKeyID:     "SOMEACCESSKEYID",
				SecretAccessKey: "supersecret",
				SessionToken:    "token",
				Created:         testTime,
				Duration:        30,
				LeaseID:         "vaultleaseid",
			},
			wantErr: false,
			want: map[string]string{
				"AWS_ACCESS_KEY_ID":          "SOMEACCESSKEYID",
				"AWS_SECRET_ACCESS_KEY":      "supersecret",
				"AWS_SESSION_TOKEN":          "token",
				"AWS_SECURITY_TOKEN":         "token",
				"AWS_SESSION_START":          "1",
				"AWS_SESSION_VAULT_LEASE_ID": "vaultleaseid",
				"AWS_SESSION_DURATION":       "30",
			},
		},
		{
			name: "empty access key",
			fields: fields{
				SecretAccessKey: "supersecret",
				SessionToken:    "token",
				Created:         testTime,
				Duration:        30,
				LeaseID:         "vaultleaseid",
			},
			wantErr: true,
		},
		{
			name: "empty secret key",
			fields: fields{
				AccessKeyID:  "SOMEACCESSKEYID",
				SessionToken: "token",
				Created:      testTime,
				Duration:     30,
				LeaseID:      "vaultleaseid",
			},
			wantErr: true,
		},
		{
			name: "empty token",
			fields: fields{
				AccessKeyID:     "SOMEACCESSKEYID",
				SecretAccessKey: "supersecret",
				Created:         testTime,
				Duration:        30,
				LeaseID:         "vaultleaseid",
			},
			wantErr: true,
		},
		{
			name: "empty lease id",
			fields: fields{
				AccessKeyID:     "SOMEACCESSKEYID",
				SecretAccessKey: "supersecret",
				SessionToken:    "token",
				Created:         testTime,
				Duration:        30,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AWSCredentials{
				AccessKeyID:     tt.fields.AccessKeyID,
				SecretAccessKey: tt.fields.SecretAccessKey,
				SessionToken:    tt.fields.SessionToken,
				Created:         tt.fields.Created,
				Duration:        tt.fields.Duration,
				LeaseID:         tt.fields.LeaseID,
			}
			err := a.buildEnv()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, a.EnvMap)
			}

		})
	}
}
