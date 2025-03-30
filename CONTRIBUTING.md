# Contributing guidelines

Welcome to our project's contribution guidelines! We're thrilled that you're interested in contributing to our repository.

This document outlines the process for making contributions to our project, whether you're fixing a bug, implementing a new feature, or suggesting improvements. Please take a moment to review these guidelines before diving into your contributions. Your participation is invaluable, and we appreciate your efforts to make our project better for everyone.

Let's build something great together!

- [Contributing guidelines](#contributing-guidelines)
  - [Reporting issues](#reporting-issues)
  - [Suggest a feature](#suggest-a-feature)
  - [Contribute with code](#contribute-with-code)
    - [Setup development envrionment](#setup-development-envrionment)
      - [Using Nix and Devenv](#using-nix-and-devenv)
      - [Using dev containers](#using-dev-containers)
      - [Install everything on host machine.](#install-everything-on-host-machine)
    - [Tools we use](#tools-we-use)
    - [Development lifecycle](#development-lifecycle)
      - [Submitting your changes for review](#submitting-your-changes-for-review)
      - [Commit guidelines](#commit-guidelines)
  - [Release process](#release-process)


## Reporting issues

If you found any issue, feel free to submit a GitHub issue.

Before submitting a new issue, we encourage you to utilize the search functionality to check if a similar issue has already been reported. This ensures that we avoid duplication and allows us to focus on addressing unique problems effectively.

When creating a new issue, please provide the most information you can like application version, operating system, logs and stack traces and anything else that you think is relevant and can help the investigation process.

## Suggest a feature

TBD


## Contribute with code

If you are a developer, and want to contribute with the code, please follow the next steps, to understand how to setup the development envrionment and our release process.

### Setup development envrionment

#### Using Nix and Devenv

The easiest way to setup this project on your local machine is to use [Devenv](https://devenv.sh/getting-started/)

Then run `devenv shell` to start a shell, with all the tools installed.

It´s also recommended to install [direnv](https://direnv.net/). Direnv integrates very well with Devbox and facilitates the management of project level envrionment variables.

There are multiple ways to install direnv. Check [Installation Guide](https://direnv.net/docs/installation.html) to see the most appropriate way, based on your Operating system.

#### Using dev containers

If you use VSCode or GitHub Codespaces, we also provide a [Devcontainer](https://containers.dev/) definition that you can use. It´s simply a wrapper for Devenv, but allows to start coding right way, without even installing Devbox on your machine.

#### Install everything on host machine.

While, not recommended, you can still install everything directly on your development machine.

You will need [Go](https://go.dev/) 1.22 and a few [Tools](#tools-we-use).


### Tools we use

- [Task](https://taskfile.dev/) a task runner / build tool, modern alternative to Make. Useful to define common tasks like build the application or run the tests. Run `task -l` or check [Taskfile.yml](Taskfile.yml) to see the available tasks.
- [lefthook](https://github.com/evilmartians/lefthook) -  Fast and powerful Git hooks manager for any type of projects. Useful to run tasks like linting and formatting, before commiting changes GitHub.
- [golangci-lint](https://golangci-lint.run/) - Go linters aggregator. Helps ensuring the quality of the code.
- [gomarkdoc](https://github.com/princjef/gomarkdoc) - Generate markdown documentation for Go code
- [gotestsum](https://github.com/gotestyourself/gotestsum) - 'go test' runner with output optimized for humans, JUnit XML for CI integration, and a summary of the test results.

### Development lifecycle

This project follows [GitHub flow](https://docs.github.com/en/get-started/using-github/github-flow) for managing changes.

When implmenting a new feature, start by creating a new branch from `main`, with a descriptive name (Ex: `feat/my-awesome-feature` or `fix/some-bug`).

Having a descriptive name helps to reason about the branches, when you have many.

Checkout to that branch and do your changes.

Some useful guidelines when working on feature branches:

- **keep it short lived** - Long running feature branches can lead to problems, like merge conflicts. You should aim to create a feature branch, for a feature than is small enough to be done in a few days.
- **rebase with main at least once a day** - this ensure you are always working with the most recent code and allows to fix any conflicts that might occurr, early in the process.

#### Submitting your changes for review

When you are ready create a Pull request to the main branch.

When creating a pull request, you should:

- Provide a descriptive PR title, following [Conventional Commits](https://www.conventionalcommits.org/en/) specification.
- Provide a short description of what changes you did, core architecture decisions you took, and link to any issue the PR might relate to.
- Ensure that any automated checks like Linting and Tests pass.

The PR will then be reviewed and changes may by requested. Keep commiting those changes, until the PR is approved.

After being approved, the maintainers will merge the PR to main branch and start the release process.

#### Commit guidelines

The project folows [Conventional Commits](https://www.conventionalcommits.org/en/) specification.

Each commit message should begin with a type, indicating the nature of the change (e.g., feat for a new feature, fix for a bug fix, docs for documentation changes), followed by a concise and descriptive message.

Additionally, providing an optional scope and further details in the commit message body is encouraged when necessary. This approach streamlines the review process, facilitates automated release notes generation, and enhances overall project maintainability.

We also recommend squashing your commits when appropriate.

## Release process

We use [Release Drafter](https://github.com/marketplace/actions/release-drafter) to automatically create draft releases with appropriate release notes, anytime a PR is merged.

When we are ready to create a new release, we simply publish the release, which will trigger GitHub actions, that will publish any related artifacts and commit a Changelog to the project repository.


