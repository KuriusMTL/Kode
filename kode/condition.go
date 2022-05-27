package kode

// ! ConditionBlock : A block of code that has a condition.
// -------------------------
// ! Condition : The condition of the block.
// -------------------------
// ! ConditionIndex : The line number of the condition.
// -------------------------
// ! Code : The code of the block.
type ConditionBlock struct {
	Condition      string
	ConditionIndex int
	Code           string
}

/**
 * Parse the conditions block(s).
 * @param tokens : *Queue - The tokens to parse.
 * @param currentLine : int - The current line number of the first "if" token.
 * @param lines : []string - The lines of the current scope.
 * @return []ConditionBlock - The parsed conditions block(s).
 * @return int - The ending index line of the condition block(s).
 * @return error - The error if any.
 */
func ParseConditionBlocks(tokens *Queue, currentLine int, lines []string, startLine int) ([]ConditionBlock, int, *ErrorStack) {

	// Get the condition as a string
	// Get the rest of the line tokens and join them to feed the condition
	condition := InlineQueueToString(tokens)

	// Array of the condition blocks. Store the current line number, the condition and the code
	conditionBlocks := append([]ConditionBlock{}, ConditionBlock{condition, currentLine, ""})
	currentLine++          // Skip to next line to avoid including the condition in the code
	foundBoundary := false // Flag to indicate if the boundary of the condition has been found
	nestedBlocksCount := 0 // Keep track of the number of nested blocks

	// Read the next line until the end of the block is reached (end if)
	for currentLine < len(lines) {

		// Parse the current line
		parsed := InlineParse(lines[currentLine], []string{" ", "\t"}, true)

		// Add nested if block
		// A counter is used to determine if the block is nested
		// nestedBlocksCount = 0 is the main outer block
		if len(parsed) > 0 && parsed[0] == "if" {
			nestedBlocksCount++

			// Same level of nesting condition "else if"
		} else if len(parsed) > 1 && parsed[0] == "else" && parsed[1] == "if" && nestedBlocksCount == 0 {

			// Get the rest of the line tokens and join them to feed the condition
			// Not so elegant, but it works
			ifElseCondition := ""
			for i := 2; i < len(parsed); i++ {
				if parsed[i] == " " {
					ifElseCondition += " "
				} else {
					// Don't add an extra space if the token is just a space
					ifElseCondition += parsed[i] + " "
				}
			}

			// Append the new condition block to the list
			conditionBlocks = append(conditionBlocks, ConditionBlock{ifElseCondition, currentLine, ""})
			currentLine++ // Skip to next line to avoid including the condition ending in the code
			continue

			// Start of an else block
		} else if len(parsed) > 0 && parsed[0] == "else" && nestedBlocksCount == 0 {

			// If its the current main block, then the block is ended
			// Else ignore the else and continue (the block is nested)
			if nestedBlocksCount != 0 {

				nestedBlocksCount--

			} else {

				// Append the condition block to the list
				conditionBlocks = append(conditionBlocks, ConditionBlock{"else", currentLine, ""})
				currentLine++
				continue
			}

		}

		// An ending block is reached
		if len(parsed) > 1 && parsed[0] == "end" && parsed[1] == "if" {

			// Check if the main block ended
			if nestedBlocksCount == 0 {

				foundBoundary = true // Cab safely break the loop and confirm the boundary
				break

			} else {

				// Decrease the nested block counter
				// Ignore the end if and continue
				nestedBlocksCount--

				// Don't forget to add the current line to the condition block
				conditionBlocks[len(conditionBlocks)-1].Code += lines[currentLine] + "\n"

			}

			// No ending block found
		} else {

			// Add the current line to the current condition block
			conditionBlocks[len(conditionBlocks)-1].Code += lines[currentLine] + "\n"
		}

		// Increment the current line
		currentLine++
	}

	// Check if the end of the block was found
	// e.g. "end if"
	if !foundBoundary {
		return []ConditionBlock{}, currentLine, CreateError("Condition block not closed with \"end if\"", conditionBlocks[len(conditionBlocks)-1].ConditionIndex+1+startLine)
	}

	// Return the condition blocks
	return conditionBlocks, currentLine, nil
}
