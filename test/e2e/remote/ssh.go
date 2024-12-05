//go:build test
// +build test

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package remote

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/aks-engine-azurestack/test/e2e/kubernetes/util"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

const (
	sshRetries = 20
	scriptsDir = "scripts"
)

var (
	khpath    string
	lineBreak string
)

func init() {
	khpath = filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	lineBreak = "\n"
}

// Connection is
type Connection struct {
	Host           string
	Port           string
	User           string
	PrivateKeyPath string
	ClientConfig   *ssh.ClientConfig
	Client         *ssh.Client
}

// NewConnectionResult is the result type for GetConfigAsync
type NewConnectionResult struct {
	Connection *Connection
	Err        error
}

// NewConnectionAsync wraps NewConnection with a struct response for goroutine + channel usage
func NewConnectionAsync(host, port, user, keyPath string, keyPathPub string) NewConnectionResult {
	connection, err := NewConnection(host, port, user, keyPath, keyPathPub)
	if connection == nil {
		connection = &Connection{}
	}
	return NewConnectionResult{
		Connection: connection,
		Err:        err,
	}
}

// NewConnection will build and return a new Connection object
func NewConnection(host, port, user, keyPath string, keyPathPub string) (*Connection, error) {
	conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		log.Printf("unable to establish net connection $SSH_AUTH_SOCK has value %s\n", os.Getenv("SSH_AUTH_SOCK"))
		return nil, err
	}
	defer conn.Close()
	ag := agent.NewClient(conn)

	privateKeyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	addKey := agent.AddedKey{
		PrivateKey: privateKey,
	}

	ag.Add(addKey)
	signers, err := ag.Signers()
	if err != nil {
		log.Println("unable to add key to agent")
		return nil, err
	}
	auths := []ssh.AuthMethod{ssh.PublicKeys(signers...)}

	hkCallback, err := knownHostsHostKeyCallback()
	if err != nil {
		return nil, err
	}
	hostKeyCallback := func(host string, remote net.Addr, pubKey ssh.PublicKey) error {
		var keyErr *kh.KeyError
		if cbErr := hkCallback(host, remote, pubKey); cbErr != nil && errors.As(cbErr, &keyErr) {
			hostname := strings.Split(host, ":")[0]
			if len(keyErr.Want) > 0 {
				log.Println("Strict host key check failed. Remote host identification has changed.")
				log.Println("Key '%v' does not match a key known for host %s.", hostKeyString(pubKey), hostname)
				return keyErr
			}
			if len(keyErr.Want) == 0 {
				if err := addHostKey(hostname, pubKey); err != nil {
					return err
				}
				log.Println("Permanently added '%s' (%s) to the list of known hosts (%s)", hostname, pubKey.Type(), khpath)
			}
		}
		return nil
	}

	cfg := &ssh.ClientConfig{
		User:            user,
		Auth:            auths,
		HostKeyCallback: hostKeyCallback,
	}

	cnctStr := fmt.Sprintf("%s:%s", host, port)
	sshClient, err := ssh.Dial("tcp", cnctStr, cfg)
	if err != nil {
		log.Println("unable to dial ssh connection, err: %s\n", err)
		return nil, err
	}

	return &Connection{
		Host:           host,
		Port:           port,
		User:           user,
		PrivateKeyPath: keyPath,
		ClientConfig:   cfg,
		Client:         sshClient,
	}, nil
}

// NewConnectionWithRetry establishes an ssh connection, allowing for retries
func NewConnectionWithRetry(host, port, user, keyPath string, keyPathPub string, sleep, timeout time.Duration) (*Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ch := make(chan NewConnectionResult)
	var mostRecentNewConnectionWithRetryError error
	var connection *Connection
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- NewConnectionAsync(host, port, user, keyPath, keyPathPub)
				time.Sleep(sleep)
			}
		}
	}()
	for {
		select {
		case result := <-ch:
			mostRecentNewConnectionWithRetryError = result.Err
			connection = result.Connection
			if mostRecentNewConnectionWithRetryError == nil {
				return connection, nil
			}
		case <-ctx.Done():
			return nil, errors.Errorf("NewConnectionWithRetry timed out: %s\n", mostRecentNewConnectionWithRetryError)
		}
	}
}

// knownHostsHostKeyCallback returns a host key callback that uses file
// ${HOME}/.ssh/known_hosts to store known host keys
func knownHostsHostKeyCallback() (ssh.HostKeyCallback, error) {
	err := ensuresKnownHosts()
	if err != nil {
		return nil, err
	}
	khCallback, err := kh.New(khpath)
	if err != nil {
		return nil, errors.Wrap(err, "creating HostKeyCallback instance")
	}
	return khCallback, nil
}

// ensuresKnownHosts creates file ${HOME}/.ssh/known_hosts if it does not exist
func ensuresKnownHosts() error {
	if err := os.MkdirAll(path.Dir(khpath), 0700); err != nil {
		return errors.Wrap(err, "creating .ssh directory")
	}
	f, err := os.OpenFile(khpath, os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "creating known_hosts file")
	}
	f.Close()
	return nil
}

// hostKeyString pretty-prints a public key struct
func hostKeyString(k ssh.PublicKey) string {
	return fmt.Sprintf("%s %s", k.Type(), base64.StdEncoding.EncodeToString(k.Marshal()))
}

// addHostKey adds an entry to ${HOME}/.ssh/known_hosts
func addHostKey(hostname string, pubKey ssh.PublicKey) error {
	f, err := os.OpenFile(khpath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "opening known_hosts file")
	}
	defer f.Close()
	// append blank line
	if _, err = f.WriteString(lineBreak); err != nil {
		return errors.Wrap(err, "appending blank line to known_hosts file")
	}
	// append host key line
	knownHosts := kh.Normalize(hostname)
	_, err = f.WriteString(kh.Line([]string{knownHosts}, pubKey))
	if err != nil {
		return errors.Wrap(err, "writing known_hosts file")
	}
	return nil
}

// Execute will execute a given cmd on a remote host
func (c *Connection) Execute(cmd string, printStdout bool) error {
	session, err := c.Client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fmt.Printf("\n$ %s\n", cmd)
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Printf("Error output:%s\n", out)
		return err
	}
	if printStdout {
		log.Printf("%s\n", out)
	}
	return nil
}

func (c *Connection) Write(data, path string) error {
	remoteCommand := fmt.Sprintf("echo %s > %s", data, path)
	connectString := fmt.Sprintf("%s@%s", c.User, c.Host)
	cmd := exec.Command("ssh", "-i", c.PrivateKeyPath, "-o", "ConnectTimeout=30", "-o", "StrictHostKeyChecking=no", connectString, "-p", c.Port, remoteCommand)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error output:%s\n", out)
		return err
	}
	return nil
}

func (c *Connection) Read(path string) ([]byte, error) {
	remoteCommand := fmt.Sprintf("cat %s", path)
	connectString := fmt.Sprintf("%s@%s", c.User, c.Host)
	cmd := exec.Command("ssh", "-i", c.PrivateKeyPath, "-o", "ConnectTimeout=30", "-o", "StrictHostKeyChecking=no", connectString, "-p", c.Port, remoteCommand)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error output:%s\n", out)
		return nil, err
	}
	return out, nil
}

// CopyTo uses this ssh connection to send files to the remote ssh listener's underlying file system
func (c *Connection) CopyTo(filename string) error {
	var scpError error
	var scpOut []byte
	for i := 0; i < sshRetries; i++ {
		connectString := fmt.Sprintf("%s@%s", c.User, c.Host)
		cmd := exec.Command("scp", "-i", c.PrivateKeyPath, "-P", c.Port, "-o", "StrictHostKeyChecking=no", filepath.Join(scriptsDir, filename), connectString+":/tmp/"+filename)
		util.PrintCommand(cmd)
		scpOut, scpError = cmd.CombinedOutput()
		if scpError != nil {
			log.Printf("Error output:%s\n", scpOut)
			continue
		} else {
			break
		}
	}
	return scpError
}

// CopyFrom uses this ssh connection to get remote files via scp
func (c *Connection) CopyFrom(hostname, path string) error {
	remoteCommand := fmt.Sprintf("scp -o StrictHostKeyChecking=no %s:%s /tmp/%s-%s", hostname, path, hostname, filepath.Base(path))
	connectString := fmt.Sprintf("%s@%s", c.User, c.Host)
	cmd := exec.Command("ssh", "-A", "-i", c.PrivateKeyPath, "-o", "ConnectTimeout=30", "-o", "StrictHostKeyChecking=no", connectString, "-p", c.Port, remoteCommand)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error output:%s\n", out)
		return err
	}
	return nil
}

// CopyToRemote uses this ssh connection to send files via scp to a remote host
func (c *Connection) CopyToRemote(hostname, path string) error {
	remoteCommand := fmt.Sprintf("scp -o StrictHostKeyChecking=no %s %s:%s", path, hostname, path)
	connectString := fmt.Sprintf("%s@%s", c.User, c.Host)
	cmd := exec.Command("ssh", "-A", "-i", c.PrivateKeyPath, "-o", "ConnectTimeout=30", "-o", "StrictHostKeyChecking=no", connectString, "-p", c.Port, remoteCommand)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error output:%s\n", out)
	}
	return err
}

// CopyToRemoteWithRetry sends files via scp to a remote host with retry tolerance
func (c *Connection) CopyToRemoteWithRetry(hostname, path string, sleep, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ch := make(chan error)
	var mostRecentCopyToRemoteWithRetry error
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- c.CopyToRemote(hostname, path)
				time.Sleep(sleep)
			}
		}
	}()
	for {
		select {
		case result := <-ch:
			mostRecentCopyToRemoteWithRetry = result
			if mostRecentCopyToRemoteWithRetry == nil {
				return nil
			}
		case <-ctx.Done():
			return errors.Errorf("CopyToRemoteWithRetry timed out: %s\n", mostRecentCopyToRemoteWithRetry)
		}
	}
}

// ExecuteRemote uses this ssh connection to run a remote command from the primary master node
func (c *Connection) ExecuteRemote(node, command string, printStdout bool) error {
	cmd := exec.Command("ssh", "-A", "-i", c.PrivateKeyPath, "-p", c.Port, "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-o", "LogLevel=ERROR", fmt.Sprintf("%s@%s", c.User, c.Host), "ssh", "-o", "ConnectTimeout=10", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-o", "LogLevel=ERROR", node, command)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error output:%s\n", out)
	} else {
		if printStdout {
			log.Printf("%s\n", out)
		}
	}
	return err
}

// ExecuteRemoteWithRetry runs a remote command with retry tolerance
func (c *Connection) ExecuteRemoteWithRetry(node, command string, printStdout bool, sleep, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ch := make(chan error)
	var mostRecentExecuteRemoteWithRetryError error
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- c.ExecuteRemote(node, command, printStdout)
				time.Sleep(sleep)
			}
		}
	}()
	for {
		select {
		case result := <-ch:
			mostRecentExecuteRemoteWithRetryError = result
			if mostRecentExecuteRemoteWithRetryError == nil {
				return nil
			}
		case <-ctx.Done():
			return errors.Errorf("ExecuteRemoteWithRetry timed out: %s\n", mostRecentExecuteRemoteWithRetryError)
		}
	}
}
