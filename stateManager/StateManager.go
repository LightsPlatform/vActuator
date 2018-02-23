package stateManager

import (
	"strconv"
)

type State struct {
	State map[string]string `json:"State"`
}

type Store struct {
	States map[string]State `json:"States"`
}

type Config struct {
	StateType map[string][]string `json:"StateTypes"`
}

var store = Store{
	States: make(map[string]State),
}

func Init(config Config) Store{
	//TODO config from read mib file

	for index := range config.StateType {
		if config.StateType[index][0] == "int" ||
			config.StateType[index][0] == "float" ||
			config.StateType[index][0] == "boolean" ||
			config.StateType[index][0] == "string" {
			if store.States[config.StateType[index][0]].State == nil {
				store.States[config.StateType[index][0]] = State{
					map[string]string{
						index: config.StateType[index][1],
					}}
			} else {
				store.States[config.StateType[index][0]].State[index] = config.StateType[index][1]
			}
		}
	}

	return store;
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

// "Set" function gets a key and a value and assigns the new value in the State
//	if State change successfully returns true and new value otherwise returns false
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
	for index := range store.States["int"].State {
		if index == key {
			result, e := strconv.ParseInt(value, 10, 64)
			if e == nil && setValue {
				store.States["int"].State[index] = value
			} else if e != nil && setValue {
				close(channel)
				return
			}
			result, e = strconv.ParseInt(store.States["int"].State[index], 10, 64)
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
	for index := range store.States["string"].State {
		if index == key {
			if setValue {
				store.States["string"].State[index] = value
			}
			result := store.States["string"].State[index]
			channel <- result
			close(channel)
			return true
		}
	}
	close(channel)
	return false
}

func findFloat(key string, channel chan float64, setValue bool, value string) {
	for index := range store.States["float"].State {
		if index == key {
			result, e := strconv.ParseFloat(value, 64)
			if e == nil && setValue {
				store.States["float"].State[index] = value
			} else if e != nil && setValue {
				close(channel)
				return
			}
			result, e = strconv.ParseFloat(store.States["float"].State[index], 64)
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
	for index := range store.States["bool"].State {
		if index == key {
			result, e := strconv.ParseBool(value)
			if e == nil && setValue {
				store.States["bool"].State[index] = value
			} else if e != nil && setValue {
				close(channel)
				return
			}
			result, e = strconv.ParseBool(store.States["bool"].State[index])
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

