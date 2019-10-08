[![Build Status](https://travis-ci.com/eco/longy.svg?token=QuNAGfYo3kcpqd58kfZs&branch=master)](https://travis-ci.com/eco/longy)

# Longy
A blockchain based game for SF Blockchain week. The game is run on it's own blockchain using the Cosmos-Sdk.

TODO: explain the game and why there is a "master node"

# Install and Test
To build the project:
```
make install
```

Before you can run the tests you'll need to set up the runtime environment so
the system can connect with Eventbrite. To do that, export the environment
variable `EVENTBRITE_AUTH` containing your Eventbrite access token. In addition,
the system needs to know which Eventbrite event to work with. Set that in
`EVENTBRITE_EVENT`.

```
export EVENTBRITE_EVENT="3414213431"
export EVENTBRITE_AUTH="ewifkjaweklfheaklj"
```

To lint and test:
```
make test
```

# Key Service
The key service runs alongside `lyd` and `lycli` to facilitate keying accounts and email onboarding. The key service hosts
two http endpoints

The usage for the key service is shown below
```
key service for the longest chain game

Usage:
  ks [flags]

Flags:
      --port int                   port to bind the rekey service (default 1337)

      --aws-dynamo-url string      dynamodb url (default "http://localhost:8000")
      --aws-region string          aws region for dynamodb (default "us-west-1")

      --eventbrite-auth string     eventbrite authorization token
      --eventbrite-event int       id associated with the eventbrite event

      --longy-chain-id string      chain-id of the running longy game (default "longychain")
      --longy-fullnode string      tcp://host:port the full node for tx submission (default "tcp://localhost:26657")
      --longy-masterkey string     hex encoded master private key (default "fc613b4dfd6736a7bd268c8a0e74ed0d1c04a959f59dd74ef2874983fd443fca")
      --longy-restservice string   scheme://host:port of the full node rest client (default "http://localhost:1317")

      --smtp-password string       password of the email account (default "2019longygame")
      --smtp-server string         host:port of the smtp server (default "smtp.gmail.com:587")
      --smtp-username string       username of the email account (default "testecolongy@gmail.com")
      ```

The configruation can also be set through environment variables. the `-` characters replaced by `_` and all uppercase.  
   i.e `STMP_SERVER` or `EVENTBRITE_AUTH`


1. `/ping [GET]` is a health check. Simply writes a Status 200 along with "pong" in the request body  
2. `/key/<email> [GET]` will retrieve the hex-encoded private key associated with the badge id  
  Status 200: The body will contain the hex-encoded private key
  Status 404: The email was not found in the database  
3. `/key [POST]` is the entry point for keying an account.  
  Request Body:  
  ```
  {
    "attendee_id": "<id>",
    "private_key": "hex-encoded private key"
  }
  ```  
  Status 200: Key transaction was successfully submitted and the email containing the redirect uri was sent
  Status 401: The attendee has already keyed their account
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
