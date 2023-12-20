package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type PulseState int

const (
	HIGH_PULSE PulseState = 1
	LOW_PULSE  PulseState = -1
)

type ModuleState int

const (
	OFF ModuleState = -1
	ON  ModuleState = 1
)

type ModuleType int

const (
	FLIP_FLOP ModuleType = iota
	CONJUNCTION
	BROADCASTER
)

type Pulse struct {
	source      string
	destination string
	state       PulseState
}

type Module struct {
	state          ModuleState
	moduleType     ModuleType
	outputs        []string
	name           string
	receivedPulses map[string]PulseState
}

func initializeConjunctions(modules map[string]Module) {
	for moduleName, module := range modules {
		for _, output := range module.outputs {
			if conjunctionModule := modules[output]; conjunctionModule.moduleType == CONJUNCTION {
				conjunctionModule.receivedPulses[moduleName] = LOW_PULSE
			}
		}
	}
}

func handlePulse(m Module, p Pulse) (Module, []Pulse) {
	// fmt.Println(p.source, ":", p.state, "->", p.destination)
	var resultingPulses []Pulse
	switch m.moduleType {
	case FLIP_FLOP:
		if p.state == LOW_PULSE {
			for _, output := range m.outputs {
				var state PulseState
				if m.state == OFF {
					state = HIGH_PULSE
				} else {
					state = LOW_PULSE
				}
				pulse := Pulse{
					source:      m.name,
					destination: output,
					state:       state,
				}
				resultingPulses = append(resultingPulses, pulse)
			}
			m.state *= -1
		}
	case CONJUNCTION:
		m.receivedPulses[p.source] = p.state
		receivedPulsesSum := 0
		for _, receivedPulse := range m.receivedPulses {
			receivedPulsesSum += int(receivedPulse)
		}
		var pulseState PulseState = HIGH_PULSE
		if int(math.Abs(float64(receivedPulsesSum))) == len(m.receivedPulses) {
			if !math.Signbit(float64(receivedPulsesSum)) {
				// all positive
				pulseState = LOW_PULSE
			}
		}
		for _, output := range m.outputs {
			pulse := Pulse{
				source:      m.name,
				destination: output,
				state:       pulseState,
			}
			resultingPulses = append(resultingPulses, pulse)
		}
	}

	return m, resultingPulses
}

func processPulse(modules map[string]Module) (map[string]Module, int, int) {
	var pulsesToProcess []Pulse

	broadcasterModule := modules["broadcaster"]
	for _, output := range broadcasterModule.outputs {
		pulsesToProcess = append(pulsesToProcess, Pulse{
			source:      "broadcaster",
			destination: output,
			state:       LOW_PULSE,
		})
	}
	lowPulseCount := 1
	highPulseCount := 0

	for len(pulsesToProcess) > 0 {
		pulse := pulsesToProcess[0]
		if pulse.state == LOW_PULSE {
			lowPulseCount++
		} else {
			highPulseCount++
		}
		pulsesToProcess = pulsesToProcess[1:]
		resultingModule, resultingPulses := handlePulse(modules[pulse.destination], pulse)
		if _, found := modules[pulse.destination]; found {
			modules[pulse.destination] = resultingModule
		}
		for _, pulse := range resultingPulses {
			pulsesToProcess = append(pulsesToProcess, pulse)
		}
	}

	// fmt.Println(lowPulseCount)
	// fmt.Println(highPulseCount)
	return modules, lowPulseCount, highPulseCount
}

func isInitialState(modules map[string]Module) bool {
	allInInitialState := true
	for _, module := range modules {
		switch module.moduleType {
		case FLIP_FLOP:
			// flip flops have an initial state of off
			if module.state != OFF {
				allInInitialState = false
			}
		case CONJUNCTION:
			// conjunctions have an initial state of having all receivedPulses set to low
			receivedPulsesSum := 0
			for _, receivedPulse := range module.receivedPulses {
				receivedPulsesSum += int(receivedPulse)
			}
			receivedPulsesSum *= -1
			if receivedPulsesSum != len(module.receivedPulses) {
				allInInitialState = false
			}
		}
		if !allInInitialState {
			break
		}
	}
	return allInInitialState
}

func main() {
	inputFile, _ := os.Open("input")
	scanner := bufio.NewScanner(inputFile)

	modules := make(map[string]Module)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " -> ")
		moduleName := splitLine[0]
		splitOutputs := strings.Split(splitLine[1], ", ")

		module := Module{
			outputs:    splitOutputs,
			moduleType: BROADCASTER,
		}
		if moduleName[0:1] == "%" {
			module.state = OFF
			module.moduleType = FLIP_FLOP
			moduleName = moduleName[1:]
		}
		if moduleName[0:1] == "&" {
			module.moduleType = CONJUNCTION
			module.receivedPulses = make(map[string]PulseState)
			moduleName = moduleName[1:]
		}
		module.name = moduleName
		modules[moduleName] = module
	}
	initializeConjunctions(modules)

	// modules = processPulse(modules)
	// fmt.Println(isInitialState(modules))
	// modules = processPulse(modules)
	// fmt.Println(isInitialState(modules))
	// modules = processPulse(modules)
	// fmt.Println(isInitialState(modules))
	// modules = processPulse(modules)
	// fmt.Println(isInitialState(modules))
	// modules = processPulse(modules)
	// fmt.Println(isInitialState(modules))

	cycleCount := 1
	var lowPulseSum, highPulseSum int
	var lowPulseCount, highPulseCount int
	modules, lowPulseSum, highPulseSum = processPulse(modules)

	for !isInitialState(modules) && cycleCount < 1000 {
		cycleCount++
		modules, lowPulseCount, highPulseCount = processPulse(modules)
		lowPulseSum += lowPulseCount
		highPulseSum += highPulseCount
	}
	fmt.Println((lowPulseSum * (1000 / cycleCount)) * (highPulseSum * (1000 / cycleCount)))
}
