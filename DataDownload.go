//TODO REPLACE ALL 176 wit h164

package main

import (
	"encoding/json"
	"fmt"
	"github.com/jzelinskie/geddit"
	"log"
	"os"
)

type SubmnComments struct {
	Subm   geddit.Submission `json:"Subm"`
	Cmmnts []geddit.Comment  `json:"Cmmnts"`
}

//Struct to hold Configuration Information,
//loaded by file as to keep secret keys and passwords
//off of github
type Configinfo struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Clientid     string `json:"cid"`
	Clientsecret string `json:"csecret"`
}

var config Configinfo

func init() {
	//load secret file
	configdata, err := os.Open("C:\\Users\\Myke\\Dropbox\\School\\CSCI164\\Project\\config.secret")
	if err != nil {
		panic(err)
	}

	//close file when program is done
	defer configdata.Close()

	//load config file from json to config struct
	scanner := bufio.NewScanner(configdata)
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), &config)
		if err != nil {
			panic(err)
		}

	}

}

func main() {
	var x string

	listingoptions := geddit.ListingOptions{Time: "all"}
	o, err := geddit.NewOAuthSession(
		config.Clientid,
		config.Clientsecret,
		"school project to determine whether or not a comment would be upvoted",
		"http://redirect.url",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create new auth token for confidential clients (personal scripts/apps).
	// err = o.LoginAuth("my_user", "my_password")
	err = o.LoginAuth(config.Username, config.Password)
	if err != nil {
		log.Fatal(err)
	}

	// subs, err := o.MySubreddits()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	subs := [...]string{"space", "programming"}
	for i, _ := range subs {
		// var subs[i].Name map[string]int
		fmt.Printf("Subreddit: ", subs[i])
		x = "data/" + subs[i] + ".json"
		tmpfile, err := os.Create(x)
		if err != nil {
			log.Fatal(err)
		}
		submissions, err := o.SubredditSubmissions(subs[i], geddit.TopSubmissions, listingoptions)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("length of submissions", len(submissions))
		tmpsubms := make([]SubmnComments, len(submissions))

		for i, _ := range submissions {

			// tmpsubmis, err := json.Marshal(submissions[i])
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// fmt.Fprintf(tmpfile, string(tmpsubmis))
			fmt.Printf("Submissions: ", submissions[i])
			fmt.Printf("\n")
			fmt.Printf("\n")
			comments, err := o.Comments(submissions[i], geddit.TopSubmissions, listingoptions)
			if err != nil {
				log.Fatal(err)
			}

			tmpcmmnts := make([]geddit.Comment, len(comments))

			for i, _ := range comments {
				// comments[i].BodyHTML = ""
				tmpcmmnts = append(tmpcmmnts, *comments[i])
			}

			tmpinfo := SubmnComments{Subm: *submissions[i], Cmmnts: tmpcmmnts}

			// for i, _ := range comments {
			// 	// tmpcomment, err := json.Marshal(comments[i])
			// 	// if err != nil {
			// 	// 	log.Fatal(err)
			// 	// }
			// 	fmt.Fprintf(tmpfile, string(tmpcomment))

			// }
			tmpsubms = append(tmpsubms, tmpinfo)
		}

		tmpinfotowrite, err := json.Marshal(tmpsubms)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("JSON INFO TO FOLLOW")
		fmt.Printf(string(tmpinfotowrite))
		fmt.Println("JSON INFO TO ENDED")

		fmt.Fprintf(tmpfile, string(tmpinfotowrite))

		fmt.Printf("\n")
		defer tmpfile.Close()
	}
	// Ready to make API calls!
}
