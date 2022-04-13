package kode

func Run(code string) error {

	// Return if the code is empty.
	if code == "" {
		return nil
	}

	// Create a new main scope.
	scope := CreateFunction(map[string]Variable{}, map[string]Variable{}, code)

	// Enter the main scope.
	_, err := scope.Run(map[string]Variable{})

	if err == nil {
		print("Result variable: ")
		println(scope.GetVariable("result").Value.(int64))
	}

	return err

}
