# Overview
### Description
rsync metrics exporter for **Prmoetheus.io**

### Requirements
- linux OS with tail + rsync installed
- Golang >= 1.12
- [OPTIONAL] Docker

### Code files
- `exporter.go`       => main code
- `exporter_test.go`  => unit [tests](#Tests)

### Design choices / Architecture
This utility is written in GO and uses `tail -f -n +1 [logfile]` STDOUTPUT to parse a live logfile.

The overall logic is:
- `tail` is invoked and its output is piped using `StdoutPipe()`
- A `bufio.Scanner` is attached to the pipe. It provides a convenient interface for reading data such as a file of newline-delimited lines of text. Created with `bufio.NewScanner` which by default assumes a `ScanLines` split function.
- A GO subroutine runs a WHILE-type loop that scans every new line in the file.
- The line is passed to a log parser that extracts meaningful data for the metrics using carefully tailored `strings.Split` functions
- Extracted values are passed to the relevant metric object using function `recordMetrics(metricType string, value int)`
- Metrics are exposed using `promhttp` handler
- For an explanation on metric choices, go to[ Metrics section](#Metrics)

### Defaults
- Metrics for `Prometheus.io` are exposed to HTTP `/metrics` endpoint on `PORT 2112`. To reach locally just issue `wget http://localhost:2112/metrics`
- Default log is `./rsync_example.log`

### Metrics
The following metrics are exposed to Prometheus.io
- "connections_to_rsync" => The total number of connections to rsync daemon
-	"rsync_executions" => The total number of rsync executions
-	"data_sent" => The total data sent (bytes)
-	"data_received" => The total data received (bytes)

All 4 metrics are `COUNTER` types as they can only increase and can be easily transformed in rates using `rate()` function in Prometheus to get interesting derived metrics, like connections/min or send-receive bytes/sec.

# Install
### Quickstart using run.sh shell script (Linux/MacOS)
Once the repository is in a `<folder>`:
1. `cd <folder>`
2. Launch the `run.sh [full path of rsync logfile to parse]` shell script with your shell. For example, using sh: 

    `sh run.sh [full path of rsync logfile to parse]`

### Building an executable in a Linux environment
Once the repository is in a `<folder>`:
1. `cd <folder>`
2. `go build -o exporter` (to avoid that the executable has the default naming `<folder>`) 
3. `exporter [PATH_TO_RSYNC_LOGFILE]` to start monitoring the rsync logfile

### Using Docker
There is a Dockerfile to build an Alpine based container with the compiled executable (called `app`).

Once the repository is in a `<folder>`:
1. `cd <folder>`
2. `docker build --rm -f "Dockerfile" -t <AN IMAGE NAME HERE>:latest .`
3. `docker run --rm -it -p 2112:2112 -v <PATH TO RSYNC LOG FILE IN HOST>:<MOUNT POINT IN THE CONTAINER>:ro <AN IMAGE NAME HERE>:latest sh`
4. Run the app with: `app [MOUNT POINT IN THE CONTAINER]`
5. Metrics will be exposed to `http://localhost:2112/metrics` or, if you use `docker-machine`, `http://DOCKER_MACHINE_IP:2112/metrics`

# Usage
Once the app is compiled and named `exporter`:

`exporter [OPTIONS] PATH_TO_RSYNC_LOGFILE`

### Options
`--help, -h`    Help: displays all available options

`-p`            TCP port to listen fo /metrics endpoint

# Tests
Some unit tests are provided.

To run tests:

`go test`

**IMPORTANT NOTICE IF USING DOCKER IMAGE:** Tests require `gcc`, so the provided Alpine based container is not a suitable option to execute them. Use the standard `golang:1.12` instead or run them in a linux environment with `gcc` and `Go 1.12`. For convenience there is a [Dockerfile for this purpose](./Dockerfile4test).

To use it: 
`docker build --rm -f "Dockerfile4test" -t <AN IMAGE NAME SPECIFIC FOR TEST HERE>:latest .`

# License
This package is released under the [GPLv3 license](./LICENSE)