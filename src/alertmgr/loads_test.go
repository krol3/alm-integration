package alertmgr

import (
	"dbservice"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestLoads(t *testing.T) {
	cfgData := `
---
- type: common
  Max_DB_Size: 10
  Delete_Old_Data: 10
  AquaServer: https://demolab.aquasec.com
- name: jira
  type: jira
  enable: true
  url: "http://localhost:2990/jira"
  user: admin
  password: admin
  tls_verify: false
  project_key: key
  description:
  summary:
  issuetype: "Bug"
  priority: Medium
  assignee: 
  Policy-Min-Vulnerability: Critical
  labels: ["label1", "label2"]
  Policy-Min-Vulnerability: high

- name: jiraWithoutPass
  type: jira
  enable: true
  url: "http://localhost:2990/jira"
  user: admin

- name: my-slack
  type: slack
  enable: true
  url: "https://hooks.slack.com/services/TT/BBB/WWWW"

- name: email
  type: email
  enable: true
  user: EMAILUSER
  password: EMAILPASS
  host: smtp.gmail.com
  port: 587
  recipients: ["demo@gmail.com"]

- name: localEmail
  type: email
  enable: true
  useMX: true
  sender: mail@alm.demo.co
  recipients: ["demo@gmail.com"]

- name: email-empty
  type: email
  enable: true

- name: email-empty-pass
  type: email
  enable: true
  user: EMAILUSER

- name: ms-team
  type: teams
  enable: true
  url: https://outlook.office.com/webhook/.... # Webhook's url

- name: faild
  enable: true
  type: nextplugin

- name: my-servicenow
  type: serviceNow
  enable: true
  user: SERVICENOWUSER
  password: SERVICENOWPASS
  instance: dev00000

- name: noname
  type: future-plugin
  enable: true
  user: user
  password: password

- name: webhook
  type: webhook
  enable: true
  url: https://postman-echo.com/post

- name: splunk
  type: splunk
  enable: true
  url: http://localhost:8088
  token: splunk-demo-token
`
	cfgName :="cfg_test.yaml"
	ioutil.WriteFile(cfgName, []byte(cfgData),0644)
	dbPathReal := dbservice.DbPath
	savedBaseForTicker := baseForTicker
	defer func() {
		baseForTicker = savedBaseForTicker
		os.Remove(cfgName)
		os.Remove(dbservice.DbPath)
		dbservice.ChangeDbPath(dbPathReal)
	}()
	dbservice.DbPath = "test_webhooks.db"
	baseForTicker = time.Microsecond

	demoCtx := Instance()
	demoCtx.Start(cfgName)
	pluginsNumber := 10
	if len(demoCtx.plugins) != pluginsNumber {
		t.Errorf("There are stopped plugins\nWaited: %d\nResult: %d", pluginsNumber, len(demoCtx.plugins))
	}

	_, ok := demoCtx.plugins["ms-team"]
	if !ok {
		t.Errorf("'ms-team' plugin didn't start!")
	}

	aquaWaiting := "https://demolab.aquasec.com/#/images/"
	if aquaServer != aquaWaiting {
		t.Errorf("Wrong init of AquaServer link.\nWait: %q\nGot: %q", aquaWaiting, aquaServer)
	}

	if _, ok := demoCtx.plugins["my-servicenow"]; !ok {
		t.Errorf("Plugin 'my-servicenow' didn't run!")
	}
	demoCtx.Terminate()
	time.Sleep(200*time.Millisecond)
}
