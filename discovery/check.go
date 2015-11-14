package discovery

import (
	log "github.com/Sirupsen/logrus"
	cli "github.com/codegangsta/cli"

	"os"
	"time"
)

func Check(c *cli.Context) error {
	var (
		hbStr  = c.String("heartbeat")
		ttlStr = c.String("ttl")

		heartbeat time.Duration
		ttl       time.Duration
	)

	if Advertise = c.String("advertise"); Advertise == "" {
		cli.ShowAppHelp(c)
		log.Error("Required flag --advertise missing")
		os.Exit(1)
	}

	heartbeat, err := time.ParseDuration(hbStr)
	if err != nil {
		log.Fatal(err)
	}
	ttl, err = time.ParseDuration(ttlStr)
	if err != nil {
		log.Fatal(err)
	}

	if pos := c.Args(); len(pos) != 1 {
		cli.ShowAppHelp(c)
		log.Error("Required arguemnt DISCOVERY_URI")
		os.Exit(1)
	} else {
		Discovery = pos[0]
	}

	// register monitor instance
	Register(heartbeat, ttl)

	log.WithFields(log.Fields{
		"advertise": Advertise,
		"discovery": Discovery,
		"heartbeat": heartbeat,
		"ttl":       ttl,
	}).Info("begin advertise")

	return nil
}
