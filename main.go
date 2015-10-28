package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/lafikl/gochrome"
)

var (
	errChromeSocketConnectionFailed = errors.New("CHROME_SOCKET_CONNECTION_FAILED")
	errCouldNotStartChrome          = errors.New("COULD_NOT_START_CHROME")
)

const chrome = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
const port = 9222
const dataDir = "./data"

func main() {
	cmd := mustStartChrome()
	tc := time.Tick(3 * time.Second)

	ccc := mustConnectToChrome()

	ccc.navigate(navigateMsg{"http://pol.dk"})

	msgChan := make(chan gochrome.Message)

	for {
		select {
		case m := <-msgChan:
			str, _ := json.Marshal(m)
			fmt.Println(string(str))
		case <-tc:
			fmt.Println(cmd.ProcessState.String())
		}

	}
}

func mustStartChrome() *exec.Cmd {
	cmd := exec.Command(
		chrome,
		fmt.Sprintf("--remote-debugging-port=%v", port),
		fmt.Sprintf("--user-data-dir=%s", dataDir),
	)

	err := cmd.Start()

	if err != nil {
		panic(errCouldNotStartChrome)
	}

	fmt.Printf("Started Chrome - pid: %v\n", cmd.Process.Pid)
	return cmd
}

func mustConnectToChrome() *cc {
	try := func() (*gochrome.Chrome, error) {
		return gochrome.New(fmt.Sprintf("http://localhost:%v", port), 0)
	}

	for i := 0; i < 10; i++ {
		gc, err := try()

		if err == nil {
			fmt.Printf("Connected to Chrome over web socket")
			return &cc{gc}
		}
		fmt.Printf("Attempting to connect - retry number %v\n", i+1)
		time.Sleep(200 * time.Millisecond)
	}

	panic(errChromeSocketConnectionFailed)
}

type cc struct {
	cl *gochrome.Chrome
}

type arg struct {
	name string
	val  interface{}
}

type navigateMsg struct {
	URL string `json:"url"`
}

func (c cc) navigate(msg navigateMsg) int {
	id, com := comm("Page.navigate", arg{"url", msg.URL})
	c.cl.Send(com)
	return id
}

var i int

func comm(method string, params ...arg) (int, gochrome.Command) {
	m := map[string]interface{}{}
	for _, a := range params {
		m[a.name] = a.val
	}
	i++
	return i, gochrome.Command{Id: i, Method: method, Params: m}
}
