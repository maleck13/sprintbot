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
    console.log(request)
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
    var issues = []
    for (var i = 0; i < resBody.Issues.length; i++) {
      var issue = resBody.Issues[i];
      if (issue.PRs && issue.PRs.length > 0) {
        for (var k = 0; k < issue.PRs.length; k++) {
          issues.push({
            "title": "PR needs review",
            "title_link": issue.PRs[k],
            "text": "The Jira " + issue.Link + " has open PR(s) that haven't been reviewed yet:\n " + issue.PRs[k] + "",
            "color": "#764FA5"
          })
        }
      }
    }
    if (issues.length == 0) {
      for (var i = 0; i < resBody.Issues.length; i++) {
        var issue = resBody.Issues[i];
        issues.push({
          "title": "Jira " + issue.Description,
          "title_link": issue.Link,
          "text": "The Jira " + issue.Link + " needs  looking at ",
          "color": "#764FA5"
        })
      }
    }
    // Return empty will proceed with the default response process
    return {
      content: {
        "username": "sprintbot",
        "text": resBody.Message,
        "attachments": issues
      }
    }

  }
}