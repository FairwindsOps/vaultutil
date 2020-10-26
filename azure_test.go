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

func TestAzureCredentials_buildEnv(t *testing.T) {
	testTime := time.Unix(1, 0)
	type fields struct {
		ClientID     string
		ClientSecret string
		Created      time.Time
		Duration     int64
		LeaseID      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		want    map[string]string
	}{
		{
			name: "simple no error",
			fields: fields{
				ClientID:     "id",
				ClientSecret: "supersecret",
				Created:      testTime,
				Duration:     30,
				LeaseID:      "vaultleaseid",
			},
			wantErr: false,
			want: map[string]string{
				"ARM_CLIENT_ID":              "id",
				"ARM_CLIENT_SECRET":          "supersecret",
				"ARM_SESSION_START":          "1",
				"ARM_SESSION_VAULT_LEASE_ID": "vaultleaseid",
				"ARM_SESSION_DURATION":       "30",
			},
		},
		{
			name: "no client id",
			fields: fields{
				ClientSecret: "supersecret",
				Created:      testTime,
				Duration:     30,
				LeaseID:      "vaultleaseid",
			},
			wantErr: true,
		},
		{
			name: "no client secret",
			fields: fields{
				ClientID: "id",
				Created:  testTime,
				Duration: 30,
				LeaseID:  "vaultleaseid",
			},
			wantErr: true,
		},
		{
			name: "no lease id",
			fields: fields{
				ClientID:     "id",
				ClientSecret: "supersecret",
				Created:      testTime,
				Duration:     30,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			az := &AzureCredentials{
				ClientID:     tt.fields.ClientID,
				ClientSecret: tt.fields.ClientSecret,
				Created:      tt.fields.Created,
				Duration:     tt.fields.Duration,
				LeaseID:      tt.fields.LeaseID,
			}
			err := az.buildEnv()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, az.EnvMap)
			}
		})
	}
}

func TestAzureCredentials_Expired(t *testing.T) {
	type fields struct {
		Created  time.Time
		Duration int64
	}

	tests := []struct {
		name   string
		fields fields
		buffer int64
		want   bool
	}{
		{
			name: "not expired",
			fields: fields{
				Created:  time.Now(),
				Duration: 100,
			},
			buffer: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AzureCredentials{
				Created:  tt.fields.Created,
				Duration: tt.fields.Duration,
			}
			got := a.Expired(tt.buffer)
			assert.Equal(t, tt.want, got)
		})
	}
}
