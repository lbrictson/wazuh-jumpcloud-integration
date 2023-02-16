# wazuh-jumpcloud-integration
A pipeline for ingesting [JumpCloud](https://jumpcloud.com/) events into [Wazuh](https://wazuh.com/)

![image](https://user-images.githubusercontent.com/8505034/219263945-23901d63-f974-4feb-8696-d759a86156a3.png)

## Overview

This integration is designed to be placed on a Wazuh Manager system to interact with the JumpCloud API in order ot pull events

Events pulled from JumpCloud

- System Events
- Directory Events
- SSO Events
- LDAP Events
- Radius Events

Events that do not match any rule are set to level 0 and therefore ignored by Wazuh.

Rules are found in `rules/jumpcloud_rules.xml`

## Requirements

- A valid [JumpCloud API Key](https://support.jumpcloud.com/support/s/article/jumpcloud-apis1)
- An installation of [Wazuh](https://wazuh.com/)
- SSH Access to your Wazuh Manager server

## Installation

Note:  Paths are examples, you can use any path you like

```bash
# Create directories
mkdir -p /opt/jumpcloud
# Download the latest release
wget https://github.com/lbrictson/wazuh-jumpcloud-integration/releases
# Setup the config file
wget https://raw.githubusercontent.com/lbrictson/wazuh-jumpcloud-integration/master/config/config.json -O /opt/jumpcloud/config.json
# Place your JumpCloud API Key in the config file
sed -i 's/this-is-not-a-real-key/YOUR-JUMPCLOUD-API-KEY-HERE/g' /opt/jumpcloud/config.json
# Setup permissions
chmod +x /opt/jumpcloud/wazuh-jumpcloud-integration
chown -R root:wazuh /opt/jumpcloud
```

Once all the components are in place it is time to modify the Wazuh configuration

Always backup your configuration before making changes

```bash
# Edit the ossec.conf file to add the JumpCloud integration
vim /var/ossec/etc/ossec.conf
```
Add the following to the `ossec.conf` file.  Change any paths if you customized the installation location.  Optionally change the interval, 5m will keep you under the JumpCloud API Rate limits
```xml
<wodle name="command">
  <disabled>no</disabled>
  <tag>jumpcloud</tag>
  <command>/bin/bash -c "/opt/jumpcloud/wazuh-jumpcloud-integration /opt/jumpcloud/config.json /opt/jumpcloud/output.log"</command>
  <interval>5m</interval>
  <ignore_output>yes</ignore_output>
  <run_on_start>yes</run_on_start>
</wodle>
```

Add a block to the `ossec.conf` file to configure the JumpCloud log file.  This instructs Wazuh to ingest the logs emitted by the integration
```xml
<localfile>
    <log_format>json</log_format>
    <location>/opt/jumpcloud/output.log</location>
</localfile>
```

Lastly add the ruleset
```bash
wget https://raw.githubusercontent.com/lbrictson/wazuh-jumpcloud-integration/master/rules/jumpcloud_rules.xml -O /var/ossec/etc/rules/jumpcloud_rules.xml
chown wazuh:wazuh /var/ossec/etc/rules/jumpcloud_rules.xml
```

Restart the Wazuh Manager
```bash
systemctl restart wazuh-manager
```

Monitor the logs to see if the integration is working
```bash
tail -f /var/ossec/logs/ossec.log
```

## Troubleshooting

If you are having issues with the integration you can run it manually to see what is happening
```bash
/opt/jumpcloud/wazuh-jumpcloud-integration /opt/jumpcloud/config.json /opt/jumpcloud/output.log
```

## How it Works

The integration program relies on the config.json file to locate the JumpCloud API key, additionally this file is automatically updated with the last successful time the integration was run.

Each time the integration runs it checks the config file, reads the last time and only gathers events since that time.

Events are emitted as JSON into the designated output file.  Wazuh will then read the output file and ingest the events.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Happy to accept requests to update and modify the rules to match more events
