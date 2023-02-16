package pkg

import (
	"fmt"
	"os"
	"time"
)

type TimeTracker interface {
	UpdateLast(newTime time.Time) error
	GetLastTime() time.Time
}

type JumpCloudConnector interface {
	GetEventsSinceTime(time.Time) (*JumpCloudEvents, error)
}

// RunService is the main entry point for the service it will run a single time and return an error if one is encountered
func RunService(timeTracker TimeTracker, j JumpCloudConnector, pathToLogFile string) error {
	f, err := os.OpenFile(pathToLogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	lastTime := timeTracker.GetLastTime()
	e, err := j.GetEventsSinceTime(lastTime)
	if err != nil {
		return err
	}
	// Before doing anything make sure there is at least one event, if there isn't we don't need to do anything
	if len(e.Directory) == 0 && len(e.LDAP) == 0 && len(e.Systems) == 0 && len(e.SSO) == 0 && len(e.Radius) == 0 {
		return nil
	}
	lastEventSeen := time.Time{}
	// Loop over all events and find the newest timestamp, we will use this to update the last time we ran the service
	for _, x := range e.Directory {
		if x.Timestamp.After(lastEventSeen) {
			lastEventSeen = x.Timestamp
		}
		_, writeErr := f.WriteString(x.convertToWazuhString() + "\n")
		if writeErr != nil {
			fmt.Printf("Error writing to file: %s", writeErr.Error())
		}
	}
	for _, x := range e.LDAP {
		if x.Timestamp.After(lastEventSeen) {
			lastEventSeen = x.Timestamp
		}
		_, writeErr := f.WriteString(x.convertToWazuhString() + "\n")
		if writeErr != nil {
			fmt.Printf("Error writing to file: %s", writeErr.Error())
		}
	}
	for _, x := range e.Systems {
		if x.Timestamp.After(lastEventSeen) {
			lastEventSeen = x.Timestamp
		}
		_, writeErr := f.WriteString(x.convertToWazuhString() + "\n")
		if writeErr != nil {
			fmt.Printf("Error writing to file: %s", writeErr.Error())
		}
	}
	for _, x := range e.SSO {
		if x.Timestamp.After(lastEventSeen) {
			lastEventSeen = x.Timestamp
		}
		_, writeErr := f.WriteString(x.convertToWazuhString() + "\n")
		if writeErr != nil {
			fmt.Printf("Error writing to file: %s", writeErr.Error())
		}

	}
	err = timeTracker.UpdateLast(lastEventSeen.Add(time.Second * 1))
	return err
}
