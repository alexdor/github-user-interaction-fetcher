# Github Contributions Fetcher
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Falexdor%2Fgithub-user-interaction-fetcher.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Falexdor%2Fgithub-user-interaction-fetcher?ref=badge_shield)


A simple go app, that fetches all the open source Github repositories that a list of users has interacted with, using Github's GraphQL API. As an interaction, it considers a contribution, a commit, a new issue or a comment on an issue.

## Prerequisites

In order to run this app you need to generate a github access token from [here](https://github.com/settings/tokens) -> Personal access tokens -> Generate new token and check the `public_repo` permission.

## Starting the app

### Option 1: Docker

* Add the github access token to the docker-compose.yml
* Run `docker-compose up --build`

### Opton 2: Native

* Make sure you have go and dep install on your machine
* Set the environment variable GITHUB_TOKEN to your Github access key for example `export GITHUB_TOKEN=your_access_key`
* Start the app `go run main.go`

## Technologies Used:

* [Gin](https://gin-gonic.github.io/gin/)
* [MUI CSS](https://www.muicss.com/)


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Falexdor%2Fgithub-user-interaction-fetcher.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Falexdor%2Fgithub-user-interaction-fetcher?ref=badge_large)