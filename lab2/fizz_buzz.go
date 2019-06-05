package main

import "fmt"

func main() {
	for i := 1; i <= 100; i++ {
		//check := 0
		if i%3 == 0 {
			//check = 1
			fmt.Println("Fizz")
		} else if i%5 == 0 {
			//check = 2
			fmt.Println("Buzz")
		} else
		/*if i%3 == 0 && i%5 == 0 {
			check = 3
		}*/
		if i%15 == 0 {
			fmt.Println("FizzBuzz")
		} else {
			fmt.Println(i)
		}

		/*switch check {
		case 1:
			fmt.Println("Fizz")
		case 2:
			fmt.Println("Buzz")
		case 3:
			fmt.Println("FizzBuzz")
		default:
			fmt.Println(i)
		}*/
	}
}

func fizzbuzz() {

}
