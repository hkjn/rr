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
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	RepoBase string
}

// repoHasIssues returns false if there were no issues in given repo, or true otherwise.
func repoHasIssues(m string) bool {
	gitWtf := "git-wtf.rb"
	os.Chdir(m)
	glog.V(1).Infof("about to run %q for %q..\n", gitWtf, m)
	cmd := exec.Command(gitWtf)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(stdout)
	if err != nil {
		panic(err)
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
		return true
	}
	return false
}

func main() {
	flag.Parse()
	var c config
	if err := envconfig.Process("rr", &c); err != nil {
		panic(err)
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
	for _, m := range matches {
		glog.Infof("Checking %q..\n", m)
		if repoHasIssues(m) {
			failed = true
		}
	}
	if failed {
		glog.Errorln("Some repos had issues.")
		os.Exit(1)
	}
	glog.Infoln("No issues, all done.")
	os.Exit(0)
}
