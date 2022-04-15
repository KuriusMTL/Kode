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
		print("string1: ")
		println(scope.GetVariable("string1").Value.(string))

		print("string2: ")
		println(scope.GetVariable("string2").Value.(string))

		print("result: ")
		println(scope.GetVariable("result").Value.(string))
	}

	return err

}
