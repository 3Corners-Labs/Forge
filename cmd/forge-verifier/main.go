package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
	verifierruntime "github.com/3corners/forge/internal/verifier"
)

func main() {
	planPath := flag.String("plan", "", "immutable Assurance Plan")
	bundlePath := flag.String("bundle", "", "immutable candidate bundle")
	workspace := flag.String("workspace", "", "fresh source checkout")
	resultPath := flag.String("result", "", "verification result output")
	flag.Parse()
	if *planPath == "" || *bundlePath == "" || *workspace == "" || *resultPath == "" {
		exit("verifier arguments are invalid")
	}
	planBytes, err := os.ReadFile(*planPath)
	if err != nil {
		exit("assurance plan is unavailable")
	}
	bundleBytes, err := os.ReadFile(*bundlePath)
	if err != nil {
		exit("candidate bundle is unavailable")
	}
	plan, err := verificationcontract.DecodePlan(planBytes)
	if err != nil {
		exit("assurance plan is invalid")
	}
	bundle, err := workercontract.DecodeBundle(bundleBytes)
	if err != nil {
		exit("candidate bundle is invalid")
	}
	engine, _ := verifierruntime.NewEngine(func() time.Time { return time.Now().UTC() })
	result, err := engine.Verify(context.Background(), plan, bundle, *workspace)
	if err != nil {
		exit("candidate verification failed")
	}
	resultBytes, _ := json.Marshal(result)
	if err := os.WriteFile(*resultPath, resultBytes, 0o644); err != nil {
		exit("verification result could not be written")
	}
}

func exit(message string) {
	_, _ = fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
