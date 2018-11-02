// A bytecube exercise to assign guests of a wedding party
// to available tables, taking into account seating preferences if possible
package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var w *Wedding

func init() {
	fmt.Println("Initializing test package")
}

func TestSeatGuestsBase(t *testing.T) {
	w := NewWedding()
	w.addTable("A", 8)
	w.addTable("B", 8)
	w.addTable("C", 7)
	w.addTable("D", 7)

	w.addReservation("Thornton", 3)
	w.addReservation("Garcia", 2)
	w.addReservation("Owens", 6, "Thornton", "Taylor")
	w.addReservation("Smith", 1, "Garcia")
	w.addReservation("Taylor", 5)
	w.addReservation("Reese", 7)

	res, err := w.seatGuests(false)
	assert.True(t, res)
	assert.Nil(t, err)

	w.printSeating()
}

func TestSeatGuestsPartyTooLarge(t *testing.T) {
	w := NewWedding()
	w.addTable("A", 8)
	w.addTable("B", 8)
	w.addTable("C", 7)
	w.addTable("D", 7)

	w.addReservation("Thornton", 3)
	w.addReservation("Garcia", 2)
	w.addReservation("Owens", 6, "Thornton", "Taylor")
	w.addReservation("Smith", 1, "Garcia")
	w.addReservation("Taylor", 9)
	w.addReservation("Reese", 7)

	_, err := w.seatGuests(false)
	fmt.Println(err)
	assert.NotNil(t, err)
}

func TestSeatGuestsTooManyGuest(t *testing.T) {
	w := NewWedding()
	w.addTable("A", 8)

	w.addReservation("Thornton", 3)
	w.addReservation("Garcia", 2)
	w.addReservation("Owens", 6, "Thornton", "Taylor")
	w.addReservation("Smith", 1, "Garcia")
	w.addReservation("Taylor", 1)
	w.addReservation("Reese", 7)

	_, err := w.seatGuests(false)
	assert.NotNil(t, err)
	fmt.Printf("%s", err)
}

func TestSeatGuestsNoTables(t *testing.T) {
	w := NewWedding()

	w.addReservation("Thornton", 3)
	w.addReservation("Garcia", 2)
	w.addReservation("Owens", 6, "Thornton", "Taylor")
	w.addReservation("Smith", 1, "Garcia")
	w.addReservation("Taylor", 1)
	w.addReservation("Reese", 7)

	_, err := w.seatGuests(false)
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestSeatGuestsNoReservations(t *testing.T) {
	w := NewWedding()

	w.addTable("A", 8)
	w.addTable("B", 8)
	w.addTable("C", 7)
	w.addTable("D", 7)

	_, err := w.seatGuests(false)
	assert.NotNil(t, err)
	fmt.Println(err)
}

func TestSeatGuestsPartyUsePreferences(t *testing.T) {
	w := NewWedding()

	w.addTable("A", 8)
	w.addTable("B", 8)
	w.addTable("C", 7)
	w.addTable("D", 7)

	w.addReservation("Thornton", 3)
	w.addReservation("Garcia", 2)
	w.addReservation("Owens", 6, "Thornton", "Taylor")
	w.addReservation("Smith", 1, "Reese")
	w.addReservation("Taylor", 5)
	w.addReservation("Reese", 7)

	_, err := w.seatGuests(true)
	assert.Nil(t, err)

	w.printSeating()
}
