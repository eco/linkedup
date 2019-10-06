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

The usage for the key service is shown below
```p
key service for the longest chain game

Usage:
  ks [flags]

Flags:
      --port int                   port to bind the rekey service (default 1337)

      --eb-auth-token string       eventbrite authorization token
      --eb-event-id int            id associated with the eventbrite event

      --longy-chain-id string      chain-id of the running longy game (default "longychain")
      --longy-fullnode string      tcp://host:port the full node for tx submission (default "tcp://localhost:26657")
      --longy-masterkey string     hex encoded master private key for the longy game (default "fc613b4dfd6736a7bd268c8a0e74ed0d1c04a959f59dd74ef2874983fd443fca")
      --longy-restservice string   scheme://host:port of the full node rest client (default "http://localhost:1317")

      --smtp-password string       password of the email account (default "2019longygame")
      --smtp-server string         host:port of the smtp server (default "smtp.gmail.com:587")
      --smtp-username string       username of the email account (default "testecolongy@gmail.com")
      ```

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
