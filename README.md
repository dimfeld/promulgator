A bot that will provide two-way integration between Slack and Jira. Currently it provides one-way integration from Jira to Slack.

## Features

* Take webhooks from Jira and post new comments to a Slack channel.

## Usage

### Bot Configuration

All configuration is done through environment variables.

* PROMULGATOR_JIRAURL -- The URL of the associated Jira instance.
* PROMULGATOR_SLACKKEY -- The token provided by Slack's bot integration.
* PROMULGATOR_SLACKUSER -- Post updates as this user. This should generally match the name of the bot.
* PROMULGATOR_SLACKDEFAULTCHANNEL -- Post updates to this channel. The bot must be a member of the channel.
* PROMULGATOR_WEBHOOKBIND -- Bind to this address/port to receive web hooks.
* PROMULGATOR_WEBHOOKKEY -- The key to use to verify incoming Jira webhooks.

### Jira Configuration

Set up a Jira webhook to `http://SERVER:PORT/jirahook?key=PROMULGATOR_WEBHOOKKEY`. The webhook should trigger on issue updates.

### Slack Configuration

Set up a Slack bot integration. The bot's username should be the value of `PROMULGATOR_SLACKUSER` and the provided API token should be the value of `PROMULGATOR_SLACKKEY`. The bot should be invited into the desired channel, and the
channel's name (with the `#` at the beginning) should be configured as `PROMULGATOR_SLACKDEFAULTCHANNEL`.

## Building

This project requires the `gb` dependency management tool to build.
If you do not have this tool, execute `go get -u github.com/constabulary/gb/...`

With the tool installed, you can pull down the dependencies using
`gb vendor restore` and then build with `gb build`.

## Upcoming Features

I'm working on this in my spare time, so can't make any promises on when these will be done.

* Real logging instead of writing to stdout
* Support HTTPS for webhooks
* Handle events from Jira for issues created, resolved, etc.
* Move to three-tier system where incoming actions are completely isolated from outgoing actions.
* Allow commenting on, resolving, closing, and assigning Jira issues from Slack. Need to decide on slash commands (`/jira blah blah`) vs. bot-type commands (`@jira blah blah`).
* Handle webhooks from Git servers to resolve issues based on pushed commits.
