# Lukcy 38

A digital croupier, _"dealing cards as a meditation"_... and for job interviews.

---

## Table of Contents

- Getting Started
    - [Requirements](#requirements)
    - [Set up](#set-up)
    - [Running tests](#running-tests)
- Usage
    - [Endpoints](#endpoints)

---

## Getting Started

### Requirements

- [Docker](https://www.docker.com/products/docker-desktop) **OR** Go ^1.15 and PostgreSQL 13.2
- REST client (i.e. [Postman](https://www.postman.com), [Paw](https://paw.cloud))

## Set up

There are a few easy steps you have to follow in order to complete the set up and get the best croupier you can ever
have up and running:

- Create and update Docker .env file for Docker to figure out where your project folder lies.
- Create and update project .env files for the digital croupier to figure out which ~~blackjack table~~ database it has
  to work with.

#### Setting up Docker environment variables

1. Create `/deployment/.env` from `/deployment/.env.dist`
2. Change **PROJECT_PATH** with the full path of the project folder.

#### Setting up project environment variables

1. Create `.env` and `.env.test` from `.env.dist` at the same level with the file.
2. If you're using project defaults:
    - Leave .env file as it is
    - Update **DB_SOURCE** in **.env.test** with _postgresql://db_admin:admin321@db/lucky_test?sslmode=disable_

### Using the Makefile

Everybody has the best of luck when it comes to running commands with Lucky 38! With the makefile included in the
project, you can run the following commands with make, without having to worry about flags and options!

|Command|Description|
|-------|-----------|
| test | Runs tests in Docker container |
| server-run | Builds docker images as well as the app, then runs it |
| server-stop | Takes everything down |
| server-restart | For your convenience, Lucky 38 comes with a restart command! Which basically just runs _" server-stop"_ and _"server-run"_, one after the other. |

### Running tests

Tests under pkg/storage include DB integration tests and require a DB connection. Hence, it's highly recommended that
you run the tests in Docker. To do so, open your favorite command line prompter, change the current directory to
project's and run `make test`.

If you want to run tests on your machine, make sure you:

1. Update **DB_SOURCE** in **.env.test** with `postgresql://db_admin:admin321@localhost/lucky_test?sslmode=disable`
2. Run the command `APP_ENV=test go test -v ./pkg/...`

---

## Usage

When everything is up and running, you'll find your digital croupier available at `localhost:3000`. Below are the endpoints and their descriptions. There
is also a Postman and Paw collection available in the project for our VIP gamblers as yourself can benefit.

### Endpoints

#### Health

Call this endpoint if you want to check on your digital croupier.

- URL: /health
- Method: GET
- Response:

```json
"Why, hello!"
```

#### Create Deck

Creates a deck shuffled or in order; full or partial.

- URL: /deck
- Method: POST
- Body: `{ "shuffled": true|false }`
- Query string: cards (optional) Ex: http://localhost:3000/deck?cards=AS,2S,3D
- Response:

```json
{
  "deck_id": "008e2cbf-5c1b-4956-b7f6-40f68792b6cb",
  "shuffled": true,
  "remaining": 4
}
```

#### Open Deck

Returns the requested deck and available cards in it.

- URL: /deck/:id
- Method: GET
- Parameters:
    - id (required): Deck ID
- Response example:

```json
{
  "deck_id": "008e2cbf-5c1b-4956-b7f6-40f68792b6cb",
  "shuffled": true,
  "remaining": 4,
  "cards": [
    {
      "code": "2D",
      "value": "2",
      "suit": "DIAMONDS"
    },
    {
      "code": "AC",
      "value": "ACE",
      "suit": "CLUBS"
    },
    {
      "code": "KH",
      "value": "KING",
      "suit": "HEARTS"
    },
    {
      "code": "3D",
      "value": "3",
      "suit": "DIAMONDS"
    }
  ]
}
```

#### Draw Card

Draws cards from the deck and returns them.

- URL: /deck/:id/draw/:amount
- Method: PUT
- Parameters:
    - id (required): Deck ID
    - amount (required): How many cards to draw from the deck
- Response:

```json
[
  {
    "value": "3",
    "suit": "DIAMONDS",
    "code": "3D"
  },
  {
    "value": "ACE",
    "suit": "CLUBS",
    "code": "AC"
  }
]
```

