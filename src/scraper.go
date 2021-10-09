// Copyright 2021 Edoardo Riggio
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"fmt"
	"time"
	"strings"
	
	"github.com/joho/godotenv"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func checkWebsite() bool {
	godotenv.Load("../.env")

	platform := os.Getenv("PLATFORM")
	driverPath := "../bin/chromedriver"

	if strings.Compare(platform, "pi") == 0 {
		driverPath = "/usr/lib/chromium-browser/chromedriver"
	}

	const (
		seleniumPath = "../bin/selenium-server-standalone-3.141.59.jar"
		port         = 8080
	)

	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(driverPath),
		selenium.Output(os.Stderr),
	}

	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)

	if err != nil {
		panic(err)
	}

	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}

	chromeCaps := chrome.Capabilities {
		Args: []string{
			"--headless",
			"--no-sandbox",
		},
	}

	caps.AddChrome(chromeCaps)

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	if err != nil {
		panic(err)
	}

	defer wd.Quit()

	// Navigate to website
	if err := wd.Get("https://usi.qualtrics.com/jfe/form/SV_9t5eb0E2ddzd1Qi"); err != nil {
		panic(err)
	}

	// Page 1
	butt1, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"QID56\"]/div[3]/div/fieldset/div/table/tbody/tr/td[1]")
	next1, err1 := wd.FindElement(selenium.ByXPATH, "//*[@id=\"NextButton\"]")

	butt1.Click()
	next1.Click()

	time.Sleep(time.Millisecond * 1000)

	if err != nil || err1 != nil {
		panic(err)
	}

	// Page 2
	butt2, err2 := wd.FindElement(selenium.ByXPATH, "//*[@id=\"QID61\"]/div[3]/div/fieldset/div/table/tbody/tr/td[1]")
	next2, err3 := wd.FindElement(selenium.ByXPATH, "//*[@id=\"NextButton\"]")

	butt2.Click()
	next2.Click()

	time.Sleep(time.Millisecond * 1000)

	if err2 != nil || err3 != nil {
		panic(err)
	}

	// Page 3
	input1, err4 := wd.FindElements(selenium.ByTagName, "input")
	next3, err5 := wd.FindElement(selenium.ByXPATH, "//*[@id=\"NextButton\"]")

	input1[0].SendKeys("Edoardo")
	input1[1].SendKeys("Riggio")
	input1[2].SendKeys("M")
	input1[3].SendKeys("15/05/2000")
	input1[4].SendKeys("edoardo.riggio@usi.ch")
	next3.Click()

	time.Sleep(time.Millisecond * 1000)

	if err4 != nil || err5 != nil {
		panic(err)
	}

	// Page 4
	butt3, err6 := wd.FindElement(selenium.ByXPATH, "//*[@id=\"QID68\"]/div[3]/div/fieldset/div/table/tbody/tr/td[1]")
	next4, err7 := wd.FindElement(selenium.ByXPATH, "//*[@id=\"NextButton\"]")

	butt3.Click()
	next4.Click()

	time.Sleep(time.Millisecond * 1000)

	if err6 != nil || err7 != nil {
		panic(err)
	}

	// Page 5
	butt4, err8 := wd.FindElement(selenium.ByXPATH, "//*[@id=\"QID64\"]/div[3]/div/fieldset/div/table/tbody/tr/td[2]")
	next5, err9 := wd.FindElement(selenium.ByXPATH, "//*[@id=\"NextButton\"]")

	butt4.Click()
	next5.Click()

	time.Sleep(time.Millisecond * 1000)

	if err8 != nil || err9 != nil {
		panic(err)
	}

	// Final Page
	fields, err10 := wd.FindElements(selenium.ByTagName, "option")

	if err10 != nil {
		panic(err)
	}

	return len(fields)-7 == 3
}