package utils

func Black(str string) string {
	return "\u001b[30m" + str + "\u001b[0m"
}

func Red(str string) string {
	return "\u001b[31m" + str + "\u001b[0m"
}

func Green(str string) string {
	return "\u001b[32m" + str + "\u001b[0m"
}

func Yellow(str string) string {
	return "\u001b[33m" + str + "\u001b[0m"
}

func Blue(str string) string {
	return "\u001b[34m" + str + "\u001b[0m"
}

func Magenta(str string) string {
	return "\u001b[35m" + str + "\u001b[0m"
}

func Cyan(str string) string {
	return "\u001b[36m" + str + "\u001b[0m"
}

func White(str string) string {
	return "\u001b[37m" + str + "\u001b[0m"
}

func Grey(str string) string {
	return "\u001b[90m" + str + "\u001b[0m"
}
