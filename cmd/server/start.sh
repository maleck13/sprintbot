JU=${JIRA_USER}
JP=${JIRA_PASS}
GT=${GITHUB_TOKEN}
RT=${ROCKET_TOKEN}

docker run -e SB_ROCKET_TOKEN="${RT}" -e SB_JIRA_BOARD="RHMAP Core Team" -e SB_JIRA_SPRINT="4.x sprint 1" --rm  --name sprintbots -it $@ server  -jira-host=https://issues.jboss.org -jira-user="${JU}" -jira-pass="${JP}" -github-token="${GT}"