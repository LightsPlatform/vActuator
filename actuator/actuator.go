package actuator

import (
	"github.com/LightsPlatform/vActuator/stateManager"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"fmt"
	"encoding/json"
)

type Actuator struct {
	id     int
	Name   string `json:"name"`
	State  stateManager.Store
	quit   chan struct{}
	trap   chan string
	triggerResult chan bool
	config stateManager.Config
}

// New creates new actuator and store its user given script
func New(name string, script []byte, config stateManager.Config) (*Actuator, error) {
	// Store user script
	path := os.TempDir() + "/actuator-%s.py"
	f, err := os.Create(fmt.Sprintf(path, name))
	if err != nil {
		return nil, err
	}
	f.Write(script)

	return &Actuator{
		Name:   name,
		config: config,

		quit:  make(chan struct{}, 0),
		trap:  make(chan string,0),
		triggerResult:  make(chan bool, 0),
		State: stateManager.Init(config),
	}, nil
}

// Stop stops running actuator
func (a *Actuator) Stop() {
	a.quit <- struct{}{}

	close(a.quit)
	close(a.trap)
	close(a.triggerResult)
}

// Stop stops running actuator
func (a *Actuator) Trigger(action string) bool{
	a.trap <- action
	for{
		select {
		case r := <-a.triggerResult:
			return r;
		}
	}
	return true
}

func (a *Actuator) Run() {

	for {
		select {
		case action := <-a.trap:

			path := os.TempDir() + "/actuator-%s.py"
			b, err := json.Marshal(a.State)
			if(err != nil){
				log.Errorf("Can't pass args")
				a.triggerResult <- false
				continue
			}
			cmd := exec.Command("runtime.py", fmt.Sprintf(path, a.Name),string(b),action)
			// run
			value, err := cmd.Output()
			if err != nil {
				if err, ok := err.(*exec.ExitError); ok {
					log.Errorf("%s: %s", err.Error(), err.Stderr)
					a.triggerResult <- false
					continue
				}
			}

			if err := json.Unmarshal(value, &a.State); err == nil {
				log.Infoln(a.State)
			} else {
				log.Errorf("%s", err)
				a.triggerResult <- false
				continue
			}

			log.Println("trap")
			a.triggerResult <- true
		case <-a.quit:
			log.Println("quit")
			return
		}
	}
}

func main() {
	config := stateManager.Config{map[string][]string{
		"temp":   []string{"int", "2"},
		"light":  []string{"float", "1.2"},
		"light2": []string{"float", "1.2"},
	}}
	actuator, error := New("test", nil, config)
	if error != nil {
		log.Println(error)
	}
	go actuator.Run()
	//actuator.Trigger()
	//actuator.Trigger()
	//actuator.Trigger()
	//actuator.Stop()
}
