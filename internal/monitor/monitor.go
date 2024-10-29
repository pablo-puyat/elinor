package monitor

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type Monitor struct {
	stats    NetworkStats
	logger   *log.Logger
	stopChan chan struct{}
}

func New(logger *log.Logger) *Monitor {
	return &Monitor{
		stats: NetworkStats{
			lastIOCounts: make(map[int32]ProcessIO),
		},
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (m *Monitor) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.updateStats()
		case <-m.stopChan:
			return
		}
	}
}

func (m *Monitor) Stop() {
	close(m.stopChan)
}

func (m *Monitor) GetStats() NetworkStats {
	m.stats.mutex.RLock()
	defer m.stats.mutex.RUnlock()
	return m.stats
}

func (m *Monitor) updateStats() {
	connections, err := net.Connections("all")
	if err != nil {
		m.logger.Printf("Error getting connections: %v", err)
		return
	}

	m.stats.mutex.Lock()
	defer m.stats.mutex.Unlock()

	m.stats.Timestamp = time.Now()
	m.stats.Connections = make([]ConnectionInfo, 0)
	processStats := make(map[int32]ProcessStats)

	for _, conn := range connections {
		if conn.Pid == 0 {
			continue
		}

		proc, err := process.NewProcess(conn.Pid)
		if err != nil {
			continue
		}

		name, err := proc.Name()
		if err != nil {
			continue
		}

		ioCounters, err := proc.IOCounters()
		if err != nil {
			continue
		}

		lastIO := m.stats.lastIOCounts[conn.Pid]
		bytesSent := ioCounters.WriteBytes - lastIO.WriteBytes
		bytesReceived := ioCounters.ReadBytes - lastIO.ReadBytes

		m.stats.lastIOCounts[conn.Pid] = ProcessIO{
			ReadBytes:  ioCounters.ReadBytes,
			WriteBytes: ioCounters.WriteBytes,
		}

		connInfo := ConnectionInfo{
			Pid:           conn.Pid,
			ProcessName:   name,
			LocalAddr:     conn.Laddr.IP,
			LocalPort:     conn.Laddr.Port,
			RemoteAddr:    conn.Raddr.IP,
			RemotePort:    conn.Raddr.Port,
			Status:        conn.Status,
			Protocol:      conn.Type,
			BytesSent:     bytesSent,
			BytesReceived: bytesReceived,
		}

		m.stats.Connections = append(m.stats.Connections, connInfo)

		pStats := processStats[conn.Pid]
		pStats.Name = name
		pStats.BytesSent += bytesSent
		pStats.BytesReceived += bytesReceived
		pStats.Connections++
		processStats[conn.Pid] = pStats
	}

	m.stats.ProcessStats = processStats
}
