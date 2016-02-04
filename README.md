A bot that will provide two-way integration between Slack and Jira. Currently it provides one-way integration from Jira to Slack.

## Features

* Take webhooks from Jira and post new comments to a Slack channel.

## Usage

### Bot Configuration

All configuration is done through environment variables.

* `JIRA_URL` -- The URL of the associated Jira instance.
* `JIRA_WEBHOOK_KEY` -- The key to use to verify incoming Jira webhooks.
* `SLACK_KEY` -- The token provided by Slack's bot integration.
* `SLACK_USER` -- Post updates as this user. This should generally match the name of the bot.
* `SLACK_DEFAULT_CHANNEL` -- Post updates to this channel. The bot must be a member of the channel.
* `WEBHOOK_BIND` -- Bind to this address/port to receive web hooks.


### Jira Configuration

Set up a Jira webhook to `http://SERVER:PORT/jirahook?key=WEBHOOK_KEY`. The webhook should trigger on issue updates.

### Slack Configuration

Set up a Slack bot integration. The bot's username should be the value of `SLACK_USER` and the provided API token should be the value of `SLACK_KEY`. The bot should be invited into the desired channel, and the
channel's name (with the `#` at the beginning) should be configured as `SLACK_DEFAULT_CHANNEL`.

## Building

This project uses the `glide` dependency management tool. Use `glide install` and then build normally, ensuring that `GO15VENDOREXPERIMENT=1` is in your environment if using Go 1.5.

For quick builds you don't have to do this, but the versions of dependencies may
not be the same as those that I have built against.

## Upcoming Features

I'm working on this in my spare time, so can't make any promises on when these will be done.

* Real logging instead of writing to stdout
* Support HTTPS for webhooks
* Handle events from Jira for issues created, resolved, etc.
* Move to three-tier system where incoming actions are completely isolated from outgoing actions.
* Allow commenting on, resolving, closing, and assigning Jira issues from Slack. Need to decide on slash commands (`/jira blah blah`) vs. bot-type commands (`@jira blah blah`). I may support both options.
* Handle webhooks from Git servers to resolve issues based on pushed commits.
