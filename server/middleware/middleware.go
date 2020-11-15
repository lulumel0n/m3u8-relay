package middleware

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"net/http"

	"github.com/lulumel0n/m3u8-relay/server/model"
)

// GetStreaming get m3u file
func GetStreaming(w http.ResponseWriter, r *http.Request) {
	payload := getFromRadio()
	sendGeneric(w, r, []byte(payload))
}

func sendGeneric(w http.ResponseWriter, r *http.Request, payload []byte) {
	w.Header().Set("Context-Type", "application/octet-stream") // send
	w.Header().Set("Content-Disposition", "attachment; filename=filename.m3u")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	bytesSent, err := w.Write(payload)
	if err != nil {
		fmt.Println("Failed to send")
		fmt.Println(err)
	} else {
		fmt.Printf("Sent %d\n", bytesSent)
	}
}

func getFromRadio() string {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	req, err := http.NewRequest("GET", model.ENDPOINT+"live_11.m3u8", nil)
	if err != nil {
		// handle err
		fmt.Println(err.Error())
		return "go die"
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		fmt.Println(err.Error())
		return "go die"
	}

	defer resp.Body.Close()

	return transformResponse(resp.Body)
}

func transformResponse(data io.ReadCloser) string {
	sb := new(strings.Builder)
	tsliner, _ := regexp.Compile(".*[.]ts")

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		if tsliner.MatchString(scanner.Text()) {
			sb.WriteString(model.ENDPOINT)
			sb.WriteString(scanner.Text())
		} else {
			sb.WriteString(scanner.Text())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
