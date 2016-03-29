A bot that will provide two-way integration between Slack and Jira. Currently it provides one-way integration from Jira to Slack.

## Features

* Take webhooks from Jira and post new comments to a Slack channel.

## Usage

### Bot Configuration

All configuration is done through environment variables. The Config object in `model/model.go` has the full set of options, with documentation. These are the essentials:

* `JIRA_URL` -- The URL of the associated Jira instance, including username and password in HTTP Basic Auth format (i.e. `https://username:password@jiraserver`).
* `JIRA_WEBHOOK_KEY` -- The key to use to verify incoming Jira webhooks.
* `SLACK_KEY` -- The token provided by Slack's bot integration.
* `SLACK_USER` -- Post updates as this user. This should generally match the name of the bot.
* `SLACK_DEFAULT_CHANNEL` -- Post updates to this channel. The bot must be a member of the channel.
* `WEBHOOK_BIND` -- Bind to this address/port to receive web hooks.


### Jira Configuration

Set up a Jira webhook to `http://SERVER:PORT/jirahook?key=JIRA_WEBHOOK_KEY`. The webhook should trigger on issue updates.

### Slack Configuration

Set up a Slack bot integration. The bot's username should be the value of `SLACK_USER` and the provided API token should be the value of `SLACK_KEY`. The bot should be invited into the desired channel, and the
channel's name (with the `#` at the beginning) should be configured as `SLACK_DEFAULT_CHANNEL`.

## Building

This project uses the `glide` dependency management tool. Use `glide install` and then build normally, ensuring that `GO15VENDOREXPERIMENT=1` is in your environment if using Go 1.5.

For quick builds you don't have to do this, but the versions of dependencies may
not be the same as those that I have built against.

## Upcoming Features

I'm working on this in my spare time, so can't make any promises on when these will be done. This list is in rough order of priority.

* Allow commenting on and assigning Jira issues from Slack. Slash commands first, RTM with @jira support later.
* Allow resolving and closing issues from Slack. Jira's workflow system makes this complex and potentially different for every installation.
* Attempted autocorrelation between Jira and Slack users
* Handle events from Jira for issues created, resolved, etc. The Slack integration does this right now so it's not too important for me.
* Option to post an blurb with issue details whenever an issue is referenced.
* Handle webhooks from Git servers to resolve issues based on pushed commits.
* Support HTTPS for webhooks
* Move to three-tier system for output where incoming actions are completely isolated from outgoing actions.
