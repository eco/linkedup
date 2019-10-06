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
make test
```

# Key Service
The key service runs alongside `lyd` and `lycli` to facilitate keying accounts and email onboarding. The key service hosts
two http endpoints

1. `/ping [GET]` is a health check. Simply writes a Status 200 along with "pong" in the request body
2. `/key [POST]` is the entry point for keying an account.  
  Request Body:  
  ```
  {
    "attendee_id": "<id>",
    "pubkey": {
      "type":"tendermint/PubKeySecp256k1",
      "value":"Aq4PU4ws0ozmtkAmKv8y9Fs5FLbXnoJPDJRBHnimin62"
      }
  }
  ```  
  Status 200: Key transaction was successfully submitted and the email containing the redirect uri was sent
  Status 403: The attendee id was not found in the eventbrite event
  Status 500: Something internal went wrong. (Communicating with eventbrite, sending an email, the transaction failed to sent).  
    - The logs will contain information about what went wrong

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
ENABLE_DB_TESTS=true
make test
```
