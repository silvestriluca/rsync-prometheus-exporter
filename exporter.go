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

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// GLOBAL VARIABLES
var (
	portPtr    *int   //TCP port to listen fo /metrics endpoint
	file2Parse string //Log file that has to be parsed
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

func init() {
	portPtr = flag.Int("p", 2112, "TCP port to listen fo /metrics endpoint")
	flag.Parse()
	if len(flag.Args()) >= 1 {
		file2Parse = flag.Arg(0)
	} else {
		file2Parse = "./rsync_example.log"
	}
	fmt.Println("Value of port declared: ", *portPtr)
	fmt.Println("File to parse: ", file2Parse)
	startMessage()
}

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
