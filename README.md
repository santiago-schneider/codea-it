# Go API Challenge

A simple Go application that uses SQLite for in-memory data storage and provides a RESTful API.

## Table of Contents

- [Description](#description)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Running Tests](#running-tests)
- [Dockerization](#dockerization)

## Description

This project is a RESTful API built with Go. It uses SQLite for data storage and is designed to run inside a Docker container. The API provides endpoints to interact with the data from Kraken API.

## Requirements

- Go 1.22 or higher
- Docker

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/santiago-schneider/codea-it.git
    cd codea-it
    ```

2. Install Go dependencies:
    ```sh
    go mod download
    ```

## Usage

1. Run the application locally:
    ```sh
    go run main.go
    ```

2. The application will start on the address defined by the `SERVER_ADDRESS` environment variable, defaulting to `:8080`.

## Running Tests

To run integration tests for the API, use the following command:
```sh
go test -v ./...
```

## Dockerization

1. Build the Docker image:
    ```sh
    docker build -t codea-it-app:latest .
    ``` 
2. Run the Docker container:
    ```sh
    docker run -p 8080:8080 -e SERVER_ADDRESS=":8080" myapi:latest
    ```
