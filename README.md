[![Build Status](https://travis-ci.com/eco/longy.svg?token=QuNAGfYo3kcpqd58kfZs&branch=master)](https://travis-ci.com/eco/longy)

# Longy
A blockchain based game for SF Blockchain week. The game is run on it's own blockchain using the Cosmos-Sdk.

TODO: explain the game and why there is a "master node"

# Install and Test
To build the project:
```
make install
```

To lint and test:
```
./scripts/lint_and_test.sh
```

## Running database tests
The database tests depend on a local DynamoDB instance running on port 8000.
To enable the tests, first launch the DynamoDB service:
```
docker run --rm -p 8000:8000 amazon/dynamodb-local
```

Then run the suite in an environment with:
```
ENABLE_DB_TESTS=true
```
eg:
```
ENABLE_DB_TESTS=true ginkgo ./...
```
