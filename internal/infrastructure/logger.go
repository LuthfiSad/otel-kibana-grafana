package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func PushLog(layer, component, status, message string, payload map[string]interface{}) {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://elasticsearch:9200"
	}

	logEntry := map[string]interface{}{
		"@timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"service":    "massive-enterprise-app",
		"layer":      layer,
		"component":  component,
		"status":     status,
		"message":    message,
		"details":    payload,
	}

	body, err := json.Marshal(logEntry)
	if err != nil {
		return
	}

	go func() {
		client := &http.Client{Timeout: 3 * time.Second}
		indexName := fmt.Sprintf("app-logs-%s", time.Now().Format("2006.01.02"))
		url := fmt.Sprintf("%s/%s/_doc", esURL, indexName)
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
		if err == nil {
			resp.Body.Close()
		}
	}()
}
