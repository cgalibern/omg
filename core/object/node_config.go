package object

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"opensvc.com/opensvc/core/keyop"
	"opensvc.com/opensvc/core/rawconfig"
	"opensvc.com/opensvc/core/xconfig"
	"opensvc.com/opensvc/util/hostname"
	"opensvc.com/opensvc/util/key"
)

func (t Node) Log() *zerolog.Logger {
	return &t.log
}

func (t *Node) ConfigFile() string {
	return filepath.Join(rawconfig.Node.Paths.Etc, "node.conf")
}

func (t *Node) ClusterConfigFile() string {
	return filepath.Join(rawconfig.Node.Paths.Etc, "cluster.conf")
}

func (t *Node) loadConfig() error {
	var err error
	if t.config, err = xconfig.NewObject(t.ConfigFile()); err != nil {
		return err
	}
	t.config.Referrer = t
	if t.mergedConfig, err = xconfig.NewObject(t.ConfigFile(), t.ClusterConfigFile()); err != nil {
		return err
	}
	t.mergedConfig.Referrer = t
	return err
}

func (t Node) Config() *xconfig.T {
	return t.config
}

func (t Node) MergedConfig() *xconfig.T {
	return t.mergedConfig
}

func (t Node) ID() uuid.UUID {
	if t.id != uuid.Nil {
		return t.id
	}
	idKey := key.Parse("id")
	if idStr := t.config.GetString(idKey); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			t.id = id
			return t.id
		}
	}
	t.id = uuid.New()
	op := keyop.T{
		Key:   key.Parse("id"),
		Op:    keyop.Set,
		Value: t.id.String(),
	}
	_ = t.config.Set(op)
	if err := t.config.Commit(); err != nil {
		t.log.Error().Err(err).Msg("")
	}
	return t.id
}

func (t Node) Env() string {
	k := key.Parse("env")
	if s := t.config.GetString(k); s != "" {
		return s
	}
	return rawconfig.Node.Node.Env
}

func (t Node) App() string {
	k := key.Parse("app")
	return t.config.GetString(k)
}

func (t Node) Dereference(ref string) (string, error) {
	switch ref {
	case "id":
		return t.ID().String(), nil
	case "name", "nodename":
		return hostname.Hostname(), nil
	case "short_name", "short_nodename":
		return strings.SplitN(hostname.Hostname(), ".", 1)[0], nil
	case "dnsuxsock":
		return t.DNSUDSFile(), nil
	case "dnsuxsockd":
		return t.DNSUDSDir(), nil
	}
	switch {
	case strings.HasPrefix(ref, "safe://"):
		return ref, fmt.Errorf("TODO")
	}
	return ref, fmt.Errorf("unknown reference: %s", ref)
}

func (t Node) PostCommit() error {
	return nil
}

func (t Node) Nodes() []string {
	return []string{hostname.Hostname()}
}

func (t Node) DRPNodes() []string {
	return []string{}
}

func (t Node) EncapNodes() []string {
	return []string{}
}
