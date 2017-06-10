## Sprint Bot

A service that watches your sprint gives useful information and integrates with rocket chat


- It will watch for open pull requests and notify the room

- It will answer when asked ```sprintbot next```

-It will allow you to create a distinct log entry that can be seen via the team at ```sprintbot log <today>```
Creating an entry would be ```sprintbot log start mylog``` 
...
... ```sprintbot log commit mylog ```

- It will prompt if all PRs haven't been closed

## Running Locally

clone the repo

```
cd cmd/server
export GOOS=linux; go build .
docker build -t sprintbot:latest .

```

next start rocket chat

... TODO go through setup of webhook and linking the bot


## Setup the server

