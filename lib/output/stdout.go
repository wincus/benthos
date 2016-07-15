/*
Copyright (c) 2014 Ashley Jeffs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package output

import (
	"bytes"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/jeffail/benthos/lib/types"
	"github.com/jeffail/util/log"
	"github.com/jeffail/util/metrics"
)

//--------------------------------------------------------------------------------------------------

func init() {
	constructors["stdout"] = typeSpec{
		constructor: NewSTDOUT,
		description: "TODO",
	}
}

//--------------------------------------------------------------------------------------------------

// STDOUT - An output type that pushes messages to STDOUT.
type STDOUT struct {
	running int32

	conf Config
	log  log.Modular

	messages     <-chan types.Message
	responseChan chan types.Response

	closedChan chan struct{}
}

// NewSTDOUT - Create a new STDOUT output type.
func NewSTDOUT(conf Config, log log.Modular, stats metrics.Aggregator) (Type, error) {
	s := STDOUT{
		running:      1,
		conf:         conf,
		log:          log.NewModule(".output.stdout"),
		messages:     nil,
		responseChan: make(chan types.Response),
		closedChan:   make(chan struct{}),
	}

	return &s, nil
}

//--------------------------------------------------------------------------------------------------

// loop - Internal loop brokers incoming messages to output pipe.
func (s *STDOUT) loop() {
	defer func() {
		close(s.responseChan)
		close(s.closedChan)
	}()

	s.log.Infoln("Sending messages through STDOUT")

	for atomic.LoadInt32(&s.running) == 1 {
		msg, open := <-s.messages
		// If the messages chan is closed we do not close ourselves as it can replaced.
		if !open {
			return
		}
		var err error
		if len(msg.Parts) == 1 {
			_, err = fmt.Fprintf(os.Stdout, "%s\n", msg.Parts[0])
		} else {
			_, err = fmt.Fprintf(os.Stdout, "%s\n\n", bytes.Join(msg.Parts, []byte("\n")))
		}
		s.responseChan <- types.NewSimpleResponse(err)
	}
}

// StartReceiving - Assigns a messages channel for the output to read.
func (s *STDOUT) StartReceiving(msgs <-chan types.Message) error {
	if s.messages != nil {
		return types.ErrAlreadyStarted
	}
	s.messages = msgs
	go s.loop()
	return nil
}

// ResponseChan - Returns the errors channel.
func (s *STDOUT) ResponseChan() <-chan types.Response {
	return s.responseChan
}

// CloseAsync - Shuts down the STDOUT output and stops processing messages.
func (s *STDOUT) CloseAsync() {
	atomic.StoreInt32(&s.running, 0)
}

// WaitForClose - Blocks until the STDOUT output has closed down.
func (s *STDOUT) WaitForClose(timeout time.Duration) error {
	select {
	case <-s.closedChan:
	case <-time.After(timeout):
		return types.ErrTimeout
	}
	return nil
}

//--------------------------------------------------------------------------------------------------
