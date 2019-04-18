/*
  <exporter.go>
  Copyright (C) <2019>  <Luca Silvestri>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// GLOBAL CONSTANTS
const (
	//Metrics
	Connections  = "connections_to_rsync"
	Executions   = "rsync_executions"
	DataSent     = "data_sent"
	DataReceived = "data_received"
)

// GLOBAL VARIABLES
//1. Command line related variables
var (
	portPtr    *int   //TCP port to listen fo /metrics endpoint
	file2Parse string //Log file that has to be parsed
)

//2. Metrics related variables
var (
	connectionsMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: Connections,
		Help: "The total number of connections to rsync daemon",
	})
	executionsMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: Executions,
		Help: "The total number of rsync executions",
	})
	dataSentMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: DataSent,
		Help: "The total data sent (bytes)",
	})
	dataReceivedMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: DataReceived,
		Help: "The total data received (bytes)",
	})
)

//Setup a http endpoint to expose metrics to Prometheus.io
func setupHTTPListener() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil)
}

//Outputs the start message
func startMessage() {
	fmt.Println("")
	fmt.Println("*****************************************************************")
	fmt.Println("Copyright (C) 2019  Luca Silvestri")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY.")
	fmt.Println("This is free software, and you are welcome to redistribute it")
	fmt.Println("under GPLv3 License or above")
	fmt.Println("****************************************************************")
}

//Analyze the log line looking for relevant metrics
func parseLogLine(logLine string) (parsed string) {
	//Takes of the timestamp + PID from the log line
	parsablePart := strings.Split(logLine, "] ")[1]

	/*
		//Extract the timestamp and converts it in a date object
		datePart := strings.Split(logLine, " [")[0]
		t, err := time.Parse("2006/01/02 15:04:05", datePart)
		if err != nil {
			fmt.Println(err)
		}
	*/

	//Parse the event and extract relevant metrics
	words := strings.Split(parsablePart, " ")
	switch words[0] {
	case "connect":
		recordMetrics(Connections, 1)
	case "rsync":
		recordMetrics(Executions, 1)
	case "sent":
		bytesSent, err1 := strconv.Atoi(words[1])
		if err1 != nil {
			fmt.Println("Error while converting sent data value")
		} else {
			recordMetrics(DataSent, bytesSent)
		}
		bytesReceived, err2 := strconv.Atoi(words[5])
		if err2 != nil {
			fmt.Println("Error while converting received data value")
		} else {
			recordMetrics(DataReceived, bytesReceived)
		}
	}

	//Returns the parsable part of the log line
	return parsablePart
}

//Updates the metrics exposed to Prometheus
func recordMetrics(metricType string, value int) {
	//Converts value to float64
	floatValue := float64(value)

	//Updates the metrics acording to the type
	switch metricType {
	case Connections:
		connectionsMetric.Add(floatValue)
	case Executions:
		executionsMetric.Add(floatValue)
	case DataSent:
		dataSentMetric.Add(floatValue)
	case DataReceived:
		dataReceivedMetric.Add(floatValue)
	}
}

//Init function
func init() {
	//Prints the start message
	startMessage()
	//Evaluates the cli parameters and assign default values:
	//Port
	portPtr = flag.Int("p", 2112, "TCP port to listen fo /metrics endpoint")
	//File to parse
	flag.Parse()
	if len(flag.Args()) >= 1 {
		file2Parse = flag.Arg(0)
	} else {
		file2Parse = "./rsync_example.log"
	}
	//Outputs the cli parameters values
	fmt.Println("Value of port declared: ", *portPtr)
	fmt.Println("File to parse: ", file2Parse)
	//
	fmt.Println("****************************************************************")
}

//Main procedure
func main() {
	//Launch HTTP Server
	go setupHTTPListener()

	//Defines a command that invokes "tail"
	cmd := exec.Command("tail", "-f", "-n", "+1", file2Parse)

	//Create a pipe for the output of "tail" command
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for tail", err)
		return
	}

	//Create a scanner on the IO Buffer
	scanner := bufio.NewScanner(cmdReader)

	//Launch a goroutine that scans the pipeline LIVE (line by line) until it gets closed
	go func() {
		i := 0
		for scanner.Scan() {
			//Outputs a single line (LIVE)
			line := scanner.Text()
			fmt.Printf("Line %d", i)
			fmt.Printf("\t > %s\n", line)
			i = i + 1
			//Parse the single line
			parseLogLine(line)
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
