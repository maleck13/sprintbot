## Sprint Bot

A service that watches your sprint gives useful information and integrates with rocket chat


- It will watch for open pull requests and notify the room

- It will answer when asked what next

- It will prompt if all PRs haven't been closed

## Running Locally

clone the repo

```
cd cmd/server
export GOOS=linux; go build .
docker build -t sprintbot:latest .

```

next start rocket chat

