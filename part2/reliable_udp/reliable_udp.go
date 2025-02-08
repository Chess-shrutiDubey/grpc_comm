package reliable_udp

import (
    "fmt"
    "math/rand"
    "net"
    "sync"
    "sync/atomic"
    "time"
)

const (
    MaxRetries    = 5
    RetryTimeout  = 100 * time.Millisecond
    MaxPacketSize = 1024
    ACK          = "ACK"
)

type Packet struct {
    SequenceNumber int64
    Data          []byte
    Timestamp     time.Time
}

type Statistics struct {
    mu            sync.Mutex
    sentPackets   int
    recvPackets   int
    lostPackets   int
    droppedPackets int  // New field for tracking initially dropped packets
    totalRTT      time.Duration
    dropRate      float64
}

var (
    stats           = Statistics{}
    currentSequence atomic.Int64
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

// createPacket creates a new packet with sequence number and timestamp
func createPacket(data []byte) Packet {
    return Packet{
        SequenceNumber: currentSequence.Add(1),
        Data:          data,
        Timestamp:     time.Now(),
    }
}

// validatePacket checks if the packet size is within limits
func validatePacket(data []byte) bool {
    return len(data) <= MaxPacketSize
}


// SendReliable sends data with retry mechanism
func SendReliable(conn *net.UDPConn, data string) (time.Duration, error) {
    if !validatePacket([]byte(data)) {
        return 0, fmt.Errorf("packet size exceeds maximum allowed size of %d bytes", MaxPacketSize)
    }

    packet := createPacket([]byte(data))
    start := time.Now()
    ackBuf := make([]byte, MaxPacketSize)

    stats.mu.Lock()
    stats.sentPackets++
    stats.mu.Unlock()

    for retry := 0; retry < MaxRetries; retry++ {
        // Send packet
        if _, err := conn.Write(packet.Data); err != nil {
            return 0, fmt.Errorf("send error: %v", err)
        }

        // Wait for ACK with timeout
        conn.SetReadDeadline(time.Now().Add(RetryTimeout))
        n, _, err := conn.ReadFrom(ackBuf)
        
        if err == nil && string(ackBuf[:n]) == ACK {
            rtt := time.Since(start)
            
            stats.mu.Lock()
            stats.recvPackets++
            stats.totalRTT += rtt
            stats.mu.Unlock()
            
            return rtt, nil
        }

        if retry < MaxRetries-1 {
            stats.mu.Lock()
            stats.droppedPackets++
            stats.mu.Unlock()
            
            fmt.Printf("Retry %d: packet %d (size: %d bytes)\n", 
                retry+1, packet.SequenceNumber, len(packet.Data))
        }
    }

    stats.mu.Lock()
    stats.lostPackets++
    stats.mu.Unlock()
    
    return 0, fmt.Errorf("max retries exceeded for packet %d", packet.SequenceNumber)
}
// ReceiveReliable handles incoming packets and sends ACKs
func ReceiveReliable(conn *net.UDPConn) ([]byte, *net.UDPAddr, error) {
    buffer := make([]byte, MaxPacketSize)
    n, addr, err := conn.ReadFromUDP(buffer)
    if err != nil {
        return nil, nil, fmt.Errorf("read error: %v", err)
    }

    stats.mu.Lock()
    stats.recvPackets++
    shouldDrop := (stats.dropRate > 0 && rand.Float64()*100 < stats.dropRate)
    stats.mu.Unlock()

    if shouldDrop {
        stats.mu.Lock()
        stats.lostPackets++
        stats.mu.Unlock()
        return nil, nil, fmt.Errorf("packet dropped (artificial loss)")
    }

    // Send ACK
    if _, err := conn.WriteToUDP([]byte(ACK), addr); err != nil {
        return nil, nil, fmt.Errorf("failed to send ACK: %v", err)
    }

    return buffer[:n], addr, nil
}

// GetStatistics returns current communication statistics
func GetStatistics() Statistics {
    stats.mu.Lock()
    defer stats.mu.Unlock()
    return stats
}

// SetDropRate sets artificial packet loss rate (0-100)
func SetDropRate(rate float64) {
    stats.mu.Lock()
    defer stats.mu.Unlock()
    if rate < 0 {
        stats.dropRate = 0
    } else if rate > 100 {
        stats.dropRate = 100
    } else {
        stats.dropRate = rate
    }
}

// ResetStatistics resets all statistics counters
func ResetStatistics() {
    stats.mu.Lock()
    defer stats.mu.Unlock()
    stats.sentPackets = 0
    stats.recvPackets = 0
    stats.lostPackets = 0
    stats.totalRTT = 0
}