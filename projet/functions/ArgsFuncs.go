package functions

import (
	"fmt"
	"os"
	"strconv"
)

// NoValueArg is a type representing an argument without a value
type NoValueArg struct {
	Name string
}

// ArgValue is a type representing the value of an argument
type ArgValue string

const (
	// ArgStringValue is a type representing a string value
	ArgStringValue ArgValue = "string"
	// ArgIntValue is a type representing an int value
	ArgIntValue ArgValue = "int"
	// ArgFloatValue is a type representing a float value
	ArgFloatValue ArgValue = "float"
)

// ValueArg is a type representing an argument with a value
// (e.g. '--name value' where 'name' is the Name and 'value' is of the ValueType type)
type ValueArg struct {
	Name      string
	ValueType ArgValue
}

// argList is a type representing a list of arguments
var argList []interface{}

// getArgs return the arguments given to the program
func getArgs() []string {
	return os.Args[1:]
}

// AddNoValueArg adds an argument of the type NoValueArg to the list of arguments
func AddNoValueArg(names ...string) {
	if len(names) == 0 {
		ErrorPrintln("No names given for the NoValueArg")
		return
	}
	for _, name := range names {
		argList = append(argList,
			NoValueArg{
				Name: name,
			})
	}
}

// AddValueArg adds an argument of the type ValueArg to the list of arguments
func AddValueArg(valueType ArgValue, names ...string) {
	if len(names) == 0 {
		ErrorPrintln("No names given for the ValueArg")
		return
	}
	for _, name := range names {
		argList = append(argList,
			ValueArg{
				Name:      name,
				ValueType: valueType,
			})
	}
}

// GetArgValue returns the value of the argument with the given name
// Returns an empty string if the argument doesn't exist.
// Returns the value of the argument if it exists.
// Returns an error if the argument exists but isn't of its expected type or if the argument is of an unknown type or if the argument is a NoValueArg.
func GetArgValue(names ...string) (any, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("no names given for the argument")
	}
	for _, name := range names {
		// Get the argument with the given name if it exists
		for _, arg := range argList {
			switch arg.(type) {
			case NoValueArg:
				// Since the argument doesn't have a value, we just check if it exists
				argWithoutValue := arg.(NoValueArg)
				if argWithoutValue.Name == name {
					return nil, fmt.Errorf("the argument %s is not a NoValueArg and cannot have a value", name)
				}
			case ValueArg:
				argWithValue := arg.(ValueArg)
				args := getArgs()
				if argWithValue.Name == name {
					if argWithValue.ValueType == ArgStringValue {
						// Get the value of the argument
						for i, arg := range args {
							if arg == "--"+name {
								if i+1 < len(args) {
									return args[i+1], nil
								}
							}
						}
					} else if argWithValue.ValueType == ArgIntValue {
						for i, arg := range args {
							if arg == "--"+name {
								if i+1 < len(args) {
									val, err := strconv.Atoi(args[i+1])
									if err != nil {
										return "", fmt.Errorf("the argument %s is of type int but the value is not an int", name)
									}
									return val, nil
								}
							}
						}
					} else if argWithValue.ValueType == ArgFloatValue {
						for i, arg := range args {
							if arg == "--"+name {
								if i+1 < len(args) {
									val, err := strconv.ParseFloat(args[i+1], 64)
									if err != nil {
										return "", fmt.Errorf("the argument %s is of type float but the value is not a float", name)
									}
									return val, nil
								}
							}
						}

					} else {
						return nil, fmt.Errorf("the argument %s is of an unknown type", name)
					}
				}
			default:
				return nil, fmt.Errorf("the argument %s is neither a NoValueArg nor a ValueArg", name)
			}
		}
	}
	return nil, nil
}

// GetArgNoValue returns true if the argument with the given name is present and doesn't have a value
// Returns false if the argument doesn't exist or if it's given a value
func GetArgNoValue(names ...string) (bool, error) {
	if len(names) == 0 {
		return false, fmt.Errorf("no names given for the argument")
	}
	args := getArgs()
	for _, name := range names {
		for _, arg := range argList {
			switch arg.(type) {
			case NoValueArg:
				argWithoutValue := arg.(NoValueArg)
				if argWithoutValue.Name == name {
					for _, arg := range args {
						if arg == "-"+name {
							return true, nil
						}
					}
				}
			case ValueArg:
				argWithoutValue := arg.(ValueArg)
				if argWithoutValue.Name == name {
					return false, fmt.Errorf("the argument %s is not a NoValueArg and cannot have a value", name)
				}
			}
		}
	}
	return false, nil
}

// ArgPresent returns true if the argument with the given name exists
func ArgPresent(name string) bool {
	for _, arg := range argList {
		switch arg.(type) {
		case NoValueArg:
			argWithoutValue := arg.(NoValueArg)
			if argWithoutValue.Name == name {
				return true
			}
		case ValueArg:
			argWithValue := arg.(ValueArg)
			if argWithValue.Name == name {
				return true
			}
		}
	}
	return false
}

// IsNoValueArgPresent returns true if the argument with the given name is present
// and doesn't have a value
func IsNoValueArgPresent(name string) bool {
	for _, arg := range argList {
		switch arg.(type) {
		case NoValueArg:
			argWithoutValue := arg.(NoValueArg)
			if argWithoutValue.Name == name {
				return true
			}
		}
	}
	return false
}

// IsValueArgPresent returns true if the argument with the given name is present
// and has a value
func IsValueArgPresent(name string) bool {
	for _, arg := range argList {
		switch arg.(type) {
		case ValueArg:
			argWithValue := arg.(ValueArg)
			if argWithValue.Name == name {
				return true
			}
		}
	}
	return false
}
