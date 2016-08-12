package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func Talk(request string, token string) map[string]interface{} {
	resp, err := http.Get("https://api.vk.com/method/" + request +
		"&access_token=" + token + "&v=5.53")

	if err != nil {
		log.Fatalf("Failed to send request (%s): %s\n", request, err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf(
			"Failed to receive response body for sent request (%s): %s\n",
			request, err,
		)
	}

	var val interface{}

	json.Unmarshal(body, &val)

	obj := val.(map[string]interface{})

	if err, has_err := obj["error"].(map[string]interface{}); has_err {
		log.Fatalln("Server returned an error: %s (%d)", err["error_msg"], err["error_code"])
	}

	return obj["response"].(map[string]interface{})
}
