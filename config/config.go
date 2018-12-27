package config

import (
	env "genialo/helpers"
)

// TODO: Add a prefix env variable
// Define environment variables here
var (
	GithubAccessToken = env.GetEnvVar("GUITARISTS_RELEASE_VERSION", "XXXXXXXXXXXXXXXXXXXXX").(string)
)
