package printer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var UseANSI = true

// helper: wrap s in code and reset if UseANSI, otherwise return s unmodified.
func wrap(s, code string) string {
	if UseANSI {
		return code + s + Reset
	}
	return s
}

func PrintError(warning string) {
	prefix := wrap("[!]: ", BoldRed)
	msg := wrap(warning, Red)
	fmt.Println(prefix + msg)
}

func PrintWarning(warning string) {
	prefix := wrap("[x]: ", BoldYellow)
	msg := wrap(warning, Yellow)
	fmt.Println(prefix + msg)
}

func PrintNormal(warning string) {
	prefix := wrap("[~]: ", BoldCyan)
	msg := wrap(warning, "")
	fmt.Println(prefix + msg)
}

func PrintSuccess(warning string) {
	prefix := wrap("[✓]: ", BoldGreen)
	msg := wrap(warning, "")
	fmt.Println(prefix + msg)
}

func PrintInfo(warning string) {
	prefix := wrap("[i]: ", BoldBlue)
	msg := wrap(warning, "")
	fmt.Println(prefix + msg)
}

func PrintObvious(msg string) {
	if UseANSI {
		fmt.Printf("\033[1;5;91m[!]: %s\033[0m\n", msg)
	} else {
		fmt.Printf("[!]: %s\n", msg)
	}
}

func PrintBanner() {
	lines := []string{
		"      ::::::::  :::::::: :::::::::  ::::::::::::::::::::::: ",
		"    :+:    :+: :+:    :+: :+:    :+: :+:    :+:   :+:       ",
		"   +:+     +:+ +:+    +:+ +:+    +:+ +:+   +:+ +:+        ",
		"  #+#     #+# #+#    #+# #+#    #+# #+#   #+# #+#         ",
		" +#+   +#+#+#+#    #+# #+#    #+# #+#   #+# +#+           ",
		"#+#    #+# #+#    #+# #+#    #+# #+#   #+# #+#            ",
		"########  ######## #########  ########    ###             ",
		"      :::       :::::::::::::    :::::::::::::::::::::::: ",
		"    +:+           +:+    +:+:+:+  +:+    +:+    +:+       ",
		"    +:+           +:+    +:+:+:+  +:+    +:+    +:+       ",
		"   +#+           +#+    +#+ +:+ +#+    +#+    +#+++:++#   ",
		"  +#+           +#+    +#+  +#+#+#+    +#+    +#+       ",
		" #+#           #+#    #+#   +#+#+#+    #+#    #+#       ",
		"########################    ####    ###    #############   ###     ",
		"",
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}

func PPrintArray(arr []string) {
	// Dim only if ANSI
	if UseANSI {
		fmt.Print(Dim)
	}
	limit, n := 5, len(arr)
	for i := 0; i < n && i < limit; i++ {
		fmt.Printf("  %s", arr[i])
		if i != limit-1 && i != n-1 {
			fmt.Println(",")
		}
	}
	if n > limit {
		fmt.Printf(",\n  + %d more lines", n-limit)
	}
	if UseANSI {
		fmt.Println("\n" + Reset)
	} else {
		fmt.Println()
	}
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
		p := wrap("[?]: ", Magenta)
		fmt.Printf("%s %s [Y/n]:", p, prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input.")
			continue
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input == "" || input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" {
			return false
		}
		fmt.Println("Please enter 'y' or 'n'.")
	}
}
