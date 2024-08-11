package main

import (
	"os"
	"os/exec"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	// constant
	const defDuration = 60

	inputVideoArg := os.Args[1] // argument: video filename without path

	// verify the filename
	_, err := os.Open(inputVideoArg)
	if err != nil {
		log.Fatal(err)
	}

	// get video duration
	out, _ := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=duration", "-of", "default=noprint_wrappers=1:nokey=1", inputVideoArg).Output()

	outStr := string(out[:])
	outSplit := strings.Split(outStr, ".")
	s, _ := strconv.Atoi(outSplit[0])

	fmt.Println("Video duration is " + getTime(s))

	counter := 1 // will be used for naming the output video file
	outputVideo := "" // output video file name
	startTime, endTime := "", "" // placeholder for starttime and endtime

	for i := 0; i < s; i=i + defDuration {
		outputVideo = strconv.Itoa(counter) + ".mp4" // assuming all file is mp4

		startTime = getTime(i)
		endTime = getTime(i + defDuration)

		fmt.Println("Processing " + outputVideo)
		exec.Command("ffmpeg", "-ss", startTime, "-to", endTime, "-i", inputVideoArg, outputVideo).Output()

		counter++
	}
}


func getTime(s int) string {
	if checkSeconds(s) {
		a, b, c := getHours(s)
		return standardFormat(a,b,c)
	} else {
		a, b := getMinutes(s)
		return standardFormat(0,a,b)
	}
}

// function to check whether the seconds is equal to "1 hour or more"
func checkSeconds(seconds int) bool {
	if seconds >= 3600 {
		return true
	} else {
		return false
	}
}

// function to get minutes, and remainder seconds
func getMinutes(seconds int) (int, int) {
	m := seconds/60
	s := seconds%60

	return m, s
}

// function to get hour, remainder minutes and remainder seconds
func getHours(seconds int) (int, int, int) {
	h := seconds/3600
	hr := seconds%3600 // h remainder
	m, s := getMinutes(hr)

	return h, m, s
}

// function to generate hh:mm:ss
func standardFormat(h, m, s int) string {
	hs, ms, ss := "00", "00", "00" // default string

	if isLessTen(h) {
		hs = "0" + strconv.Itoa(h)
	} else {
		hs = strconv.Itoa(h)
	}

	if isLessTen(m) {
		ms = "0" + strconv.Itoa(m)
	} else {
		ms = strconv.Itoa(m)
	}

	if isLessTen(s) {
		ss = "0" + strconv.Itoa(s)
	} else {
		ss = strconv.Itoa(s)
	}

	return hs + ":" + ms + ":" + ss
}

//
func isLessTen(num int) bool {
	if num<10 {
		return true
	} else {
		return false
	}
}