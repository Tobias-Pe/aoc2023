package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type InstructionSet struct {
	instructions []rune
}

func (instructionSet InstructionSet) peopleVersionRunOn(startNode *Node, targetNodeName string) int {
	counterTillTarget := 0
	currentNode := startNode

	for currentNode.name != targetNodeName {
		instruction := instructionSet.instructions[counterTillTarget%len(instructionSet.instructions)]
		if instruction == 'L' {
			currentNode = currentNode.left
		} else {
			currentNode = currentNode.right
		}

		counterTillTarget++
	}

	return counterTillTarget
}

func (instructionSet InstructionSet) ghostVersionRunOn(startingNodes []*Node, targetNodeSuffix string) (int, int) {
	counterFromStartTillTarget := 0
	var currentNodes []*Node
	currentNodes = append(currentNodes, startingNodes...)

	suffix, countHavingSuffix := haveAllNodesSuffix(currentNodes, targetNodeSuffix)
	for ; !suffix; suffix, countHavingSuffix = haveAllNodesSuffix(currentNodes, targetNodeSuffix) {
		instruction := instructionSet.instructions[counterFromStartTillTarget%len(instructionSet.instructions)]
		for i := range currentNodes {
			if instruction == 'L' {
				currentNodes[i] = currentNodes[i].left
			} else {
				currentNodes[i] = currentNodes[i].right
			}
		}
		counterFromStartTillTarget++
	}
	fmt.Println("From Start till Target", counterFromStartTillTarget, listNodeNames(currentNodes), countHavingSuffix)

	counterFromTargetTillTarget := 0
	currentCounter := counterFromStartTillTarget
	_, countHavingSuffix = haveAllNodesSuffix(currentNodes, targetNodeSuffix)
	suffix = false
	for ; !suffix; suffix, countHavingSuffix = haveAllNodesSuffix(currentNodes, targetNodeSuffix) {
		instruction := instructionSet.instructions[currentCounter%len(instructionSet.instructions)]
		for i := range currentNodes {
			if instruction == 'L' {
				currentNodes[i] = currentNodes[i].left
			} else {
				currentNodes[i] = currentNodes[i].right
			}
		}
		counterFromTargetTillTarget++
		currentCounter++
	}
	fmt.Println("From Target till Target", counterFromTargetTillTarget, listNodeNames(currentNodes), countHavingSuffix)

	return counterFromStartTillTarget, counterFromTargetTillTarget
}

func listNodeNames(nodes []*Node) []string {
	var names []string
	for _, node := range nodes {
		names = append(names, node.name)
	}
	return names
}

func haveAllNodesSuffix(nodes []*Node, suffix string) (bool, int) {
	count := 0
	result := true
	for _, node := range nodes {
		hasSuffix := strings.HasSuffix(node.name, suffix)
		if hasSuffix {
			count++
		}
		result = result && hasSuffix
	}
	return result, count
}

type Node struct {
	name  string
	left  *Node
	right *Node
}

func main() {
	start := time.Now()

	lines := readFile("day8/input.txt")
	instructionSet, startNodePtr := parseTreePeopleVersion(lines)
	stepCount := instructionSet.peopleVersionRunOn(startNodePtr, "ZZZ")
	fmt.Println("[Part1] Steps from AAA till ZZZ:", stepCount)

	lines = readFile("day8/input.txt")
	instructionSet, startNodePtrs := parseTreeGhostVersion(lines)

	var targetCounts []int64
	var iterationCounts []int64
	for _, ptr := range startNodePtrs {
		var nodes []*Node
		nodes = append(nodes, ptr)
		stepsTillTarget, repeatingSteps := instructionSet.ghostVersionRunOn(nodes, "Z")
		targetCounts = append(targetCounts, int64(stepsTillTarget))
		iterationCounts = append(iterationCounts, int64(repeatingSteps))
		fmt.Println("[Part2] Steps from", ptr.name, " till second time of *Z:", stepsTillTarget, repeatingSteps)
	}

	fmt.Println("[Part2] Steps till all are at target nodes", lcmArray(targetCounts), targetCounts, iterationCounts)

	fmt.Println("Finished in", time.Since(start))
}

func lcmArray(numbers []int64) int64 {
	lcm := findLcm(numbers[0], numbers[1])
	for i := 2; i < len(numbers); i++ {
		lcm = findLcm(lcm, numbers[i])
	}
	return lcm
}

func findLcm(first, second int64) int64 {
	return first * second / findGcd(first, second)
}

func findGcd(first, second int64) int64 {
	if first == 0 {
		return second
	}
	// recursively call findGCD
	return findGcd(second%first, first)
}

func readFile(file string) []string {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n")
	return split
}

func parseTreePeopleVersion(lines []string) (InstructionSet, *Node) {
	instructionSet := InstructionSet{instructions: []rune(lines[0])}
	var startNode *Node = nil

	var nameToNodeMap = map[string]*Node{}
	for i := 2; i < len(lines); i++ {
		nodePtr := readoutAndSaveNode(lines[i], nameToNodeMap)

		if nodePtr.name == "AAA" {
			startNode = nodePtr
		}
	}

	return instructionSet, startNode
}

func parseTreeGhostVersion(lines []string) (InstructionSet, []*Node) {
	instructionSet := InstructionSet{instructions: []rune(lines[0])}
	var startingNodes []*Node

	var nameToNodeMap = map[string]*Node{}
	for i := 2; i < len(lines); i++ {
		nodePtr := readoutAndSaveNode(lines[i], nameToNodeMap)

		if strings.HasSuffix(nodePtr.name, "A") {
			startingNodes = append(startingNodes, nodePtr)
		}
	}

	return instructionSet, startingNodes
}

func readoutAndSaveNode(line string, nameToNodeMap map[string]*Node) *Node {
	fields := strings.Fields(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(line, "=", ""), "(", ""), ")", ""), ",", ""))

	nodePtr := checkCreateOrGetNodeName(nameToNodeMap, fields[0])
	nodePtr.left = checkCreateOrGetNodeName(nameToNodeMap, fields[1])
	nodePtr.right = checkCreateOrGetNodeName(nameToNodeMap, fields[2])
	nameToNodeMap[fields[0]] = nodePtr
	return nodePtr
}

func checkCreateOrGetNodeName(nameToNodeMap map[string]*Node, nodeName string) *Node {
	nodePtr, ok := nameToNodeMap[nodeName]
	if !ok {
		nodePtr = &Node{
			name:  nodeName,
			left:  nil,
			right: nil,
		}
		nameToNodeMap[nodeName] = nodePtr
	}
	return nodePtr
}
