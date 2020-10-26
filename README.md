# vaultutil

[![PkgGoDev][doc-image]][doc-link] [![GitHub release (latest SemVer)][release-image]][release-link] [![GitHub go.mod Go version][version-image]][version-link] [![CircleCI][circleci-image]][circleci-link] [![Code Coverage][codecov-image]][codecov-link] [![Go Report Card][goreport-image]][goreport-link]

[doc-image]: https://pkg.go.dev/badge/fairwindsops/vaultutil
[doc-link]: https://pkg.go.dev/fairwindsops/vaultutil

[version-image]: https://img.shields.io/github/go-mod/go-version/FairwindsOps/vaultutil
[version-link]: https://github.com/FairwindsOps/vaultutil

[release-image]: https://img.shields.io/github/v/release/FairwindsOps/vaultutil
[release-link]: https://github.com/FairwindsOps/vaultutil

[goreport-image]: https://goreportcard.com/badge/github.com/FairwindsOps/vaultutil
[goreport-link]: https://goreportcard.com/report/github.com/FairwindsOps/vaultutil

[circleci-image]: https://circleci.com/gh/FairwindsOps/vaultutil/tree/master.svg?style=svg
[circleci-link]: https://circleci.com/gh/FairwindsOps/vaultutil

[codecov-image]: https://codecov.io/gh/FairwindsOps/vaultutil/branch/master/graph/badge.svg
[codecov-link]: https://codecov.io/gh/FairwindsOps/vaultutil

This library provides utilities for utilizing Vault in various user workflows and environments.

```
import "github.com/fairwindsops/vaultutil"
```

## AWS

There are helpers for:

- Getting and refreshing STS credentials from a vault aws backend
- Generating AWS Console login links from STS credentials

## Azure

There are helpers for:

- Getting and refreshing service principals from a vault azure backend
