package kode

func Run(code string) {

	// Return if the code is empty.
	if code == "" {
		return
	}

	// Split the code into lines.
	lines := lineParse(code)

	// Loop through the lines.
	for _, line := range lines {

		result, err := EvaluateExpression(line)

		// Check for any runtime errors.
		if err != nil {
			panic(err)
		}

		// Print the result.
		println(result)

	}

}
