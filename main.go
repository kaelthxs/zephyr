package main

func isPalindrome(x int) bool {
    var ourNumber[]int
    if x == 0 {
        return false
    }

    for x != 0 {
		digit := x % 10
		ourNumber = append([]int{digit}, ourNumber...)
		x /= 10
	}

    log.Print(ourNumber)

    for i := 0; i < len(ourNumber)-1; i++ {
        for j := len(ourNumber)-1; j >= 0; j-- {
            if (ourNumber[i] == ourNumber[j]) && (i != j) {
                return true
            }
        }
    }
    return false
}