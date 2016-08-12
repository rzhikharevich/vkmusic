package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func eq(str, search string) bool {
	if len(search) == 0 {
		return true
	}

	return strings.EqualFold(str, search)
}

func songFail(desc string) {
	fmt.Fprintf(
		os.Stderr,
		"\r[%s✗%s] %s\n",
		RedColor, ResetColor,
		desc,
	)
}

func songOk(desc string) {
	fmt.Printf(
		"\r[%s✓%s] %s\n",
		GreenColor, ResetColor,
		desc,
	)
}

func main() {
	log.SetFlags(0)

	token := Auth()

	resp := Talk("audio.get?count=6000", token)

	artist_dirs := make(map[string]bool)

	for _, arg := range os.Args[1:] {
		comp := strings.Split(arg, ":")

		artist := comp[0]
		if len(comp) == 1 {
			artist = ""
		}

		title := comp[len(comp)-1]

		for _, item := range resp["items"].([]interface{}) {
			item := item.(map[string]interface{})

			item_artist := item["artist"].(string)
			item_title := item["title"].(string)
			url := item["url"].(string)

			if !eq(item_artist, artist) || !eq(item_title, title) {
				continue
			}

			desc := item_artist + " – " + item_title

			fmt.Print("[*] ", desc)

			if _, dir_exists := artist_dirs[item_artist]; !dir_exists {
				if err := os.MkdirAll(item_artist, 0777); err != nil {
					songFail(desc)
					continue
				}
			}

			artist_dirs[item_artist] = true

			path := item_artist + string(filepath.Separator) + item_title + ".mp3"

			if _, err := os.Stat(path); !os.IsNotExist(err) {
				songOk(desc)
				continue
			}

			resp, err := http.Get(url)
			if err != nil {
				songFail(desc)
				continue
			}

			out, err := os.Create(path)
			if err == nil {
				_, err := io.Copy(out, resp.Body)
				if err != nil {
					songFail(desc)
					os.Remove(path)
				} else {
					songOk(desc)
				}

				out.Close()
			} else {
				songFail(desc)
			}

			resp.Body.Close()
		}
	}
}
