package cmd

import (
	// System
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	// 3rd party
	"github.com/mauhftw/genialo/config"
	"github.com/mauhftw/genialo/helpers"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Subobject of tagPayload
type TagObject struct {
	CommitType string `json:"type"`
	Sha        string `json:"sha"`
	Url        string `json:"url"`
}

// Payload for github tags. Keep in mind this is an array of objects (collection)
type TagPayloadJson struct {
	Ref    string    `json:"ref"`
	NodeId string    `json:"node_id"`
	Url    string    `json:"url"`
	Object TagObject `json:"object"`
}

// Github release options
type ReleaseOptions struct {
	Organization  string
	FutureRelease string
	Application   string
	Token         string
}

// Github release payload
type ReleasePayload struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	Draft           bool   `json:"draft"`
	PreRelease      bool   `json:"prerelease"`
}

// Create changelogs and git releases
func GithubReleaseHandler(cmd *cobra.Command, args []string) {

	// Set release options for generating the changelog
	opt := &ReleaseOptions{
		Organization: Organization,
		Application:  Application,
		Token:        config.GithubAccessToken,
	}

	// Defines type of release X.X.X [major, minor or patch]
	releaseType := cmd.Use

	// Check repository and release status
	log.Info("Starting...")
	CheckRepoStatus(opt)

	// Get latest release created
	log.Info("Calling to getLatestTag handler...")
	GetLatestTag(opt)

	// Increment version depending on semver type: major,minor,patch
	log.Info("Calling to GetFutureRelease handler...")
	GetFutureRelease(opt, releaseType)

	// Generate changelog
	log.Info("Calling to BuildChangelog handler...")
	BuildChangelog(opt)

	// Create new release version
	log.Info("Calling to CreateRelease handler...")
	CreateRelease(opt)
	log.Info("Process finished sucessfully!")

}

// Check Repo and releases status
func CheckRepoStatus(o *ReleaseOptions) {

	// Check if repo exist
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v?access_token=%v", o.Organization, o.Application, o.Token)
	req, err := http.Get(url)
	if err != nil {
		log.Errorf("%v", err)
	}

	// Check if request has been success [2XX]
	log.Infof("Checking if repo %v exists", o.Application)
	if !strings.HasPrefix(req.Header.Get("Status"), "2") {
		r, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Errorf("\t %v", err)
		}

		log.Fatal(string(r))
	}

	defer req.Body.Close()

	// Check if there's at least a minor version
	url = fmt.Sprintf("https://api.github.com/repos/%v/%v/releases/latest?access_token=%v", o.Organization, o.Application, o.Token)
	req, err = http.Get(url)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Check if releases exists [2XX]
	log.Infof("Checking releases for %v ", o.Application)
	if !strings.HasPrefix(req.Header.Get("Status"), "2") {
		r, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Errorf("\t %v", err)

		}

		log.Warnf("\t %v", string(r))
		log.Warn("No releases had been detected!")

		o.FutureRelease = "0.1.0"

		// TODO: make a function
		// Attempt to create default initial release 0.1.0
		reader := bufio.NewReader(os.Stdin)
		log.Warnf("We are going to attempt to create a default initial release tagged as %v", o.FutureRelease)
		log.Warn("Continue? [y/n]")
		opt, _, _ := reader.ReadRune()

		switch opt {
		case 'y':
			// Generate changelog
			log.Info("Calling to BuildChangelog handler...")
			BuildChangelog(o)

			// Create new release version
			log.Info("Calling to CreateRelease handler...")
			CreateRelease(o)
			log.Info("Process finished sucessfully!")
			os.Exit(0)

		case 'n':
			log.Error("Halting execution by user selection")
			os.Exit(0)
		default:
			log.Warn("Incorrect option!")
		}
	}
}

// Generates changelog.md
func BuildChangelog(o *ReleaseOptions) {

	// Define path to checkout applications and changelog
	var dst string
	var changeLogDst string

	// Sets temporal directory for checking out application
	dst = "/tmp/applications"
	changeLogDst = dst + "/" + o.Application

	log.Infof("\t Attempting to create changelog...")

	// Checks if temporal applications directory exists
	_, err := os.Stat(dst)
	if err != nil {
		log.Warnf("\t Application directory %v doesn't exist. Creating directory...", dst)
		err := os.Mkdir(dst, 0777)
		if err != nil {
			log.Fatalf("\t %v", err)
		}
		log.Infof("\t Application directory %v created successfully!", dst)
	}

	// Checks if repo directory exists
	_, err = os.Stat(changeLogDst)
	if err != nil {
		log.Warnf("\t Repository directory %v doesn't exist. Creating directory...", o.Application)

		// Perform git clone command
		repo := fmt.Sprintf("git@github.com:%s/%s.git", o.Organization, o.Application)
		c := &helpers.BashCmd{
			Cmd:      "git",
			Args:     []string{"clone", repo, "--single-branch"},
			ExecPath: dst,
		}

		helpers.ExecBashCmd(c)
		log.Infof("\t Application directory %v created successfully!", dst)

	} else {

		// Perform git checkout and git pull commands
		c := &helpers.BashCmd{
			Cmd:      "/bin/sh",
			Args:     []string{"-c", "git checkout . && git clean -fd && git pull origin master"},
			ExecPath: changeLogDst,
		}

		helpers.ExecBashCmd(c)

	}

	// Perform github_changelog_generator command
	c := &helpers.BashCmd{
		Cmd: "github_changelog_generator",
		Args: []string{
			"--future-release", o.FutureRelease,
			"--user", o.Organization,
			"--project", o.Application,
			"--token", o.Token,
			"--date-format", "%d-%m-%Y"},
		ExecPath: changeLogDst,
	}

	helpers.ExecBashCmd(c)

	// Perform git add, commit and push commands
	push := fmt.Sprintf("git add -A && git commit -m 'release %v' && git push origin master", o.FutureRelease)
	c = &helpers.BashCmd{
		Cmd:      "/bin/sh",
		Args:     []string{"-c", push},
		ExecPath: changeLogDst,
	}

	helpers.ExecBashCmd(c)
	log.Info("\t Changelog created succesfully!")
}

// Set latest github release/tag
func GetLatestTag(o *ReleaseOptions) {

	// TODO: check that version is x.x.x

	// We build https://api.github.com/repos/{owner}/{repo}/git/refs/tags?access_token={access_token} for requesting tags
	log.Infof("\t Attempting to get latest release of %v...", o.Application)
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/git/refs/tags?access_token=%v", o.Organization, o.Application, o.Token)

	// Perform get request
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("\t %v", err)
	}

	// Check if request has been success [2XX]
	if !strings.HasPrefix(resp.Header.Get("Status"), "2") {
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("\t %v", err)
		}

		log.Fatal(string(r))
	}

	defer resp.Body.Close()

	// Decode request stream into json
	p := make([]TagPayloadJson, 0)
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&p); err != nil {
		log.Fatalf("Can't decode json. %v", err)
	}

	// Get github's latest tag and parse it
	l := p[len(p)-1].Ref
	s := strings.Split(l, "/")
	latestTag := s[len(s)-1]

	// Sets future release
	o.FutureRelease = latestTag
	log.Infof("\t Latest release is: %v", o.FutureRelease)

}

// Create a new software release
func CreateRelease(o *ReleaseOptions) {

	// Build github release URL
	url := "https://api.github.com/repos/" + o.Organization +
		"/" + o.Application +
		"/releases?access_token=" +
		o.Token

	// Set request's payload
	log.Info("\t Attempting to create new release...")
	now := time.Now().Format(time.RFC1123) // NOTE: Monday, 02-Jan-06 15:04:05 MST format
	body := &ReleasePayload{
		TagName:         o.FutureRelease,
		TargetCommitish: "master",
		Name:            "Release " + o.FutureRelease + " on " + now,
		Body:            "full changelog at: " + "https://github.com/" + o.Organization + "/" + o.Application + "/blob/master" + "/CHANGELOG.md",
		Draft:           false,
		PreRelease:      false,
	}

	// Prepare payload's format and set request
	bJson, err := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bJson))
	if err != nil {
		log.Errorf("\t %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Open a new http client and perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Check if request has been success [2XX]
	if !strings.HasPrefix(resp.Header.Get("Status"), "2") {
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("\t Error in reading request: %v", err)
		}

		log.Errorf("\t Response not OK: %v", resp.Header.Get("Status"))
		log.Fatalf("\t Error in response %v", string(r))
	}

	defer resp.Body.Close()
	log.Infof("\t Release: %v created successfully!", o.FutureRelease)
}
