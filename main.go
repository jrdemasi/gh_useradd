package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

func parseArgs() []string {
	// Check for at least one arg, bail if none
	if len(os.Args) < 2 {
		log.Fatalln("You must provide at least one GitHub username.  Usernames can be provided either as a comma-separated list or as separate arguments.")
	} else if len(os.Args) == 2 {
		//  One argument, comma-separated list
		if strings.Contains(os.Args[1], ",") {
			// We have a comma, so assume multiple users
			usernames := strings.Split(os.Args[1], ",")
			return usernames
		} else {
			// No comma, one username provided
			usernames := os.Args[1:]
			return usernames
		}
	} else {
		// Here we assume usernames were provided as separate args, maybe with commas
		usernames := os.Args[1:]
		j := 0
		for i := 0; i < len(usernames); i++ {
			usernames[j] = strings.Replace(usernames[i], ",", "", -1)
			j++
		}
		usernames = usernames[:j]
		return usernames
	}
	return make([]string, 0)
}

func checkUsername(username string) bool {

	valid := false

	log.Printf("Checking for GitHub user %s", username)
	url := fmt.Sprintf("https://github.com/%s.keys", username)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode == 200 {
		log.Printf("Found valid user %s", username)
		valid = true
		return valid
	} else {
		log.Printf("%s is an invalid user", username)
	}

	return valid
}

func addUser(username string) {

	useradd, lookErr := exec.LookPath("useradd")
	if lookErr != nil {
		log.Fatal(lookErr)
	}

	log.Printf("Adding Unix user %s", username)
	/*args := []string{"useradd", "-m", username}
	  env := os.Environ()

	  execErr := syscall.Exec(useradd, args, env)
	  if execErr != nil {
	      log.Fatal(execErr)
	  } */
	addCmd := exec.Command(useradd, "-m", username)
	_, err := addCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func fetchKeys(username string) {

    log.Printf("Fetching keys for user %s", username)
	path := fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)
	out, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	url := fmt.Sprintf("https://github.com/%s.keys", username)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println(err)
	}
	return
}

func addSshDir(username string) {

	log.Printf("Adding .ssh directory for user %s", username)
	path := fmt.Sprintf("/home/%s/.ssh", username)
	os.Mkdir(path, 0700)
	return
}

func fixPerms(username string) {
    log.Printf("Cleaning up, fixing permissions for %s", username)
	user, _ := user.Lookup(username)
	sshdir := fmt.Sprintf("%s/.ssh", user.HomeDir)
	keyfile := fmt.Sprintf("%s/.ssh/authorized_keys", user.HomeDir)
	uid, _ := strconv.Atoi(user.Uid)
	gid, _ := strconv.Atoi(user.Gid)

	os.Chmod(keyfile, 0600)
	os.Chown(sshdir, uid, gid)
	os.Chown(keyfile, uid, gid)
	return
}

func main() {
	usernames := parseArgs()
	for _, username := range usernames {
		if checkUsername(username) {
			addUser(username)
			addSshDir(username)
			fetchKeys(username)
			fixPerms(username)
		}
	}
	return
}
