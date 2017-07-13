// Copyright 2017 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package log

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	mt "github.com/mendersoftware/mendertesting"
	"github.com/sirupsen/logrus"
)

var logrus_quotes string = ""

func TestSetup(t *testing.T) {
}

func TestLogging(t *testing.T) {
	checkLogging(t, "log_test")
}

func setupLogging(t *testing.T) {
	fd, err := os.Create("output.log")
	mt.AssertTrue(t, err == nil)
	SetOutput(fd)

	// How you run the test influences whether the output is connected to a
	// terminal or not (-v gives you a terminal, otherwise not). So disable
	// that moving target for this test.
	Log.Formatter.(*logrus.TextFormatter).DisableColors = true

	// Also disable timestamps, which are not predictable.
	Log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true

	SetLevel(logrus.DebugLevel)

	// Older versions of logrus add quotes to fields with underscores, while
	// newer ones do not. We cannot avoid underscores because of the name of
	// our test file, so we need to account for that in the string
	// comparisons.
	Log.WithField("test_field", "test_value").Println("Testing fields")
	fd.Close()
	fd, err = os.Open("output.log")
	mt.AssertTrue(t, err == nil)
	var buf [4096]byte
	n, err := fd.Read(buf[:])
	mt.AssertTrue(t, err == nil)
	mt.AssertTrue(t, n < 4096)

	if string(buf[0:n]) == "level=info msg=\"Testing fields\" test_field=\"test_value\" \n" {
		logrus_quotes = "\""
	} else if string(buf[0:n]) == "level=info msg=\"Testing fields\" test_field=test_value \n" {
		logrus_quotes = ""
	} else {
		t.Fatalf("Unexpected log format: %s", string(buf[0:n]))
	}

	fd, err = os.Create("output.log")
	mt.AssertTrue(t, err == nil)
	SetOutput(fd)
}

func checkLogging(t *testing.T, module string) {
	setupLogging(t)

	Print("Print")
	Printf("%s : %s", "Printf", "Printf")
	Println("Println")
	Debug("Debug")
	Debugf("%s : %s", "Debugf", "Debugf")
	Debugln("Debugln")
	Info("Info")
	Infof("%s : %s", "Infof", "Infof")
	Infoln("Infoln")
	Warn("Warn")
	Warnf("%s : %s", "Warnf", "Warnf")
	Warnln("Warnln")
	Warning("Warning")
	Warningf("%s : %s", "Warningf", "Warningf")
	Warningln("Warningln")
	Error("Error")
	Errorf("%s : %s", "Errorf", "Errorf")
	Errorln("Errorln")
	func() {
		defer func() {
			recover()
		}()
		Panic("Panic")
	}()
	func() {
		defer func() {
			recover()
		}()
		Panicf("%s : %s", "Panicf", "Panicf")
	}()
	func() {
		defer func() {
			recover()
		}()
		Panicln("Panicln")
	}()

	Log.Print("Print")
	Log.Printf("%s : %s", "Printf", "Printf")
	Log.Println("Println")
	Log.Debug("Debug")
	Log.Debugf("%s : %s", "Debugf", "Debugf")
	Log.Debugln("Debugln")
	Log.Info("Info")
	Log.Infof("%s : %s", "Infof", "Infof")
	Log.Infoln("Infoln")
	Log.Warn("Warn")
	Log.Warnf("%s : %s", "Warnf", "Warnf")
	Log.Warnln("Warnln")
	Log.Warning("Warning")
	Log.Warningf("%s : %s", "Warningf", "Warningf")
	Log.Warningln("Warningln")
	Log.Error("Error")
	Log.Errorf("%s : %s", "Errorf", "Errorf")
	Log.Errorln("Errorln")
	func() {
		defer func() {
			recover()
		}()
		Log.Panic("Panic")
	}()
	func() {
		defer func() {
			recover()
		}()
		Log.Panicf("%s : %s", "Panicf", "Panicf")
	}()
	func() {
		defer func() {
			recover()
		}()
		Log.Panicln("Panicln")
	}()

	// Would have tested Fatal() calls here, but it's tricky because they
	// call os.Exit(). Skipping those tests for now, they're not likely to
	// be used much.

	checkString := "level=info msg=Print module=%% \n" +
		"level=info msg=\"Printf : Printf\" module=%% \n" +
		"level=info msg=Println module=%% \n" +
		"level=debug msg=Debug module=%% \n" +
		"level=debug msg=\"Debugf : Debugf\" module=%% \n" +
		"level=debug msg=Debugln module=%% \n" +
		"level=info msg=Info module=%% \n" +
		"level=info msg=\"Infof : Infof\" module=%% \n" +
		"level=info msg=Infoln module=%% \n" +
		"level=warning msg=Warn module=%% \n" +
		"level=warning msg=\"Warnf : Warnf\" module=%% \n" +
		"level=warning msg=Warnln module=%% \n" +
		"level=warning msg=Warning module=%% \n" +
		"level=warning msg=\"Warningf : Warningf\" module=%% \n" +
		"level=warning msg=Warningln module=%% \n" +
		"level=error msg=Error module=%% \n" +
		"level=error msg=\"Errorf : Errorf\" module=%% \n" +
		"level=error msg=Errorln module=%% \n" +
		"level=panic msg=Panic module=%% \n" +
		"level=panic msg=\"Panicf : Panicf\" module=%% \n" +
		"level=panic msg=Panicln module=%% \n" +
		"level=info msg=Print module=%% \n" +
		"level=info msg=\"Printf : Printf\" module=%% \n" +
		"level=info msg=Println module=%% \n" +
		"level=debug msg=Debug module=%% \n" +
		"level=debug msg=\"Debugf : Debugf\" module=%% \n" +
		"level=debug msg=Debugln module=%% \n" +
		"level=info msg=Info module=%% \n" +
		"level=info msg=\"Infof : Infof\" module=%% \n" +
		"level=info msg=Infoln module=%% \n" +
		"level=warning msg=Warn module=%% \n" +
		"level=warning msg=\"Warnf : Warnf\" module=%% \n" +
		"level=warning msg=Warnln module=%% \n" +
		"level=warning msg=Warning module=%% \n" +
		"level=warning msg=\"Warningf : Warningf\" module=%% \n" +
		"level=warning msg=Warningln module=%% \n" +
		"level=error msg=Error module=%% \n" +
		"level=error msg=\"Errorf : Errorf\" module=%% \n" +
		"level=error msg=Errorln module=%% \n" +
		"level=panic msg=Panic module=%% \n" +
		"level=panic msg=\"Panicf : Panicf\" module=%% \n" +
		"level=panic msg=Panicln module=%% \n"
	checkString = strings.Replace(checkString, "%%", logrus_quotes+module+logrus_quotes, -1)

	verifyLogging(t, checkString)
	cleanupLogging(t)
}

func verifyLogging(t *testing.T, checkString string) {
	fd, err := os.Open("output.log")
	mt.AssertTrue(t, err == nil)
	defer fd.Close()
	var buf [4096]byte
	n, err := fd.Read(buf[:])
	mt.AssertTrue(t, err == nil)
	mt.AssertTrue(t, n < 4096)

	mt.AssertStringEqual(t, string(buf[0:n]), checkString)
}

func cleanupLogging(t *testing.T) {
	Log.Formatter.(*logrus.TextFormatter).DisableColors = false
	Log.Formatter.(*logrus.TextFormatter).DisableTimestamp = false

	Log.Out.(*os.File).Close()

	os.Remove("output.log")
}

func TestModules(t *testing.T) {
	PushModule("test_1")
	checkLogging(t, "test_1")
	PushModule("test_2")
	checkLogging(t, "test_2")
	PopModule()
	checkLogging(t, "test_1")
	PopModule()
	checkLogging(t, "log_test")
}

func TestSyslog(t *testing.T) {
	setupLogging(t)

	if AddSyslogHook() != nil {
		// If we cannot connect to syslog we have no choice but to skip
		// the test. It's perfectly legitimate that it's not available.
		cleanupLogging(t)
		t.Skip("Skip syslog test because syslog is not available")
	}

	Log.Formatter.(*logrus.TextFormatter).ForceColors = true

	rand.Seed(time.Now().UTC().UnixNano())

	// In order to not get false passes because the syslog has entries from
	// the previous run.
	testrand := rand.Int()

	SetLevel(DebugLevel)

	Log.Errorf("For syslog testing: Error with no module: %d", testrand)
	Log.PushModule("test_1")
	Log.Warnf("For syslog testing: Warning with test1 module: %d", testrand)
	Log.Debugf("For syslog testing: Debug with test1 module: %d", testrand)
	Log.PopModule()

	syslog := "/var/log/syslog"
	if _, err := os.Stat(syslog); err != nil {
		syslog = "/var/log/messages"
		if _, err = os.Stat(syslog); err != nil {
			t.Fatal("Could not locate syslog, cannot test")
		}
	}
	cmd := exec.Command("tail", syslog)
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Error tailing syslog: %s", err.Error())
	}

	// Make sure there are no colors in the syslog, even if we forced them
	// for the console.
	var checkString string
	// Should show.
	checkString = fmt.Sprintf("level=error msg=\"For syslog testing: Error with no module: "+
		"%d\" module=%slog_test%s", testrand, logrus_quotes, logrus_quotes)
	mt.AssertTrue(t, strings.Index(string(output), checkString) >= 0)
	checkString = fmt.Sprintf("level=warning msg=\"For syslog testing: Warning with test1 module: "+
		"%d\" module=%stest_1%s", testrand, logrus_quotes, logrus_quotes)
	mt.AssertTrue(t, strings.Index(string(output), checkString) >= 0)
	// Should not show.
	checkString = fmt.Sprintf("level=debug msg=\"For syslog testing: Debug with test1 module: "+
		"%d\" module=%stest_1%s", testrand, logrus_quotes, logrus_quotes)
	mt.AssertTrue(t, strings.Index(string(output), checkString) < 0)

	cleanupLogging(t)

	// Get rid of the syslog hook again.
	Log = New()
}

func TestModuleFilter(t *testing.T) {
	setupLogging(t)

	SetModuleFilter([]string{"test"})
	Debug("Should not show")
	SetModuleFilter([]string{"log_test"})
	Debug("Should show")
	SetModuleFilter([]string{"test", "log_test"})
	Debug("Should also show")
	PushModule("test")
	SetModuleFilter([]string{"test", "log_test"})
	Debug("Should show as well")
	PushModule("test2")
	Debug("Should not show again")
	PopModule()
	Debug("Should show after module reappeared")
	PopModule()
	Debug("Should show after file reappeared")

	checkString := "level=debug msg=\"Should show\" module=%% \n" +
		"level=debug msg=\"Should also show\" module=%% \n" +
		"level=debug msg=\"Should show as well\" module=test \n" +
		"level=debug msg=\"Should show after module reappeared\" module=test \n" +
		"level=debug msg=\"Should show after file reappeared\" module=%% \n"
	checkString = strings.Replace(checkString, "%%", logrus_quotes+"log_test"+logrus_quotes, -1)

	verifyLogging(t, checkString)
	cleanupLogging(t)
}

func TestLogLevels(t *testing.T) {
	setupLogging(t)

	SetLevel(logrus.DebugLevel)
	Debug("Debug log level should show")
	SetLevel(logrus.InfoLevel)
	Debug("Debug log level should not show")
	Info("Info log level should show")
	SetLevel(logrus.WarnLevel)
	Debug("Debug log level should not show")
	Info("Info log level should not show")
	Debug("Debug log level should not show")
	Warn("Warn log level should show")

	checkString := "level=debug msg=\"Debug log level should show\" module=%% \n" +
		"level=info msg=\"Info log level should show\" module=%% \n" +
		"level=warning msg=\"Warn log level should show\" module=%% \n"
	checkString = strings.Replace(checkString, "%%", logrus_quotes+"log_test"+logrus_quotes, -1)

	verifyLogging(t, checkString)
	cleanupLogging(t)
}
