//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func main() {
	nodes := []struct {
		name string
		addr string
	}{
		{"Node A", "127.0.0.1:3001"},
		{"Node B", "127.0.0.1:3002"},
		{"Node C", "127.0.0.1:3003"},
	}

	fmt.Println("🚀 Starting HiveFS Load Generator for Grafana screen recording...")
	fmt.Println("This script will generate multiple traffic and connection spikes across all nodes.")
	fmt.Println("Press Ctrl+C to stop the load generator.\n")

	rand.Seed(time.Now().UnixNano())
	iteration := 1

	for {
		fmt.Printf("--- Iteration %d ---\n", iteration)

		// Choose a random node
		target := nodes[rand.Intn(len(nodes))]
		fmt.Printf("Connecting to %s on %s...\n", target.name, target.addr)

		conn, err := net.Dial("tcp", target.addr)
		if err != nil {
			fmt.Printf("❌ Failed to connect to %s: %v\n", target.name, err)
			time.Sleep(3 * time.Second)
			continue
		}

		// Keep connection open to show active peer count = 1
		fmt.Printf("✅ Peer connected to %s. Sending messages...\n", target.name)
		
		// Send a few messages over the same connection with small delays
		msgCount := rand.Intn(4) + 2 // 2 to 5 messages
		for i := 1; i <= msgCount; i++ {
			payloadSize := rand.Intn(800) + 100 // 100 to 900 bytes payload
			payload := make([]byte, payloadSize)
			rand.Read(payload)

			_, err = conn.Write(payload)
			if err != nil {
				fmt.Printf("❌ Error writing payload %d: %v\n", i, err)
				break
			}
			fmt.Printf("   👉 Sent message %d (%d bytes) to %s\n", i, payloadSize, target.name)
			time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond) // wait 500-1500ms
		}

		// Keep peer connected for a bit longer to stabilize the gauge in Grafana
		holdTime := time.Duration(rand.Intn(2)+2) * time.Second // 2 to 3 seconds
		fmt.Printf("Holding connection for %v to show active peer gauge...\n", holdTime)
		time.Sleep(holdTime)

		conn.Close()
		fmt.Printf("🔒 Connection to %s closed.\n", target.name)

		// Sleep between iterations to create distinct valleys/spikes in the Grafana graphs
		idleTime := time.Duration(rand.Intn(3)+3) * time.Second // 3 to 5 seconds
		fmt.Printf("Sleeping for %v to create valleys in the graphs...\n\n", idleTime)
		time.Sleep(idleTime)

		iteration++
	}
}
