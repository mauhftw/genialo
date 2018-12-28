package cmd

import (
	"bytes"
	"encoding/json"
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

// github release payload
type ReleasePayload struct {
	TagName         string `json:tag_name`
	TargetCommitish string `json:target_commitish`
	Name            string `json:name`
	Body            string `json:body`
	Draft           bool   `json:draft`
	PreRelease      bool   `json:prerelease`
}

// TASKS
// get latest tag  X
// increment version X
// generate changelog X
// commit changelog
// create release X

// create changelogs and git releases
func GithubHandler(cmd *cobra.Command, args []string) {

	// set release options for generating the changelog
	opt := &ReleaseOptions{
		Organization: Organization,
		Application:  Application,
		Token:        config.GithubAccessToken,
	}

	// get latest release created
	log.Infof("Getting latest release of %v", opt.Application)
	getLastestTag(opt)
	log.Infof("Latest release is: %v", opt.FutureRelease)

	// increment version depending major,minor,patch
	getFutureRelease(opt, "major")
	log.Infof("New release is: %v", opt.FutureRelease)

	// generate changelog. Works if i use sudo, ruby problem
	log.Infof("Creating changelog in /tmp/ %v", opt.Application)
	buildChangelog(opt)
	log.Infof("Changelog created succesfully in /tmp/ %v", opt.Application)

	// commit changelog

	// create new release version
	log.Infof("Creating new release: %v", opt.FutureRelease)
	createRelease(opt)
	log.Infof("Release: %v has been created successful %v", opt.FutureRelease)

}

// generates software changelog
func buildChangelog(o *ReleaseOptions) error {

	// sets temporal directory for building changelogs
	dst := "/tmp/" + o.Application + "/"

	// check for temporal directory existence. REFACTOR use switch case
	if _, err := os.Stat(dst); err != nil {
		log.Infof("Directory %v doesn't exist. Creating directory...", o.Application)
		if err := os.Mkdir(dst, 0777); err != nil {
			log.Fatalf("There was an error creating the directory: %v ", err)
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
		log.Fatalf("There was an error running the %v command: %v ", cmd, err)
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
		log.Fatalf("There was a problem trying to reach %v. %v:", url, err)
	}

	defer resp.Body.Close()

	// decode request stream into json
	p := make([]TagPayloadJson, 0)
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&p); err != nil {
		log.Fatalf("Can't decode json. %v", err)
	}

	// get github latest tag and parse it
	l := p[len(p)-1].Ref
	s := strings.Split(l, "/")
	latestTag := s[2]

	// sets future release
	o.FutureRelease = latestTag

}

// create a new software release
func createRelease(o *ReleaseOptions) {

	// we build https://api.github.com/repos/{owner}/{repo}/releases?access_token={access_token} for creating a new release
	url := "https://api.github.com/repos/" + o.Organization +
		"/" + o.Application +
		"/" + "/releases?access_token=" +
		o.Token

	contentType := "application/json"
	body := &ReleasePayload{
		TagName:         "release" + o.FutureRelease,
		TargetCommitish: "master",
		Name:            "release" + o.FutureRelease,
		Body:            "example",
		Draft:           false,
		PreRelease:      false,
	}

	// encode json to stream
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(body); err != nil {
		log.Fatalf("Can't encode json. %s", err)
	}

	resp, err := http.Post(url, contentType, &buf)
	if err != nil {
		log.Fatalf("There was a problem when trying to create a release. %v", err)
	}

	defer resp.Body.Close()

}
