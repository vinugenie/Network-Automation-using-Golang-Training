package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("*.html"))
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/restconf", restConf)
	http.ListenAndServe("192.168.0.75:8080", nil)
}

func restConf(w http.ResponseWriter, r *http.Request) {
	target := r.FormValue("dns")

	url := "https://" + target + ":443/restconf/data/ietf-interfaces:interfaces"

	ignoreCert := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	content := []byte(`{
		"ietf-interfaces:interface": {
			"name": "Loopback101",
			"description": "Configured using Golang",
			"type": "iana-if-type:softwareLoopback",
			"enabled": true,
			"ietf-ip:ipv4": {
				"address": {
					"ip": "101.101.101.101",
					"netmask": "255.255.255.255"
				}
			}
		}
	}`)

	client := &http.Client{Transport: ignoreCert}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(content))
	if err != nil {
		fmt.Println(err)
	}

	req.SetBasicAuth("admin", "admin")
	req.Header.Set("Access-Control-Allow-Origin", "*")
	req.Header.Add("Content-Type", "application/yang-data+json")
	req.Header.Add("Accept", "application/yang-data+json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Fprint(w, string(body))
}
