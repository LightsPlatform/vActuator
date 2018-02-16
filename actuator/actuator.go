package main

import (
	"strconv"
)

type State struct {
	state map[string]string
}

type Store struct {
	states map[string]State
}

type Config struct {
	stateType map[string][]string
}

var store = Store{
	states: make(map[string]State),
}

func _init(config Config) {
	//TODO config from read mib file

	for index := range config.stateType {
		if config.stateType[index][0] == "int" ||
			config.stateType[index][0] == "float" ||
			config.stateType[index][0] == "boolean" ||
			config.stateType[index][0] == "string" {
			if store.states[config.stateType[index][0]].state == nil {
				store.states[config.stateType[index][0]] = State{
					map[string]string{
						index: config.stateType[index][1],
					}}
			} else {
				store.states[config.stateType[index][0]].state[index] = config.stateType[index][1]
			}
		}
	}
}

//
func Set(key string, value string) (bool, string) {
	result, status := find(key, true, value)
	return status, result
}

// "Get" function gets a key and returns the value in string format
func Get(key string) string {
	result, _ := find(key, false, "")
	return result
}

// "Set" function gets a key and a value and assigns the new value in the state
//	if state change successfully returns true and new value otherwise returns false
func find(key string, setValue bool, value string) (string, bool) {
	intChan, stringChan, floatChan, boolChan :=
		make(chan int64), make(chan string), make(chan float64), make(chan bool)
	go findBool(key, boolChan, setValue, value)
	go findFloat(key, floatChan, setValue, value)
	go findInt(key, intChan, setValue, value)
	go findString(key, stringChan, setValue, value)

	resInt, okInt := <-intChan
	resFloat, okFloat := <-floatChan
	resStr, okStr := <-stringChan
	resBool, okBool := <-boolChan

	if okInt {
		return strconv.FormatInt(resInt, 10), true
	} else if okFloat {
		return strconv.FormatFloat(resFloat, 'f', 6, 64), true
	} else if okStr {
		return resStr, true
	} else if okBool {
		return strconv.FormatBool(resBool), true
	}
	return "", false
}

func findInt(key string, channel chan int64, setValue bool, value string) {
	for index := range store.states["int"].state {
		if index == key {
			result, e := strconv.ParseInt(value, 10, 64)
			if e == nil && setValue {
				store.states["int"].state[index] = value
			} else if e != nil && setValue {
				close(channel)
				return
			}
			result, e = strconv.ParseInt(store.states["int"].state[index], 10, 64)
			if e == nil {
				channel <- result
				close(channel)
				return
			}
		}
		close(channel)
		return
	}
}

func findString(key string, channel chan string, setValue bool, value string) bool {
	for index := range store.states["string"].state {
		if index == key {
			if setValue {
				store.states["string"].state[index] = value
			}
			result := store.states["string"].state[index]
			channel <- result
			close(channel)
			return true
		}
	}
	close(channel)
	return false
}

func findFloat(key string, channel chan float64, setValue bool, value string) {
	for index := range store.states["float"].state {
		if index == key {
			result, e := strconv.ParseFloat(value, 64)
			if e == nil && setValue {
				store.states["float"].state[index] = value
			} else if e != nil && setValue {
				close(channel)
				return
			}
			result, e = strconv.ParseFloat(store.states["float"].state[index], 64)
			if e == nil {
				channel <- result
				close(channel)
				return
			}
		}
	}
	close(channel)
	return
}

func findBool(key string, channel chan bool, setValue bool, value string) {
	for index := range store.states["bool"].state {
		if index == key {
			result, e := strconv.ParseBool(value)
			if e == nil && setValue {
				store.states["bool"].state[index] = value
			} else if e != nil && setValue {
				close(channel)
				return
			}
			result, e = strconv.ParseBool(store.states["bool"].state[index])
			if e == nil {
				channel <- result
				close(channel)
				return
			}
		}
	}
	close(channel)
	return
}

func main() {
	config := Config{map[string][]string{
		"temp":   []string{"int", "2"},
		"light":  []string{"float", "1.2"},
		"light2": []string{"float", "1.2"},
	}}
	_init(config)
}
