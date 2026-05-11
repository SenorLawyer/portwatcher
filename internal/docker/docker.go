package docker

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/SenorLawyer/portwatcher/internal/scanner"
)

type Enricher struct {
	enabled bool
}

func New(enabled bool) *Enricher {
	return &Enricher{enabled: enabled}
}

func (e *Enricher) Enrich(ctx scanner.Context, ports []scanner.PortEntry) ([]scanner.PortEntry, []string) {
	if !e.enabled {
		return ports, nil
	}
	stdCtx, ok := ctx.(context.Context)
	if !ok {
		stdCtx = context.Background()
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return ports, []string{fmt.Sprintf("docker: %v", err)}
	}
	defer cli.Close()

	containers, err := cli.ContainerList(stdCtx, container.ListOptions{All: false})
	if err != nil {
		return ports, []string{fmt.Sprintf("docker: %v", err)}
	}

	type portKey struct {
		protocol string
		port     uint32
	}
	byPort := map[portKey]types.Container{}
	privateByPort := map[portKey]uint16{}
	for _, c := range containers {
		for _, p := range c.Ports {
			if p.PublicPort == 0 {
				continue
			}
			key := portKey{protocol: strings.ToLower(p.Type), port: uint32(p.PublicPort)}
			byPort[key] = c
			privateByPort[key] = p.PrivatePort
		}
	}

	for i := range ports {
		key := portKey{protocol: ports[i].Protocol, port: ports[i].Port}
		c, ok := byPort[key]
		if !ok {
			continue
		}
		ports[i].ContainerID = shortID(c.ID)
		ports[i].Container = cleanName(c.Names)
		ports[i].Image = c.Image
		ports[i].PrivatePort = privateByPort[key]
	}
	return ports, nil
}

func cleanName(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return strings.TrimPrefix(names[0], "/")
}

func shortID(id string) string {
	if len(id) <= 12 {
		return id
	}
	if _, err := strconv.ParseUint(id[:12], 16, 64); err == nil {
		return id[:12]
	}
	return id
}
