package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"genialo/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// subobject of tagPayload
type TagObject struct {
	CommitType string `json:type`
	Sha        string `json:sha`
	Url        string `json:url`
}

// payload for github tags. Keep in mind this is an array of objects (collection)
type TagPayloadJson struct {
	Ref    string    `json:ref`
	NodeId string    `json:node_id`
	Url    string    `json:url`
	Object TagObject `json:object`
}

// github release options
type ReleaseOptions struct {
	Organization  string
	FutureRelease string
	Application   string
	Token         string
}

// TASKS
// get latest tag  X
// increment version X
// generate changelog X
// create release

// create changelogs and git releases
func GithubHandler(cmd *cobra.Command, args []string) {

	// set release options for generating the changelog
	opt := &ReleaseOptions{
		Organization: Organization,
		Application:  Application,
		Token:        config.GithubAccessToken,
	}

	log.Info("Getting latest release of ", opt.Application)
	getLastestTag(opt)
	log.Info("Latest release is: ", opt.FutureRelease)

	// increment version depending major,minor,patch
	getFutureRelease(opt, "major")
	log.Info("New release is: ", opt.FutureRelease)

	// generate changelog. Works if i use sudo, ruby problem
	log.Info("Creating changelog in /tmp/", opt.Application)
	buildChangelog(opt)
	log.Info("Changelog created succesfully in /tmp/", opt.Application)

	// commit changelog

	// create release

}

// generates software changelog
func buildChangelog(o *ReleaseOptions) error {

	// sets temporal directory for building changelogs
	dst := "/tmp/" + o.Application + "/"

	// check for temporal directory existence. REFACTOR use switch case
	if _, err := os.Stat(dst); err != nil {
		log.Info("Creating directory...")
		if err := os.Mkdir(dst, 0777); err != nil {
			fmt.Fprintln(os.Stdout, "There was an error creating the directory: ", err)
		}
	}

	// execute github_changelog_generator command
	cmd := "github_changelog_generator"
	cmdArgs := []string{
		"--future-release", o.FutureRelease,
		"--user", o.Organization,
		"--project", o.Application,
		"--token", o.Token,
		"--date-format", "%d-%m-%Y"}

	cmdRun := exec.Command(cmd, cmdArgs...)
	cmdRun.Dir = dst
	err := cmdRun.Run()

	if err != nil {
		fmt.Fprintln(os.Stdout, "There was an error running command: ", err)
		log.Fatal("The command couldn't be executed")
		return err
	}

	return nil
}

// test if env var are set
// get latest github tag
func getLastestTag(o *ReleaseOptions) {

	// we build https://api.github.com/repos/{owner}/{repo}/git/refs/tags?access_token={access_token} for requesting tags
	url := "https://api.github.com/repos/" + o.Organization + "/" +
		Application + "/git/refs/tags?access_token=" + o.Token

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("error:", err)
	}

	defer resp.Body.Close()

	// transofrm request stream into json
	p := make([]TagPayloadJson, 0)
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&p); err != nil {
		log.Fatal("Error", err)
	}

	// get github latest tag and parse it
	l := p[len(p)-1].Ref
	s := strings.Split(l, "/")
	latestTag := s[2]

	// sets future release
	o.FutureRelease = latestTag

}
