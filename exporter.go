package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//Setup a http endpoint to expose metrics to Prometheus.io
func setupHTTPListener() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func main() {
	//Launch HTTP Server
	go setupHTTPListener()

	//Defines a command that invokes "tail"
	cmd := exec.Command("tail", "-f", "-n", "+1", "./rsync_example.log")

	//Create a pipe for the output of "tail" command
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for tail", err)
		return
	}

	//Create a scanner on the IO Buffer
	scanner := bufio.NewScanner(cmdReader)

	//Launch a goroutine that scans the pipeline until it gets closed
	go func() {
		i := 0
		for scanner.Scan() {
			fmt.Printf("Line %d", i)
			fmt.Printf("\t > %s\n", scanner.Text())
			i = i + 1
		}
	}()

	//Starts the "tail" command
	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting tail", err)
		return
	}

	//Waits for the "tail" command to end
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for tail", err)
		return
	}

}
