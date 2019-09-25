[![Build Status](https://travis-ci.com/eco/longy.svg?token=QuNAGfYo3kcpqd58kfZs&branch=master)](https://travis-ci.com/eco/longy)

# Longy
A blockchain based game for SF Blockchain week. The game is run on it's own blockchain using the Cosmos-Sdk. In addition to the validator nodes
that are run for the duration of the game, the "master" node deploys an http `re-key service` to faciliate rekeying users.  


TODO: explain the game and why there is a "master node"

# The Rekey service
The `rekey-service` module defines a http service. There is a corresponding binary in `cmd/rks` that will start the service and defines the commandline
flags needed to start the service. The usage is copied below:
```
Usage:
  rks [flags]

Flags:
      --port int                 port to bind the rekey service (default 1337)
      --longy-masterkey string   master private key for the longy game.
      --eb-auth-token string     eventbrite authorization token
      --eb-event-id int          id associated with the eventbrite event
      --smtp-server string       host:port of the smtp server
      --smtp-username string     username of the email account
      --smtp-password string     password of the email account
```

The service defines two endpoints

1. `/ping` which simple returns a status 200 and a "pong" response. This endpoint' intended use is health checks.
2. `/rekey/<attendee_id>?nonce=<nonce>` main url path which will generate the rekey signature and send the email containing the appropriate redirect. The
attendee's email is fetched from eventbrite using `attendee_id`. The signature is over `SHA256("resetkey(id=<id>, nonce=<nonce>)")
  Response Codes:
  1. 200 - All went well!
  2. 403 - The attentee ID was not found in the event
  3. 500 - something went wrong. Check the logs. Email probably didn't send

## Install & Test
To build the project:
```
make install
```

To lint and test:
```
./scripts/lint_and_test.sh
```
