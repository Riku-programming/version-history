package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const latestFileName string = "-latest.txt"
const compareFileName string = "-compare.txt"
const versionHistoryFileName string = "-version_history.txt"
const dateLayout string = "2006-01-02"

//const URL string = "https://test-nesic-cp.axlbox.biz/common/versions.txt"

func main() {
	file, err := os.Open("client_host.txt")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		clientData := fileScanner.Text()
		arr := strings.Split(clientData, ",")
		clientName := arr[0]
		URL := arr[1]
		if err := writeLine(clientName+compareFileName, fetchVersion(URL)); err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
		if err := fileScanner.Err(); err != nil {
			log.Fatalf("Error while reading file %s", err)
		}
		if deepCompare(clientName+latestFileName, clientName+compareFileName) == false {
			fmt.Println("Difference")
			removeFile(clientName + latestFileName)
			renameFile(clientName+compareFileName, clientName+latestFileName)
			appendHistory(clientName+versionHistoryFileName, URL)
		} else {
			fmt.Println("Same")
			removeFile(clientName + compareFileName)
		}
	}
	err = file.Close()
	if err != nil {
		return
	}
	return
}

func fetchVersion(URL string) string {
	curl := exec.Command("curl", URL)
	out, err := curl.Output()
	if err != nil {
		fmt.Println("error", err)
		return "error"
	}
	return string(out)
}

//
func writeLine(fileName, lines string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	for _, line := range lines {
		_, err := file.WriteString(string(line))
		if err != nil {
			return err
		}
	}
	return nil
}

func appendHistory(fileName, URL string) {
	nowTime := time.Now().Format(dateLayout)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	_, err = fmt.Fprintln(file, nowTime)
	if err != nil {
		return
	}
	_, err = fmt.Fprintln(file, fetchVersion(URL))
	if err != nil {
		return
	}
}

func deepCompare(file1, file2 string) bool {
	sf, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	df, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	sscan := bufio.NewScanner(sf)
	dscan := bufio.NewScanner(df)

	for sscan.Scan() {
		dscan.Scan()
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			return false
		}
	}
	return true
}

func removeFile(filename string) {
	if err := os.Remove(filename); err != nil {
		fmt.Println(err)
	}
}

func renameFile(oldFilename, newFilename string) {
	if err := os.Rename(oldFilename, newFilename); err != nil {
		fmt.Println(err)
	}
}

// ファイルから1行ずつ読み込む　それをURLとして読み込む
// nesic,https://test-nesic-cp.axlbox.biz/common/versions.txt
// こんな感じで値が出るから nesic = Arr[0] https: Arr[1]
// forでまわす
