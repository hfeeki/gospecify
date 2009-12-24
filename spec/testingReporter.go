/*
Copyright (c) 2009 Samuel Tesla <samuel.tesla@gmail.com>

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
package main

import (
	"container/list";
	"os";

	t "../src/testspecify";
)

type TestingReporter interface {
	t.Reporter;
	FailingExamples() int;
	PassingExamples() int;
	PendingExamples() int;
	EachFailure() <-chan os.Error;
	EachPending() <-chan string;
}

func testRun(name string, block func(t.Runner)) (reporter TestingReporter) {
	runner := t.NewRunner();
	runner.Describe(name, func() { block(runner) });
	reporter = newTestingReporter();
	runner.Run(reporter);
	return;
}

func newTestingReporter() (result *testingReporter) {
	result = &testingReporter{};
	result.failures = list.New();
	result.pending = list.New();
	return;
}

type testingReporter struct {
	passing		int;
	failures	*list.List;
	pending		*list.List;
}

func (self *testingReporter) Fail(err os.Error) {
	self.failures.PushBack(err)
}
func (self *testingReporter) Finish()	{}
func (self *testingReporter) Pass()	{ self.passing++ }
func (self *testingReporter) Pending(name string) {
	self.pending.PushBack(name)
}

func (self *testingReporter) FailingExamples() int {
	return self.failures.Len()
}
func (self *testingReporter) PassingExamples() int {
	return self.passing
}
func (self *testingReporter) PendingExamples() int {
	return self.pending.Len()
}
func (self *testingReporter) EachFailure() <-chan os.Error {
	ch := make(chan os.Error, self.failures.Len());
	for val := range self.failures.Iter() {
		if err, ok := val.(os.Error); !ok {
			panic("typecast error");
		} else {
			ch <- err;
		}
	}
	close(ch);
	return ch;
}
func (self *testingReporter) EachPending() <-chan string {
	ch := make(chan string, self.pending.Len());
	for val := range self.pending.Iter() {
		if name, ok := val.(string); !ok {
			panic("typecast error");
		} else {
			ch <- name;
		}
	}
	close(ch);
	return ch;
}
