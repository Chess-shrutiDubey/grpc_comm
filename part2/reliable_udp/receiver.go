package reliable_udp

import (
    "fmt"
    "math/rand"
    "net"
    "os"
    "strconv"
    "time"
)

func RunReceiver() {
    rand.Seed(time.Now().UnixNano())

    // Parse drop rate from command line
    dropRate := 0.0
    if len(os.Args) > 1 {
        if rate, err := strconv.ParseFloat(os.Args[1], 64); err == nil {
            dropRate = rate
        }
    }

    // Create UDP listener
    addr, err := net.ResolveUDPAddr("udp", ":8080")
    if err != nil {
        fmt.Printf("Error resolving address: %v\n", err)
        return
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Printf("Error starting listener: %v\n", err)
        return 
    }
    defer conn.Close()


    fmt.Printf("Receiver started (drop rate: %.1f%%)\n", dropRate)
    SetDropRate(dropRate)

    buffer := make([]byte, MaxPacketSize)
    for {
        n, remoteAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Printf("Error reading: %v\n", err)
            continue
        }

        // Apply artificial packet loss
        if rand.Float64() * 100 < dropRate {
            fmt.Printf("Dropping packet from %v (%.1f%% drop rate)\n", 
                remoteAddr, dropRate)
            continue
        }

        // Send ACK
        if _, err := conn.WriteToUDP([]byte(ACK), remoteAddr); err != nil {
            fmt.Printf("Error sending ACK: %v\n", err)
            continue
        }
        
        fmt.Printf("Received %d bytes from %v\n", n, remoteAddr)
    }
    
}