package tls

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"sync"

	"github.com/sagernet/fswatch"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/ntp"
)

var errInsecureUnused = E.New("tls: insecure unused")

type STDServerConfig struct {
	access          sync.RWMutex
	config          *tls.Config
	logger          log.Logger
	certificate     []byte
	key             []byte
	certificatePath string
	keyPath         string
	echKeyPath      string
	watcher         *fswatch.Watcher
}

func (c *STDServerConfig) ServerName() string {
	c.access.RLock()
	defer c.access.RUnlock()
	return c.config.ServerName
}

func (c *STDServerConfig) SetServerName(serverName string) {
	c.access.Lock()
	defer c.access.Unlock()
	config := c.config.Clone()
	config.ServerName = serverName
	c.config = config
}

func (c *STDServerConfig) NextProtos() []string {
	c.access.RLock()
	defer c.access.RUnlock()
	return c.config.NextProtos
}

func (c *STDServerConfig) SetNextProtos(nextProto []string) {
	c.access.Lock()
	defer c.access.Unlock()
	config := c.config.Clone()
	config.NextProtos = nextProto
	c.config = config
}

func (c *STDServerConfig) Config() (*STDConfig, error) {
	return c.config, nil
}

func (c *STDServerConfig) Client(conn net.Conn) (Conn, error) {
	return tls.Client(conn, c.config), nil
}

func (c *STDServerConfig) Server(conn net.Conn) (Conn, error) {
	return tls.Server(conn, c.config), nil
}

func (c *STDServerConfig) Clone() Config {
	return &STDServerConfig{
		config: c.config.Clone(),
	}
}

func (c *STDServerConfig) Start() error {
	err := c.startWatcher()
	if err != nil {
		c.logger.Warn("create fsnotify watcher: ", err)
	}
	return nil
}

func (c *STDServerConfig) startWatcher() error {
	var watchPath []string
	if c.certificatePath != "" {
		watchPath = append(watchPath, c.certificatePath)
	}
	if c.keyPath != "" {
		watchPath = append(watchPath, c.keyPath)
	}
	if c.echKeyPath != "" {
		watchPath = append(watchPath, c.echKeyPath)
	}
	if len(watchPath) == 0 {
		return nil
	}
	watcher, err := fswatch.NewWatcher(fswatch.Options{
		Path: watchPath,
		Callback: func(path string) {
			err := c.certificateUpdated(path)
			if err != nil {
				c.logger.Error(E.Cause(err, "reload certificate"))
			}
		},
	})
	if err != nil {
		return err
	}
	err = watcher.Start()
	if err != nil {
		return err
	}
	c.watcher = watcher
	return nil
}

func (c *STDServerConfig) certificateUpdated(path string) error {
	if path == c.certificatePath || path == c.keyPath {
		if path == c.certificatePath {
			certificate, err := os.ReadFile(c.certificatePath)
			if err != nil {
				return E.Cause(err, "reload certificate from ", c.certificatePath)
			}
			c.certificate = certificate
		} else if path == c.keyPath {
			key, err := os.ReadFile(c.keyPath)
			if err != nil {
				return E.Cause(err, "reload key from ", c.keyPath)
			}
			c.key = key
		}
		keyPair, err := tls.X509KeyPair(c.certificate, c.key)
		if err != nil {
			return E.Cause(err, "reload key pair")
		}
		c.access.Lock()
		config := c.config.Clone()
		config.Certificates = []tls.Certificate{keyPair}
		c.config = config
		c.access.Unlock()
		c.logger.Info("reloaded TLS certificate")
	} else if path == c.echKeyPath {
		echKey, err := os.ReadFile(c.echKeyPath)
		if err != nil {
			return E.Cause(err, "reload ECH keys from ", c.echKeyPath)
		}
		err = c.setECHServerConfig(echKey)
		if err != nil {
			return err
		}
		c.logger.Info("reloaded ECH keys")
	}
	return nil
}

func (c *STDServerConfig) Close() error {
	if c.watcher != nil {
		return c.watcher.Close()
	}
	return nil
}

func NewSTDServer(ctx context.Context, logger log.Logger, options option.InboundTLSOptions) (ServerConfig, error) {
	if !options.Enabled {
		return nil, nil
	}
	var err error
	tlsConfig := &tls.Config{}
	tlsConfig.Time = ntp.TimeFuncFromContext(ctx)
	if options.ServerName != "" {
		tlsConfig.ServerName = options.ServerName
	}
	if len(options.ALPN) > 0 {
		tlsConfig.NextProtos = append(options.ALPN, tlsConfig.NextProtos...)
	}
	if options.MinVersion != "" {
		minVersion, err := ParseTLSVersion(options.MinVersion)
		if err != nil {
			return nil, E.Cause(err, "parse min_version")
		}
		tlsConfig.MinVersion = minVersion
	}
	if options.MaxVersion != "" {
		maxVersion, err := ParseTLSVersion(options.MaxVersion)
		if err != nil {
			return nil, E.Cause(err, "parse max_version")
		}
		tlsConfig.MaxVersion = maxVersion
	}
	if options.CipherSuites != nil {
	find:
		for _, cipherSuite := range options.CipherSuites {
			for _, tlsCipherSuite := range tls.CipherSuites() {
				if cipherSuite == tlsCipherSuite.Name {
					tlsConfig.CipherSuites = append(tlsConfig.CipherSuites, tlsCipherSuite.ID)
					continue find
				}
			}
			return nil, E.New("unknown cipher_suite: ", cipherSuite)
		}
	}
	var certificate []byte
	var key []byte
	var echKeyPath string
	if options.ECH != nil && options.ECH.Enabled {
		err = parseECHServerConfig(ctx, options, tlsConfig, &echKeyPath)
		if err != nil {
			return nil, err
		}
	}
	serverConfig := &STDServerConfig{
		config:          tlsConfig,
		logger:          logger,
		certificate:     certificate,
		key:             key,
		certificatePath: options.CertificatePath,
		keyPath:         options.KeyPath,
		echKeyPath:      echKeyPath,
	}
	serverConfig.config.GetConfigForClient = func(info *tls.ClientHelloInfo) (*tls.Config, error) {
		serverConfig.access.Lock()
		defer serverConfig.access.Unlock()
		return serverConfig.config, nil
	}
	return serverConfig, nil
}
