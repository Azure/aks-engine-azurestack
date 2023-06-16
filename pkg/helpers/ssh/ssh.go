// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package ssh

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

// ExecuteRemote executes a script in a remote host.
//
// Context ctx is only enforced during the process that stablishes
// the SSH connection and creates the SSH client.
func ExecuteRemote(ctx context.Context, host *RemoteHost, script string) (combinedOutput string, err error) {
	c, err := clientWithRetry(ctx, host)
	if err != nil {
		return "", errors.Wrap(err, "creating SSH client")
	}
	defer c.Close()
	s, err := c.NewSession()
	if err != nil {
		return "", errors.Wrap(err, "creating SSH session")
	}
	defer s.Close()
	if co, err := s.CombinedOutput(script); err != nil {
		return string(co), errors.Wrapf(err, "executing script")
	}
	return "", nil
}

// PublicKeyAuth returns an AuthMethod that uses a ssh key pair
func PublicKeyAuth(sshPrivateKeyPath string) (ssh.AuthMethod, error) {
	b, err := os.ReadFile(sshPrivateKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "reading ssh private key file")
	}
	k, err := ssh.ParsePrivateKey(b)
	if err != nil {
		return nil, errors.Wrap(err, "parsing ssh private key content")
	}
	return ssh.PublicKeys(k), nil
}

// ValidateConfig checks the JumpBox configuration
func ValidateConfig(host *JumpBox) error {
	jbConfig, err := config(host.AuthConfig)
	if err != nil {
		return errors.Wrap(err, "creating ssh client config")
	}
	_, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.URI, host.Port), jbConfig)
	if err != nil {
		return errors.Wrapf(err, "dialing ssh (%s)", host.URI)
	}
	return nil
}

func clientWithRetry(ctx context.Context, host *RemoteHost) (*ssh.Client, error) {
	// TODO Granular retry func
	retryFunc := func(err error) bool {
		select {
		case <-ctx.Done():
			return false
		default:
			return true
		}
	}
	backoff := wait.Backoff{Steps: 300, Duration: 10 * time.Second}
	var c *ssh.Client
	var err error
	err = retry.OnError(backoff, retryFunc, func() error {
		c, err = client(host)
		return err
	})
	return c, err
}

func client(host *RemoteHost) (*ssh.Client, error) {
	jbConfig, err := config(host.Jumpbox.AuthConfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating jumpbox client config")
	}
	jbConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.Jumpbox.URI, host.Jumpbox.Port), jbConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "dialing jumpbox (%s)", host.Jumpbox.URI)
	}
	hostConn, err := jbConn.Dial("tcp", fmt.Sprintf("%s:%d", host.URI, host.Port))
	if err != nil {
		return nil, errors.Wrapf(err, "dialing host (%s)", host.URI)
	}
	hostConfig, err := config(host.AuthConfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating host client config")
	}
	ncc, chans, reqs, err := ssh.NewClientConn(hostConn, host.URI, hostConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "starting new client connection to host (%s)", host.URI)
	}
	c, err := ssh.NewClient(ncc, chans, reqs), nil
	if err != nil {
		return nil, errors.Wrapf(err, "creating new ssh client for host (%s)", host.URI)
	}
	return c, nil
}

func config(authConfig *AuthConfig) (*ssh.ClientConfig, error) {
	authMethod, err := clientConfigAuth(authConfig)
	if err != nil {
		return nil, err
	}
	hkCallback, err := knownHostsHostKeyCallback()
	if err != nil {
		return nil, err
	}
	hostKeyCallback := func(host string, remote net.Addr, pubKey ssh.PublicKey) error {
		var keyErr *kh.KeyError
		if cbErr := hkCallback(host, remote, pubKey); cbErr != nil && errors.As(cbErr, &keyErr) {
			hostname := strings.Split(host, ":")[0]
			if len(keyErr.Want) > 0 {
				log.Errorf("Strict host key check failed. Remote host identification has changed.")
				log.Errorf("Key '%v' does not match a key known for host %s.", hostKeyString(pubKey), hostname)
				return keyErr
			}
			if len(keyErr.Want) == 0 {
				if err := addHostKey(hostname, pubKey); err != nil {
					return err
				}
				log.Warnf("Permanently added '%s' (%s) to the list of known hosts (%s)", hostname, pubKey.Type(), khpath)
			}
		}
		return nil
	}
	return &ssh.ClientConfig{
		HostKeyCallback: hostKeyCallback,
		User:            authConfig.User,
		Auth:            authMethod,
	}, nil
}

// clientConfigAuth returns the ssh authentication method
func clientConfigAuth(authConfig *AuthConfig) ([]ssh.AuthMethod, error) {
	// TODO we may be able to return both methods, not changing behavior for now
	// TODO this can be reworked so it is executed once per command
	if authConfig.PrivateKeyPath != "" {
		keyAuth, err := PublicKeyAuth(authConfig.PrivateKeyPath)
		if err != nil {
			return nil, errors.Wrap(err, "creating public key authentication method")
		}
		return []ssh.AuthMethod{keyAuth}, nil
	}
	return []ssh.AuthMethod{ssh.Password(authConfig.Password)}, nil
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
