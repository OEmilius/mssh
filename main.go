package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	//	"sync"

	"golang.org/x/crypto/ssh"
)

var cfg Config

type Config struct {
	Hosts []string
}

const (
	CERT_PASSWORD        = 1
	CERT_PUBLIC_KEY_FILE = 2
	DEFAULT_TIMEOUT      = 3 // second
)

type SSH struct {
	Ip      string
	User    string
	Cert    string //password or key file path
	Port    int
	session *ssh.Session
	client  *ssh.Client
}

func (ssh_client *SSH) readPublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func (ssh_client *SSH) Connect(mode int) error {
	var ssh_config *ssh.ClientConfig
	var auth []ssh.AuthMethod
	if mode == CERT_PASSWORD {
		auth = []ssh.AuthMethod{ssh.Password(ssh_client.Cert)}
	} else if mode == CERT_PUBLIC_KEY_FILE {
		auth = []ssh.AuthMethod{ssh_client.readPublicKeyFile(ssh_client.Cert)}
	} else {
		log.Println("does not support mode: ", mode)
		//var err error
		return fmt.Errorf("does not support mode:", mode)
	}

	ssh_config = &ssh.ClientConfig{
		User: ssh_client.User,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * DEFAULT_TIMEOUT,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ssh_client.Ip, ssh_client.Port), ssh_config)
	if err != nil {
		//fmt.Println(err)
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		//fmt.Println(err)
		client.Close()
		return err
	}

	ssh_client.session = session
	ssh_client.client = client
	return nil
}

func (ssh_client *SSH) RunCmd(cmd string) {
	out, err := ssh_client.session.CombinedOutput(cmd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func (ssh_client *SSH) RunCmdGetOut(cmd string) (result string) {
	out, err := ssh_client.session.CombinedOutput(cmd)
	if err != nil {
		result = fmt.Sprintln(err)
		return result
	}
	result = fmt.Sprintln(string(out))
	return result
}

func (ssh_client *SSH) Close() {
	ssh_client.session.Close()
	ssh_client.client.Close()
}

func (client *SSH) ConnectRunCmdGetOut(mode int, cmd string) (result string, err error) {
	err = client.Connect(CERT_PASSWORD)
	if err != nil {
		return fmt.Sprintln(err), err
	} else {
		result = client.RunCmdGetOut(cmd)
		client.Close()
		return result, nil
	}
}

func main() {
	cfg = ReadConfig("config.json")
	//fmt.Println(cfg)
	for _, host := range cfg.Hosts {
		client := &SSH{
			User: login,
			Port: 22,
			Cert: password,
		}
		client.Ip = host
		go func() {
			result, _ := client.ConnectRunCmdGetOut(CERT_PASSWORD, cmd)
			fmt.Println(client.Ip, result)
		}()
	}
	fmt.Scanln()
}
