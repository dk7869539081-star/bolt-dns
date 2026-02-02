package main

import (
	"fmt"
	"net"
	"strings"
	"os"
	"os/signal"
	"syscall"
)

// Global Counters for Stats
var totalQueries = 0
var blockedQueries = 0

// Blacklist Domains
var blacklist = map[string]bool{
	"doubleclick.net":      true,
	"google-analytics.com": true,
	"facebook.com":         true,
	"ads.twitter.com":      true,
	"telemetry.main.com":   true,
}

func main() {
	// Step 1: Listen on Localhost Port 53
	// TIP: Agar "Permission Denied" aaye toh Port 5353 try karna testing ke liye
	addr := net.UDPAddr{
		Port: 53,
		IP:   net.ParseIP("127.0.0.1"),
	}
	
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("❌ ERROR: Port 53 bind nahi ho saka. \n💡 Tip: Run as Admin/Sudo ya Port change karo.\nDetails: %v\n", err)
		return
	}
	defer conn.Close()

	// Exit handle karne ke liye (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("🛡️  SHIELD-CLI STARTED SUCCESSFULLY")
	fmt.Println("📍 Listening on 127.0.0.1:53")
	fmt.Println("-------------------------------------------")

	buffer := make([]byte, 512)

	// Background loop to process queries
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}

		totalQueries++
		queryData := string(buffer[:n])
		isBlocked := false

		// Check if any blacklisted domain is in the query
		for domain := range blacklist {
			if strings.Contains(queryData, domain) {
				isBlocked = true
				break
			}
		}

		if isBlocked {
			blockedQueries++
			fmt.Printf("🚫 [BLOCKED] Query from %s\n", remoteAddr)
		} else {
			fmt.Printf("✅ [ALLOWED] Query from %s\n", remoteAddr)
		}

		// LIVE DASHBOARD PRINT
		showStats()
	}
}

func showStats() {
	// Terminal clear karke stats dikhane ka magic
	// \033[H\033[2J terminal clear karta hai
	fmt.Print("\033[H\033[2J")
	fmt.Println("===========================================")
	fmt.Println("      🛡️  SHIELD-CLI LIVE DASHBOARD       ")
	fmt.Println("===========================================")
	fmt.Printf("  TOTAL QUERIES   : %d\n", totalQueries)
	fmt.Printf("  ADS BLOCKED     : %d\n", blockedQueries)
	if totalQueries > 0 {
		efficiency := (float64(blockedQueries) / float64(totalQueries)) * 100
		fmt.Printf("  PROTECTION RATE : %.2f%%\n", efficiency)
	}
	fmt.Println("===========================================")
	fmt.Println(" (Press Ctrl+C to stop the proxy)")
}
