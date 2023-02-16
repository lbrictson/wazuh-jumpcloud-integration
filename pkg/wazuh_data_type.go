package pkg

// Simple version to text JSON strings for Wazuh to ingest, might need to customize these later

import "encoding/json"

func (d *JumpCloudSystemEvent) convertToWazuhString() string {
	d.JumpCloudEventType = "system"
	b, _ := json.Marshal(d)
	return string(b)
}

func (d *JumpCloudLDAPEvent) convertToWazuhString() string {
	d.JumpCloudEventType = "ldap"
	b, _ := json.Marshal(d)
	return string(b)
}

func (d *JumpCloudDirectoryEvent) convertToWazuhString() string {
	d.JumpCloudEventType = "directory"
	b, _ := json.Marshal(d)
	return string(b)
}

func (d *JumpCloudRadiusEvent) convertToWazuhString() string {
	d.JumpCloudEventType = "radius"
	b, _ := json.Marshal(d)
	return string(b)
}

func (d *JumpCloudSSOEvent) convertToWazuhString() string {
	d.JumpCloudEventType = "sso"
	b, _ := json.Marshal(d)
	return string(b)
}
