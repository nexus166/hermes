package hermes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// POSTMessage /api/v1/chat.postMessage
func POSTMessage(t POSTTemplate) (int, error) {
	if t.Channel == "" && t.RoomID == "" {
		t.Channel = RCChannel
		t.RoomID = RCRoomID
	}
	b, err := MakePayload(t)
	if err != nil {
		return -1, err
	} else if Verbose {
		fmt.Println(string(b.Bytes()))
	}

	req, err := http.NewRequest("POST", RCURL+"/api/v1/chat.postMessage", b)
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", RCToken)
	req.Header.Set("X-User-Id", RCTokenID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	if Verbose {
		var response map[string]interface{}
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return -1, err
		}
		j, err := json.MarshalIndent(response, "", "\t")
		if err != nil {
			return -1, err
		}
		fmt.Println(string(j))
	}
	if resp.StatusCode > 200 {
		return resp.StatusCode, fmt.Errorf("bad response from RC server")
	}
	return resp.StatusCode, nil
}
