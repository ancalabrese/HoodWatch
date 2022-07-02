package core

import "time"

type Crime struct {
	ID           uint
	ReportedTime time.Time
	Type         string
	Payload      map[string]string
}

const CRIME_PAYLOAD_DATE_KEY = "DATE_KEY"
const CRIME_PAYLOAD_MSG_KEY = "MSG_KEY"

func NewCrime(crimeType string, payload map[string]string) *Crime {
	//TODO: generate unique ID for each crime.
	reportedTime := time.Now()
	return &Crime{
		ReportedTime: reportedTime,
		Payload:      payload,
		Type:         crimeType,
	}
}
