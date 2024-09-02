# api-service
This repository contains the user-facing backend service for the minireipaz project. It serves as the core API that the frontend communicates with, providing essential functionalities and data management. Built with Go, it ensures high performance, scalability, and reliability to support various user interactions and workflows.

## Project Description

This project implements a workflow automation platform using architecture that combines Vercel Functions with Kafka Confluent and an HTTP connector. This solution is designed to overcome the limitations of Vercel's serverless functions, allowing for the execution of complex and long-running workflows.

## Architecture

### Main Components

1. **Frontend and Initial API**: Hosted on Vercel
2. **Vercel Functions**: For quick operations and initial request handling
3. **Kafka Confluent**: As a central event bus
4. **Confluent HTTP Connector**: To receive events from Vercel Functions
5. **Processing Services**: Kafka consumers to execute workflows
6. **ClickHouse**: As a database to store results and workflow data

### Workflow

1. Users interact with the frontend hosted on Vercel.
2. Quick requests are handled directly by Vercel Functions.
3. For tasks that require more than 10 seconds:
   - The Vercel Function sends an event to the Kafka HTTP connector.
   - The event is stored in a Kafka topic.
   - A consumer service processes the event and executes the workflow.
4. Results and states are stored in ClickHouse.

## Architecture Justification

### Vercel Functions Limitations

- **Maximum execution time**: 10 seconds
- **Lack of native triggers**: No support for scheduled or complex event-based executions

### Solution: Kafka Confluent with HTTP Connector

This architecture allows us to:
1. Simulate triggers by sending events to Kafka through an HTTP endpoint.
2. Handle long-running processing outside of Vercel Functions.

### Pros and Cons

#### Pros:

1. **Decoupling**: Separates initial business logic from long-running processing.
2. **Scalability**: Kafka can handle large volumes of events.
3. **Persistence**: Events are stored in Kafka, providing reliability.
4. **Flexibility**: Easy to add new event producers or consumers.
5. **Asynchronous processing**: Allows handling of long-running tasks.
6. **Compatibility**: Integrates well with existing Vercel-based infrastructure.

#### Cons:

1. **Additional complexity**: Introduces new components that need to be managed.
2. **Cost**: Kafka Confluent has associated costs.
3. **Latency**: There may be a small increase in latency.
4. **Learning curve**: Requires knowledge of Kafka and its ecosystem.
5. **Maintenance**: Needs additional configuration and maintenance.


```
graph TD
    A[user interacts] --> B[CDN Frontend React]
    B --> C[Frontend Vercel Function]
    C --> D[Backend Vercel Function]
    D --> E[REST Proxy]

    subgraph KAFKA BROKER
    F[Kafka Topic Broker Task Queue] --> G[Kafka HTTP Sink Connector]
    end

    subgraph KAFKA AS WORKFLOW DB
    H[Kafka Topic workflow Workflow DB] --> I[Kafka HTTP Sink Connector]
    end

    E --> H
    J[Scheduler Vercel Function] --> F
    J --> H
    G --> K[workers]
    I --> J
    K --> L[Results DB]
    L --> J

    L --> D
    
    %% Long Polling
    B -.->|Long Polling| C
    C -.->|Long Polling| D
    D -.->|Long Polling| L
```

### Errors
TODO

### Retries

#### Retries with General Error Classification

- **HTTP 4xx errors:** Generally not candidates for automatic retries, with specific exceptions.
- **HTTP 5xx errors:** Potential candidates for retries, but with caution.
- **Network or timeout errors:** Usually appropriate for retry attempts.

#### Retries with Error-Specific Retry Strategies

- **HTTP 429 (Too Many Requests):** Implement retries with exponential backoff, respecting "Retry-After" headers if present.
- **HTTP 500, 502, 503, 504:** Retry with exponential backoff.
- **Connection or timeout errors:** Immediate retry for the first couple of attempts, then apply backoff.
- **HTTP 400, 401, 403, 404, etc.:** No automatic retries; log for review.

#### Exponential Backoff Implementation

1. Initialize with a base delay (e.g., 1 second).
2. Apply exponential factor on each attempt: 1s, 2s, 4s, 8s, etc.
3. Incorporate random jitter to prevent retry synchronization.
4. Set a maximum delay threshold (e.g., 60 seconds).

#### Retry Attempts

- Define a maximum retry count (e.g., 3-5) based on the operation type.
- Consider the total acceptable time window for the operation.

#### Persistent Failure Handling

- After exhausting retry attempts, log detailed error information.
- Implement a circuit breaker pattern to prevent system overload.

#### Monitoring and Logging

- Log each retry attempt with relevant data (error type, attempt number, delay).
- Set up alerts for recurring error patterns.


## Makefile

This project utilizes a Makefile to automate common development and testing tasks.
Available commands are:

```
make openapi_http
make lint
make test
make fmt
```

`make openapi_http`: Generates HTTP API client code from the OpenAPI specification.
`make lint`: Analyzes the source code to identify and report stylistic errors and programming mistakes.
`make test`: Runs the test suite to ensure that the code is functioning correctly.
`make fmt`: Formats the source code according to the project's style.


### openapi_http

To generate the OpenAPI documentation, simply run the following command:

```
make openapi_http
```

makefile details:

```
.PHONY: openapi_http
openapi_http:
	@echo "Generating OpenAPI documentation for the HTTP API..."
	@./scripts/generate_openapi.sh
```

### lint

Lints the Go codebase using golangci-lint.

Installation (optional):

If you need to install `golangci-lint` and `go-cleanarch`, run the script with the -install option:
```
./scripts/lint.sh -install
```

This will:
- Download and install `golangci-lint` version v1.59.1.
- Display the installed version of `golangci-lint`.
- Install the `go-cleanarch` tool.

Script details:

This will execute `golangci-lint` with the specified options:
-v: Enable verbose output.
--tests=false: Exclude tests from linting.
--timeout=2m: Set the timeout to 2 minutes.
--config ./.golangci.yaml: Use the specified configuration file for `golangci-lint`.


### fmt
This command will format all Go files in the current directory and its subdirectories using `goimports`.
```
make fmt
```

makefile details:

```
goimports -l -w -d -v ./
```

-l: Lists the files that would be formatted.
-w: Writes the formatted content directly to the files.
-d: Displays the diffs instead of applying them.
-v: Verbose mode, providing detailed output.


### test

By running the make test command, both the unit and E2E test suites will be executed sequentially. 

To run the tests, execute the following command. The tests will automatically retrieve the necessary environment variables:
```
make test
```

makefile details:

```
test:
	@./scripts/test.sh .env
	@./scripts/test.sh .e2e.env
```

Environment Variables:

.env: This file contains the environment variables for the main test suite.
.e2e.env: This file contains the environment variables for the end-to-end (E2E) test suite.

`./scripts/test.sh`: This script is responsible for executing the tests. It reads the environment variables from the specified files and runs the tests accordingly.
