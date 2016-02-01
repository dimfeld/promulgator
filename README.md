A bot that will provide two-way integration between Slack and Jira. Currently it provides one-way integration from Jira to Slack. 

This project requires the `gb` dependency management tool to build.
If you do not have this tool, execute `go get -u github.com/constabulary/gb/...`

With the tool installed, you can pull down the dependencies using
`gb vendor restore` and then build with `gb build`.

## What's Working:

* Take webhooks from Jira and post new comments to a Slack channel.

## Upcoming Features

I'm working on this in my spare time, so can't make any promises on when these will be done.

* Documentation on how to use it. For now, see `src/model/model.go`. All config is through the environment.
* Tests
* Support HTTPS for webhooks
* Handle events from Jira for issues created, resolved, etc.
* Move to three-tier system where incoming actions are completely isolated from outgoing actions.
* Allow commenting on, resolving, closing, and assigning Jira issues from Slack. Need to decide on slash commands (`/jira blah blah`) vs. bot-type commands (`@jira blah blah`).
* Handle webhooks from Git servers to resolve issues based on pushed commits.
