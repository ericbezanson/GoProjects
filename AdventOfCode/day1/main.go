// AoC Day 1

// pair up the numbers and measure how far apart they are. Pair up the smallest number in the left list with the smallest number in the right list, then the second-smallest left number with the second-smallest right number, and so on.

// Within each pair, figure out how far apart the two numbers are; you'll need to add up all of those distances. For example, if you pair up a 3 from the left list with a 7 from the right list, the distance apart is 4; if you pair up a 9 with a 3, the distance apart is 6.

package main

import (
	"fmt"  // used for formatted I/O (Fscanf)
	"os"   // provide os functionality (open)
	"sort" // used to sort slices (Ints)
)

func main() {
	// Read input
	leftList, rightList, err := readInput("input.txt")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Calculate total distance
	totalDistance := calculateTotalDistance(leftList, rightList)
	fmt.Println("Total Distance:", totalDistance)
}

// Read the input and return two slices of integers
func readInput(filename string) ([]int, []int, error) {
	// uses the os package ( A platform-independent interface to operating system functionality.) that opens the input file for reading, assigning it to "file"
	file, err := os.Open(filename)

	//error handling
	if err != nil {
		return nil, nil, err
	}

	// make sure file is properly closed, even if there are no errors
	defer file.Close()

	// declare two slices of type int that will be used to store the left/right side values of the list
	var leftList, rightList []int

	// declare variables that will be used to temp store the values from each side
	var left, right int

	// a loop that will read the values from filename (ex: input.txt) line by line until an error occurs or end of the file is reached
	for {

		// NOTE: Fscanf is a function iside of gos fmt package. it will read a formatted input from an io.Reader and parse it based on the specified format string
		// in this case im assigning the left int from input.text to &left and the right to &right
		_, err := fmt.Fscanf(file, "%d %d\n", &left, &right)

		// error handling
		if err != nil {
			break // Exit loop when input ends
		}

		// add the numbers to their specified lists
		leftList = append(leftList, left)
		rightList = append(rightList, right)
	}

	// return both lists after the loop concludes
	return leftList, rightList, nil
}

// Calculate the total distance between two lists
func calculateTotalDistance(leftList, rightList []int) int {
	// sort is a default go package that profides primatives for sorting slices, "Ints" will sort the values in increasing order (smallest to largest)
	sort.Ints(leftList)
	sort.Ints(rightList)

	// initialize total distance var at 0, will accumulate all the differences between left and right list entries
	totalDistance := 0

	// loop that itterates over the left list (both lists have the same length)
	for i := range leftList {

		// before adding the differece between left/right to totalDistance we account for negative results (right side entry larger than left)
		totalDistance += accountForNegativeNumbers(leftList[i] - rightList[i])
	}

	// return the total distance
	return totalDistance
}

// if the result is negaitve, we reutrn it as a positive number to keep the difference calculation accurate.
func accountForNegativeNumbers(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
