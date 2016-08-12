package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/toqueteos/webbrowser"
)

func Auth() string {
	portn := 22091
	ports := strconv.Itoa(portn)

	auth := "https://oauth.vk.com/oauth/authorize?" +
		// token receiver URI
		"redirect_uri=http://localhost:" + ports + "/rewrap&" +
		// ask for audio access
		"scope=8&" +
		// get an access token
		"response_type=token&" +
		// API version
		"v=5.53&" +
		// app identifier
		"client_id=5581005"

	ch := make(chan string)

	http.HandleFunc("/rewrap", func(w http.ResponseWriter, r *http.Request) {
		c := "<script>" +
			"location.href = '/auth?'+location.hash.substr(1)" +
			"</script>"

		fmt.Fprintln(w, c)
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<script>window.close()</script>")

		q := r.URL.Query()

		token := q.Get("access_token")

		if len(token) == 0 {
			log.Fatalf("Authentification failure: %s (%s)\n",
				q.Get("error_description"), q.Get("error"))
		}

		ch <- token
	})

	listener, err := net.Listen("tcp", ":"+ports)
	if err != nil {
		log.Fatalln("Can't listen on port", ports+":", err)
	}

	defer listener.Close()

	go http.Serve(listener, nil)

	webbrowser.Open(auth)

	token := <-ch

	return token
}
