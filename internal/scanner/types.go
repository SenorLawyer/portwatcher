package scanner

import "time"

type PortEntry struct {
	ID          string    `json:"id"`
	Protocol    string    `json:"protocol"`
	Address     string    `json:"address"`
	Host        string    `json:"host"`
	Port        uint32    `json:"port"`
	Remote      string    `json:"remote,omitempty"`
	State       string    `json:"state"`
	PID         int32     `json:"pid"`
	Process     string    `json:"process,omitempty"`
	Command     string    `json:"command,omitempty"`
	User        string    `json:"user,omitempty"`
	ContainerID string    `json:"container_id,omitempty"`
	Container   string    `json:"container,omitempty"`
	Image       string    `json:"image,omitempty"`
	PrivatePort uint16    `json:"private_port,omitempty"`
	SeenAt      time.Time `json:"seen_at"`
	Warnings    []string  `json:"warnings,omitempty"`
}

type Snapshot struct {
	Ports    []PortEntry `json:"ports"`
	At       time.Time   `json:"at"`
	Duration string      `json:"duration"`
	Warnings []string    `json:"warnings,omitempty"`
}

type DockerMapper interface {
	Enrich(ctx Context, ports []PortEntry) ([]PortEntry, []string)
}

type Context interface {
	Done() <-chan struct{}
	Err() error
}
