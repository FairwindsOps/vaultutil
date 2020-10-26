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
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"k8s.io/klog"
)

// AWSCredentials holds the AWS Credential JSON
type AWSCredentials struct {
	// AccessKeyId is the access key ID of the sts credentials
	AccessKeyID string `json:"sessionId"`
	// SecretAccessKey is the secret access key of the sts credentials
	SecretAccessKey string `json:"sessionKey"`
	// SessionToken is the token of the sts credentials
	SessionToken string `json:"sessionToken"`
	// Created is the time that the credentials were issued
	Created time.Time `json:"created,omitempty"`
	// Duration is the time in seconds that the credentials are valid for
	Duration int64 `json:"duration,omitempty"`
	// LeaseID is the vault lease id. Can be usd to revoke the credentials
	LeaseID string `json:"lease_id,omitempty"`
	// EnvMap is a map of environment variables to the values above. It can be used
	// to populate necessary CLI environement variables for using the credentials. In
	// addition, this tool adds the vault lease and the duration/creation in order to
	// reduce the number of times that new credentials need to be generated.
	// The environment variables are:
	//  AWS_ACCESS_KEY_ID=AccessKeyID
	//  AWS_SECRET_ACCESS_KEY=SecretAccessKey
	//  AWS_SESSION_TOKEN=SessionToken
	//  AWS_SECURITY_TOKEN=SessionToken
	//  AWS_SESSION_START=Created (in Unix time)
	//  AWS_SESSION_DURATION=Duration
	//  AWS_SESSION_VAULT_LEASE_ID=LeaseID
	EnvMap map[string]string `json:"environment"`
}

type vaultAWSCredentials struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	LeaseDuration int64  `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`
	Data          struct {
		AccessKey     string `json:"access_key"`
		SecretKey     string `json:"secret_key"`
		SecurityToken string `json:"security_token"`
	} `json:"data"`
	Warnings interface{} `json:"warnings"`
}

// Expired returns true if the AWS credentials are expired
func (a *AWSCredentials) Expired(buffer int64) bool {
	return expired(buffer, a.Duration, a.Created)
}

// ReadFromEnv populates a set of aws credentials that were previously exported
func (a *AWSCredentials) ReadFromEnv() error {
	created, err := strconv.ParseInt(os.Getenv("AWS_SESSION_START"), 10, 64)
	if err != nil {
		return err
	}

	duration, err := strconv.ParseInt(os.Getenv("AWS_SESSION_DURATION"), 10, 64)
	if err != nil {
		return err
	}

	a.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	a.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	a.SessionToken = os.Getenv("AWS_SESSION_TOKEN")
	a.LeaseID = os.Getenv("AWS_SESSION_VAULT_LEASE_ID")
	a.Created = time.Unix(created, 0)
	a.Duration = duration

	if err := a.buildEnv(); err != nil {
		return err
	}

	return nil
}

// Revoke revokes the vault lease associated with the credentials
func (a *AWSCredentials) Revoke() error {
	return revokeLease(a.LeaseID)
}

// buildEnv populates the environment variable of the credentials struct
// This allows consumers of the credentials to reliably and consistently export
// the correct environment variables. The variables are documented in the
// AWSCredentials type.
func (a *AWSCredentials) buildEnv() error {
	a.EnvMap = make(map[string]string)

	if a.AccessKeyID != "" {
		a.EnvMap["AWS_ACCESS_KEY_ID"] = a.AccessKeyID
	} else {
		return fmt.Errorf("cannot set env: access key id was empty")
	}

	if a.SecretAccessKey != "" {
		a.EnvMap["AWS_SECRET_ACCESS_KEY"] = a.SecretAccessKey
	} else {
		return fmt.Errorf("cannot set env: secret access key was empty")
	}

	if a.SessionToken != "" {
		a.EnvMap["AWS_SESSION_TOKEN"] = a.SessionToken
		a.EnvMap["AWS_SECURITY_TOKEN"] = a.SessionToken
	} else {
		return fmt.Errorf("cannot set env: session token was empty")
	}

	if a.LeaseID != "" {
		a.EnvMap["AWS_SESSION_VAULT_LEASE_ID"] = a.LeaseID
	} else {
		return fmt.Errorf("cannot set env: vault lease id was empty")
	}
	a.EnvMap["AWS_SESSION_DURATION"] = strconv.FormatInt(a.Duration, 10)
	a.EnvMap["AWS_SESSION_START"] = strconv.FormatInt(a.Created.Unix(), 10)

	klog.V(10).Infof("environment vars: %v", a.EnvMap)
	return nil
}

// NewAWSCredentials returns existing ones from env if they are not expired
// if they are expired, or if we can't get any from env, return a new set
func (c Config) NewAWSCredentials() (*AWSCredentials, error) {
	creds := &AWSCredentials{}
	err := creds.ReadFromEnv()
	if err == nil {
		klog.V(3).Infof("credentials found in environment - checking expiration")
		if creds.Expired(c.BufferSeconds) {
			klog.V(3).Info("credentials were expired - getting new ones")
			newCreds, err := c.AWSLogin()
			if err != nil {
				return nil, err
			}
			return newCreds, nil
		}
		klog.V(3).Infof("credentials were valid - returning them")
		return creds, nil
	}
	klog.V(3).Infof("error getting credentials from env: %s", err.Error())
	klog.V(2).Info("no existing credentials found - getting new ones")

	newCreds, err := c.AWSLogin()
	if err != nil {
		return nil, err
	}
	return newCreds, nil
}

// BuildConsoleLogin returns a new console login
func (c Config) BuildConsoleLogin() (string, error) {
	creds, err := c.getAWSCredentials()
	if err != nil {
		return "", err
	}

	token, err := c.getSigninToken(creds)
	if err != nil {
		return "", err
	}

	klog.V(9).Info(token)

	url, err := c.generateSigninURL(token)
	if err != nil {
		return "", err
	}
	fmt.Println(url)

	return "", nil
}

// getSigninToken gets a federation token for signin
func (c Config) getSigninToken(creds *AWSCredentials) (*string, error) {
	credsData, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}
	klog.V(4).Info(string(credsData))

	req, err := http.NewRequest("GET", fmt.Sprintf("https://signin.%s/federation", c.AWSBaseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("Action", "getSigninToken")
	q.Add("Session", string(credsData))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get signin token failed with code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	klog.V(10).Info(string(body))

	var respParsed map[string]string

	err = json.Unmarshal([]byte(body), &respParsed)
	if err != nil {
		return nil, err
	}

	token, ok := respParsed["SigninToken"]
	if !ok {
		return nil, fmt.Errorf("could not get signin token from body")
	}

	return &token, nil
}

// generateSigninURL returns a string that is the url to sign in to the console
func (c Config) generateSigninURL(token *string) (string, error) {
	dest := url.PathEscape(fmt.Sprintf("https://console.%s/", c.AWSBaseURL))
	issuer := url.PathEscape("https://github.com/fairwindsops/vault-util")

	url := fmt.Sprintf("https://signin.%s/federation"+
		"?Action=login"+
		"&Issuer=%s"+
		"&Destination=%s"+
		"&SigninToken=%s",
		c.AWSBaseURL, issuer, dest, *token)

	return url, nil
}

// getAWSCredentials retrieves the AWS credentials from the current environment
// This function only returns the access key, secret key and token.
// Used by the aws-console to build JSON for generating links
func (c Config) getAWSCredentials() (*AWSCredentials, error) {
	creds := credentials.NewEnvCredentials()

	// Retrieve the credentials value
	credValue, err := creds.Get()
	if err != nil {
		return nil, err
	}
	klog.V(10).Info(credValue)

	ret := &AWSCredentials{
		AccessKeyID:     credValue.AccessKeyID,
		SecretAccessKey: credValue.SecretAccessKey,
		SessionToken:    credValue.SessionToken,
	}
	klog.V(4).Info(*ret)

	return ret, nil
}

// AWSLogin calls vault write on an sts endpoint and generates the necessary environment variables
func (c Config) AWSLogin() (*AWSCredentials, error) {
	endpoint := fmt.Sprintf("%s/sts/%s", c.Path, c.Role)
	klog.V(3).Infof("attempting to get aws credentials from vault at %s", endpoint)
	cmd := exec.Command("vault", "write", endpoint, "-format=json")

	data, err := cmd.CombinedOutput()
	if err != nil {
		output := strings.TrimSpace(string(data))
		return nil, fmt.Errorf("vault aws credentials failed with status %d: %s", cmd.ProcessState.ExitCode(), output)
	}

	creds := &vaultAWSCredentials{}
	err = json.Unmarshal(data, creds)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling vault token: %s", err.Error())
	}

	ret := &AWSCredentials{
		AccessKeyID:     creds.Data.AccessKey,
		SecretAccessKey: creds.Data.SecretKey,
		SessionToken:    creds.Data.SecurityToken,
		Created:         time.Now(),
		Duration:        creds.LeaseDuration,
		LeaseID:         creds.LeaseID,
	}

	klog.V(10).Infof("got credentials: %v", ret)

	if err := ret.buildEnv(); err != nil {
		return nil, err
	}
	return ret, nil
}
