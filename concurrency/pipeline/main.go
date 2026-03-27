package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Service   string
	Action    string
	Payload   any
	Username  string
}

type UserPayload struct {
	UserID int
}

type PaymentPayload struct {
	OrderID int
	Amount  int
}

type OrderPayload struct {
	OrderID int
	UserID  int
}

type StockPayload struct {
	SKU   string
	Delta int
}

type EmailPayload struct {
	Template string
	UserID   int
}

func ParseLogLine(line string) (LogEntry, error) {
	parts := strings.SplitN(strings.TrimSpace(line), "|", 5)
	if len(parts) != 5 {
		return LogEntry{}, fmt.Errorf("invalid log format: %q", line)
	}

	ts, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return LogEntry{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	payload, err := parsePayloadByAction(parts[3], parts[4])
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		Timestamp: ts,
		Level:     parts[1],
		Service:   parts[2],
		Action:    parts[3],
		Payload:   payload,
	}, nil
}

func parsePayloadByAction(action, payload string) (any, error) {
	kv := parseKV(payload)

	switch action {
	case "user_login", "user_logout", "token_refresh", "password_reset":
		userID, err := mustInt(kv, "user_id")
		if err != nil {
			return nil, err
		}
		return UserPayload{UserID: userID}, nil

	case "payment_success", "payment_failed":
		orderID, err := mustInt(kv, "order_id")
		if err != nil {
			return nil, err
		}
		amount, err := mustInt(kv, "amount")
		if err != nil {
			return nil, err
		}
		return PaymentPayload{OrderID: orderID, Amount: amount}, nil

	case "order_created", "order_cancelled":
		orderID, err := mustInt(kv, "order_id")
		if err != nil {
			return nil, err
		}
		userID, err := mustInt(kv, "user_id")
		if err != nil {
			return nil, err
		}
		return OrderPayload{OrderID: orderID, UserID: userID}, nil

	case "stock_updated":
		sku, ok := kv["sku"]
		if !ok {
			return nil, errors.New("missing field sku")
		}
		delta, err := mustInt(kv, "delta")
		if err != nil {
			return nil, err
		}
		return StockPayload{SKU: sku, Delta: delta}, nil

	case "email_sent":
		template, ok := kv["template"]
		if !ok {
			return nil, errors.New("missing field template")
		}
		userID, err := mustInt(kv, "user_id")
		if err != nil {
			return nil, err
		}
		return EmailPayload{Template: template, UserID: userID}, nil

	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}
}

func parseKV(payload string) map[string]string {
	out := map[string]string{}
	for _, token := range strings.Fields(payload) {
		parts := strings.SplitN(token, "=", 2)
		if len(parts) == 2 {
			out[parts[0]] = parts[1]
		}
	}
	return out
}

func mustInt(kv map[string]string, key string) (int, error) {
	v, ok := kv[key]
	if !ok {
		return 0, fmt.Errorf("missing field %s", key)
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid int %s=%q", key, v)
	}
	return n, nil
}

func extractUserID(payload any) (int, bool) {
	switch p := payload.(type) {
	case UserPayload:
		return p.UserID, true
	case OrderPayload:
		return p.UserID, true
	case EmailPayload:
		return p.UserID, true
	default:
		return 0, false
	}
}

func main() {
	start := time.Now()

	generatorQueue := make(chan string, 256)
	parserQueue := make(chan LogEntry, 256)
	filterQueue := make(chan LogEntry, 256)
	enricherQueue := make(chan LogEntry, 256)

	file, err := os.Open("./log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	userMap := map[int]string{
		101: "alice",
		102: "bob",
		103: "charlie",
		104: "diana",
		105: "eric",
		106: "farah",
		107: "george",
		108: "hana",
		109: "ivan",
		110: "jane",
		111: "kyle",
		112: "lara",
		113: "mike",
		114: "nina",
		115: "omar",
		116: "priya",
		117: "quinn",
		118: "ravi",
		119: "sana",
		120: "tina",
	}

	var emitted int64
	var parsed int64
	var parseFailed int64
	var filteredDebug int64
	var filteredStale int64
	var enriched int64

	go func() {
		for scanner.Scan() {
			generatorQueue <- scanner.Text()
			atomic.AddInt64(&emitted, 1)
		}
		close(generatorQueue)
	}()

	var parserWG sync.WaitGroup
	for i := 1; i <= 3; i++ {
		parserWG.Add(1)
		go func() {
			defer parserWG.Done()
			for line := range generatorQueue {
				entry, err := ParseLogLine(line)
				if err != nil {
					atomic.AddInt64(&parseFailed, 1)
					continue
				}
				atomic.AddInt64(&parsed, 1)
				parserQueue <- entry
			}
		}()
	}

	go func() {
		parserWG.Wait()
		close(parserQueue)
	}()

	go func() {
		cutoff := time.Now().UTC().Add(-1 * time.Hour)
		for entry := range parserQueue {
			if entry.Level == "DEBUG" {
				atomic.AddInt64(&filteredDebug, 1)
				continue
			}
			if entry.Timestamp.Before(cutoff) {
				atomic.AddInt64(&filteredStale, 1)
				continue
			}
			filterQueue <- entry
		}
		close(filterQueue)
	}()

	var enricherWG sync.WaitGroup
	for i := 1; i <= 2; i++ {
		enricherWG.Add(1)
		go func() {
			defer enricherWG.Done()
			for entry := range filterQueue {
				if userID, ok := extractUserID(entry.Payload); ok {
					if username, found := userMap[userID]; found {
						entry.Username = username
					} else {
						entry.Username = "unknown"
					}
				}
				atomic.AddInt64(&enriched, 1)
				enricherQueue <- entry
			}
		}()
	}

	go func() {
		enricherWG.Wait()
		close(enricherQueue)
	}()

	serviceCount := map[string]int{}
	levelCount := map[string]int{}
	for entry := range enricherQueue {
		serviceCount[entry.Service]++
		levelCount[entry.Level]++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[Generator] Emitted: %d\n", emitted)
	fmt.Printf("[Parser] Parsed: %d, Failed: %d\n", parsed, parseFailed)
	fmt.Printf("[Filter] Dropped DEBUG: %d, Dropped stale: %d\n", filteredDebug, filteredStale)
	fmt.Printf("[Enricher] Forwarded: %d\n", enriched)

	fmt.Println("\n========== PIPELINE RESULTS ==========")
	fmt.Println("Events per service:")
	for service, count := range serviceCount {
		fmt.Printf("  %s: %d\n", service, count)
	}
	fmt.Println("Events per log level:")
	for level, count := range levelCount {
		fmt.Printf("  %s: %d\n", level, count)
	}
	fmt.Printf("Pipeline duration: %s\n", time.Since(start))
	fmt.Println("======================================")

}
