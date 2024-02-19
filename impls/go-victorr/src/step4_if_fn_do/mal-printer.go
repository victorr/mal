package main

func PrintString(object MalObject) (string, error) {
	return object.Print(true), nil
}
