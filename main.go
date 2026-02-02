package main

import (
	"fmt"
	"net"
	"strings"
)

// List of domains to block (Blacklist)
var blacklist = map[string]bool{
	"doubleclick.net": true,
	"google-analytics.com": true,
	"facebook.com": true, // Testing ke liye
}

func main() {
	// Step 1: Listen on UDP Port 53 (DNS Port)
	// Note: Phone/PC par iske liye sudo/admin permissions chahiye hongi
	addr := net.UDPAddr{
		Port: 53,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("❌ Error: Port 53 bind nahi ho paya: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("🛡️ Shield-CLI is active on 127.0.0.1:53")
	fmt.Println("Press Ctrl+C to stop.")

	buffer := make([]byte, 512)
	for {
		// Step 2: Read incoming DNS packet
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP:", err)
			continue
		}

		// Simple logic to extract domain (Minimal parsing for demo)
		// Real world mein hum "miekg/dns" library use karte hain logic ke liye
		query := string(buffer[:n])
		
		fmt.Printf("🔍 Query received from %s\n", remoteAddr)

		// Step 3: Check Blacklist and Logic
		// Abhi ke liye hum simple console print kar rahe hain
		// Asli implementation mein yahan binary header parse hoga
		blocked := false
		for domain := range blacklist {
			if strings.Contains(query, domain) {
				blocked = true
				break
			}
		}

		if blocked {
			fmt.Println("🚫 BLOCKED: Ad/Tracker domain detected!")
		} else {
			fmt.Println("✅ ALLOWED: Forwarding to Upstream DNS...")
			// TODO: Forward to 8.8.8.8
		}
	}
}
