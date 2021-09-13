# vaultutil

[![PkgGoDev][doc-image]][doc-link] [![GitHub release (latest SemVer)][release-image]][release-link] [![GitHub go.mod Go version][version-image]][version-link] [![CircleCI][circleci-image]][circleci-link] [![Go Report Card][goreport-image]][goreport-link]

[doc-image]: https://pkg.go.dev/badge/fairwindsops/vaultutil
[doc-link]: https://pkg.go.dev/github.com/fairwindsops/vaultutil

[version-image]: https://img.shields.io/github/go-mod/go-version/FairwindsOps/vaultutil
[version-link]: https://github.com/FairwindsOps/vaultutil

[release-image]: https://img.shields.io/github/v/release/FairwindsOps/vaultutil
[release-link]: https://github.com/FairwindsOps/vaultutil

[goreport-image]: https://goreportcard.com/badge/github.com/FairwindsOps/vaultutil
[goreport-link]: https://goreportcard.com/report/github.com/FairwindsOps/vaultutil

[circleci-image]: https://circleci.com/gh/FairwindsOps/vaultutil/tree/master.svg?style=svg
[circleci-link]: https://circleci.com/gh/FairwindsOps/vaultutil

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


<!-- Begin boilerplate -->
## Join the Fairwinds Open Source Community

The goal of the Fairwinds Community is to exchange ideas, influence the open source roadmap,
and network with fellow Kubernetes users.
[Chat with us on Slack](https://join.slack.com/t/fairwindscommunity/shared_invite/zt-e3c6vj4l-3lIH6dvKqzWII5fSSFDi1g)
[join the user group](https://www.fairwinds.com/open-source-software-user-group) to get involved!

<a href="https://www.fairwinds.com/t-shirt-offer?utm_source=vaultutil&utm_medium=vaultutil&utm_campaign=vaultutil-tshirt">
  <img src="https://www.fairwinds.com/hubfs/Doc_Banners/Fairwinds_OSS_User_Group_740x125_v6.png" alt="Love Fairwinds Open Source? Share your business email and job title and we'll send you a free Fairwinds t-shirt!" />
</a>

## Other Projects from Fairwinds

Enjoying Vaultutil? Check out some of our other projects:
* [Polaris](https://github.com/FairwindsOps/Polaris) - Audit, enforce, and build policies for Kubernetes resources, including over 20 built-in checks for best practices
* [Goldilocks](https://github.com/FairwindsOps/Goldilocks) - Right-size your Kubernetes Deployments by compare your memory and CPU settings against actual usage
* [Pluto](https://github.com/FairwindsOps/Pluto) - Detect Kubernetes resources that have been deprecated or removed in future versions
* [Nova](https://github.com/FairwindsOps/Nova) - Check to see if any of your Helm charts have updates available
* [rbac-manager](https://github.com/FairwindsOps/rbac-manager) - Simplify the management of RBAC in your Kubernetes clusters
