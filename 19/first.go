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

func parseItem(line string) map[string]int {
	item := make(map[string]int)
	splitParts := strings.Split(line[1:len(line)-1], ",")
	for _, splitPart := range splitParts {
		parts := strings.Split(splitPart, "=")
		partIdentifier := parts[0]
		partRank, _ := strconv.Atoi(parts[1])
		item[partIdentifier] = partRank
	}

	return item
}

func evaluateCondition(condition Condition, categoryRating int) bool {
	switch condition.comparisonOperator {
	case ">":
		return categoryRating > condition.limit
	case "<":
		return categoryRating < condition.limit
	}
	panic(fmt.Sprintf("Didn't match %s", condition.comparisonOperator))
}

func processWorkflow(workflow Workflow, item map[string]int) string {
	conditionCounter := map[string]int{
		"a": -1,
		"m": -1,
		"s": -1,
		"x": -1,
	}
	for _, conditionKey := range workflow.conditionKeys {
		conditionCounter[conditionKey]++
		conditions := workflow.conditions[conditionKey]
		if evaluateCondition(conditions[conditionCounter[conditionKey]], item[conditionKey]) {
			return conditions[conditionCounter[conditionKey]].result
		}
	}
	return workflow.noneMatched
}

func sumItemRanks(item map[string]int) uint64 {
	var sum uint64 = 0
	for _, rank := range item {
		sum += uint64(rank)
	}

	return sum
}

func process(workflows map[string]Workflow, items []map[string]int) uint64 {
	var acceptedSum uint64 = 0
	var rejectedItems []map[string]int
	for _, item := range items {
		workflow := workflows[STARTING_WORKFLOW]
		for true {
			nextWorkflow := processWorkflow(workflow, item)
			if nextWorkflow == ACCEPT {
				acceptedSum += sumItemRanks(item)
				break
			}
			if nextWorkflow == REJECT {
				rejectedItems = append(rejectedItems, item)
				break
			}
			workflow = workflows[nextWorkflow]
		}
	}
	return acceptedSum
}

func main() {
	inputFile, _ := os.Open("input")
	scanner := bufio.NewScanner(inputFile)

	workflows := make(map[string]Workflow)
	var items []map[string]int
	readingWorkflows := true
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			// handle empty line separating workflows and items
			readingWorkflows = false
			continue
		}
		if readingWorkflows {
			workflowName, workflow := parseWorkflow(line)
			workflows[workflowName] = workflow
		} else {
			item := parseItem(line)
			items = append(items, item)
		}
	}
	fmt.Println(process(workflows, items))
}
