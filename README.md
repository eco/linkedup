# LinkedUp
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-green.svg)](https://github.com/RichardLitt/standard-readme)
[![Build Status](https://travis-ci.com/eco/longy.svg?token=QuNAGfYo3kcpqd58kfZs&branch=master)](https://travis-ci.com/eco/longy)
> A blockchain based game for SF Blockchain week. The game is run on it's own blockchain using the Cosmos-Sdk.

TODO: explain the game and why there is a "master node"

## Table of Contents
 - [Background](#background)
 - [Usage](#usage)
 - [Contributing](#contributing)
 - [License](#license)

## Background
LinkedUp is a networking game built for SF Blockchain Week 2019. The game runs
on a Cosmos blockchain, and rewards players for establishing social connections
between each other at the conference.

The game addresses two challenges - player onboarding and wallet installation.
Players use ephemeral wallets to play the game, and a key escrow service
provides key recovery as well as facilitating player onboarding. This escrow
service also provides administrative functions useful for day-of operations.
These functions should become less necessary as we learn more about how the
game works in practice.

## Usage
### Running in Docker
To get up and running with docker:
```
docker-compose up --build
```

LocalStack may depend on the presence of a `/tmp/localstack` directory on your
system. If needed, ensure its presence:
```
mkdir -p /tmp/localstack
```

Make sure to edit `/etc/hosts` to make localstack an alias for 127.0.0.1
```
127.0.0.1	localhost
127.0.0.1	localstack
255.255.255.255	broadcasthost
```
### Email Functionality
SES-based email sending functionality depends on AWS API keys. When run deployed
in AWS these will be provided by the instance role. For development purposes
you'll need to provide your own. See the
[linkedup-content](https://github.com/eco/linkedup-content) repository for the
email templates - they can be installed using:
```
aws ses create-template --template "`cat path/to/rekey.json`" --region us-west-2
```

## Contributing
### Build and Test
To build the project:
```
make all
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

### Running the Key Service
The key service runs alongside `lyd` and `lycli` to facilitate keying accounts and email onboarding. The key service hosts
two http endpoints

The usage for the key service is shown below
```
Usage:
  ks [flags]

Flags:
      --port int                    port to bind the rekey service (default 1337)
      --longy-chain-id string       chain-id of the running longy game (default "longychain")
      --longy-masterkey string      hex encoded master private key (default "fc613b4dfd6736a7bd268c8a0e74ed0d1c04a959f59dd74ef2874983fd443fca")
      --longy-restservice string    scheme://host:port of the full node rest client (default "http://localhost:1317")
	  --longy-app-url              scheme://host of the client web app

      --eventbrite-auth string      eventbrite authorization token
      --eventbrite-event int        id associated with the eventbrite event

      --aws-content-bucket string   content bucket for user uploads (default "linkedup-user-content")
      --email-mock                  print email URLs instead of emailing
      --localstack                  use localstack instead of aws; implies --email-mock
```

The configruation can also be set through environment variables. the `-` characters replaced by `_` and all uppercase.  
   i.e `STMP_SERVER` or `EVENTBRITE_AUTH`

#### Email Data Testing
Running the key service with the `--email-mock` flag will cause email template
parameters to be logged instead of sent to an email system.

#### Using LocalStack
Running the key service with the `--localstack` flag will cause AWS-backed
drivers to look for LocalStack services on the `localstack` host instead. It
implies `--email-mock`.

#### Key Service API
**Quick & Convenient Endpoints**  
`/ping [GET]` health check. Simply writes a Status 200 along with "pong" in the request body  
`/id/<id> [GET]` converts an badge id into a determinsistic cosmos address  

**Main Endpoints**  
`/key [POST]` is the entry point for keying an account.  
    **Request Body**:  
    ```
    {
      "attendee_id": "<id>",
      "cosmos_private_key": "hex-encoded private key"
      "cosmos_rsa_key": "string representation of the rsa key"
      "use_verification": true|false
    }
    ```  
    If `use_verification` is set, the email will contain a 6-digit verification code that can then be used to retrieve all the stored
    attendee information using `/recover/<id>/<token>`  
  
    Status 200: Key transaction was successfully submitted and the email (redirect or verification).  
    Status 400: Bad input to this endpoint. Request body could be malformed. failed conversion for `cosmos_private_key` to a valid secp256k1 key.  
    Status 404: The attendee id was not found in the eventbrite event.  
    Status 409: This attendee has already gone through the onboarding flow and has information stored.  
    Status 503: Another external component outside of the key service went wrong. (i.e email / backend storage / etc).  
      - The logs will contain information about what went wrong  
    Status 500: Something internal went wrong. (i.e marshalling data)  
      - The logs will contain information about what went wrong  

`/recover [POST]` will start the process of recovering an account  
    **Request Body**:  
    ```
    {
      "attendee_id": "<id>",
      "use_verification": true|false
    }
    ```  
    If `use_verification` is set, the email will contain a 6-digit verification code that can then be used to retrieve all the stored
    attendee information using `/recover/<id>/<token>`  
  
    Status 200: Recovery email with the token was sent succesfully.  
    Status 400: Bad request body. Check the error response.  
    Status 404: The attendee for the corresponding id was not found in eventbrite.  
    Status 503: Another external component outside the service is down. Backend storage or email.  
      - Check the logs  

`/recover/<id>/<authtoken> [GET]` will retrieve complete attendee information that is stored  
    **Response Body**:
    ```
    {
      "address": "<account adddress>",
      "attendee": {
          "id": <badgeID>,
          "first_name": "first name",
          "last_name": "last name",
          "email": "email address"
      },
  
      "cosmos_private_key": "hex encoded private key bytes"
      "rsa_private_key": ".."
  
      "commitment": "commitment set in the key message",
      "commitment_secret": "pre-image to the commitment",
      "image_upload_url": "url to upload avatar to"
    }
    ```  
    Status 200: Check the response body for the result.  
    Status 404: No information found on this attendee.  
    Status 401: Incorrect authentication token OR there is not authentication token set for this attendee id.  

### Running database tests
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

## License
MIT
