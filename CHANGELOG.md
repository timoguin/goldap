# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

<!--
Any notes about merges to master that haven't been added to a Git tag should go
here. When it's time to cut a new release, create a header for the new version
below and move the content of this section down to the new version.

### Added

### Changed

### Removed
-->

No notes about an upcoming release.

## [0.0.1]

Initial release with support for LDAP searches

### Added

- Clone of the `ldapsearch` tool with the `goldap search` command
- Bash completion available by sourcing the output of `goldap bash-completion`
- Customizable config via CLI flag, ENV vars, or config file

<!--
This section should be updated with every release. It contains a sequence of
links to GitHub that show the full Git diff between each release. The brackets
allow us to render the version headers as links by adding brackets to any
matching headers. Any commits that don't yet belong to a Git tag will show the
Git diff from the last tag to the master branch HEAD.
-->
[Unreleased]: https://github.com/timoguin/goldap/compare/v0.0.1..HEAD
[0.0.1]: https://github.com/timoguin/goldap/releases/tag/v0.0.1
