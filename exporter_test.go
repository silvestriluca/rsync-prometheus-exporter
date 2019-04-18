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
	"testing"
)

func Test_setupHTTPListener(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupHTTPListener()
		})
	}
}

func Test_startMessage(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startMessage()
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_parseLogLine(t *testing.T) {
	type args struct {
		logLine string
	}
	tests := []struct {
		name       string
		args       args
		wantParsed string
	}{
		{"Log line", args{"2017/09/01 03:10:18 [16551] sent 6092862 bytes  received 55750 bytes  total size 302215840"}, "sent 6092862 bytes  received 55750 bytes  total size 302215840"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotParsed := parseLogLine(tt.args.logLine); gotParsed != tt.wantParsed {
				t.Errorf("parseLogLine() = %v, want %v", gotParsed, tt.wantParsed)
			}
		})
	}
}
