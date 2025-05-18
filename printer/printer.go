package printer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PrintError(warning string) {
	fmt.Println(BoldRed + "[!]: " + Reset + Red + warning + Reset)
}

func PrintWarning(warning string) {
	fmt.Println(BoldYellow + "[x]: " + Reset + Yellow + warning + Reset)
}

func PrintNormal(warning string) {
	fmt.Println(BoldCyan + "[~]: " + Reset + warning + Reset)
}

func PrintInfo(warning string) {
	fmt.Println(BoldBlue + "[i]: " + Reset + warning + Reset)
}

func PrintObvious(msg string) {
	fmt.Printf("\033[1;5;91m[!]: %s\033[0m\n", msg)
}

func PrintBanner() {
	fmt.Println("      ::::::::  :::::::: :::::::::  ::::::::::::::::::: ")
	fmt.Println("    :+:    :+::+:    :+::+:    :+::+:    :+:   :+:      ")
	fmt.Println("   +:+       +:+    +:++:+    +:++:+    +:+   +:+       ")
	fmt.Println("  :#:       +#+    +:++#+    +:++#+    +:+   +#+        ")
	fmt.Println(" +#+   +#+#+#+    +#++#+    +#++#+    +#+   +#+         ")
	fmt.Println("#+#    #+##+#    #+##+#    #+##+#    #+#   #+#          ")
	fmt.Println("########  ######## #########  ########    ###           ")
	fmt.Println("      :::       :::::::::::::::    ::::::::::::::::::::::::::::::::: ")
	fmt.Println("    +:+           +:+    :+:+:+  +:+    +:+    +:+       +:+    +:+  ")
	fmt.Println("    +:+           +:+    :+:+:+  +:+    +:+    +:+       +:+    +:+  ")
	fmt.Println("   +#+           +#+    +#+ +:+ +#+    +#+    +#++:++#  +#++:++#:    ")
	fmt.Println("  +#+           +#+    +#+  +#+#+#    +#+    +#+       +#+    +#+    ")
	fmt.Println(" #+#           #+#    #+#   #+#+#    #+#    #+#       #+#    #+#     ")
	fmt.Println("########################    ####    ###    #############    ###      ")
	fmt.Println("")
}

func PPrintArray(arr []string) {
	fmt.Printf(Dim)
	limit := 5
	n := len(arr)

	for i := 0; i < n && i < limit; i++ {
		fmt.Printf("  %s", arr[i])
		if i != limit-1 && i != n-1 {
			fmt.Printf(",\n")
		}
	}

	if n > limit {
		fmt.Printf(",\n  + %d more lines", n-limit)
	}

	fmt.Printf("\n" + Reset)
}

func DebugPrintArray(arr []string) {
	fmt.Print("[")
	for i, s := range arr {
		fmt.Printf("'%s'", s)

		if i != len(arr)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Print("]\n")
}

func AskConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s[?]: %s [Y/n]: %s", Magenta, prompt, Reset)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input.")
			continue
		}

		// Clean and normalize input
		input = strings.TrimSpace(strings.ToLower(input))

		// Handle default (Enter)
		if input == "" || input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" {
			return false
		} else {
			fmt.Println("Please enter 'y' or 'n'.")
		}
	}
}
