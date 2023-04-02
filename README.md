# Hashcash - Proof of Work

## Test task for Server Engineer

### Description

Design and implement “Word of Wisdom” tcp server.

### Requirements

- TCP server should be protected from DDOS attacks with the Prof of Work
  (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other
  collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge.

---

### Solution

#### Proof of Work

Proof of Work (PoW) is a piece of data which is difficult (costly, time-consuming) to produce but easy for others to
verify and which satisfies certain requirements. Producing a proof of work can be a random process with low probability
so that a lot of trial and error is required on average before a valid proof of work is generated. Bitcoin uses the
Hashcash proof of work system.

For the solution of the task I chose the following algorithm:

- Hashcash (https://en.wikipedia.org/wiki/Hashcash) - a proof-of-work algorithm that uses a cryptographic hash function
  to create a one-way function. The hashcash algorithm is used to prevent spam in email and Usenet newsgroups, and to
  prevent denial-of-service attacks on websites.
- SHA-256 (https://en.wikipedia.org/wiki/SHA-2) - a set of cryptographic hash functions designed by the United States
  National Security Agency (NSA) and published by the NIST as a U.S. Federal Information Processing Standard (FIPS).
  SHA-256 is one of the four algorithms in the SHA-2 set, all of which produce 256-bit hashes.

#### Server

The server is written in Go. It uses the standard library for the network and the crypto/sha256 package for the hash
function. The server is started with the following command:

##### Environment variables

| Name                  | Description                              | Default        |
|-----------------------|------------------------------------------|----------------|
| LOG_LEVEL             | Logger level, enum: debug, info, error   | info           |
| SERVER_ADDR           | Address of the server                    | localhost:8080 |
| SERVER_TTL_CONNECTION | How long server will be using connection | 1m             |
| HEADER_DIFFICULTY     | Difficulty to calculating header         | 5              |
| HEADER_TTL            | How long header will be alive            | 10m            |

##### Run

```bash
make server
```

The server listens on port 8080. The server accepts connections and reads the data from the client. The server checks
the correctness of the data and sends the quote to the client. The server is protected from DDOS attacks with the Prof
of Work. The server sends a challenge to the client and waits for the response. If the response is correct, the server
sends the quote to the client. If the response is incorrect, the server broke the connection with client. The server
sends the quote to the client only after the correct response is received.

#### Client

The client is written in Go. It uses the standard library for the network and the crypto/sha256 package for the hash
function. The client is started with the following command:

##### Environment variables

| Name                     | Description                                                 | Default        |
|--------------------------|-------------------------------------------------------------|----------------|
| LOG_LEVEL                | Logger level, enum: debug, info, error                      | info           |
| SERVER_ADDR              | Address of the server                                       | localhost:8080 |                    
| CLIENT_TIMEOUT           | How long clients will be using connection                   | 5s             |
| HEADER_MAX_ITERATIONS    | Count of iterations for calculating header (<0 is infinity) | -1             |
| HEADER_TTL               | Time to live header                                         | 10m            |
| CLIENT_DDOS_MODE         | Is need to ddos to the server?                              | false          |
| CLIENT_DDOS_CLIENT_COUNT | How many clients will be in DDOS                            | 100            |
| CLIENT_DDOS_TIMEOUT      | How long DDOS will be?                                      | 1m             |
| CLIENT_DDOS_WAITING      | If receive an error, how long need to wait for repeat       | 1s             |

##### Run

```bash
make client
```

The client connects to the server and sends the challenge to the server. The client calculates the hash of the challenge
and the nonce. If the hash of the challenge and the nonce is less than the target, the client sends the response to the
server. The client calculates the hash of the challenge and the nonce until the hash is less than the target. The client
sends the response to the server. The client receives the quote from the server and displays it on the screen.

##### DDOS attack

For the test of the server, I used the ddos attack. I used the following command:

```bash
make ddos
```

#### Docker

##### Server

The Dockerfile for the server is located in the server directory.

The Dockerfile for the server is started with the following command:

```bash
make docker-server
```

The Dockerfile for the client is located in the client directory.

The Dockerfile for the client is started with the following command:

```bash
make docker-client
```

### Testing

You can test the server and the client with the following command:

```bash
make test
```

## TODO

- [x] Add Dockerfile
- [x] Add README.md
- [ ] Refactor code
- [ ] Use lint
- [ ] Use go-releaser
- [ ] Add tests

