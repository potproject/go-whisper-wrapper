package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sqweek/dialog"
)

var CHATGPT_WHISPER_API string
var CHATGPT_API_KEY string

var EXPORT_EXTENTION string       // csv or srt
var TRANSCRIPTION_LANGUAGE string // example: ja
var SERVICE string                // openai or oss

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	CHATGPT_WHISPER_API = os.Getenv("CHATGPT_WHISPER_API")
	CHATGPT_API_KEY = os.Getenv("CHATGPT_API_KEY")
	EXPORT_EXTENTION = os.Getenv("EXPORT_EXTENTION")
	TRANSCRIPTION_LANGUAGE = os.Getenv("TRANSCRIPTION_LANGUAGE")
	SERVICE = os.Getenv("SERVICE")
}

func main() {

	defer bufio.NewReader(os.Stdin).ReadBytes('\n')
	defer fmt.Print("Press 'Enter' to continue...")

	fmt.Println("Start Processing...")

	filePaths := []string{}
	if len(os.Args) < 2 {
		filePath, err := dialog.File().Load()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		filePaths = append(filePaths, filePath)
	} else {
		filePaths = os.Args[1:]
	}

	// 複数ファイル対応
	// os.Argsが空になるまで標準入力を読み込む
	for index, filePath := range filePaths {

		fmt.Println("Input: " + filePath)

		fileDirectory := filepath.Dir(filePath)
		fileBaseName := filepath.Base(filePath[:len(filePath)-len(filepath.Ext(filePath))])

		var text string
		var err error
		if SERVICE == "openai" {
			text, err = whisper(filePath)
		} else {
			text, err = whisperOss(filePath)
		}

		if err != nil {
			fmt.Printf("Error: %s\n", err)
			// 次のファイルがあるか確認
			if index+1 < len(os.Args) {
				fmt.Println("File failed to process. Continue to next file.")
			}
			continue
		}

		if EXPORT_EXTENTION == "csv" {
			text = verboseJsonToCSV(text)
		}

		fmt.Println("Done.")
		fmt.Println("Output: " + fileBaseName + "." + EXPORT_EXTENTION)

		filePath := filepath.Join(fileDirectory, fileBaseName+"."+EXPORT_EXTENTION)
		err = os.WriteFile(filePath, []byte(text), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("All Done.")
}

func whisperOss(file string) (string, error) {
	// exec: whisper %1 --output_format ${format}  --language ${lang} --model large
	var format string
	if EXPORT_EXTENTION == "csv" {
		format = "json"
	}
	if EXPORT_EXTENTION == "srt" {
		format = "srt"
	}
	cmd := exec.Command("whisper", file, "--output_format", format, "--language", TRANSCRIPTION_LANGUAGE, "--model", "large", "--output_dir", filepath.Dir(file))

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	fileDirectory := filepath.Dir(file)
	fileBaseName := filepath.Base(file[:len(file)-len(filepath.Ext(file))])
	outputFilePath := filepath.Join(fileDirectory, fileBaseName+"."+format)
	fileData, err := os.ReadFile(outputFilePath)
	if err != nil {
		return "", err
	}
	return string(fileData), nil
}

func whisper(file string) (string, error) {
	fileData, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	var format string
	if EXPORT_EXTENTION == "csv" {
		format = "verbose_json"
	}
	if EXPORT_EXTENTION == "srt" {
		format = "srt"
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("model", "whisper-1")
	writer.WriteField("language", TRANSCRIPTION_LANGUAGE)
	writer.WriteField("response_format", format)
	part, err := writer.CreateFormFile("file", "test.mp3")
	if err != nil {
		return "", err
	}
	part.Write(fileData)
	writer.Close()

	req, err := http.NewRequest("POST", CHATGPT_WHISPER_API, body)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+CHATGPT_API_KEY)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type WhisperVerboseJson struct {
	//Task     string                      `json:"task"`
	Language string `json:"language"`
	//Duration float64                     `json:"duration"`
	Text     string                      `json:"text"`
	Segments []WhisperVerboseJsonSegment `json:"segments"`
}

type WhisperVerboseJsonSegment struct {
	Id               int64   `json:"id"`
	Seek             int64   `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int64 `json:"tokens"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

func verboseJsonToCSV(text string) string {
	var whisperVerboseJson WhisperVerboseJson
	err := json.Unmarshal([]byte(text), &whisperVerboseJson)
	if err != nil {
		panic(err)
	}
	// to csv
	records := [][]string{}
	// header
	records = append(records, []string{
		"StartTime",
		"EndTime",
		"Text",
	})
	for _, segment := range whisperVerboseJson.Segments {
		records = append(records, []string{
			fmt.Sprintf("%g", segment.Start),
			fmt.Sprintf("%g", segment.End),
			segment.Text,
		})
	}

	var csvString strings.Builder
	writer := csv.NewWriter(&csvString)
	err = writer.WriteAll(records)
	if err != nil {
		panic(err)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		panic(err)
	}
	return csvString.String()
}
