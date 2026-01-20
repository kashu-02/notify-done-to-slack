package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RequestBody struct {
	Text string `json:"text"`
}

type threadSafeBuffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *threadSafeBuffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}

func (b *threadSafeBuffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.String()
}

func NotifyDoneToSlack(cmd *cobra.Command, args []string) error {
	webhookURL := viper.GetString("webhook-url")
	if webhookURL == "" {
		return fmt.Errorf("please configure the Slack Webhook URL at configure command")
	}

	headNumber, _ := cmd.Flags().GetInt32("head")
	tailNumber, _ := cmd.Flags().GetInt32("tail")

	var lines []string
	var err error

	if len(args) > 0 {
		lines, err = runCommand(cmd, args)
	} else {
		lines, err = readStdIn(cmd)
	}

	if err != nil {
		return err
	}

	if headNumber > 0 {
		lines = printHead(lines, int(headNumber))
	} else if tailNumber > 0 {
		lines = printTail(lines, int(tailNumber))
	}

	result := strings.Join(lines, "")

	requestBody := &RequestBody{
		Text: "Command done.\nResult:```\n" + result + "```",
	}

	err = sendToSlack(webhookURL, requestBody)
	if err != nil {
		return err
	}

	return nil
}

func runCommand(cmd *cobra.Command, args []string) ([]string, error) {
	var outputBuf threadSafeBuffer
	cmdName := args[0]
	cmdArgs := args[1:]

	c := exec.Command(cmdName, cmdArgs...)
	c.Stdin = cmd.InOrStdin()
	c.Stdout = io.MultiWriter(os.Stdout, &outputBuf)
	c.Stderr = io.MultiWriter(os.Stderr, &outputBuf)

	err := c.Run()

	var lines []string
	fullOutput := outputBuf.String()
	if len(fullOutput) > 0 {
		lines = strings.SplitAfter(fullOutput, "\n")
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
	}

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return lines, nil
		}
		return lines, err
	}
	return lines, nil
}

func readStdIn(cmd *cobra.Command) ([]string, error) {
	var lines []string
	var inputReader io.Reader = cmd.InOrStdin()
	reader := bufio.NewReader(inputReader)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func printHead(lines []string, n int) []string {
	var result []string
	for i := 0; i < n && i < len(lines); i++ {
		result = append(result, lines[i])
	}
	return result
}

func printTail(lines []string, n int) []string {
	var result []string
	start := len(lines) - n
	if start < 0 {
		start = 0
	}
	for i := start; i < len(lines); i++ {
		result = append(result, lines[i])
	}
	return result
}

func sendToSlack(webhookURL string, requestBody *RequestBody) error {
	requestBodyJSON, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(
		"POST",
		webhookURL,
		bytes.NewBuffer(requestBodyJSON),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}