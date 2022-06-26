# Golang Project

## Getting started

1. run default configuration

    ```bash
        go run main.go
    ```

2. run custom configuration

    ```bash
        go run main.go -c template.json
    ```

3. run from environment variable

    ```bash
        TEST_SECRETKEY="test" go run main.go -c template.json
    ```