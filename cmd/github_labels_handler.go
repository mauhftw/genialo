package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"genialo/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type LabelsOptions struct {
	Organization string
	LabelFile    string
	Application  string
	Token        string
}

// Define label stuct
type Label struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

// Create labels from file
func GithubLabelCreatorHandler(cmd *cobra.Command, args []string) {

	opt := &LabelsOptions{
		Organization: Organization,
		LabelFile:    LabelFile,
		Application:  Application,
		Token:        config.GithubAccessToken,
	}

	// Open label file
	log.Info("Starting...")
	log.Infof("Trying to read label file at: %v", opt.LabelFile)
	f, err := os.Open(opt.LabelFile)
	if err != nil {
		log.Fatalf("\t Error opening label file: %v", err)
	}
	defer f.Close()

	// Read label file
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("\t Error reading label file: %v", err)
	}

	// Prepare Unmarshal (to struct type)
	var labels []Label
	err = json.Unmarshal([]byte(byteValue), &labels)
	if err != nil {
		log.Info(err)
	}

	// Attempt to create labels
	log.Info("\t Calling to CreateLabels handler...")
	for _, l := range labels {
		CreateLabels(opt, l)
	}

	log.Infof("Labels created successfully!")

}

// Destroy all labels
func GithubLabelDestroyerHandler(cmd *cobra.Command, args []string) {

	opt := &LabelsOptions{
		Organization: Organization,
		LabelFile:    LabelFile,
		Application:  Application,
		Token:        config.GithubAccessToken,
	}

	// Get list of github labels
	log.Info("Starting...")
	log.Info("Trying to get labels...")
	log.Info("Calling to ListLabels handler...")
	r := ListLabels(opt)
	// TODO: Refactor. This compares to empty string, due to the []bytes to string convertion :/
	if r == "[]" {
		log.Fatalf("\t There are no labels created in that repo")
	}

	log.Info("Labels found!...")

	// Unmarshal json to interface
	var labels []map[string]interface{}
	err := json.Unmarshal([]byte(r), &labels)
	if err != nil {
		log.Fatalf("\t Error Unmarshalling %v", err)
	}

	// Parse label
	for _, v := range labels {
		// Delete labels
		DeleteLabel(v["name"].(string), opt)
	}

	log.Info("Labels destroyed successfully!")

}

// Attempt to DELETE github label
func DeleteLabel(l string, o *LabelsOptions) {

	// Prepare request for creating labels
	var byteValue []byte
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/labels/%v?access_token=%v", o.Organization, o.Application, l, o.Token)
	req, err := http.NewRequest("DELETE", url, bytes.NewReader(byteValue))
	if err != nil {
		log.Fatalf("\t Error when attempting to perform request to %v: %v", url, err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Attempt request to github
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("\t Error opening a connection to %v: %v", url, err)
	}

	// Check if request has been success [2XX]
	if !strings.HasPrefix(resp.Header.Get("Status"), "2") {
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("\t Error in reading request %v: %v", url, err)
		}

		log.Errorf("\t Response not OK: %v", resp.Header.Get("Status"))
		log.Fatalf("\t Error in response %v", string(r))
	}

	defer resp.Body.Close()
}

// Attempt to POST github labels
func CreateLabels(o *LabelsOptions, l Label) {

	// Prepare Marshal json (to []bytes/stream type)
	byteValue, err := json.Marshal(l)
	if err != nil {
		log.Fatalf("\t There was a problem when trying to Marshal json: %v", err)
	}

	// Prepare request for creating labels
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/labels?access_token=%v", o.Organization, o.Application, o.Token)
	req, err := http.NewRequest("POST", url, bytes.NewReader(byteValue))
	if err != nil {
		log.Fatalf("\t Error when attempting to perform request to %v: %v", url, err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Attempt request to github
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("\t Error opening a connection to %v: %v", url, err)
	}

	// Check if request has been success [2XX]
	if !strings.HasPrefix(resp.Header.Get("Status"), "2") {
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("\t Error in reading request %v: %v", url, err)
		}

		log.Errorf("\t Response not OK: %v", resp.Header.Get("Status"))
		log.Fatalf("\t Error in response %v", string(r))
	}

	defer resp.Body.Close()

}

func ListLabels(o *LabelsOptions) string {

	// Prepare request for creating labels
	var byteValue []byte
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/labels?access_token=%v", o.Organization, o.Application, o.Token)
	req, err := http.NewRequest("GET", url, bytes.NewReader(byteValue))
	if err != nil {
		log.Fatalf("\t Error when attempting to perform request to %v: %v", url, err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Attempt request to github
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("\t Error opening a connection to %v: %v", url, err)
	}

	// Check if request has been success [2XX]
	if !strings.HasPrefix(resp.Header.Get("Status"), "2") {
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("\t Error in reading request %v: %v", url, err)
		}

		log.Errorf("\t Response not OK: %v", resp.Header.Get("Status"))
		log.Fatalf("\t Error in response %v", string(r))

	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("\t Error in reading request %v: %v", url, err)
	}

	defer resp.Body.Close()
	return string(r)

}
