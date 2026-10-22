package main

import (
	"strings"
	"testing"

	"github.com/kevinadi/mlgo/internal/topology"
)

func TestMongodRoleIsExplicit(t *testing.T) {
	m := new(Mongod)
	if err := m.Init(27017, false, "configuration", false, topology.RoleData); err != nil {
		t.Fatal(err)
	}
	args := strings.Join(m.Cmd_mongod(), " ")
	if strings.Contains(args, "--configsvr") || strings.Contains(args, "--shardsvr") {
		t.Fatalf("set name must not select a server role: %s", args)
	}

	shard := new(Mongod)
	if err := shard.Init(27018, false, "not-a-shard-name", false, topology.RoleShard); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(shard.Cmd_mongod(), " "), "--shardsvr") {
		t.Fatal("expected --shardsvr")
	}

	cfg := new(Mongod)
	if err := cfg.Init(27019, false, "cfg", false, topology.RoleConfig); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(cfg.Cmd_mongod(), " "), "--configsvr") {
		t.Fatal("expected --configsvr")
	}
}

func TestReplSetInitUsesPlanner(t *testing.T) {
	rs := new(ReplSet)
	err := rs.Init(2, 27017, "PSA", true, true, "replset", false, true, true, topology.RoleData)
	if err == nil {
		t.Fatal("expected -cfg and -num to conflict")
	}

	rs = new(ReplSet)
	if err := rs.Init(1, 27017, "PSH", true, false, "replset", false, false, true, topology.RoleData); err != nil {
		t.Fatal(err)
	}
	if rs.Num != 3 || rs.Config != "PSH" {
		t.Fatalf("num=%d cfg=%s", rs.Num, rs.Config)
	}
	joined := strings.Join(rs.Cmdlines[len(rs.Cmdlines)-1], " ")
	if !strings.Contains(joined, `"hidden":true`) || !strings.Contains(joined, `"votes":0`) {
		t.Fatalf("initiate command = %s", joined)
	}
}
