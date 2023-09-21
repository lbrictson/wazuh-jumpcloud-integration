package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Valid JumpCloud service types are:
// all: Logs from all services.
// directory: Logs activity in the Admin Portal and User Portal, including admin changes in the directory and admin/user authentications to the Admin Portal and User Portal.
// ldap: Logs user authentications to LDAP, including LDAP Bind and Search event types.
// mdm: Logs MDM command results.
// password_manager: Logs activity related to JumpCloud password manager.
// radius: Logs user authentications to RADIUS, used for Wi-Fi and VPNs.
// software: Logs application activity when software is added, removed, or changed on a macOS, Windows, or Linux device. Events are logged based on changes to an application version during each device check-in.
// sso: Logs user authentications to SAML applications.
// systems: Logs user authentications to MacOS, Windows, and Linux systems, including agent-related events on lockout, password changes, and File Disk Encryption key updates.

// JumpCloudAPI can be used to interact with the JumpCloud API
type JumpCloudAPI struct {
	apiKey  string
	baseURL string
	orgID   string
}

// NewJumpCloudAPIOptions are the options for creating a new JumpCloudAPI object
type NewJumpCloudAPIOptions struct {
	APIKey  string
	BaseURL string
	OrgID   string
}

// NewJumpCloudAPI returns a new JumpCloudAPI object, if you do not provide a base URL, it will default to the JumpCloud API
func NewJumpCloudAPI(options NewJumpCloudAPIOptions) *JumpCloudAPI {
	a := JumpCloudAPI{
		apiKey:  options.APIKey,
		baseURL: options.BaseURL,
		orgID:   options.OrgID,
	}
	if options.BaseURL == "" {
		a.baseURL = "https://api.jumpcloud.com"
	}
	return &a
}

// GetEventsSinceTime returns all JumpCloud events since the given time
func (a *JumpCloudAPI) GetEventsSinceTime(startTime time.Time) (*JumpCloudEvents, error) {
	url := a.baseURL + "/insights/directory/v1/events"
	method := "POST"
	// JumpCloud API requires a time in RFC3339 format
	starterTime := startTime.Format(time.RFC3339)
	payload := strings.NewReader(fmt.Sprintf(`{"service": ["all"], "start_time": "%v", "limit": 10000}`, starterTime))
	// Default Go HTTP client, might need to customize this later
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("x-api-key", a.apiKey)
	req.Header.Add("Content-Type", "application/json")
	if a.orgID != "" {
		req.Header.Add("x-org-id", a.orgID)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v | %v | %v", res.Status, res.StatusCode, err)
	}
	// JumpCloud API returns a 200 even if there are no events
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error response from JumpCloud: %v | %v | %v", res.Status, res.StatusCode, string(body))
	}
	events, err := decodeJumpCloudEvents(body)
	if err != nil {
		return nil, fmt.Errorf("error decoding JumpCloud response: %v", err)
	}
	return &events, nil
}

type JumpCloudEvents struct {
	LDAP      []JumpCloudLDAPEvent      `json:"ldap_events"`
	Systems   []JumpCloudSystemEvent    `json:"systems"`
	Directory []JumpCloudDirectoryEvent `json:"directory"`
	Radius    []JumpCloudRadiusEvent    `json:"radius"`
	SSO       []JumpCloudSSOEvent       `json:"sso"`
	Admin     []JumpCloudAdminEvent     `json:"admin"`
}

type BaseJumpCloudEvent struct {
	Service string `json:"service"`
}

// decodeJumpCloudEvents decodes the raw JumpCloud API response into a JumpCloudEvents object that contains events
// of the varying types
func decodeJumpCloudEvents(raw []byte) (JumpCloudEvents, error) {
	finished := JumpCloudEvents{}
	generic := []map[string]interface{}{}
	err := json.Unmarshal(raw, &generic)
	if err != nil {
		return JumpCloudEvents{}, err
	}
	var events []BaseJumpCloudEvent
	err = json.Unmarshal(raw, &events)
	for i, x := range events {
		fmt.Println(x.Service)
		switch x.Service {
		case "ldap":
			b, err := json.Marshal(generic[i])
			if err != nil {
				return JumpCloudEvents{}, err
			}
			var e JumpCloudLDAPEvent
			err = json.Unmarshal(b, &e)
			if err != nil {
				return JumpCloudEvents{}, err
			}
			finished.LDAP = append(finished.LDAP, e)
		case "systems":
			b, err := json.Marshal(generic[i])
			if err != nil {
				return JumpCloudEvents{}, err
			}
			var e JumpCloudSystemEvent
			err = json.Unmarshal(b, &e)
			if err != nil {
				return JumpCloudEvents{}, err
			}
			finished.Systems = append(finished.Systems, e)
		case "directory":
			b, err := json.Marshal(generic[i])
			if err != nil {
				return JumpCloudEvents{}, err
			}
			var e JumpCloudDirectoryEvent
			err = json.Unmarshal(b, &e)
			if err != nil {
				return JumpCloudEvents{}, err
			}
			finished.Directory = append(finished.Directory, e)
		case "radius":
			b, err := json.Marshal(generic[i])
			if err != nil {
				return JumpCloudEvents{}, err
			}
			var e JumpCloudRadiusEvent
			err = json.Unmarshal(b, &e)
			if err != nil {
				return JumpCloudEvents{}, err
			}
			finished.Radius = append(finished.Radius, e)
		case "sso":
			b, err := json.Marshal(generic[i])
			if err != nil {
				return JumpCloudEvents{}, err
			}
			var e JumpCloudSSOEvent
			err = json.Unmarshal(b, &e)
			if err != nil {
				return JumpCloudEvents{}, err
			}
			finished.SSO = append(finished.SSO, e)
		case "admin":
			b, err := json.Marshal(generic[i])
			if err != nil {
				return JumpCloudEvents{}, err
			}
			var e JumpCloudAdminEvent
			err = json.Unmarshal(b, &e)
			if err != nil {
				return JumpCloudEvents{}, err
			}
			finished.Admin = append(finished.Admin, e)

		}
	}
	return finished, err
}
