/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"

	mlog "mosn.io/pkg/log"

	//"github.com/autsu/holmes1/reporters/http_reporter"
	"math/rand"
	"net/http"
	"time"

	"github.com/autsu/holmes1"
)

// run `curl http://localhost:10024/rand` after 15s(warn up)
// run `curl http://localhost:10024/spike` after 15s(warn up)
func init() {
	http.HandleFunc("/rand", randAlloc)
	http.HandleFunc("/spike", spikeAlloc)
	go http.ListenAndServe(":10024", nil)
}

func main() {
	// reporter := http_reporter.NewReporter("TOKEN", "URL")
	h, _ := holmes.New(
		holmes.WithDumpPath("./tmp"),
		holmes.WithLogger(holmes.NewFileLog("./tmp/holmes.log", mlog.DEBUG)),
		holmes.WithBinaryDump(),
		holmes.WithMemoryLimit(100*1024*1024), // 100MB
		holmes.WithGCHeapDump(10, 20, 40, time.Minute),
		// holmes.WithProfileReporter(reporter),
	)
	h.EnableGCHeapDump().Start()
	time.Sleep(time.Hour)
}

var (
	base = make([]byte, 1024*1024*10) // 10 MB long live memory.
)

func randAlloc(wr http.ResponseWriter, req *http.Request) {
	var s = make([][]byte, 0) // short live
	for i := 0; i < 1024; i++ {
		len := rand.Intn(1024)
		bytes := make([]byte, len)

		s = append(s, bytes)

		if len == 0 {
			s = make([][]byte, 0)
		}
	}
	time.Sleep(time.Millisecond * 10)
	fmt.Fprintf(wr, "slice current length: %v\n", len(s))
}

func spikeAlloc(wr http.ResponseWriter, req *http.Request) {
	var s = make([][]byte, 0, 1024) // spike, 10MB
	for i := 0; i < 10; i++ {
		bytes := make([]byte, 1024*1024)
		s = append(s, bytes)
	}
	// live for a while
	time.Sleep(time.Millisecond * 500)
	fmt.Fprintf(wr, "spike slice length: %v\n", len(s))
}
