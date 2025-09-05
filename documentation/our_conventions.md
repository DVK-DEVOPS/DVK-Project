# Our Conventions
Below will be a list of conventions that the team has agreed upon to adhere during development.

## General Conventions
We will follow the conventions as defined from Go documentation which can be seen [here.](https://go.dev/doc/effective_go)

Go conventions emphasize short, lowercase, meaningful, singular package names and use PascalCase for exported types and camelCase for unexported types and variables. File names are lowercase and use no underscores.


## Package Structure
`go/db` will contain database configuration.

`go/handlers` will contain the underlying controller code for routing the user.

`go/models` will contain the internal "business logic" of our system.

`go/templates` will contain the .html files that we will serve our users based on their routing within our system.

## Pull Request Template
Our repository contains a pull request template that will be inserted into future pull requests making sure that the developer has a checklist that they must adhere to before finalising their pull request which can be seen **below**. ![Image of pull request template within documentation assets](/documentation/assets/pull_request_template_img.png)
