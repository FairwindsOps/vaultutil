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
// limitations under the License

// Package vaultutil is a go library that provides functions and helpers for working
// with vault-provided cloud credentials.
//
// Example of Getting AWS Credentials:
//  import (
//      "fmt"
//      "os"
//
//      "github.com/fairwindsops/vaultutil"
//  )
//
//  func main() {
//      c := vaultutil.NewConfig("aws", "admin", "aws-account", 120)
//      creds, err := c.AWSLogin()
//      if err != nil {
//         fmt.Println(err)
//         os.Exit(1)
//      }
//  }
//
package vaultutil
