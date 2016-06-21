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
	glog.Infof("rr tool checking repos under %q..\n", c.RepoBase)
	gitWtf := "git-wtf.rb"
	glog.V(1).Infof("about to run %q..\n", gitWtf)
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
	glog.Infof("FIXMEH: stdout=%s\n", b)
	b, err = ioutil.ReadAll(stderr)
	if err != nil {
		panic(err)
	}
	glog.Infof("FIXMEH: stderr=%s\n", b)
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
	//if err := cmd.Run(); err != nil {
	//glog.Fatalf("git-wtf.rb command failed: %v\n", err)
	//}
}
