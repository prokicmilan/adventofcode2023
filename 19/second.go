package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Condition struct {
	comparisonOperator string
	limit              int
	result             string
}

type Workflow struct {
	conditionKeys []string
	conditions    map[string][]Condition
	noneMatched   string
}

type Range struct {
	low, high uint16
}

type WorkflowRanges struct {
	workflowName string
	ranges       map[string]Range
}

const ACCEPT = "A"
const REJECT = "R"
const STARTING_WORKFLOW = "in"

func parseWorkflow(line string) (string, Workflow) {
	var workflow Workflow
	var workflowName string

	workflowName = line[0:strings.Index(line, "{")]
	conditionDefinitions := strings.Split(line[strings.Index(line, "{")+1:strings.Index(line, "}")], ",")

	workflow.conditionKeys = make([]string, len(conditionDefinitions)-1)
	workflow.conditions = make(map[string][]Condition)

	for conditionIx := 0; conditionIx < len(conditionDefinitions)-1; conditionIx++ {
		definitionSplit := strings.Split(conditionDefinitions[conditionIx], ":")
		conditionDefinition := definitionSplit[0]
		result := definitionSplit[1]
		conditionKey := conditionDefinition[0:1]
		comparisonOperator := conditionDefinition[1:2]
		limit, _ := strconv.Atoi(conditionDefinition[2:])
		condition := Condition{
			comparisonOperator: comparisonOperator,
			limit:              limit,
			result:             result,
		}
		workflow.conditionKeys[conditionIx] = conditionKey

		if _, found := workflow.conditions[conditionKey]; found {
			workflow.conditions[conditionKey] = append(workflow.conditions[conditionKey], condition)
		} else {
			workflow.conditions[conditionKey] = []Condition{condition}
		}
	}
	workflow.noneMatched = conditionDefinitions[len(conditionDefinitions)-1]

	return workflowName, workflow
}

func determineNewRange(condition Condition, rng Range, invert bool) Range {
	low := rng.low
	high := rng.high
	switch condition.comparisonOperator {
	case ">":
		if invert {
			high = min(rng.high, uint16(condition.limit))
		} else {
			low = max(uint16(condition.limit)+1, low)
		}
		return Range{
			low:  low,
			high: high,
		}
	case "<":
		if invert {
			low = max(uint16(condition.limit), low)
		} else {
			high = min(uint16(condition.limit)-1, rng.high)
		}
		return Range{
			low:  low,
			high: high,
		}
	}

	return rng
}

func process(workflows map[string]Workflow, items []map[string]int) uint64 {
	var acceptedSum uint64 = 0

	var queue []WorkflowRanges

	queue = append(queue, WorkflowRanges{
		workflowName: "in",
		ranges: map[string]Range{
			"x": {low: 1, high: 4000},
			"m": {low: 1, high: 4000},
			"a": {low: 1, high: 4000},
			"s": {low: 1, high: 4000},
		},
	})
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]
		if item.workflowName == "A" {
			var rangeSum uint64 = 1
			for _, rng := range item.ranges {
				rangeSum *= uint64(rng.high) - uint64(rng.low) + 1
			}
			acceptedSum += rangeSum
			continue
		}
		if item.workflowName == "R" {
			continue
		}
		workflow := workflows[item.workflowName]
		conditionCounter := map[string]int{
			"x": -1,
			"m": -1,
			"a": -1,
			"s": -1,
		}
		for _, conditionKey := range workflow.conditionKeys {
			conditionCounter[conditionKey]++
			conditions := workflow.conditions[conditionKey]
			condition := conditions[conditionCounter[conditionKey]]
			newRanges := make(map[string]Range)
			for k, v := range item.ranges {
				newRanges[k] = v
			}
			newRanges[conditionKey] = determineNewRange(condition, newRanges[conditionKey], false)
			newItem := WorkflowRanges{
				workflowName: condition.result,
				ranges:       newRanges,
			}
			queue = append(queue, newItem)
			item.ranges[conditionKey] = determineNewRange(condition, item.ranges[conditionKey], true)
		}
		item.workflowName = workflow.noneMatched
		queue = append(queue, item)

	}
	return acceptedSum
}

func main() {
	inputFile, _ := os.Open("input")
	scanner := bufio.NewScanner(inputFile)

	workflows := make(map[string]Workflow)
	var items []map[string]int
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		workflowName, workflow := parseWorkflow(line)
		workflows[workflowName] = workflow

	}
	fmt.Println(process(workflows, items))
}
