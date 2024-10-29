package monitor

import (
	"sync"
	"time"
)

type NetworkStats struct {
	Timestamp     time.Time              `json:"timestamp"`
	Connections   []ConnectionInfo       `json:"connections"`
	ProcessStats  map[int32]ProcessStats `json:"process_stats"`
	mutex         sync.RWMutex
	lastIOCounts  map[int32]ProcessIO
	totalReceived uint64
	totalSent     uint64
}

type ConnectionInfo struct {
	Pid           int32  `json:"pid"`
	ProcessName   string `json:"process_name"`
	LocalAddr     string `json:"local_addr"`
	LocalPort     uint32 `json:"local_port"`
	RemoteAddr    string `json:"remote_addr"`
	RemotePort    uint32 `json:"remote_port"`
	Status        string `json:"status"`
	Protocol      string `json:"protocol"`
	BytesSent     uint64 `json:"bytes_sent"`
	BytesReceived uint64 `json:"bytes_received"`
}

type ProcessStats struct {
	Name          string `json:"name"`
	BytesSent     uint64 `json:"bytes_sent"`
	BytesReceived uint64 `json:"bytes_received"`
	Connections   int    `json:"connections"`
}

type ProcessIO struct {
	ReadBytes  uint64
	WriteBytes uint64
}
