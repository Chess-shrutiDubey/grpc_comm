package reliable_udp

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type PacketStatus struct {
	sent     bool
	received bool
	dropped  bool
	rtt      time.Duration
}

func RunSender() {
	// Parse command line arguments
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <packet_count> <drop_rate> <packet_size>\n", os.Args[0])
		return
	}

	numPackets, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid packet count: %v", err)
	}

	dropRate, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		log.Fatalf("Invalid drop rate: %v", err)
	}

	msgSize, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Invalid packet size: %v", err)
	}

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Initialize metrics
	start := time.Now()
	totalSent := 0
	totalReceived := 0
	totalBytes := 0
	var totalRTT time.Duration
	packetStatus := make([]PacketStatus, numPackets)

	// Create test data
	testData := make([]byte, msgSize)
	rand.Read(testData)

	// Connect to receiver
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Failed to resolve address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Send packets
	for i := 0; i < numPackets; i++ {
		// Simulate packet loss based on drop rate
		if rand.Float64() < dropRate/100.0 {
			packetStatus[i].dropped = true
			fmt.Fprintf(os.Stderr, "Packet %d dropped (simulated)\n", i)
			continue
		}

		sendTime := time.Now()
		_, err := conn.Write(testData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send packet %d: %v\n", i, err)
			continue
		}

		packetStatus[i].sent = true
		totalSent++
		totalBytes += len(testData)

		// Wait for ACK
		ackBuffer := make([]byte, 64)
		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		_, _, err = conn.ReadFromUDP(ackBuffer)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Packet %d ACK timeout\n", i)
		} else {
			rtt := time.Since(sendTime)
			packetStatus[i].received = true
			packetStatus[i].rtt = rtt
			totalRTT += rtt
			totalReceived++
			fmt.Fprintf(os.Stderr, "Packet %d sent successfully, RTT: %v\n", i, rtt)
		}

		time.Sleep(time.Millisecond)
	}

	// Calculate final metrics
	duration := time.Since(start)
	var avgRTT float64
	if totalReceived > 0 {
		avgRTT = float64(totalRTT.Microseconds()) / float64(totalReceived) / 1000.0 // Convert to ms
	}
	lossRate := 100.0 * float64(totalSent-totalReceived) / float64(totalSent)
	bandwidthMBps := float64(totalBytes) / duration.Seconds() / (1024 * 1024)

	// Write CSV metrics to stdout
	fmt.Printf("Metric,Value\n")
	fmt.Printf("Packets_Sent,%d\n", totalSent)
	fmt.Printf("Packets_Received,%d\n", totalReceived)
	fmt.Printf("Dropped_Packets,%d\n", totalSent-totalReceived)
	fmt.Printf("Packet_Loss_Rate,%.2f\n", lossRate)
	fmt.Printf("Bandwidth_MBps,%.5f\n", bandwidthMBps)
	fmt.Printf("Average_RTT_ms,%.3f\n", avgRTT)
}
