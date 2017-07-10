JU=${JIRA_USER}
JP=${JIRA_PASS}
GHT=${GITHUB_TOKEN}
GLT=${GITLAB_TOKEN}
RT=${ROCKET_TOKEN}

docker run -e SB_JIRA_HOST="https://issues.jboss.org" -e SB_JIRA_USER="${JU}" -e SB_JIRA_PASS="${JP}" -e SB_GITHUB_TOKEN="${GHT}" -e SB_GITLAB_TOKEN="${GLT}" -e SB_ROCKET_TOKEN="${RT}" -e SB_JIRA_BOARD="RHMAP Core Team" -e SB_JIRA_SPRINT="4.x sprint 2" --rm  --name sprintbots -it $@ server --log-level=debug
