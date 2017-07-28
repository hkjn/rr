// The rr tool reports status of repos.
//
// Checks that all files under REPOBASE are:
// 1. Directories
// 2. With git repos inside
// 3. With a clean working tree
// TODO(hkjn): Reimplement git-wtf.rb calls in lib/ bash scripts, to
// avoid wrapping it in shell call here (and having dependency on
// ruby).
//
// TODO(hkjn): Also look for LICENSE, README.md?
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
)

type (
	config struct {
		RepoBase string
	}
	result struct {
		Repo          string `json:"repo"`
		Error, Output string `json:"error"`
		ExitStatus    int    `json:"exit_status"`
	}
	results []result
)

// repoHasIssues returns false if there were no issues in given repo, or true otherwise.
func repoHasIssues(m string) (*result, error) {
	gitWtf := "git-wtf.rb"
	os.Chdir(m)
	glog.V(1).Infof("about to run %q for %q..\n", gitWtf, m)
	cmd := exec.Command(gitWtf)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	out := string(b)
	os.Chdir("..")
	b, err = ioutil.ReadAll(stderr)
	if err != nil {
		panic(err)
	}
	eout := string(b)
	if err := cmd.Wait(); err != nil {
		glog.Infof("--- start of %q output ---\n", m)
		glog.Infof(out)
		glog.Infof("--- end of %q output ---\n", m)
		glog.Errorf("--- start of %q errors ---\n", m)
		glog.Errorf(eout)
		glog.Errorf("--- end of %q errors ---\n", m)
		return &result{Repo: m, Error: eout, Output: out, ExitStatus: 1}, nil
	}
	return &result{Repo: m, Error: eout, Output: out, ExitStatus: 0}, nil
}

func main() {
	flag.Parse()
	var c config
	if err := envconfig.Process("rr", &c); err != nil {
		log.Fatalf("failed to process config: %v\n", err)
	}
	if c.RepoBase == "" {
		c.RepoBase = filepath.Join(
			os.Getenv("GOPATH"),
			"src/hkjn.me",
		)
	}
	glog.Infof("Checking repos under %q..\n", c.RepoBase)
	matches, err := filepath.Glob(filepath.Join(c.RepoBase, "*"))
	if err != nil {
		panic(err)
	}
	failed := false
	r := results{}
	for _, m := range matches {
		glog.Infof("Checking %q..\n", m)
		result, err := repoHasIssues(m)
		if err != nil {
			log.Fatalf("failed to check repos: %v\n", err)
		}
		r = append(r, *result)
		if result.ExitStatus != 0 {
			failed = true
		}
	}
	if err := json.NewEncoder(os.Stdout).Encode(r); err != nil {
		log.Fatalf("failed to encode json: %v\n", err)
	}
	if failed {
		os.Exit(1)
	}
	os.Exit(0)
}
