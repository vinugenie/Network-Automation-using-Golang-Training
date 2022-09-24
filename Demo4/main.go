package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

const tmpl_config = `
{{- range $rtr := .Routers }}
config t
hostname {{ $rtr.Hostname }}
!
interface Loopback100
ip address {{ $rtr.Loopback }} 255.255.255.255
!
{{- end }}
end
`

type Router struct {
	Hostname string `yaml:"hostname"`
	IP       string `yaml:"ip"`
	Loopback string `yaml:"lo100"`
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
}

type Inventory struct {
	Routers []Router `json:"routers"`
}

func main() {
	file, err := os.Open("data.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	src := yaml.NewDecoder(file)

	var inv Inventory

	err = src.Decode(&inv)
	if err != nil {
		panic(err)
	}

	cfg, err := genConfig(inv)
	if err != nil {
		panic(err)
	}

	// fmt.Println(cfg)
	var ctr int = 0
	for ctr < len(inv.Routers) {
		pushConfigs(inv.Routers[ctr].IP, inv.Routers[ctr].User, inv.Routers[ctr].Pwd, cfg)
		ctr += 1
	}

}

func genConfig(inv Inventory) (b bytes.Buffer, err error) {
	t, err := template.New("config").Parse(tmpl_config)
	if err != nil {
		return b, fmt.Errorf("Failed to create Template: %w", err)
	}

	err = t.Execute(&b, inv)
	if err != nil {
		return b, fmt.Errorf("Failed to create Template: %w", err)
	}
	return b, nil
}

func pushConfigs(IP string, user string, pwd string, cfg bytes.Buffer) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial(
		"tcp", fmt.Sprintf("%s:%d", IP, 22),
		config,
	)

	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	stdin, err := session.StdinPipe()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Shell()
	log.Print("Connected to the device.. Configuring it now..")
	cfg.WriteTo(stdin)
	session.Close()
}
