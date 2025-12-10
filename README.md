# DVK-Project - What is it?

The DVK Project is a mirror repository of Whoknows Variations for use during the 4th semester DevOps elective(2025) at EK(formerly KEA). This readme file is tasked with providing an easy overview of the directories and files that the team believes to be important in order to understand the motivation of the team behind the rewrite.

The goal is to rewrite a legacy project using traditional DevOps tools and adhering to DevOps principles.

## Requirements

* Go 1.25
* Gorilla Mux 1.8.1
* Postgres driver: lib/pq 1.10.9


## Notable Content

[Overview of the documentation for this rewrite project](/documentation/)

[Our Choices - Programming Language, Repo Structure & More](/documentation/our_choices.md)

[Our Conventions](/documentation/our_conventions.md)

[Security Issues With Legacy Project](/documentation/legacy_codebase/Legacy_Codebase_Problems.md)

[Service Level Agreement](/documentation/SLA.md)

## Github Actions

[![Continuous Delivery](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/continuous_delivery.yml/badge.svg)](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/continuous_delivery.yml)
[![Continuous Deployment](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/continuous_deployment.yml/badge.svg)](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/continuous_deployment.yml)
[![Scheduled Health Check](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/health.yml/badge.svg)](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/health.yml)
[![Tests](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/tests.yml/badge.svg)](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/tests.yml)
[![Golangci Lint](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/golangci_lint.yml/badge.svg)](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/golangci_lint.yml)
[![Hadolint Dockerfile](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/hadolint.yml/badge.svg)](https://github.com/DVK-DEVOPS/DVK-Project/actions/workflows/hadolint.yml)

## Quality Analysis

[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=DVK-DEVOPS_DVK-Project&metric=bugs)](https://sonarcloud.io/summary/new_code?id=DVK-DEVOPS_DVK-Project)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=DVK-DEVOPS_DVK-Project&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=DVK-DEVOPS_DVK-Project)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=DVK-DEVOPS_DVK-Project&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=DVK-DEVOPS_DVK-Project)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=DVK-DEVOPS_DVK-Project&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=DVK-DEVOPS_DVK-Project)

## The Legacy Project

Peruse the documentation regarding running and installing the legacy project [here.](documentation\legacy_codebase\README.md)

**Below is the original readme file from the WhoKnows repository.**

## Whoknows Variations - Repository ReadMe

This is the Whoknows variations repository. It is not meant for production as it contains several security vulnerabilities and problematic parts on purpose.

### How to get started

Each branch is a tutorial in a different topic based on the same Flask application as in the `main` branch.

One way to follow along is by:

1. Forking the repository to your own account.

2. Cloning the repository to your local machine.

3. Checking out the branch you are interested in (e.g. `git checkout <branch_name>`).

4. Following the instructions in the README of the branch.

5. You can now push changes to your own repository.

### Pull requests

If you have any suggestions or improvements to the tutorials, feel free to open a pull request.

## Developing in a Dev Container

This repository includes a VS Code devcontainer configuration that uses the project's
`Dockerfile.dev` so you can run the Go app in a containerised environment similar to production.

To open the project in the devcontainer from VS Code: use the Command Palette -> "Remote-Containers: Open Folder in Container..." and choose the repository root.

What the devcontainer provides:

- Uses `Dockerfile.dev` to build the development image
- Forwards port 8080 from the container to the host
- Mounts the workspace into the container so code edits are live

If you prefer to run the same environment with docker-compose locally, there's a `docker-compose.dev.yml` which starts a service named `whoknows` bound to port 8080.
