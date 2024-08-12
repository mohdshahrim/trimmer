package main

import (
	"os"
	"os/exec"
	"fmt"
	"log"
	"strconv"
	"strings"
	"flag"
)

func main() {
	// constant
	const defDuration = 60

	frontCutoffFlag := flag.String("min", "00:00:00", "cutoff at the beginning of the main video")
	backCutoffFlag := flag.String("max", "00:00:00", "cutoff at the ending of the main video")
	flag.Parse()

	//inputVideoArg := os.Args[1] // argument: video filename without path
	flagTails := flag.Args()
	inputVideoArg := flagTails[0]

	// validate front and back flag
	front := standardToSeconds(*frontCutoffFlag)
	back := standardToSeconds(*backCutoffFlag)

	// verify the filename
	_, err := os.Open(inputVideoArg)
	if err != nil {
		log.Fatal(err)
	}

	// get video duration
	out, _ := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=duration", "-of", "default=noprint_wrappers=1:nokey=1", inputVideoArg).Output()

	outStr := string(out[:])
	outSplit := strings.Split(outStr, ".")
	s, _ := strconv.Atoi(outSplit[0]) // duration in seconds

	fmt.Println("Original video duration is " + getTime(s))

	if back!=0 && back<s {
		s = back
	}

	fmt.Println("Cutoff video duration is " + getTime(s-front))

	counter := 1 // will be used for naming the output video file
	outputVideo := "" // output video file name
	startTime, endTime := "", "" // placeholder for starttime and endtime, to be used within loop


	for i := front; i < s; i=i + defDuration {
		outputVideo = strconv.Itoa(counter) + ".mp4" // assuming all file is mp4

		startTime = getTime(i)
		if i + defDuration > s {
			endTime = getTime(s)
		} else {
			endTime = getTime(i + defDuration)
		}

		fmt.Println("processing " + outputVideo)
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

// function to convert hh:mm:ss to seconds
func standardToSeconds(standardFormat string) int {
	splitted := strings.Split(standardFormat, ":")

	h, _ := strconv.Atoi(splitted[0]) //
	m, _ := strconv.Atoi(splitted[1]) //
	s, _ := strconv.Atoi(splitted[2]) //

	return (3600 * h) + (60 * m) + s
}