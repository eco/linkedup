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
2. `/id/<id> [GET]` is a convenience endpoint to convert an badge id into a determinsistic cosmos address  
3. `/key [POST]` is the entry point for keying an account.  
  Request Body:  
  ```
  {
    "attendee_id": "<id>",
    "cosmos_private_key": "hex-encoded private key"
    "cosmos_rsa_key": "string representation of the rsa key"
  }
  ```  
  Status 200: Key transaction was successfully submitted and the email containing the redirect uri was sent  
  Status 400: Bad request body. Check the returned response  
  Status 401: The attendee has already keyed their account  
  Status 404: The attendee id was not found in the eventbrite event  
  Status 503: Another external component outside of the key service went wrong. (i.e email / backend storage / etc)  
    - The logs will contain information about what went wrong
  Status 500: Something internal went wrong. (i.e marshalling data)  
    - The logs will contain information about what went wrong

4. `/recover [POST]` will start the process of recovering an account  
  Request Body:  
  `badge id number`

  Status 200: An email was sent with a redirect uri to hit the following endpoint below with the authentication token to retrieve attendee information.  
  Status 400: Bad request body. Check the error response  
  Status 404: The attendee for the corresponding id was not found in eventbrite  
  Status 503: Another external component outside the service is down. Backend storage or email  
    - Check the logs

5. `/recover/<id>/<authtoken> [GET]` will retrieve complete attendee information that is stored  

  Status 200: Check the response body for the result  
  Status 404: No information found on this attendee  
  Status 401: Incorrect authentication token OR there is not authentication token set for this attendee id. Must start from the `/recover` endpoint above  

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

## Email Functionality
Email sending functionality depends on AWS API keys. When run deployed in AWS
these will be provided by the instance role. For development purposes you'll
need to provide your own. See the
[linkedup-content](https://github.com/eco/linkedup-content) repository for the
email templates - they can be installed using:
```
aws ses create-template --template "`cat path/to/rekey.json`" --region us-west-2
```
