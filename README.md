## Sprint Bot

A service that watches your sprint, gives useful information and integrates with rocket chat.

Current Commands Available:

| Command            | Response                                                                              |
|--------------------|---------------------------------------------------------------------------------------|
| sprintbot next     | Tickets that are available next in the Sprint. i.e. Tickets with prs for review       |
| sprintbot status   | Sprint data e.g. Points Completed/remaining, Sprint velocity, No. of open issues      |
| sprintbot help     |  Available commands and responses       							                     |

In the future it will:

- It will watch for open pull requests and notify the room

- It will allow you to create a distinct log entry that can be seen via the team at ```sprintbot log <today>```
Creating an entry would be ```sprintbot log start mylog``` 
...
... ```sprintbot log commit mylog ```

- It will prompt if all PRs haven't been closed

## Running Locally

### Start rocket chat

```
docker run --name db -d mongo:3.0 --smallfiles
docker run --name rocketchats -p 3001:3000 --env ROOT_URL=http://localhost --link db --link sprintbots -d rocket.chat
```
See original steps @ [dockerhub](https://hub.docker.com/_/rocket.chat/)

### Set up rocket web hook

Take the following steps to set-up the rocket web hook:

- go to localhost:3001
- create user
- click on username on top left
- Administration -> Integrations -> New Integration -> Outgoing Webhook -> Event Trigger -> Message Sent
- Set the following
  - Enabled to True
  - Select the channel where the sprintbot will be used
  - Paste to URLs: http://sprintbots:3000/chat/message?source=rocket
  - Make note of the rocket token for use below
  - Copy `integrations/rocket/script.js` to Script

### Set Up sprintbot

* Clone the repo and run:

```
cd cmd/server
export GOOS=linux; go build .
docker build -t sprintbot:latest .
```
* export the required env vars or add them to your `~/.bashrc` file

```
export JIRA_USER=your_jira_username
export JIRA_PASS=your_jira_password
export GITHUB_TOKEN=github_token
export ROCKET_TOKEN=rocket_token
```

* Run sprintbot

```
./start.sh <image-hash>
```

You will now have a local development environment. To view logs from rocketchat run:
```
 docker logs -f <container-id>
```

![Local Development Logs](images/debug.png?raw=true)