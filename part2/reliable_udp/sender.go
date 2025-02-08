package reliable_udp

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "time"
)

func RunSender() {
    // Validate command line args
    if len(os.Args) < 3 {
        fmt.Println("Usage: sender <num_packets> <loss_rate>")
        return
    }

    numPackets, err := strconv.Atoi(os.Args[1])
    if err != nil || numPackets <= 0 {
        fmt.Println("Invalid number of packets")
        return
    }

    lossRate, err := strconv.ParseFloat(os.Args[2], 64) 
    if err != nil || lossRate < 0 || lossRate > 100 {
        fmt.Println("Invalid loss rate (must be 0-100)")
        return
    }

    // Connect to receiver
    serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
    if err != nil {
        fmt.Printf("Error resolving address: %v\n", err)
        return
    }

    conn, err := net.DialUDP("udp", nil, serverAddr) 
    if err != nil {
        fmt.Printf("Error connecting: %v\n", err)
        return
    }
    defer conn.Close()

    // Test RTT with small message
    fmt.Println("Testing RTT with small message...")
    msg := "ping"
    rtt, err := SendReliable(conn, msg)
    if err != nil {
        fmt.Printf("Error sending ping: %v\n", err)
        return
    }
    fmt.Printf("RTT for small message: %v\n", rtt)

    
    // ... existing validation code ...

    fmt.Printf("\nMeasuring bandwidth with %d packets...\n", numPackets)
    start := time.Now()
    successCount := 0
    totalBytes := 0
    msgSize := 100 // Fixed size messages for consistent bandwidth measurement

    for i := 0; i < numPackets; i++ {
        // Create fixed-size packet
        msg := fmt.Sprintf("%-*d", msgSize-1, i) // Pad with spaces to achieve fixed size
        if rtt, err := SendReliable(conn, msg); err == nil {
            successCount++
            totalBytes += len(msg)
            fmt.Printf("Packet %d sent successfully, RTT: %v\n", i, rtt)
        } else {
            fmt.Printf("Packet %d failed: %v\n", i, err)
        }
        time.Sleep(1 * time.Millisecond) // Add small delay between packets
    }

    duration := time.Since(start)
    stats := GetStatistics()
    
    // Calculate actual bandwidth (bytes/second)
    bandwidthMBps := float64(totalBytes) / duration.Seconds() / (1024 * 1024)
    
    fmt.Printf("\nResults:\n")
    fmt.Printf("Packets sent: %d\n", numPackets)
    fmt.Printf("Packets received: %d\n", successCount)
    fmt.Printf("Initial drops: %d\n", stats.droppedPackets)
    fmt.Printf("Final packet loss: %.2f%%\n", 100-float64(successCount)/float64(numPackets)*100)
    fmt.Printf("Bandwidth: %.5f MB/s\n", bandwidthMBps)
    fmt.Printf("Average RTT: %v\n", stats.totalRTT/time.Duration(successCount))
    fmt.Printf("Total duration: %v\n", duration)
    
}