package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jzelinskie/geddit"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	// "sync"
)

//Struct to hold a single submission, and all of the comments within it
type SubmnComments struct {
	Subm   geddit.Submission `json:"subm"`
	Cmmnts []geddit.Comment  `json:"cmmnts"`
}

//Struct to hold User information
type UserData struct {
	User      string
	Subs      []string
	Words     map[string]SubsWordData
	UserScore float64
}

//Struct to hold data regarding words
type SubsWordData struct {
	Word     string
	Numoccur float64
	Avgscore float64
	Heur     float64
}

type SubsCommentData struct {
	SubName      string
	Numcomments  float64
	AvgUpVote    float64
	NumUpVotes   float64
	AvgDownVotes float64
	NumDownVotes float64
	AvgPostScore float64
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

//global vars
//empty structs referenced
var config Configinfo
var reg regexp.Regexp
var o geddit.OAuthSession

//maxnumsubs, needs to be twice the size of the number of subreddits the user can have
const MAXNUMSUBS int = 100

// //track number of upvotes and downvotes made by the bot
var SubCommentMap map[string]SubsCommentData

// var wg sync.WaitGroup

//maps of data, key is the word, value is subwordadat strat
var usagecmmts map[string](map[string]SubsWordData)
var usagesubmstitle map[string](map[string]SubsWordData)
var usagesubmsbody map[string](map[string]SubsWordData)

//maps of userdata and stopwords
var usermap map[string]UserData
var stopwords map[string]struct{}

//init automatically runs before main
//sets up global data, loads stopwords and config
func init() {
	//allocate memory to maps
	usermap = make(map[string]UserData)
	stopwords = make(map[string]struct{})
	SubCommentMap = make(map[string]SubsCommentData)
	usagecmmts = make(map[string](map[string]SubsWordData))
	usagesubmstitle = make(map[string](map[string]SubsWordData))
	usagesubmsbody = make(map[string](map[string]SubsWordData))

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

	//load stopwords
	stopdata, err := os.Open("C:\\Users\\Myke\\Dropbox\\School\\CSCI164\\Project\\stopwords")
	if err != nil {
		panic(err)
	}

	//place stop words into a map
	scanner = bufio.NewScanner(stopdata)
	for scanner.Scan() {
		var z struct{}
		stopwords[scanner.Text()] = z
	}

	//close file when program is done
	defer stopdata.Close()

}

func main() {

	//prepare filenames, load data
	dirname := "C:\\Users\\Myke\\Dropbox\\School\\CSCI164\\Project\\data/"
	d, err := os.Open(dirname)
	if err != nil {
		panic(err)
	}

	defer d.Close()

	//get slice of pointers to file
	files, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}

	//range through files
	for _, files := range files {
		//remove .json from filename
		tmpfilename := strings.ToLower(files.Name())
		fname := strings.TrimSuffix(tmpfilename, filepath.Ext(tmpfilename))

		//if file is valid
		if files.Mode().IsRegular() {

			//allocate memory to nested map for specific subreddit
			usagecmmts[fname] = make(map[string]SubsWordData)
			usagesubmstitle[fname] = make(map[string]SubsWordData)
			usagesubmsbody[fname] = make(map[string]SubsWordData)

			//debug info to show that data has been loaded correctly
			fmt.Println(files.Name(), files.Size(), "bytes")

			//Get File location and read file
			filename := dirname + files.Name()
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				panic(err)
			}

			//empty struct to place data into
			jsondata := []SubmnComments{}
			//unmarshal json to struct
			err = json.Unmarshal(data, &jsondata)

			fmt.Println(len(jsondata))
			//range through subreddit submissions
			for i, _ := range jsondata {
				//if it is valid
				if jsondata[i].Subm.Title != "" {
					//Information on Words of Subreddits
					// wordcount(jsondata[i].Subm.Title, usagesubmstitle[l])
					// wordcount(jsondata[i].Subm.Selftext, usagesubmsbody[l])
					//range through comments
					for j, _ := range jsondata[i].Cmmnts {
						//add words to map for subreddit for the parent comment
						wordcountcomment(jsondata[i].Cmmnts[j], usagecmmts[fname])
						//send replies to parent comment to recursive func to get word data
						parsereplies(jsondata[i].Cmmnts[j], usagecmmts[fname])
					}

				}
			}
		}

		//struct information to help read output
		fmt.Println("Word Data Struct Has the following fields '\n' type SubsWordData struct  '\n' Word     string '\n' Numoccur float64'\n' Avgscore float64'\n' Heur     float64'\n' '\n'")
		fmt.Println("Comments Usage for", files.Name(), ": \n")

		//print out words that occur at least 30 times in the top all time submissions of this particular subreddit
		cnt := 0
		for _, v := range usagecmmts[fname] {
			if v.Numoccur > 30 {
				fmt.Println(v)
				cnt++
			}
		}
		fmt.Println("\n Number of words used at least 50 times (2x per thread)", cnt, "\n")
	}

	// //analyzing user information
	// for i, _ := range usermap {
	// 	tmpuser := usermap[i]
	// 	tmpscore := tmpuser.UserScore
	// 	for k, _ := range usermap[i].Words {
	// 		tmpscore += usermap[i].Words[k].Avgscore
	// 	}

	// 	tmpuser.UserScore = tmpscore
	// 	usermap[i] = tmpuser

	// 	tmpuser.Words = nil
	// 	// fmt.Println(tmpuser)

	// 	if usermap[i].UserScore > 5000 {
	// 		fmt.Println("User Info: ", usermap[i].User, "\n")
	// 		fmt.Println("Score: ", usermap[i].UserScore, "\n")
	// 		fmt.Println("Num words: ", len(usermap[i].Words), "\n")
	// 		fmt.Println("Num Subs: ", len(usermap[i].Subs), "\n")
	// 		fmt.Println("Subs: ", usermap[i].Subs, "\n")
	// 	}
	// }

	//determine listings type and amount for comment and submissions
	sublistingoptions := geddit.ListingOptions{Time: "all", Limit: 3}
	cmmntlistingoptions := geddit.ListingOptions{Time: "all"}

	//prepare struct for authentication with reddit with my oauth keys
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
	err = o.LoginAuth(config.Username, config.Password)
	if err != nil {
		log.Fatal(err)
	}

	//get my subreddits
	subs, err := o.MySubreddits()
	if err != nil {
		log.Fatal(err)
	}

	//range through subs
	for _, val := range subs {
		//prepare struct and place it in map
		SubCommentData := SubsCommentData{SubName: val.Name}
		SubCommentMap[val.Name] = SubCommentData

		//get submissions from that subreddit
		submissions, err := o.SubredditSubmissions(val.Name, geddit.HotSubmissions, sublistingoptions)
		if err != nil {
			log.Fatal(err)
		}

		//range through the submissions
		for _, submvalue := range submissions {
			fmt.Println("Submission: ", submvalue, "\n")

			if len(submvalue.Title) != 0 {

				//get comments
				comments, err := o.Comments(submvalue, geddit.TopSubmissions, cmmntlistingoptions)
				if err != nil {
					log.Fatal(err)
				}

				//range through comments
				for _, comvalue := range comments {
					//determine the vote of the parent
					SubCommentData = determinevote(comvalue, submvalue, SubCommentData)
					//if reply is valid, range through them and determine their votes
					if len(comvalue.Replies) != 0 {
						for _, replyvalue := range comvalue.Replies {
							SubCommentData = determinevote(replyvalue, submvalue, SubCommentData)

						}
					}
				}
			}

		}

		//place the data back into the map
		SubCommentMap[val.Name] = SubCommentData

		fmt.Println("\n")
		fmt.Println("Subreddit: ", val.Name, " is done")
		fmt.Println(SubCommentData)

	}

	//when done, range through all data
	for _, val := range subs {
		tmpdata := SubCommentMap[val.Name]
		fmt.Println("\n")
		fmt.Println("Subreddit: ", val.Name, " \n")
		fmt.Println("Number of Comments: ", tmpdata.Numcomments, " \n")
		fmt.Println("Average Number of Upvotes on Upvoted posts: ", tmpdata.AvgUpVote, " \n")
		fmt.Println("Number of Upvotes: ", tmpdata.NumUpVotes, " \n")
		fmt.Println("Average Number of Downvotes on Downvoted posts: ", tmpdata.AvgDownVotes, " \n")
		fmt.Println("Number of Downvotes: ", tmpdata.NumDownVotes, " \n")
		fmt.Println("Average posts score on the subreddit: ", tmpdata.AvgPostScore, "\n")
		fmt.Println("Performance metric: ", tmpdata.AvgUpVote/tmpdata.AvgPostScore, "\n")
		fmt.Println("\n")

	}

}

//heuristic func to determine action
func determinevote(x *geddit.Comment, y *geddit.Submission, subdata SubsCommentData) SubsCommentData {
	//prepare oauth struct
	o = geddit.OAuthSession{}
	o, err := geddit.NewOAuthSession(
		config.Clientid,
		config.Clientsecret,
		"school project to determine whether or not a comment would be upvoted",
		"http://redirect.url",
	)
	if err != nil {
		log.Fatal(err)
	}

	//make a http client for voting purposes
	err = o.LoginAuth(config.Username, config.Password)
	if err != nil {
		log.Fatal(err)
	}

	//get each individual word from the body of the comment
	substrs := strings.Fields(x.Body)

	//prepare score variables
	var score, timestopped float64
	score = 0
	timestopped = 0

	//range through words in comment
	for i, _ := range substrs {
		//remove caps
		substrs[i] = strings.ToLower(substrs[i])
		//if word exists in the stopwords, keep track
		if _, stopped := stopwords[substrs[i]]; stopped {
			timestopped += 1
		} else {
			if word, exists := usagecmmts[subdata.SubName][substrs[i]]; exists {
				//if the wor dexists in the map, add to the score
				score += word.Avgscore
			}
		}

	}

	//get the average score per word
	score = score / (float64(len(substrs)) - timestopped)

	//apply amt of time calcs
	score = score / (x.Created - y.DateCreated) * 60
	//track number of comments
	subdata.Numcomments += 1
	//get the average score
	subdata.AvgPostScore = ((subdata.AvgPostScore * (subdata.Numcomments - 1)) + (x.UpVotes - x.DownVotes)) / (subdata.Numcomments)

	//if score is good enough to upvote, do it and log it
	if score > .5 {
		err = o.Vote(x, geddit.UpVote)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("---------------New Vote--------------------\n")
		fmt.Println("MyScore: ", score, "\n")
		fmt.Println("Upvotes: ", x.UpVotes, "\n")
		fmt.Println("avgUpvotes: ", subdata.AvgUpVote, "\n")
		fmt.Println("numcomments: ", subdata.Numcomments, "\n")
		fmt.Println("avgscore: ", subdata.AvgPostScore)

		//update stats
		subdata.AvgUpVote = ((subdata.AvgUpVote * subdata.NumUpVotes) + x.UpVotes) / (subdata.NumUpVotes + 1)
		subdata.NumUpVotes += 1

		fmt.Println("avgup: ", subdata.AvgUpVote)
		fmt.Println("ups: ", subdata.NumUpVotes)
		fmt.Println("Upvoted! Score: ", score, " Post: ", x.Body)
		fmt.Println("Expected sub(cmmnt): ", x.Subreddit, " Expected sub (subm): ", y.Subreddit, " Subvoted in: ", subdata.SubName, " \n")
		fmt.Println(subdata)
		fmt.Println("------------------------------------------\n")

	} else if score < -.5 {
		err = o.Vote(x, geddit.DownVote)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("---------------New Vote -------------------\n")
		fmt.Println("MyScore: ", score, "\n")
		fmt.Println("Upvotes: ", x.UpVotes, "\n")
		fmt.Println("AvgDownVotes: ", subdata.AvgDownVotes, "\n")
		fmt.Println("numcomments: ", subdata.Numcomments, "\n")

		//update stats
		subdata.AvgDownVotes = ((subdata.AvgDownVotes * subdata.NumDownVotes) + (x.UpVotes - x.DownVotes)) / (subdata.NumDownVotes + 1)
		subdata.NumDownVotes += 1

		fmt.Println("avgdown: ", subdata.AvgDownVotes)
		fmt.Println("downs: ", subdata.NumDownVotes)
		fmt.Println("Downvoted! Score: ", score, " Post: ", x.Body)

		fmt.Println("---------------------------------------------\n")

	}
	return subdata

}

//function to determine word usage from a comment,
//and update user map and word map
func wordcountcomment(x geddit.Comment, words map[string]SubsWordData) {
	//prepare vars
	substrs := strings.Fields(x.Body)
	tmpdata := SubsWordData{}
	tmpUser := UserData{}
	//regex to remove trailing and leading punctuation
	reg, err := regexp.Compile(`[^0-9a-zA-Z-]`)
	if err != nil {
		panic(err)
	}

	//log user information
	tmpUser, uexists := usermap[x.Author]

	//if no user exists
	if !uexists {

		tmpUser = UserData{
			User:  x.Author,
			Words: make(map[string]SubsWordData),
		}

	}

	//range through individual words
	for _, word := range substrs {
		//remove anything but alphanumeric
		word = reg.ReplaceAllString(word, "")
		//get rid of words like "I"
		if len(word) > 1 {
			//determine if word is stopword
			if _, stopped := stopwords[word]; !stopped {

				tmpdata = SubsWordData{}
				_, ok := words[strings.ToLower(word)]

				if ok == true {
					//if that worddata exists in the map
					tmpdata = words[word]

					tmpdata.Avgscore = ((tmpdata.Avgscore * tmpdata.Numoccur) + x.UpVotes) / (tmpdata.Numoccur + 1)
					tmpdata.Numoccur += 1
					tmpdata.Heur += x.UpVotes
					// tmpdata.TimePassed =

				} else {
					//if no worddata exists
					tmpdata = SubsWordData{
						Word:     strings.ToLower(word),
						Numoccur: 1,
						Avgscore: x.UpVotes,
						Heur:     x.UpVotes,
					}

				} //endelse

				//add word to map
				words[word] = tmpdata

				//empty word data for user
				tmpword := SubsWordData{}

				if userword, wordexists := tmpUser.Words[word]; wordexists {
					//check if data exists for author, if so update
					tmpword.Avgscore = ((userword.Avgscore * userword.Numoccur) + x.UpVotes) / (userword.Numoccur + 1)
					tmpword.Numoccur += 1
					tmpword.Heur += x.UpVotes

				} else {
					//create the data for the word
					tmpword.Avgscore = x.UpVotes
					tmpword.Numoccur = 1
					tmpword.Heur = x.UpVotes
				}

				//update word in user's word map
				tmpUser.Words[word] = tmpword
				// fmt.Println(tmpword)

			}
		}

	}
	//update user in global usermap
	tmpUser.Subs = append(tmpUser.Subs, x.Subreddit)
	usermap[x.Author] = tmpUser

}

//recursive function to determine if there are more replies to be evaluated
//takes in a reddit comment and a map
func parsereplies(x geddit.Comment, words map[string]SubsWordData) {

	//determine word usage from comment
	wordcountcomment(x, words)

	//if there are replies, call this function on each of them.
	if len(x.Replies) != 0 {
		for i, _ := range x.Replies {

			parsereplies(*x.Replies[i], words)

		}
	}
}
