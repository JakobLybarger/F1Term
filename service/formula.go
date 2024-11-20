package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/JakobLybarger/formula/models"
)

const baseUri = "https://api.openf1.org/v1/"

type FormulaOneEvent struct {
	Meeting models.Meeting
	Session models.Session
	Drivers []models.Driver
}

func GetLiveData() FormulaOneEvent {
	meeting := GetMeeting()
	session := GetSession(meeting.MeetingKey)
	drivers := GetDrivers()

	return FormulaOneEvent{
		Meeting: meeting,
		Session: session,
		Drivers: drivers,
	}
}

func GetMeeting() models.Meeting {

	uri := baseUri + "meetings?meeting_key=latest"

	body, err := GetBody(uri)
	if err != nil {
		panic(err)
	}

	var meetings []models.Meeting
	if err := json.Unmarshal(body, &meetings); err != nil {
		panic(err)
	}

	if len(meetings) == 0 {
		panic("len 0")
	}

	return meetings[0]
}

func GetSession(meetingKey int) models.Session {

	uri := fmt.Sprintf("%s%s%d", baseUri, "sessions?session_key=latest&meeting_key=", meetingKey)

	body, err := GetBody(uri)
	if err != nil {
		panic(err)
	}

	var sessions []models.Session
	if err := json.Unmarshal(body, &sessions); err != nil {
		panic(err)
	}

	if len(sessions) == 0 {
		panic("len 0")
	}

	return sessions[0]
}

func GetDrivers() []models.Driver {

	uri := fmt.Sprintf("%s%s", baseUri, "/drivers?meeting_key=latest&session_key=latest")

	body, err := GetBody(uri)
	if err != nil {
		panic(err)
	}

	var drivers []models.Driver
	if err := json.Unmarshal(body, &drivers); err != nil {
		panic(err)
	}

	return drivers
}

type Position struct {
	DriverNumber int `json:"driver_number"`
	Pos          int `json:"postion"`
}

func GetPosition(meetingKey, sessionKey int) []Position {
	uri := fmt.Sprintf("%s%s%d%s%d", baseUri, "postion?meeting_key=", meetingKey, "&session_key=", sessionKey)

	body, err := GetBody(uri)
	if err != nil {
		panic(err)
	}

	var positions []Position
	if err := json.Unmarshal(body, &positions); err != nil {
		panic(err)
	}

	return positions
}

func GetBody(uri string) ([]byte, error) {
	http := http.Client{Timeout: time.Second * 10}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
