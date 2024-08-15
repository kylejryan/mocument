# Mocument

Mocument is a mock implementation of AWS DocumentDB for use in unit and integration testing. It allows developers to simulate DocumentDB operations locally, enabling offline testing without the need for a real DocumentDB instance.

## Purpose

The purpose of mocument is to provide a robust, easy-to-use mock version of AWS DocumentDB. This allows developers to:

- Perform unit and integration tests without needing access to a real DocumentDB instance
- Simulate various scenarios including latency and errors
- Ensure code that interacts with DocumentDB can be tested thoroughly in a controlled environment

## Features

- Mock implementation of DocumentDB operations
- Configurable to simulate latency and errors
- Supports CRUD operations for documents within collections
- Easy to integrate into existing projects for testing purposes

## Installation

To install mocument, use the following command:

```sh
go get github.com/kylejryan/mocument
```

## Usage

### Basic Usage

Here's a basic example of how to use mocument in your project:

```go
package main

import (
    "context"
    "fmt"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/kylejryan/mocument/mock"
)

type MyEvent struct {
    Name string `json:"name"`
}

var dbClient *mock.MockDocDB

func init() {
    mockConfig := &mock.MockConfig{SimulateLatency: false, ErrorMode: false}
    dbClient = mock.NewMockDocDB(mockConfig)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event MyEvent) (string, error) {
    // Insert a document
    doc := map[string]interface{}{"name": event.Name}
    err := dbClient.InsertDocument("collection", doc)
    if err != nil {
        return "", fmt.Errorf("failed to insert document: %w", err)
    }

    // Find the document
    filter := map[string]interface{}{"name": event.Name}
    results, err := dbClient.FindDocument("collection", filter)
    if err != nil {
        return "", fmt.Errorf("failed to find document: %w", err)
    }

    return fmt.Sprintf("Found documents: %+v", results), nil
}

func main() {
    lambda.Start(Handler)
}
```

### Testing

You can use mocument to write unit tests for your code that interacts with DocumentDB. Here is an example of a test file:

```go
package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/kylejryan/mocument/mock"
	"github.com/stretchr/testify/assert"
)

var mockDBClient *mock.MockDocDB

func init() {
	mockConfig := &mock.MockConfig{SimulateLatency: false, ErrorMode: false}
	mockDBClient = mock.NewMockDocDB(mockConfig)
}

func mockHandler(ctx context.Context, event MyEvent) (string, error) {
	// Use the mock clients instead of the real AWS clients
	collection := mockDBClient

	// Insert a document
	doc := map[string]interface{}{"name": event.Name}
	err := collection.InsertDocument("collection", doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert document: %w", err)
	}

	// Find the document
	filter := map[string]interface{}{"name": event.Name}
	results, err := collection.FindDocument("collection", filter)
	if err != nil {
		return "", fmt.Errorf("failed to find document: %w", err)
	}

	return fmt.Sprintf("Found documents: %+v", results), nil
}

func TestHandler(t *testing.T) {
	// Create a fake event
	event := MyEvent{Name: "test"}

	// Call the handler with the mock client
	result, err := mockHandler(context.Background(), event)
	assert.NoError(t, err)
	assert.Contains(t, result, "Found documents")
}
```

## Contributing

We welcome contributions to mocument! If you'd like to contribute, please follow these steps:

1. Fork the repository on GitHub.
2. Clone your fork locally:

```sh
git clone https://github.com/your-username/mocument.git
```

3. Create a branch for your feature or bug fix:

```sh
git checkout -b feature-or-bugfix-branch
```

4. Commit your changes with a clear message:

```sh
git commit -am "Add new feature or fix bug"
```

5. Push to the branch:

```sh
git push origin feature-or-bugfix-branch
```

6. Create a pull request on GitHub, describing your changes.
   
Please make sure to write tests for your changes and ensure all existing tests pass.

## License

mocument is licensed under the MIT License. See the LICENSE file for more information.

## Acknowledgements

This project is inspired by the need to have reliable and efficient local testing environments for applications that interact with AWS DocumentDB.