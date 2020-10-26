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

package util

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"k8s.io/klog"
)

// AzureCredentials are a set of azure credentials from vault
type AzureCredentials struct {
	// ClientID is the username of the azure service principal
	ClientID string `json:"client_id"`
	// ClientSecret is the password of the azure service principal
	ClientSecret string `json:"client_secret"`
	// Created is the date/time that the credentials were requested
	Created time.Time `json:"created"`
	// Duration is the number of seconds the credentials are valid
	Duration int64 `json:"duration"`
	// LeaseID is the Vault Lease ID of the requested credentials.
	// This can be used to revoke the lease when the credentials are no longer needed.
	LeaseID string `json:"lease_id"`
	// EnvMap is a map of environment variables to the values above. It can be used
	// to populate necessary CLI environement variables for using the credentials. In
	// addition, this tool adds the vault lease and the duration/creation in order to
	// reduce the number of times that new credentials need to be generated.
	// The environment variables are:
	//  ARM_CLIENT_ID=ClientID
	//  ARM_CLIENT_SECRET=ClientSecret
	//  ARM_SESSION_START=Created (in Unix time)
	//  ARM_SESSION_DURATION=Duration
	//  ARM_SESSION_VAULT_LEASE_ID=LeaseID
	EnvMap map[string]string `json:"environment"`
}

// vaultAzureCredentials is the response from Vault for azure backends
// This is used internally to parse the response from Vault.
type vaultAzureCredentials struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	LeaseDuration int64  `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`
	Data          struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"data"`
	Warnings interface{} `json:"warnings"`
}

// Expired checks to see if the azure credentials are expired
func (az *AzureCredentials) Expired(buffer int64) bool {
	return expired(buffer, az.Duration, az.Created)
}

// ReadFromEnv populates a set of azure credentials that were previously exported
func (az *AzureCredentials) ReadFromEnv() error {
	created, err := strconv.ParseInt(os.Getenv("ARM_SESSION_START"), 10, 64)
	if err != nil {
		klog.V(3).Infof("error getting arm session start: %s", err.Error())
		return err
	}

	duration, err := strconv.ParseInt(os.Getenv("ARM_SESSION_DURATION"), 10, 64)
	if err != nil {
		klog.V(3).Infof("error getting arm session duration: %s", err.Error())
		return err
	}

	az.ClientID = os.Getenv("ARM_CLIENT_ID")
	az.ClientSecret = os.Getenv("ARM_CLIENT_SECRET")
	az.LeaseID = os.Getenv("ARM_SESSION_VAULT_LEASE_ID")
	az.Created = time.Unix(created, 0)
	az.Duration = duration

	if err := az.buildEnv(); err != nil {
		return err
	}

	return nil
}

// Revoke revokes the vault lease associated with the credentials
func (az *AzureCredentials) Revoke() error {
	return revokeLease(az.LeaseID)
}

// buildEnv populates the environment variable of the credentials struct
// This allows consumers of the credentials to reliably and consistently export
// the correct environment variables. The variables are documented in the
// AzureCredentials type.
func (az *AzureCredentials) buildEnv() error {
	az.EnvMap = make(map[string]string)

	if az.ClientID != "" {
		az.EnvMap["ARM_CLIENT_ID"] = az.ClientID
	} else {
		return fmt.Errorf("cannot set env: client id was empty")
	}

	if az.ClientSecret != "" {
		az.EnvMap["ARM_CLIENT_SECRET"] = az.ClientSecret
	} else {
		return fmt.Errorf("cannot set env: client secret was empty")
	}

	if az.LeaseID != "" {
		az.EnvMap["ARM_SESSION_VAULT_LEASE_ID"] = az.LeaseID
	} else {
		return fmt.Errorf("cannot set env: vault least id was empty")
	}
	az.EnvMap["ARM_SESSION_DURATION"] = strconv.FormatInt(az.Duration, 10)
	az.EnvMap["ARM_SESSION_START"] = strconv.FormatInt(az.Created.Unix(), 10)

	return nil
}

// NewAzureCredentials returns existing ones from env if they are not expired
// if they are expired, or if we can't get any from env, return a new set.
func (c Config) NewAzureCredentials() (*AzureCredentials, error) {
	creds := &AzureCredentials{}
	if err := creds.ReadFromEnv(); err == nil {
		if creds.Expired(c.BufferSeconds) {
			klog.V(3).Infof("found expired credentials - getting new ones")
			newCreds, err := c.AzureLogin()
			if err != nil {
				return nil, err
			}
			return newCreds, nil
		}
		klog.V(3).Infof("credentials were not expired - returning them")
		return creds, nil
	}
	klog.V(3).Infof("unable to retrieve existing credentials - getting new ones")
	newCreds, err := c.AzureLogin()
	if err != nil {
		return nil, err
	}

	return newCreds, nil
}

// AzureLogin calls vault read on a credentials endpoint and generates the necessary environment variables.
// If no environment variable support is desired, and renewing credentials is not needed, then this function
// can be used to get just a simple set of credentials.
func (c Config) AzureLogin() (*AzureCredentials, error) {
	endpoint := fmt.Sprintf("%s/creds/%s", c.Path, c.Role)
	klog.V(3).Infof("attempting to get azure credentials from vault at %s", endpoint)
	cmd := exec.Command("vault", "read", endpoint, "-format=json")

	data, err := cmd.CombinedOutput()
	if err != nil {
		output := strings.TrimSpace(string(data))
		return nil, fmt.Errorf("vault azure credentials failed with status %d: %s", cmd.ProcessState.ExitCode(), output)
	}

	creds := &vaultAzureCredentials{}
	err = json.Unmarshal(data, creds)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling vault token: %s", err.Error())
	}

	ret := &AzureCredentials{
		ClientID:     creds.Data.ClientID,
		ClientSecret: creds.Data.ClientSecret,
		Created:      time.Now(),
		Duration:     creds.LeaseDuration,
		LeaseID:      creds.LeaseID,
	}

	if err := ret.buildEnv(); err != nil {
		return nil, err
	}

	return ret, nil
}
