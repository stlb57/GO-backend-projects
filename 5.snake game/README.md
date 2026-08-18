# Snake in Go

A small Snake game that runs in the terminal, written in Go.

I built this mainly to practice Go by making something playable instead of another CLI CRUD project.

## Features

* WASD controls
* Continuous movement
* Snake body and growth
* Random food
* Wall collision
* Self collision
* Score
* Game over

## Run

Make sure Go is installed, then:

```bash
go run main.go
```

## Controls

```text
W - Up
A - Left
S - Down
D - Right
Q - Quit
```

## Project structure

```text
main.go

internal/
  game/
  terminal/
```

`game` contains the game logic and state.

`terminal` handles keyboard input and rendering.

## Why I made this

The goal wasn't to make a huge or production-ready game. I wanted a small project that would force me to work with things like structs, slices, pointers, game loops, terminal input, goroutines and channels.

Still working on it.
