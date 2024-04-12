package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)

	p.quit <- ch

	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order no %d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("Making pizza no %d, it will take %d seconds.\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** Ran out of ingredients for the pizza no %d!***", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making the pizza no %d!***", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza no %d is ready to be served!", pizzaNumber)
		}

		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func pizzaria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	i := 0

	// run forever until we receive a quit notification
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			// We tried to make the pizza by sending something to the data channel
			case pizzaMaker.data <- *currentPizza:
			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func main() {
	// seed random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// print out the message
	color.Cyan("The Pizzaria is open for business!")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzaria(pizzaJob)

	// create and run consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order no %d is out for delivery!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad for order no %d!", i.pizzaNumber)
			}
		} else {
			color.Cyan("The Pizzaria is done making pizzas!")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("***Error closing channel***", err)
			}
		}
	}

	// print out the ending message
	color.Cyan("---------------------------------")
	color.Cyan("The Pizzaria is done for the day!")
	color.Cyan("We made %d pizzas, but failed to make %d, with a total of %d attempts.", pizzasMade, pizzasFailed, total)
	switch {
	case pizzasFailed > 9:
		color.Red("It was an awful day!")
	case pizzasFailed >= 6:
		color.Red("It was not a good day!")
	case pizzasFailed >= 4:
		color.Yellow("It was an okay day!")
	case pizzasFailed >= 2:
		color.Yellow("It was a pretty good day!")
	default:
		color.Green("It was a gread day!")
	}
}
