/* exported Script */
/* globals console, _, s, HTTP */

/** Global Helpers
 *
 * console - A normal console instance
 * _       - An underscore instance
 * s       - An underscore string instance
 * HTTP    - The Meteor HTTP object to do sync http calls
 */

class Script {
  /**
   * @params {object} request
   */
  prepare_outgoing_request({
    request
  }) {
    // request.params            {object}
    // request.method            {string}
    // request.url               {string}
    // request.auth              {string}
    // request.headers           {object}
    // request.data.token        {string}
    // request.data.channel_id   {string}
    // request.data.channel_name {string}
    // request.data.timestamp    {date}
    // request.data.user_id      {string}
    // request.data.user_name    {string}
    // request.data.text         {string}
    // request.data.trigger_word {string}

    let match;

    // Change the URL and method of the request
    
    request.headers["Content-type"] = "application/json";
    // Prevent the request and return a new message
    match = request.data.text.match(/^sprintbot help$/);
    if (match) {
      return {
        message: {
          text: [
            '**commands**',
            '```',
            '  next : will search the available issues and let you know what is best to take next',
            '  status : will return a sprint status report',
            '```'
          ].join('\n')
        }
      };
    }

    match = request.data.text.match(/^sprintbot/);
    if (match) {
      request.method = 'POST';
      return request;
    }
  }


handleIssues(issues){
var issueList = []
for (var i = 0; i < issues.length; i++) {
      var issue = issues[i];
      if (issue.PRs && issue.PRs.length > 0) {
        for (var k = 0; k < issue.PRs.length; k++) {
          issueList.push({
            "title": "PR needs review",
            "title_link": issue.PRs[k],
            "text": "The Jira " + issue.Link + " has open PR(s) that haven't been reviewed yet:\n " + issue.PRs[k] + "",
            "color": "#764FA5"
          })
        }
      }
    }
    if (issueList.length == 0) {
      for (var i = 0; i < issues.length; i++) {
        var issue = issues[i];
        issueList.push({
          "title": "Jira " + issue.Description,
          "title_link": issue.Link,
          "text": "The Jira " + issue.Link + " needs  looking at ",
          "color": "#764FA5"
        })
      }
    }
    return issueList
}

handleStatus(status){
  /*
  { pointsCompleted: 0,
     pointsRemaining: 0,
     velocity: 0,
     issuesRemaining: 0,
     daysRemaining: 0,
     estimatedDaysWorkRemaining: 0,
     issuesOfInterest: null }
  */
  return [{
    "title":"Points Completed: " + status.pointsCompleted
  },{
    "title":"Points Remaining: " + status.pointsRemaining
  },
  {
    "title":"Work Days Remaining: " + status.daysRemaining
  },
  {
    "title":"Sprint Velocity: " + status.velocity
  },{
    "title":"Estimated Days of Work Left: " + status.estimatedDaysWorkRemaining
  },{
    "title":"Number of Issues Still Open: " + status.issuesRemaining
  }]

}

  /**
   * @params {object} request, response
   */
  process_outgoing_response({
    request,
    response
  }) {
    console.log(response);
    // request              {object} - the object returned by prepare_outgoing_request

    // response.error       {object}
    // response.status_code {integer}
    // response.content     {object}
    // response.content_raw {string/object}
    // response.headers     {object}

    // Return false will abort the response
    // return false;
    var resBody = JSON.parse(response.content_raw);
    console.log("resBody", resBody)
    var content = []
    if(resBody.CMD && resBody.CMD === "next"){
      content = this.handleIssues(resBody.Data);
    }else if (resBody.CMD && resBody.CMD === "status"){
       content = this.handleStatus(resBody.Data);
    }
    
    // Return empty will proceed with the default response process
    return {
      content: {
        "username": "sprintbot",
        "text": resBody.Message,
        "attachments": content
      }
    }

  }
}