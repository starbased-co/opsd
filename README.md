# Docker OnePassword Secrets Plugin
[![](https://badgen.net/docker/pulls/mrmarble/opsd?icon=docker)](https://hub.docker.com/r/mrmarble/opsd)
[![license](https://badgen.net/static/license/MIT/blue)](LICENSE)

This project provides a Docker plugin to use OnePassword as a secrets provider. It allows Docker containers to securely access secrets stored in OnePassword. The plugin uses the OnePassword Connect API to retrieve secrets and provide them to Docker containers.

## Prerequisites

- OnePassword (Obviously)
- Docker swarm mode (for using Docker [secrets](https://docs.docker.com/engine/swarm/secrets/))

## Configuration

The plugin can be configured using the following environment variables:

- `OP_CONNECT_HOST`: The OnePassword Connect  (hostdefault: `http://localhost:8080`)
- `OP_CONNECT_TOKEN`: The OnePassword Connect token
- `OP_VAULT_NAME`: The OnePassword vault used for secrets (default: `docker`)


## Installation

1. Set up the OnePassword Connect and Sync services using Docker Compose, follow the instructions in the [OnePassword Connect documentation](https://developer.1password.com/docs/connect/get-started?deploy=docker).

    ```sh
    docker-compose up -d # There is a docker-compose.yml file in the root of this repository
    ```

2. Install the plugin

    ```sh
    docker plugin install mrmarble/opsd:latest OP_CONNECT_HOST=<one password connect api host> OP_CONNECT_TOKEN=<your_token> OP_VAULT_NAME=<vault where secrets are stored>
    ```

## Usage

1. Create a secret in OnePassword

    ```sh
    op item create --category=password --title=my-app-secrets --vault=docker 'MY_SECRET[password]=supersecretpassword'
    ```
2. Create a Docker secret using the plugin

    ```sh
    docker secret create --driver mrmarble/opsd:latest -l item=my-app-secrets MY_SECRET
    ```
3. Use the secret in a service

    ```sh
    docker service create --secret MY_SECRET --name my-app my-app-image
    ```


## License

This project is licensed under the MIT License. See the LICENSE file for details.
