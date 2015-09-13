package gcm

import (
	"net/http"
	"bytes"
	"fmt"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const(
	gcmUrl = "https://gcm-http.googleapis.com/gcm/send"
	gcmTable = "gcm_registration"
)

type SendBody struct{
	To string `json:"to"`
	Data SendData `json:"data"`
}

type SendAll struct{
	RegIds []string `json:"registration_ids"`
	Data SendData `json:"data"`
}

type SendData struct{
	Title string `json:"title"`
	Message string `json:"message"`
}

func SendToAll(w http.ResponseWriter, r *http.Request) {
	var d SendData
	decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&d)
    if err != nil {
        panic(err.Error())
    }

	db, err := sql.Open("mysql", "/test")
	if err != nil {
   		panic(err.Error())
    }
    defer db.Close()

	s := make([]string, 1)
	rows, err := db.Query("SELECT gcmToken FROM "+gcmTable + " LIMIT 1000")
    for rows.Next() {
        var gcmToken string
        err = rows.Scan(&gcmToken)
        fmt.Println(gcmToken)
        s = append(s, gcmToken) 
    }

    sendAllObj := SendAll{s, d}

    fmt.Println(sendAllObj)

	jsonStr, _ := json.Marshal(sendAllObj)

	fmt.Println("bytes: "+string(jsonStr))

	req, _ := http.NewRequest("POST", gcmUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key=AIzaSyA1zwySLkbBiiG50wLr5sA8jMxxu2_C9fE")

	client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)

//    body, _ := ioutil.ReadAll(resp.Body)
//   fmt.Println("response Body:", string(body))

    var data Response
    if(resp.Status == "200 OK"){
    	data = Response{true, "Sent Successfully"}
    }else{
    	data = Response{false, "Could not be sent"}
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

type Gcm struct{
	RegId string `json:"regId"`
	Imei string `json:"imei"`
	DeviceId string `json:"deviceId"`
	Other OtherItem `json:"other"`
}

type OtherItem struct{
	Location string `json:"location"`
}

func InsertIntoGcm() {
	db, err := sql.Open("mysql", "/test")
	if err != nil {
   		panic(err.Error())
    }
    defer db.Close()

    query := "INSERT IGNORE INTO "+gcmTable+" (gcmToken, deviceId, imei, other, createdAt) VALUES(? ,? , ? , ?, ?)"
    stmtIns, err := db.Prepare(query)
    if err != nil {
    	panic(err.Error())
    }
    defer stmtIns.Close()

    cur_date := time.Now()

    b, err := json.Marshal(g.Other)
    if err != nil {
        fmt.Println(err)
    }    
    _, err = stmtIns.Exec(g.RegId, g.DeviceId, g.Imei, string(b), cur_date)
    if err != nil {
    	panic(err.Error())
    }
}

var g Gcm 
func Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&g)
    if err != nil {
        panic(err.Error())
    }
    go InsertIntoGcm()
    data := &Response{true, "Registered Successfully"}
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

type Response struct{
	Success bool `json:"success"`
	Message string `json:"message"`
}

type GcmResponse struct{
  MulticastId string `json:"multicast_id"`
  Success int `json:"success"`
  Failure int `json:"failure"`
  CanonicalIds int `json:"canonical_ids"`
  Results []GcmResult `json:"results"`
}

type GcmResult struct{
	MessageId string `json:"message_id"`
}