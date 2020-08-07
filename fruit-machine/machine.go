package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/briandowns/spinner"

	pretty "github.com/inancgumus/prettyslice"
)

//CreditCard contains meta for player money calculation
type CreditCard struct {
	ID     int
	Amount int
}

//Player contains human meta
type Player struct {
	Name    string
	Credits *CreditCard
}

//Machine type contains actual chips balance/rate and amount of jackpot money
type Machine struct {
	ChipAmount       int
	ChipExchangeRate int
	CashAmount       int
}

//Prize contains meta for result presenting after each spin: combination and prize(money)
type Prize struct {
	Combination []string
	Gainings    int
}

func shuffle(vals []string) []string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]string, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

var spinValues = []string{"F", "C", "U", "K"}
var winValue = []string{"F", "U", "C", "K"}
var jackpotAmount = 1000

func (m *Machine) spin() (Prize, error) {
	var prize Prize
	if m.ChipAmount < 1 {
		return prize, errors.New("Not enough chips to try, please buy more")
	}
	m.ChipAmount--

	s := spinner.New(spinner.CharSets[15], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                    // Start the spinner
	rand.Seed(time.Now().Unix())
	time.Sleep(3 * time.Second) // Run for some time to simulate work
	prize.Combination = shuffle(spinValues)
	pretty.Show("Combination", prize.Combination)

	if reflect.DeepEqual(prize.Combination, winValue) {
		m.CashAmount -= jackpotAmount
		prize.Gainings = jackpotAmount

		fmt.Println("!🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰*Jackpot*🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰🎰!")
		time.Sleep(4 * time.Second)
	}
	s.Stop()

	fmt.Printf("You won: %v USD\nSpins left:%v\n", prize.Gainings, m.ChipAmount)
	return prize, nil
}

func (p *Player) buyChips(amount int, machine *Machine) error {
	chipsCost := amount / machine.ChipExchangeRate
	if p.Credits.Amount < chipsCost {
		return errors.New("Not enough money to buy chips")
	}
	// TODO calculate floats as well
	p.Credits.Amount -= chipsCost
	machine.ChipAmount += amount

	fmt.Printf("Chips bought, machine balance: %v chips\n", machine.ChipAmount)
	return nil
}

var cardMap = map[int]int{
	112345: 1000,
	100500: 500,
}

func getCardByID(id int) (CreditCard, error) {
	card := CreditCard{
		ID: id,
	}
	if cardMap[card.ID] == 0 {
		// separate non-found errors from empty balance
		return card, errors.New("Credit card not found or empty")
	}
	card.Amount = cardMap[card.ID]

	return card, nil
}

func main() {
	// TODO grab data from stdin
	playerName := "Vladoska"
	playerCardID := 112345

	// creating player object
	fmt.Printf("Welcome %v, registering your account and verifying credit card...\n", playerName)
	creditCard, err := getCardByID(playerCardID)
	if err != nil {
		panic(err)
	}
	player := Player{
		Name:    playerName,
		Credits: &creditCard,
	}
	fmt.Printf("Good news %v, your account created, credit card authenticated successfully, ready to win?\n", playerName)

	// creating machine object
	oneHandMachine := Machine{
		ChipExchangeRate: 2,
		CashAmount:       10000,
	}
	fmt.Printf("Unbelivable, todays jackpot is %v, you have all chances to win!\n\n", oneHandMachine.CashAmount)

	scanner := bufio.NewScanner(os.Stdin)
	helpMessage := `Welcome to the 🎰FUCK🎰 machine! Everything you need is to spin 'FUCK' word!
Press number to start:
💰 1. Buy chips
--------------
🕹  2. Spin
--------------
🤑 3. Check credits balance
--------------
👋 4. exit game
`
	fmt.Println(helpMessage)
	for {
		if scanner.Scan() {

			switch scanner.Text() {
			case "1":
				fmt.Printf("%v, chips rate: 1 dollar = %v chip(s), how much chips do you need?\n", player.Name, oneHandMachine.ChipExchangeRate)
				actionsScanner := bufio.NewScanner(os.Stdin)
				if actionsScanner.Scan() {
					amount, err := strconv.Atoi(actionsScanner.Text())
					if err != nil {
						panic(err)
					}

					err = player.buyChips(amount, &oneHandMachine)
					if err != nil {
						panic(err)
					}
				} else {
					fmt.Println("please enter valid number for chips amount")
				}
			case "2":
				prize, err := oneHandMachine.spin()
				if err != nil {
					fmt.Println(err)
				} else {
					player.Credits.Amount += prize.Gainings
				}
			case "3":
				fmt.Printf("%v balance: %v USD", player.Name, player.Credits.Amount)
			case "4":
				os.Exit(0)
			default:
				fmt.Println(helpMessage)
			}
		}
		time.Sleep(2000 * time.Millisecond)
		print("\033[H\033[2J")
		fmt.Println(helpMessage)
	}
}
