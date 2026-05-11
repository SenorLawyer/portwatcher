package scanner

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	gonet "github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type Scanner struct {
	docker DockerMapper
}

func New(docker DockerMapper) *Scanner {
	return &Scanner{docker: docker}
}

func (s *Scanner) Snapshot(ctx context.Context) (Snapshot, error) {
	start := time.Now()
	conns, err := gonet.ConnectionsWithContext(ctx, "all")
	if err != nil {
		return Snapshot{}, err
	}

	processes := map[int32]procInfo{}
	ports := make([]PortEntry, 0, len(conns))
	warnings := make([]string, 0)
	now := time.Now()

	for _, conn := range conns {
		if conn.Laddr.Port == 0 {
			continue
		}
		protocol := protocolName(conn.Type)
		state := strings.ToLower(conn.Status)
		if state == "" || state == "none" {
			state = "open"
		}
		remote := remoteEndpoint(conn.Raddr.IP, conn.Raddr.Port)

		info, ok := processes[conn.Pid]
		if !ok && conn.Pid > 0 {
			info = loadProcess(ctx, conn.Pid)
			processes[conn.Pid] = info
			warnings = append(warnings, info.warnings...)
		}

		entry := PortEntry{
			ID:       fmt.Sprintf("%s/%s/%d/%s/%d/%d", protocol, conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port, conn.Pid),
			Protocol: protocol,
			Host:     conn.Laddr.IP,
			Port:     conn.Laddr.Port,
			Address:  endpoint(conn.Laddr.IP, conn.Laddr.Port),
			Remote:   remote,
			State:    state,
			PID:      conn.Pid,
			Process:  info.name,
			Command:  info.command,
			User:     info.user,
			SeenAt:   now,
		}
		ports = append(ports, entry)
	}

	sort.SliceStable(ports, func(i, j int) bool {
		if ports[i].Port == ports[j].Port {
			return ports[i].Protocol < ports[j].Protocol
		}
		return ports[i].Port < ports[j].Port
	})

	if s.docker != nil {
		enriched, dockerWarnings := s.docker.Enrich(ctx, ports)
		ports = enriched
		warnings = append(warnings, dockerWarnings...)
	}

	return Snapshot{
		Ports:    ports,
		At:       now,
		Duration: time.Since(start).String(),
		Warnings: compactWarnings(warnings),
	}, nil
}

type procInfo struct {
	name     string
	command  string
	user     string
	warnings []string
}

func loadProcess(ctx context.Context, pid int32) procInfo {
	p, err := process.NewProcessWithContext(ctx, pid)
	if err != nil {
		return procInfo{warnings: []string{fmt.Sprintf("pid %d: %v", pid, err)}}
	}

	var info procInfo
	if name, err := p.NameWithContext(ctx); err == nil {
		info.name = name
	}
	if cmd, err := p.CmdlineWithContext(ctx); err == nil {
		info.command = cmd
	}
	if user, err := p.UsernameWithContext(ctx); err == nil {
		info.user = user
	}
	if info.command == "" {
		info.command = info.name
	}
	return info
}

func protocolName(kind uint32) string {
	switch kind {
	case 1:
		return "tcp"
	case 2:
		return "udp"
	default:
		return fmt.Sprintf("ip%d", kind)
	}
}

func endpoint(host string, port uint32) string {
	if host == "" && port == 0 {
		return ""
	}
	return net.JoinHostPort(displayHost(host), fmt.Sprintf("%d", port))
}

func remoteEndpoint(host string, port uint32) string {
	if port == 0 {
		return ""
	}
	return endpoint(host, port)
}

func displayHost(host string) string {
	if host == "" || host == "::" {
		return "*"
	}
	return host
}

func compactWarnings(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, warning := range in {
		if warning == "" {
			continue
		}
		if _, ok := seen[warning]; ok {
			continue
		}
		seen[warning] = struct{}{}
		out = append(out, warning)
	}
	if len(out) > 5 {
		return append(out[:5], fmt.Sprintf("%d more warnings", len(out)-5))
	}
	return out
}
