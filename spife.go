package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Version ...
type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (v Version) String() string {
	s := make([]byte, 0, 5)
	s = strconv.AppendUint(s, v.Major, 10)
	s = append(s, '.')
	s = strconv.AppendUint(s, v.Minor, 10)
	s = append(s, '.')
	s = strconv.AppendUint(s, v.Patch, 10)

	return string(s)
}

func bumpVersion(v Version, bumpLevel string) (Version, error) {

	if bumpLevel == "patch" {
		v.Patch++
	} else if bumpLevel == "minor" {
		v.Minor++
		v.Patch = 0
	} else if bumpLevel == "major" {
		v.Major++
		v.Minor = 0
		v.Patch = 0
	} else {
		return Version{}, fmt.Errorf("Bump level '%s' unknown", bumpLevel)
	}

	return v, nil
}

func gitCommit(path, bumpLevel, newVersion, gitRemotes string, gitPush bool, gitUseFollowTags bool) {

	// git add
	fmt.Println("Adding changes to Git index")
	os.Chdir(path)
	_, err := exec.Command("git", "add", "metadata.rb").Output()
	if err != nil {
		log.Fatal(err)
	}

	if gitPush {

		// git commit
		fmt.Println("Committing changes")
		commitMessage := fmt.Sprintf("%s bump to version %s", bumpLevel, newVersion)
		_, err := exec.Command("git", "commit", "-m", commitMessage, "-n").Output()
		if err != nil {
			log.Fatal(err)
		}

		// git tag
		fmt.Printf("Adding tag '%s'\n", newVersion)
		tagMessage := fmt.Sprintf("%s", newVersion)
		_, err = exec.Command("git", "tag", "-a", newVersion, "-m", tagMessage).Output()
		if err != nil {
			log.Fatal(err)
		}

		// git push to gitRemotes
		remotes := strings.Split(gitRemotes, ",")
		for i := range remotes {
			remote := strings.TrimSpace(remotes[i])

      if gitUseFollowTags {
			  fmt.Printf("Pushing changes to '%s'\n", remote)
			  _, err = exec.Command("git", "push", remote, "master", "--follow-tags").Output()
      } else {
        fmt.Printf("Pushing changes to '%s'\n", remote)
        _, err = exec.Command("git", "push", remote, "master").Output()
        fmt.Printf("Pushing tags to '%s'\n", remote)
        _, err = exec.Command("git", "push", remote, "master", "--tags").Output()
      }

			if err != nil {
				log.Fatal(err)
			}
		}

	}
}

func parseVersion(version string) (Version, error) {
	versionArray := strings.SplitN(version, ".", 3)
	if len(versionArray) != 3 {
		return Version{}, fmt.Errorf("No major/minor/patch elements found")
	}

	var err error
	v := Version{}

	// major
	v.Major, err = strconv.ParseUint(versionArray[0], 10, 64)
	if err != nil {
		return Version{}, err
	}

	// minor
	v.Minor, err = strconv.ParseUint(versionArray[1], 10, 64)
	if err != nil {
		return Version{}, err
	}

	// patch
	v.Patch, err = strconv.ParseUint(versionArray[2], 10, 64)
	if err != nil {
		return Version{}, err
	}

	return v, nil
}

func metadata(path, bumpLevel string) string {
	path += "/metadata.rb"
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	newVersion := Version{}

	for i, line := range lines {
		if strings.Contains(line, "version") {
			lineArray := strings.Split(line, "'")
			fmt.Printf("Current version: %s\n", lineArray[1])

			version, err := parseVersion(lineArray[1])
			if err != nil {
				log.Fatalln(err)
			}

			newVersion, err = bumpVersion(version, bumpLevel)
			if err != nil {
				log.Fatalln(err)
			}

			lineArray[1] = newVersion.String()
			fmt.Printf("New version: %s\n", lineArray[1])

			lines[i] = strings.Join(lineArray, "'")
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	return newVersion.String()
}

func main() {
	path := flag.String("path", "", "Full or relative path to the cookbook directory. REQUIRED.")
	bumpLevel := flag.String("bump-level", "patch", "Version level to bump the cookbook")
	gitPush := flag.Bool("git-push", true, "Whether or not changes should be committed.")
  gitUseFollowTags := flag.Bool("git-use-follow-tags", true, "Use the directive --follow-tags.")
	gitRemotes := flag.String("git-remotes", "upstream, origin", "Comma separated list of Git remotes")
	flag.Parse()

	if *path == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	newVersion := metadata(*path, *bumpLevel)
	gitCommit(*path, *bumpLevel, newVersion, *gitRemotes, *gitPush, *gitUseFollowTags)
}
