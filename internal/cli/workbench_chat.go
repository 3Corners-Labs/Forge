package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// runWorkbenchChat is deliberately an API client, even when it starts an
// ephemeral local API. This keeps terminal, web, and future messaging adapters
// on one interaction contract.
func runWorkbenchChat(args []string, input io.Reader, output, errorOutput io.Writer, journeys journeyFactory) error {
	return runWorkbenchChatWithRequester(args, input, output, errorOutput, journeys, workbenchRequest)
}

type workbenchRequester func(string, string, any) (*workbenchRun, error)

func runWorkbenchChatWithRequester(args []string, input io.Reader, output, errorOutput io.Writer, journeys journeyFactory, requestRun workbenchRequester) error {
	flags := flag.NewFlagSet("workbench chat", flag.ContinueOnError)
	flags.SetOutput(errorOutput)
	apiURL := flags.String("api", "", "existing Forge Workbench API URL")
	if err := flags.Parse(args); err != nil {
		return err
	}
	var stop func()
	if *apiURL == "" {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return fmt.Errorf("start local Workbench API: %w", err)
		}
		server := &http.Server{Handler: newWorkbenchAPI(journeys), ReadHeaderTimeout: 5 * time.Second}
		go func() { _ = server.Serve(listener) }()
		*apiURL = "http://" + listener.Addr().String()
		stop = func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = server.Shutdown(ctx)
		}
		defer stop()
	}

	reader := bufio.NewReader(input)
	fmt.Fprintln(output, "FORGE WORKBENCH · PRODUCT TEAM TRIAL")
	fmt.Fprintln(output, "Free text becomes a proposal—never execution authority.")
	fmt.Fprintln(output)
	fmt.Fprint(output, "What should we build? ")
	message, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	message = strings.TrimSpace(message)
	if message == "" {
		return errors.New("a product outcome is required")
	}
	run, err := requestRun(http.MethodPost, *apiURL+"/api/v1/runs", map[string]string{"message": message})
	if err != nil {
		return err
	}
	printWorkbenchRun(output, run)
	for len(run.Gates) > 0 {
		gate := run.Gates[0]
		fmt.Fprintf(output, "\nDECISION · %s\n%s\nExact record: %s\n", gate.Label, gate.Summary, gate.Digest)
		fmt.Fprint(output, "Type approve to continue, or stop: ")
		answer, readErr := reader.ReadString('\n')
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return readErr
		}
		if strings.TrimSpace(strings.ToLower(answer)) != "approve" {
			fmt.Fprintln(output, "Stopped safely. No additional authority was granted.")
			return nil
		}
		run, err = requestRun(http.MethodPost, *apiURL+"/api/v1/runs/"+run.ID+"/actions", map[string]string{"action": gate.Action})
		if err != nil {
			return err
		}
		printWorkbenchRun(output, run)
	}
	fmt.Fprintln(output, "\nReviewable result ready. Merge, release, deployment, and production remain human-controlled.")
	return nil
}

func workbenchRequest(method, url string, body any) (*workbenchRun, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(method, url, bytes.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := (&http.Client{Timeout: 30 * time.Second}).Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		var failure struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.NewDecoder(response.Body).Decode(&failure)
		return nil, fmt.Errorf("Workbench API rejected request: %s", failure.Error.Message)
	}
	var run workbenchRun
	if err := json.NewDecoder(response.Body).Decode(&run); err != nil {
		return nil, err
	}
	return &run, nil
}

func printWorkbenchRun(output io.Writer, run *workbenchRun) {
	fmt.Fprintf(output, "\n%s · %d%%\n", strings.ToUpper(strings.ReplaceAll(run.State, "-", " ")), run.Progress)
	for _, agent := range run.Agents {
		marker := "○"
		switch agent.State {
		case "completed":
			marker = "✓"
		case "waiting":
			marker = "Ⅱ"
		}
		detail := agent.Summary
		if agent.WaitingOn != "" {
			detail = "waiting on " + agent.WaitingOn
		}
		fmt.Fprintf(output, "  %s %-26s %s\n", marker, agent.Label, detail)
	}
}
