package config

import (
	env "github.com/mauhftw/genialo/helpers"
)

// TODO: Add a prefix env variable
// Define environment variables here
var (
	GithubAccessToken = env.GetEnvVar("GITHUB_CHANGELOG_TOKEN", "github_token").(string)
)
