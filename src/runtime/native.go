package runtime

func handleArgs(argsCount int, paramCount int) {
	if argsCount > paramCount {
		panic("There are more arguments than parameters for the function.")
	} else if argsCount < paramCount {
		panic("there are fewer arguments than parameters for the function.")
	}
}

func nativeSprintf(args []RuntimeVal, env Environment) RuntimeVal {
	return StringVal{
		Value: sprint_format(args),
	}
}

func nativePrint(args []RuntimeVal, env Environment) RuntimeVal {
	print(sprint_format(args))
	return NullVal{}
}

func nativePrintln(args []RuntimeVal, env Environment) RuntimeVal {
	println(sprint_format(args))
	return NullVal{}
}

func nativeInt(args []RuntimeVal, env Environment) RuntimeVal {
	handleArgs(len(args), 1)
	return ToInt(args[0])
}

func nativeFloat(args []RuntimeVal, env Environment) RuntimeVal {
	handleArgs(len(args), 1)
	return ToFloat(args[0])
}

func nativeString(args []RuntimeVal, env Environment) RuntimeVal {
	handleArgs(len(args), 1)
	return ToString(args[0])
}

func nativeBool(args []RuntimeVal, env Environment) RuntimeVal {
	handleArgs(len(args), 1)
	return ToBool(args[0])
}

func nativeLen(args []RuntimeVal, env Environment) RuntimeVal {
	handleArgs(len(args), 1)
	return IntVal{
		Value: int64(len(args[0].(ArrayVal).Elements)),
	}
}
