# Kev

> Develop Kubernetes apps iteratively with Docker-Compose

![Stability:Beta](https://img.shields.io/badge/stability-beta-orange)
![CircleCI](https://img.shields.io/circleci/build/github/appvia/kev/master)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/release/appvia/kev)](https://github.com/appvia/kev/releases/tag/v0.2.0)
![License: Apache-2.0](https://img.shields.io/github/license/appvia/kev)


Kev helps developers port and iterate Docker Compose apps onto Kubernetes. It understands the Docker Compose application topology and prepares it for deployment in (multiple) target environments, with minimal user input.

We leverage the [Docker Compose](https://docs.docker.com/compose/compose-file/) specification and allow for target-specific configurations to be applied to each component of the application stack, simply.

Kev is opinionated in its choice of Kubernetes elements you should be able to control. It automatically infers key config parameters by analysing and reconciling changes in the project source compose file(s). The configuration parameters can be manually overridden to allow for better control of a cloud application deployment on Kubernetes.

Kev reduces the need for Kubernetes expertise in the team. The generated Kubernetes deployment configuration follows best industry practices, with a thin layer of config options to enable further control. See [kev reference documentation](docs/reference/config-params.md) for a list of available options.

## Features

* **Simplicity** - Based on the familiar Docker Compose specification. There is no new framework to learn, no new specification to embrace, and vastness of Kubernetes reduced to a limited set of easy to follow configuration parameters. You focus on the app development. Kev will prepare it for deployment in Kubernetes.

* **Multi-environment support** - Parameterisation enabled with the same configuration primitives you're already familiar with. Each defined environment gets its own docker-compose override file, which is there to control the behaviour of your application in Kubernetes in a simple and consistent way.

* **Best practice out of the box** - Best practice is codified and embedded in the translation layer, so you don't have to think about what's required to run your project application on Kubernetes.

* **Secure** - Kev is opinionated about the secret management. At this stage of its relatively short life it delegates that responsibility to the user, to remove the risk of potential uncontrolled secrets leak. No secrets == No leaks!

* **No vendor lock-in** - Because you already use docker-compose, you can keep using it, even if Kev turns out to be not your cup of tea.

* **Easy integrations** - You may use generated Kubernetes manifests with any tool / framework of your choice. Check out our [Skaffold Integration](docs/tutorials/kev-dev-with-skaffold.md).

## Contents

- **[Installation](#installation)**
- **[Quickstart](docs/tutorials/quickstart-guide.md)**
    * **[Initialise project](docs/tutorials/quickstart-guide.md#initialise-project)**
    * **[Generate Kubernetes manifests](docs/tutorials/quickstart-guide.md#generate-kubernetes-manifests)**
- **[Tutorials & guides](#tutorials-and-guides)**
- **[Configuration](#configuration)**
- **[Similar tools](#similar-tools)**
- **[Contributing to Kev](#contributing-to-kev)**
- **[Roadmap](#roadmap)**
- **[License](#license)**

## Installation

All you need to get started quickly is the [kev](https://github.com/appvia/kev/releases) binary added to your PATH, and one or more docker compose files.

## Tutorials and guides

- [How does Kev differ from Kompose?](docs/tutorials/how-kev-differs-from-kompose.md)
- [Getting started with Kev](docs/tutorials/getting-started-with-kev.md)
- [Develop the app with Kev and Skaffold](docs/tutorials/kev-dev-with-skaffold.md)
- [Simple Node.js app development workflow example](docs/tutorials/simple-nodejs-app-workflow.md)
- [Simple Node.js app CI workflow example](docs/tutorials/simple-nodejs-app-ci-workflow.md)

## Configuration

As mentioned in the [Quickstart](docs/tutorials/quickstart-guide.md), the environment specific configuration lives in a set of docker-compose override files. Each environment override file holds simplified Kubernetes configuration parameters for each of the application components.

Project components (aka services) are configured via a set of labels attached to them, and optionally environment variables section which allows for localised adjustments - the same exact way you'd control those in a regular docker-compose file.

Volumes come with their own set of labels to control Kubernetes storage specific parameters.

See the [configuration reference](docs/reference/config-params.md) for details.

## Similar tools

Kev is inspired by the simple, easy to use and well adopted Docker Compose specification, as well as several other tools in the Kubernetes manifests generation and templating space such as Kompose, Ksonnet and Kustomize, to name a few.

See how [Kev differs from Kompose.](docs/tutorials/how-kev-differs-from-kompose.md)

Each of the solutions above, however, come with their own set of challenges and are lacking in various areas. Some have been discontinued, some see very few contributions or updates, others require a great deal of prior Kubernetes expertise to find them useful.

Kev bridges the gaps in the existing tooling, helping developers familiar with Docker & Compose to easily get up and running with Kubernetes.

## Contributing to Kev

We welcome any contributions from the community! Have a look at our [contribution](CONTRIBUTING.md) guide for more information on how to get started. If you use Kev, find it useful, or are generally interested in improving Developer Experience on Kubernetes then please let us know by **Starring** and **Watching** this repo. Thanks!

## Roadmap

See our [Roadmap][roadmap] for details about our plans for the project.

## License

Copyright (c) 2020-2021 [Appvia Ltd](https://appvia.io)

This project is distributed under the [Apache License, Version 2.0](./LICENSE).

[roadmap]: https://github.com/appvia/kev/issues
