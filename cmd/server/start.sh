JU=${JIRA_USER}
JP=${JIRA_PASS}
GT=${GITHUB_TOKEN}
RT=${ROCKET_TOKEN}

docker run -e SB_JIRA_HOST="https://issues.jboss.org" -e SB_JIRA_USER="${JU}" -e SB_JIRA_PASS="${JP}" -e SB_GITHUB_TOKEN="${GT}" -e SB_ROCKET_TOKEN="${RT}" -e SB_JIRA_BOARD="RHMAP Core Team" -e SB_JIRA_SPRINT="4.x sprint 1" --rm  --name sprintbots -it $@ server --log-level=debug