package kode

import "errors"

type LoopBlock struct {
	Condition string
	Code      string
	LoopIndex int
}

func ParseLoopBlock(tokens *Queue, currentLine int, lines []string) (LoopBlock, int, error) {

	condition := InlineQueueToString(tokens)

	currentLine++          // Skip to next line to avoid including the condition in the code
	foundBoundary := false // Flag to indicate if the boundary of the condition has been found
	nestedBlocksCount := 0 // Keep track of the number of nested blocks
	code := ""             // The code of the loop block
	startIndex := currentLine

	for currentLine < len(lines) {

		// Parse the current line
		parsed := InlineParse(lines[currentLine], []string{" "}, true)

		if len(parsed) > 0 && parsed[0] == "for" {

			nestedBlocksCount++

		} else if len(parsed) > 1 && parsed[0] == "end" && parsed[1] == "for" {

			if nestedBlocksCount == 0 {
				foundBoundary = true
				break
			} else {
				nestedBlocksCount--
			}

		}

		code += lines[currentLine] + "\n"
		currentLine++

	}

	if !foundBoundary {
		return LoopBlock{}, currentLine, errors.New("Could not find the end of the for loop")
	}

	return LoopBlock{condition, code, startIndex}, currentLine, nil

}
