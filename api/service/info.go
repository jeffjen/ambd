package service

import (
	disc "github.com/jeffjen/docker-ambassador/discovery"

	_ "github.com/Sirupsen/logrus"

	"encoding/json"
	"net/http"
	"os"
	"time"
)

var (
	VERSION = os.Getenv("VERSION")

	BUILD = os.Getenv("BUILD")

	NODE_NAME = os.Getenv("NODE_NAME")

	NODE_REGION = os.Getenv("NODE_REGION")

	NODE_AVAIL_ZONE = os.Getenv("NODE_AVAIL_ZONE")
)

type serverInfo struct {
	Version   string `json:"version"`
	Build     string `json:"build"`
	Node      string `json:"node"`
	Region    string `json:"region"`
	Zone      string `json:"avail_zone"`
	Discovery string `json:"discovery"`
	Hearbeat  string `json:"heartbeat"`
	TTL       string `json:"ttl"`
	Timestamp string `json:"current_time"`
}

func Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(serverInfo{
		Version:   VERSION,
		Build:     BUILD,
		Node:      NODE_NAME,
		Region:    NODE_REGION,
		Zone:      NODE_AVAIL_ZONE,
		Discovery: disc.Discovery,
		Hearbeat:  disc.Hearbeat.String(),
		TTL:       disc.TTL.String(),
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
