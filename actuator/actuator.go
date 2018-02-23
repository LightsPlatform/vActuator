package actuator

import (
	"github.com/LightsPlatform/vActuator/stateManager"
	"log"
)

type Actuator struct {
	id     int
	Name   string `json:"name"`
	state  stateManager.Store
	quit   chan struct{}
	trap   chan struct{}
	config stateManager.Config
}

// New creates new actuator and store its user given script
func New(name string, script []byte, config stateManager.Config) (*Actuator, error) {
	// Store user script
	//path := os.TempDir() + "/actuator-%s.py"
	//f, err := os.Create(fmt.Sprintf(path, name))
	//if err != nil {
	//	return nil, err
	//}
	//f.Write(script)

	return &Actuator{
		Name:   name,
		config: config,

		quit:  make(chan struct{}, 0),
		trap:  make(chan struct{}, 0),
		state: stateManager.Init(config),
	}, nil
}

// Stop stops running actuator
func (a *Actuator) Stop() {
	a.quit <- struct{}{}

	close(a.quit)
	close(a.trap)
}

// Stop stops running actuator
func (a *Actuator) trigger() {
	a.trap <- struct{}{}
}

func (a *Actuator) Run() {

	for {
		select {
		case <-a.trap:
			//for i := 0; i < c; i++ {
			//	path := os.TempDir() + "/sensor-%s.py"
			//	cmd := exec.Command("runtime.py", fmt.Sprintf(path, s.Name))
			//
			//	// run
			//	value, err := cmd.Output()
			//	if err != nil {
			//		if err, ok := err.(*exec.ExitError); ok {
			//			log.Errorf("%s: %s", err.Error(), err.Stderr)
			//			continue
			//		}
			//	}
			//
			//	d := Data{
			//		Time: time.Now(),
			//	}
			//
			//	if err := json.Unmarshal(value, &d.Value); err == nil {
			//		log.Infoln(d)
			//		s.Buffer <- d
			//	} else {
			//		log.Errorf("%s", err)
			//	}
			//
			//}
			log.Println("trap")
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
	actuator.trigger()
	actuator.trigger()
	actuator.trigger()
	actuator.Stop()
}
