[![Build Status](https://travis-ci.org/nmrshll/winter-is-coming-backend.svg?branch=master)](https://travis-ci.org/nmrshll/winter-is-coming-backend)

# winter-is-coming-backend
Submission for the [backend challenge](https://github.com/mysteriumnetwork/winter-is-coming/blob/master/quests/Talk_to_Zombies.md#communication-channel-specification) for mysterium

## Quick start

Clone with

```sh
git clone git@github.com:nmrshll/winter-is-coming-backend
```

Launch with:

```sh
make dev
```

Then from a second terminal:

```sh
nc localhost 7777
```

And follow the instructions on the screen

====

### Automated tests

Automated tests get run on each push in travis. To run the test suite on your machine run:

```sh
make test
```

====

## Usage

### Requirements

- Go 1.9+

### Installation

Clone with

```sh
git clone git@github.com:nmrshll/winter-is-coming-backend
```

or install using go with

```sh
go get github.com/nmrshll/winter-is-coming-backend
```

### Launch

In one terminal run the server with:

```sh
make dev
```

In another terminal run a client with `nc` (netcat):

```sh
nc localhost 7777
```

### Playing the game

You can send two types of commands to the server (from the terminal that is connected to your server with `nc`):

1. In the first phase, you'll be asked to enter your player name. Expected format is any string (spaces and newlines will be trimmed on each side).
2. In the second phase, the zombie starts walking, you'll need to shoot (as many times as you want) by sending the coordinates where you want to shoot to attempt to kill the zombie. Expected format is a string containing two integers separated by a space (e.g. `1 5`,`0 3`,or `4 10`). 

If you enter the current coordinates of the zombie fast enough, you'll hit it and win the game. If the zombie reaches the wall before you manage to hit it, you lose the game.
