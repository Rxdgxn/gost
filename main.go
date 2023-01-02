package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type GistRequest struct {
	Files       map[string]File `json:"files"`
	Description string          `json:"description"`
	Public      bool            `json:"public"`
}
type File struct {
	Content string `json:"content"`
}

func main() {
	var public bool
	bgpoint := 1
	tok, err := os.ReadFile(".env")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	if os.Args[1] == "--pb" {
		public = true
		bgpoint = 2
	} else if os.Args[1] == "--pv" {
		public = false
		bgpoint = 2
	}
	for i := bgpoint; i < len(os.Args); i++ {
		filename := os.Args[i]
		fc, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("%s", err)
			return
		}
		files := map[string]File{string(filename)[strings.LastIndex(filename, "/")+1:]: {string(fc)}}

		gistRequest := GistRequest{
			Files:       files,
			Description: "",
			Public:      public,
		}
		gistRequestJson, err := json.Marshal(gistRequest)
		if err != nil {
			fmt.Printf("%s", err)
			return
		}

		c := http.Client{Timeout: time.Duration(1) * time.Second}

		req, err := http.NewRequest("POST", "https://api.github.com/gists", bytes.NewBuffer(gistRequestJson))
		if err != nil {
			fmt.Printf("Error %s", err)
			return
		}
		req.Header.Add("Accept", `application/json`)
		req.Header.Add("Authorization", fmt.Sprintf("token %s", string(tok)))

		resp, err := c.Do(req)
		if err != nil {
			fmt.Printf("Error %s", err)
			return
		}

		defer resp.Body.Close()
		fmt.Printf("Response status for file \"%s\": %s\n", filename, resp.Status)
	}
}
