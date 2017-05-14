package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	//get time
	t := time.Now()
	// appendstr := "go run C:\\Users\\Myke\\Dropbox\\School\\CSCI164\\Project\\DataUsage.go > C:\\Users\\Myke\\Dropbox\\School\\CSCI164\\Project\\runs\\" + t.Format("Jan2-15_04-2006") + ".txt"
	// fmt.Println(appendstr)

	// cmd := exec.Command("cmd", "/C", appendstr)
	// output, err := cmd.CombinedOutput()
	// fmt.Printf("%s\n", output)

	// // err := cmd.Run()
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", output)
	// fmt.Println(string(cmd))
	fmt.Println("entering loop")

	//while loop
	for {
		t = time.Now()
		appendstr := "go run C:\\Users\\Myke\\Dropbox\\School\\CSCI164\\Project\\DataUsage.go > C:\\Users\\Myke\\Dropbox\\School\\CSCI164\\Project\\runs\\" + t.Format("Jan2-15_04-2006") + ".txt"

		//if it is a fresh hour divisible by 2
		if t.Hour()%2 == 0 {
			if t.Minute() == 0 {
				cmd := exec.Command("cmd", "/C", appendstr)
				fmt.Println("Ran Program!")
				err := cmd.Start()
				if err != nil {
					panic(err)
				}
			}
		}

		time.Sleep(1 * time.Minute)

	}

}
