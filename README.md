A bot that provides two-way integration between Slack and Jira.

The components are relatiely orthogonal to each other to provide
ease of integration with other systems in the future.

This project requires the `gb` dependency management tool to build.
If you do not have this tool, execute `go get -u github.com/constabulary/gb/...`

With the tool installed, you can pull down the dependencies using
`gb vendor restore` and then build with `gb build`.
