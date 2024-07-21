package download

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kkdai/youtube/v2"
)

// Helper method to manage downloading a youtube stream
func downloadYTStream(client *youtube.Client, video *youtube.Video, format youtube.Format, filename string) error {
	stream, _, err := client.GetStream(video, &format) // get youtube video stream
	if err != nil {
		return fmt.Errorf("error getting yt stream: %w", err)
	}
	file, err := os.Create(filename) //Create output file
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()
	_, err = io.Copy(file, stream) //copy yt stream to output file
	if err != nil {
		return fmt.Errorf("error copying youtube stream to output file: %w", err)
	}
	return nil //no error; run through was successfull
}

// Downloads Youtube video
func DownloadYT() error {
	client := youtube.Client{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the YouTube video URL: ")
	videoURL, _ := reader.ReadString('\n')
	videoURL = strings.TrimSpace(videoURL)

	video, err := client.GetVideo(videoURL)
	if err != nil {
		log.Fatalf("Error getting video: %v", err)
	}

	fmt.Println("Available Formats:")
	for i, format := range video.Formats {
		fmt.Printf("%d: Quality: %s, MimeType: %s, AudioQuality: %s\n", i, format.QualityLabel, format.MimeType, format.AudioQuality)
	}

	var videoFormats []youtube.Format
	var audioFormats []youtube.Format

	for _, format := range video.Formats {
		if strings.Contains(format.MimeType, "audio") {
			audioFormats = append(audioFormats, format)
		} else if strings.Contains(format.MimeType, "video") {
			videoFormats = append(videoFormats, format)
		}
	}

	if len(audioFormats) == 0 {
		log.Println("No separate audio formats found. Attempting to find combined streams.")
		for _, format := range video.Formats {
			if strings.Contains(format.MimeType, "video") && format.AudioChannels > 0 {
				audioFormats = append(audioFormats, format)
			}
		}
	}

	if len(audioFormats) == 0 {
		log.Fatalf("No audio formats found.")
	}

	fmt.Println("Available video formats:")
	for i, format := range videoFormats {
		fmt.Printf("%d: %s (Quality: %s, Itag: %d)\n", i, format.MimeType, format.QualityLabel, format.ItagNo)
	}

	fmt.Println("Available audio formats:")
	for i, format := range audioFormats {
		fmt.Printf("%d: %s (AudioQuality: %s, Itag: %d)\n", i, format.MimeType, format.AudioQuality, format.ItagNo)
	}

	fmt.Print("Enter the number of the desired video quality: ")
	videoChoiceStr, _ := reader.ReadString('\n')
	videoChoiceStr = strings.TrimSpace(videoChoiceStr)
	var videoChoice int
	fmt.Sscanf(videoChoiceStr, "%d", &videoChoice)

	if videoChoice < 0 || videoChoice >= len(videoFormats) {
		log.Fatalf("Invalid video choice")
	}
	selectedVideoFormat := videoFormats[videoChoice]

	fmt.Print("Enter the number of the desired audio quality: ")
	audioChoiceStr, _ := reader.ReadString('\n')
	audioChoiceStr = strings.TrimSpace(audioChoiceStr)
	var audioChoice int
	fmt.Sscanf(audioChoiceStr, "%d", &audioChoice)

	if audioChoice < 0 || audioChoice >= len(audioFormats) {
		log.Fatalf("Invalid audio choice")
	}
	selectedAudioFormat := audioFormats[audioChoice]

	videoFilename := "video.mp4"
	audioFilename := "audio.mp3"
	outputFilename := "output.mp4"

	deleteRemenants(videoFilename, audioFilename, outputFilename)

	err = downloadYTStream(&client, video, selectedVideoFormat, videoFilename)
	if err != nil {
		log.Fatalf("Error downloading video stream: %v", err)
	}
	err = downloadYTStream(&client, video, selectedAudioFormat, audioFilename)
	if err != nil {
		log.Fatalf("Error downloading audio stream: %v", err)
	}

	cmd := exec.Command("ffmpeg", "-i", videoFilename, "-i", audioFilename, "-c:v", "copy", "-c:a", "aac", outputFilename)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error merging video and audio: %v", err)
	}

	os.Remove(videoFilename)
	os.Remove(audioFilename)

	fmt.Println("Video and audio downloaded and merged successfully")

	return nil
}

// Checks to see if file exists at a path
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Delete any remenants of previous runs like audio.mp3, video.mp4, or output.mp4
func deleteRemenants(videoFilename string, audioFilename string, outputFilename string) error {
	//Delete any remenants of previous runs
	var remenantRemovalArray []string
	remenantRemovalArray = append(remenantRemovalArray, videoFilename)
	remenantRemovalArray = append(remenantRemovalArray, audioFilename)
	remenantRemovalArray = append(remenantRemovalArray, outputFilename)

	for _, remenant := range remenantRemovalArray {
		if fileExists(remenant) {
			fmt.Printf("Deleting remenant: %s\n", remenant)
			os.Remove(remenant)
		}
	}

	return nil
}
