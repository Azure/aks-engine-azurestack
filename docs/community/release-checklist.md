# Releases

AKS Engine uses a [continuous delivery][] approach for creating releases. Every merged commit that passes
testing results in a deliverable that can be given a [semantic version][] tag and shipped.

## Master Is Always Releasable

The master `git` branch of a project should always work. Only changes considered ready to be
released publicly are merged.

AKS Engine depends on components that release new versions as often as needed. Fixing
a high priority bug requires the project maintainer to create a new patch release.
Merging a backward-compatible feature implies a minor release.

By releasing often, each release becomes a safe and routine event. This makes it faster
and easier for users to obtain specific fixes. Continuous delivery also reduces the work
necessary to release a product such as AKS Engine, which depends on several external projects.

## AKS Engine Releases As Needed

AKS Engine releases new versions when the team of maintainers determine it is needed. This usually
amounts to one or more releases each month.

Minor versions—for example, v0.**65**.0—are created from the master branch whenever
important features or changes have been merged and CI testing shows it to be stable over time.

Patch versions—for example, v0.64.**1**—are based on the previous release and created on demand
whenever important bug fixes arrive.

See "[Creating a New Release](#creating-a-new-release)" for more detail.

## Semantic Versioning

Releases of the `aks-engine-azurestack` binary comply with [semantic versioning][semantic version], with the "public API" broadly
defined as:

- REST, gRPC, or other API that is network-accessible
- Library or framework API intended for public use
- "Pluggable" socket-level protocols users can redirect
- CLI commands and output formats
- Integration with Azure public APIs such as ARM

In general, changes to anything a user might reasonably link to, customize, or integrate with should
be backward-compatible, or else require a major release. `aks-engine-azurestack` users can be confident that upgrading
to a patch or to a minor release will not break anything.

## Creating a New Release

Let's go through the process of creating a new release of the [aks-engine][] binary.

We will use **v0.65.0** as an example herein. You should replace this with the new version you're releasing.

### Release Branch

For a minor release, we will release from master. For a patch, we will create a new branch at the `Azure/` git origin from the previous release branch and use `git cherry-pick` to apply specific commits.

Once your source branch is prepared for a release, we run the "Create Release Branch" GitHub Action to automatically validate and create the destination release branch:

- [Create Release Branch action](https://github.com/Azure/aks-engine-azurestack/actions/workflows/create-release-branch.yaml)

Use the full "v"-prefixed semver release string in the field with the description "Which version are we creating a release branch for?", for example `v0.65.0`.

Use the source branch (`master` for minor releases, or a curated branch with cherry-picked commits for patch releases, for example `patch-release-v0.64.1`) in the field with the description "Which branch to source release branch from?".

Click "Run Workflow" to initiate the process of validating and creating our release branch. This automation will perform the following:

- Checkout the source commit.
- Run well-known "no egress" tests to validate that the base set of default components are pre-installed onto the default Linux and Windows VHDs.
- Create a new branch at the `Azure/` git origin named "release-<release version>", for example `release-v0.65.0` for the `v0.65.0` release.
- Generate automated release notes using the `git-chglog` tool.
- Create a branch with the generated release notes as a potential commit to the destination release branch.

#### Minor Release Branch using GitHub CLI

See example of how to prepare a minor branch:

```bash
# New release version
RELEASE_VERSION=v0.76.0
# Find your Azure/aks-engine-azurestack remote using git remote -v`
GH_ORG_REPO=Azure/aks-engine-azurestack
REMOTE=$(git remote -v | grep ${GH_ORG_REPO} | head -n 1 | cut -f 1)

# Run create-release-action GH action
WORKFLOW_FROM_BRANCH=master \
RELEASE_REPOSITORY=${GH_ORG_REPO} \
RELEASE_VERSION=${RELEASE_VERSION} \
RELEASE_FROM_BRANCH=master \
make create-release-branch

# Wait for action to finish
gh run watch -R ${GH_ORG_REPO}

# Validate remote branches
git fetch ${REMOTE}
git log ${REMOTE}/CHANGELOG-${RELEASE_VERSION}
git log ${REMOTE}/release-${RELEASE_VERSION}
```

#### Patch Release Branch using GitHub CLI

See example of how to prepare a patch release branch:

```bash
# Base release version
PREVIOUS_VERSION=v0.76.0
# New release version
RELEASE_VERSION=v0.76.1
# Find your Azure/aks-engine-azurestack remote using git remote -v`
GH_ORG_REPO=Azure/aks-engine-azurestack
REMOTE=$(git remote -v | grep ${GH_ORG_REPO} | head -n 1 | cut -f 1)

# Create local work branch
git checkout -b patch-release-${RELEASE_VERSION}
# Move it to previous release
git reset --hard ${REMOTE}/release-${PREVIOUS_VERSION}
# Cherry-pick commits to include
git cherry-pick ${COMMIT_1} ${COMMIT_2} ...
# Push work branch to remote
git push ${REMOTE} patch-release-${RELEASE_VERSION}

# Run create-release-action GH action
WORKFLOW_FROM_BRANCH=master \
RELEASE_REPOSITORY=${GH_ORG_REPO} \
RELEASE_VERSION=${RELEASE_VERSION} \
RELEASE_FROM_BRANCH=patch-release-${RELEASE_VERSION} \
scripts/gh-create-release-branch.sh

# Wait for action to finish
gh run watch -R ${GH_ORG_REPO}

# Validate remote branches
git fetch ${REMOTE}
git log ${REMOTE}/CHANGELOG-${RELEASE_VERSION}
git log ${REMOTE}/release-${RELEASE_VERSION}
```

### Release Notes

If the "Create Release Branch" GitHub Action ran successfully, there will be a new branch in the `Azure/aks-engine-azurestack` repository named "CHANGELOG-${RELEASE_VERSION}", for example `CHANGELOG-v0.65.0` for the `v0.65.0` release.

At this time project maintainers should review file `releases/CHANGELOG-${RELEASE_VERSION}.md` for the following:

- Does the generated list of changes meet with the desired set of changes to include in this release?
- Is there anything that can be improved with manual, human changes to the CHANGELOG markdown?
  - If so, edit the `.md` file and push a new commit to branch "CHANGELOG-${RELEASE_VERSION}"

After this, create a pull request for the changelog file:

```bash
gh pr create \
  --base release-${RELEASE_VERSION} \
  --head CHANGELOG-${RELEASE_VERSION} \
  --title "release: ${RELEASE_VERSION} CHANGELOG" \
  --body "Add CHANGELOG for upcoming ${RELEASE_VERSION} release" \
  --repo Azure/aks-engine-azurestack

# Merge if you have admin permissions
gh pr merge <PR_ID> --admin --squash --repo ${GH_ORG_REPO}
```

Ensure that at least two maintainers lgtm the final proposed CHANGELOG. Once this PR is merged to the release branch, a GitHub Action will perform the actual release publication.

### Publish the Release

By merging the CHANGELOG PR, you will initiate the final stage in the release process:

- Validate that the release branch has the expected CHANGELOG.
- Validate well-known "no egress" tests against the final release commit build.
- Validate well-known E2E tests against the final release commit build.
- Tag the release commit (for example, `v0.65.0` for the `v0.65.0` release).
- Build binaries for Linux, MacOS, and Windows.
- Create the release using the generated CHANGELOG, and upload the binaries.

Note: because the test validations above may be subject to environmental failures, it may be appropriate to retry the "Release" GitHub Action job if it fails for this reason. It's critical to investigate the failure and ensure that it's appropriate for retrying — failures that indicate a regression in the AKS Engine-generated ARM template should definitely block a release!

### Update Package Managers

Finally, let's make the new aks-engine-azurestack release easy to install.

#### The `brew` package manager

Create a pull request to add the new release to [brew][] through our [homebrew tap repository][brew-tap]. Update the macOS sha256 checksum with this command:

```
$ shasum -a 256  _dist/aks-engine-$RELEASE-darwin-amd64.tar.gz
```

The PR will look very similar to [this one][brew-pr].

#### The `gofish` package manager

The [gofish][] package manager has automation in place to create an update when AKS Engine creates a release. Check the [fish-food repository][gofish-food] to see that a pull request was created.

#### The `choco` package manager

Adding new versions to [choco][] is automated, but you can check the status of package approval and publishing at the [aks-engine-azurestack chocolatey page][choco-status].


[aks-engine]: https://github.com/Azure/aks-engine-azurestack/releases
[brew]: https://brew.sh/
[brew-pr]: https://github.com/Azure/homebrew-aks-engine/pull/34
[brew-tap]: https://github.com/Azure/homebrew-aks-engine/
[choco]: https://chocolatey.org/
[choco-status]: https://chocolatey.org/packages/aks-engine/
[conventional-commit]: https://www.conventionalcommits.org/en/v1.0.0-beta.3/
[git-chglog]: https://github.com/git-chglog/git-chglog
[gofish]: https://github.com/fishworks/gofish
[gofish-food]: https://github.com/fishworks/fish-food/
[gofish-pr]: https://github.com/fishworks/fish-food/pull/141
[new-release]: https://github.com/Azure/aks-engine-azurestack/releases/new
[continuous delivery]: https://en.wikipedia.org/wiki/Continuous_delivery
[semantic version]: http://semver.org
