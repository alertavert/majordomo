# Majordomo - Code generation Bot
*Created by M. Massenzio &copy; 2023-2025 AlertAvert.com All Rights Reserved*

[![Go](https://github.com/alertavert/majordomo/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/alertavert/majordomo/actions/workflows/test.yml)

![Bots, bots everywhere](docs/images/team-of-robots.jpeg)

# Overview

This project uses OpenAI Assistants to provide a locally-running code assistant, with access to project files, and allows the Assistant to modify the files and return the modified version in a pre-defined location; we don't overwrite directly the original source files because, well, hallucinations.

The project is described in [Using OpenAI GPT to build a Coding Assistant that uses OpenAI GPT to build apps (including itself)](https://codetrips.com/2023/11/19/using-openai-gpt-to-build-a-coding-assistant-that-uses-openai-gpt-to-build-apps-including-itself/).


# Majordomo API

## Build & run Server

`make help` shows all the available commands.

`make build` and `make test` do what one expects they would, the server is built in the `build/bin` folder, and tagged with the current `version` (derived from the `settings.yaml` and the current git SHA).

Run a development instance using `make dev`, the server is available at `http://localhost:5000`.

## Run Integration Tests

To run the integration tests, you need to create a `.env.test.local` file in the project root with your OpenAI API key:

```
OPENAI_API_KEY=your_api_key_here
```

Then run the tests with:

```shell
make integration_tests
```

This will execute all integration tests that interact with the OpenAI API. If the `.env.test.local` file is missing or doesn't contain a valid API key, the tests will fail with an appropriate error message.

## Docker & Kubernetes

The container can be created with `make container` the image name and version are determined automatically (the version will match what is in [`settings.yaml`](settings.yaml)), something like:

    alertavert/majordomo:0.6.1

and can be run via the `make start` command.

To deploy the container to a Kind cluster running locally, run the following commands:

```shell
kind create cluster --name dev
kind --name=dev load docker-image alertavert/majordomo:0.6.1
kubectl create ns majo
kubectl create -n majo configmap majordomo-config --from-file=$HOME/.majordomo/config.yaml
kubectl apply -f deploy.yaml
```

The service is then accessible from within the cluster at:

    http://majordomo-service.majo.svc.cluster.local

## API

`TODO`

These are currently the endpoints:

```
GET    /health
POST   /command
POST   /parse
POST   /prompt
GET    /projects
GET    /projects/:project_name
GET    /projects/:project_name/sessions
POST   /projects
PUT    /projects
PUT    /projects/:project_name
DELETE /projects/:project_name
GET    /assistants
```

Load [the Postman collection](docs/Majordomo.postman_collection.json) into [Postman]() to see example API calls and the format of the JSON body.

## OpenAI Interface

`TODO`

## Backend Architecture

`TODO`

# Majordomo UI

The web app is created using [Streamlit](https://streamlit.io) and is in the [`webapp`](webapp) folder.

<img src="docs/images/streamlit.png" alt="Streamlit UI" height="500px">

and can be run using:

    streamlit run webapp/app.py debug

`streamlit` needs to be installed in a virtualenv using `pip install streamlit`.
