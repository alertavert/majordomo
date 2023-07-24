# Majordomo - Code generation Bot
*Created by M. Massenzio &copy; 2023 AlertAvert.com All Rights Reserved*

# Build & run Server

`make help` shows all the available commands.

`make build` and `make test` do what one expects they would, the server is built in the `build/bin` folder, and tagged with the current `version` (derived from the `settings.yaml` and the current git SHA).

The server itself can be run in `dev` mode via `make run` and is reachable at `http://localhost:5000`; the `UI` is served off the `/web` endpoint.

## API

`TODO`

## OpenAI Interface

## Backend Architecture

# App UI

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

You will need `npm` to run all this, on MacOS it's just a matter of `brew install npm` and on Linux `sudo apt install npm` should work - although, with this junk, it's always a coin toss.

## Available Scripts

### Development

To install the necessary modules run:

```shell
make ui-setup
```


A `dev` version (which reloads when any of the Javascript files in the `ui/app/src` folder are modified) can be run with:

```
make ui-dev
```
and is served off `http://localhost:3000` (we use `CORS` to allow it to connect to the backend Gin server running on port `5000`).

The page will reload when you make changes.
You may also see any lint errors in the console.

### Distribution

The React App can be built via `make ui` (which will install a `webpack` minified distribution in the `build/ui` folder of the project.

The build is minified and the filenames include the hashes.
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

`TODO: update the following -- there are currently no tests on the UI`

### `npm test`

Launches the test runner in the interactive watch mode.
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.


