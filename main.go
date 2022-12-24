package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	tok, err := os.ReadFile(".env")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	filename := os.Args[1]
	fc, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	files := map[string]File{string(filename)[strings.LastIndex(filename, "/")+1:]: {string(fc)}}

	gistRequest := GistRequest{
		Files:       files,
		Description: "",
		Public:      false,
	}
	gistRequestJson, err := json.Marshal(gistRequest)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	c := http.Client{Timeout: time.Duration(1) * time.Second}

	req, err := http.NewRequest("POST", "https://api.github.com/gists", bytes.NewBuffer(gistRequestJson))
	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", string(tok)))

	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	fmt.Printf("Body:  %s\n ", body)
	fmt.Printf("Response status:  %s\n", resp.Status)
}
