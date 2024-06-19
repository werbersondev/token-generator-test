# Technical Challenge - SonarQube Token Generation Microservice

## Overview

This project is a technical challenge, focused on generating SonarQube scoped access tokens via a microservice architecture using a messaging bus. The microservice consists of two main components:
1. An HTTP API that accepts requests to generate tokens.
2. A consumer that processes these requests and generates tokens using the SonarQube API.

## Architecture

It consists of:

1. **SonarQube Community Edition**: Deployed locally using Docker.
2. **GCP Pub/Sub Emulator**: Deployed locally using Docker to simulate Google Cloud Pub/Sub.
3. **HTTP Service**: An API to receive token generation requests.
4. **Consumer Service**: Subscribes to Pub/Sub topics to process token generation requests.

## Setup and Installation

### Prerequisites

- Docker
- Go 1.22+

### Steps

1. **Clone the repository**

   ```sh
   git clone https://github.com/werbersondev/token-generator-test.git
   cd token-generator-test
   ```

2. **Setup Local Dependencies**

   ```sh
   make setup/local-dep
   ```

   This command will bring up the SonarQube and Pub/Sub emulator containers.

3. **Install Go to be ready to run**

   ```sh
   make setup
   ```
4. **Running the Services**

   > Any environment variable can be included or changed in the `.env` file at the root path.

   **HTTP Service**
   
   To run the HTTP service locally:
   
   ```sh
   make run/http
   ```

   **Consumer Service**
   
   To run the consumer service locally, you need to provide the `SONAR_AUTH_TOKEN`:
   
   ```sh
   make run/worker SONAR_AUTH_TOKEN=your_sonar_auth_token
   ```
   
   This ensures that the `SONAR_AUTH_TOKEN` is required and is passed correctly when running the worker.

## Configuration

### HTTP Service

| Environment Variable        | Description                        | Default Value            |
|-----------------------------|------------------------------------|--------------------------|
| `SERVER_ADDR`               | Address for the HTTP server        | `0.0.0.0:3000`           |
| `SERVER_READ_TIMEOUT`       | Read timeout for the HTTP server   | `30s`                    |
| `SERVER_WRITE_TIMEOUT`      | Write timeout for the HTTP server  | `30s`                    |
| `PUBSUB_EMULATOR_HOST`      | Host for the Pub/Sub emulator      | (required)               |
| `GCP_PROJECT_ID`            | GCP project ID                     | `my_project_key`         |
| `GCP_TOKEN_GENERATOR_TOPIC` | Pub/Sub topic for token generation | `token_generation_topic` |

### Consumer Service

| Environment Variable                    | Description                                 | Default Value                   |
|-----------------------------------------|---------------------------------------------|---------------------------------|
| `PUBSUB_EMULATOR_HOST`                  | Host for the Pub/Sub emulator               | (required)                      |
| `GCP_PROJECT_ID`                        | GCP project ID                              | `my_project_key`                |
| `GCP_TOKEN_GENERATOR_TOPIC`             | Pub/Sub topic for token generation          | `token_generation_topic`        |
| `GCP_TOKEN_GENERATOR_SUBSCRIPTION`      | Pub/Sub subscription for token generation   | `token_generation_subscription` |
| `SONAR_API_ADDRESS`                     | Address for the SonarQube API               | `http://localhost:9000`         |
| `SONAR_API_TIMEOUT`                     | Timeout for SonarQube API requests          | `30s`                           |
| `SONAR_AUTH_TOKEN`                      | Authentication token for SonarQube API      | (required)                      |

## HTTP API Documentation

### Request Token Generation Endpoint

#### Endpoint

`POST /generate-token`

#### Request Body

The endpoint expects a JSON object with the following structure:

```json
{
  "project_id": "your_project_id"
}
```

- `project_id` (string): The ID of the SonarQube project for which the token is to be generated. This field is required.

#### Response

- **202 Accepted**: The request to generate a token has been accepted and will be processed.
- **400 Bad Request**: The request body is invalid.
- **422 Unprocessable Entity**: The `project_id` parameter is missing.
- **500 Internal Server Error**: Failed to publish the message to the Pub/Sub topic.

#### Example

```sh
curl -X POST http://localhost:3000/generate-token \
     -H "Content-Type: application/json" \
     -d '{"project_id": "your_project_id"}'
```

## Testing

### Unit Tests

To run the unit tests:

```sh
make test
```

### Test Coverage

To run the tests with coverage:

```sh
make test-coverage
```

## Makefile Commands

| Command                | Description                                        |
|------------------------|----------------------------------------------------|
| `make download`        | Download Go module dependencies                    |
| `make install-tools`   | Install required tools like `moq`                  |
| `make setup`           | Install tools and tidy up Go modules               |
| `make test`            | Run all tests                                      |
| `make test-coverage`   | Run tests with coverage analysis                   |
| `make generate`        | Run `go generate`                                  |
| `make setup/local-dep` | Setup local environment with Docker Compose        |
| `make run/http`        | Run the HTTP service                               |
| `make run/worker`      | Run the worker service (requires SONAR_AUTH_TOKEN) |

---
