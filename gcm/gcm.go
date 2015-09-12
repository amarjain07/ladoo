package gcm

import (
	"net/http"
	"bytes"
	"fmt"
)

func Send(w http.ResponseWriter, r *http.Request) {
	url := "https://gcm-http.googleapis.com/gcm/send"

	param := []byte(`{
	    "to" : "fZ3MmQOBHik:APA91bG7ty1MnRmO1d0xFDga_jDgRF3_WDHzYnnVSE9Pf9OCYmcdamzX9SfYaiU6b3kDsF8yluITfTxNDVOdjZZUTmLwihfErxn_P294mPNd_ZcfKOu4flKjlxzN429ssUoKUBk1968Q",
	    "notification" : {
	      "body" : "Keep it up",
	      "title" : "great work!",
	      "icon" : "myicon"
	    },
	    "data" : {
	      "title" : "Great!",
	      "message" : "Keep up the great work."
	    }
	  }`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key=AIzaSyA1zwySLkbBiiG50wLr5sA8jMxxu2_C9fE")
}

func Printg(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", r.URL.Path)
}