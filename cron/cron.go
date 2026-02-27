package cron

import (
	"BlockMe/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/netresearch/go-cron"
)

type JobConfig struct {
	ResetKey *string
}

var resetKey *string

func InitCronJobs(jobConfig *JobConfig) *cron.Cron {
	fmt.Println("\nInitializing cron jobs...")
	c := cron.New()
	resetKey = jobConfig.ResetKey

	fmt.Printf("cron: adding ResetKeyRotator cron job with schedule %s\n", config.Env.IAmALittleBitchCron)
	_, err := c.AddFunc(config.Env.IAmALittleBitchCron, resetKeyRotator)
	if err != nil {
		fmt.Printf("cron: Failed to add ResetKeyRotator cron job. %s\n", err.Error())
	}

	c.Start()
	fmt.Println("cron: Starting cron jobs...")
	return c
}

func resetKeyRotator() {
	fmt.Println("cron.resetKeyRotator: Creating new resetKey...")
	newKey := uuid.New().String()
	*resetKey = newKey
	fmt.Printf("cron.resetKeyRotator: New resetKey created: %s\n", *resetKey)

	fmt.Println("cron.resetKeyRotator: Sending resetKey to discord")

	payload := map[string]interface{}{
		"content": fmt.Sprintf("New resetKey: %s\n\nhttp://localhost:8080/reset?resetKey=%s", *resetKey, *resetKey),
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("cron.resetKeyRotator: Error marshaling JSON: %v", err)
	}

	resp, err := http.Post(
		config.Env.IAmALittleBitchUrl,
		"application/json",
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		fmt.Printf("cron.resetKeyRotator: Error sending webhook: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		fmt.Println("cron.resetKeyRotator: Webhook sent successfully!")
	} else {
		fmt.Printf("cron.resetKeyRotator: Failed to send webhook. Status code: %d\n", resp.StatusCode)
	}
}
