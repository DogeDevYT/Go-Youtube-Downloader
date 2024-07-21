package main

import (
	"fmt"

	Download "github.com/DogeDevYT/go-commandline-youtube/download"
)

func main() {
	fmt.Println("Welcome to Command Line Youtube!")

	//initalize variable to store choice of what user would like to do
	var choice int

	//Create menu for user to pick what to do
	fmt.Println("What would you like to do? (1-9)")
	fmt.Println("(1) Download a YouTube video")

	//Handle selection of input choice
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		//insert download logic here
		err := Download.DownloadYT()

		if err != nil {
			fmt.Printf("Error downloading video: %v", err)
		}
	default:
		fmt.Println("Invalid choice selected! Goodbye!")
	}
}
