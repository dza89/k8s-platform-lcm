// Package xray is used to access Xray to find vulnerabilities for images
package xray

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/arminc/k8s-platform-lcm/pkg/vulnerabilities"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/target/go-arty/xray"
)

// Config contains all the information to talk to Xray
type Config struct {
	Username string   `koanf:"username"`
	Password string   `koanf:"password"`
	URL      string   `koanf:"url"`
	Prefixes []Prefix `koanf:"prefixes"`
}

// Prefix information about the index used by Xray
type Prefix struct {
	Prefix string   `koanf:"prefix"`
	Images []string `koanf:"images"`
}

// Scanner is an interface that wraps calls to Xray
type Scanner interface {
	GetVulnerabilities(name, version string, prefixes []Prefix) ([]vulnerabilities.Vulnerability, error)
	GetXrayResults(request xray.SummaryArtifactRequest) (xray.SummaryArtifact, error)
}

type xrayClient struct {
	client *xray.Client
}

// NewXray constructs access to Xray
// It returns an implementation of the Xray client represented as the Scanner interface
func NewXray(config Config) (Scanner, error) {
	client, err := xray.NewClient(config.URL, nil)
	if err != nil {
		return &xrayClient{}, err
	}
	client.Authentication.SetBasicAuth(config.Username, config.Password)
	return &xrayClient{
		client: client,
	}, nil
}

// GetVulnerabilities returns Xray results as generic Vulnerabilities instead of in the Xray format
// It returns empty Image list on error
func (x *xrayClient) GetVulnerabilities(name, version string, prefixes []Prefix) ([]vulnerabilities.Vulnerability, error) {
	log.WithField("image", name).WithField("version", version).Debug("Fetching vulnerabilities")
	path := fmt.Sprintf("%s/%s/%s", findPrefix(name, prefixes), name, version)
	xrayVulnerabilities, err := x.GetXrayResults(xray.SummaryArtifactRequest{
		Paths: &[]string{path},
	})
	if err != nil {
		return []vulnerabilities.Vulnerability{}, err
	}

	allVulnerabilities := []vulnerabilities.Vulnerability{}
	for _, issue := range xrayVulnerabilities.GetIssues() {
		for _, cve := range issue.GetCves() {
			cve := cve.GetCve()
			if cve == "" {
				cve = hashString(issue.GetDescription())
			}
			vulnerability := vulnerabilities.Vulnerability{
				Identifier:  cve,
				Description: issue.GetDescription(),
				Severity:    issue.GetSeverity(),
			}
			allVulnerabilities = append(allVulnerabilities, vulnerability)
		}
	}
	return allVulnerabilities, nil
}

// GetXrayResults returns results as they come from Xray
func (x *xrayClient) GetXrayResults(request xray.SummaryArtifactRequest) (xray.SummaryArtifact, error) {
	sum, response, err := x.client.Summary.Artifact(&request)
	if err != nil {
		return xray.SummaryArtifact{}, err
	}
	if response.StatusCode != 200 {
		return xray.SummaryArtifact{}, errors.Wrapf(err, "Error fetching xray vulnerabilities, http-status: %s", response.Status)
	}
	if len(sum.GetErrors()) >= 1 {
		return xray.SummaryArtifact{}, errors.Wrapf(err, "Got an error from xray for [%v] error [%s]", request, *sum.GetErrors()[0].Error)
	}

	// Currently we only fetch one image a a time, therefore we can just grab the first artifact
	if len(sum.GetArtifacts()) > 0 {
		return sum.GetArtifacts()[0], nil
	}
	return xray.SummaryArtifact{}, nil
}

// findPrefix returns the prefix used by Xray for the image
func findPrefix(imageName string, prefixes []Prefix) string {
	if len(prefixes) == 1 {
		return prefixes[0].Prefix
	}
	for _, prefix := range prefixes {
		for _, image := range prefix.Images {
			match, err := regexp.MatchString(image, imageName)
			if err != nil {
				log.WithError(err).Warn("Image regexp not valid")
			}
			if match {
				return prefix.Prefix
			}
		}
	}

	return ""
}

func hashString(text string) string {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(text))
	if err != nil {
		log.WithError(err).Warn("Could not hash")
		return "HASH_ERROR"
	}
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
