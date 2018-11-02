// Package wp is a bytecube exercise to assign guests of a wedding party
// to available tables, taking into account seating preferences if possible
//
package main

import (
	"fmt"
	"sort"
)

type table struct {
	name           string
	size           int32
	remainingSeats int32
	seatMap        map[string]int32 // maping of reservation name to number of seats used
}

type reservation struct {
	name     string
	size     int32
	dislikes []string // list of people preferred not to be seated next to
	seated   bool
}

// Wedding struct holds all wedding data
type Wedding struct {
	tables             []*table
	reservations       []*reservation
	maxTableSize       int32 // largest table available
	maxReservationSize int32 // largest reservation party
	sorted             bool
}

// NewWedding creates an instance of wedding struct
func NewWedding() *Wedding {
	return &Wedding{sorted: false}
}

func (w *Wedding) addTable(name string, size int32) {
	w.tables = append(w.tables, &table{name: name, size: size, remainingSeats: size, seatMap: make(map[string]int32)})

	// keep max table size up to date
	if w.maxTableSize < size {
		w.maxTableSize = size
	}
}

func (w *Wedding) addReservation(name string, size int32, dislikes ...string) {
	var dislikesList []string
	for _, val := range dislikes {
		dislikesList = append(dislikesList, val)
	}
	w.reservations = append(w.reservations, &reservation{name: name, size: size, dislikes: dislikesList, seated: false})

	// keep max reservation size up to date
	if w.maxReservationSize < size {
		w.maxReservationSize = size
	}
}

func (w *Wedding) addReservationToTable(name string, size int32) (bool, error) {
	w.tables = append(w.tables, &table{name: name, size: size, remainingSeats: size, seatMap: make(map[string]int32)})
	return true, nil
}

func (w *Wedding) seatGuests(usePreference bool) (bool, error) {
	// check if it's not possible to seat a table and if so exist immediately
	if w.maxReservationSize > w.maxTableSize {
		return false, fmt.Errorf("reservation of size [%d] is too large, max reservation size is [%d].  Note: All guests must be seated at a single table. Exiting", w.maxReservationSize, w.maxTableSize)
	}

	if w.maxReservationSize == 0 {
		return false, fmt.Errorf("no reservations to process. Exiting")
	}

	fmt.Println("Attempting to seat guests...")

	if !w.sorted {
		// sort the reservations by size from largest to smallest
		sort.Slice(w.reservations, func(i, j int) bool { return w.reservations[i].size > w.reservations[j].size })

		// sort the tables by size from largest to smallest
		sort.Slice(w.tables, func(i, j int) bool { return w.tables[i].size > w.tables[j].size })
		w.sorted = true
	}

	// go through each reservation, starting with the largest and attempt to
	// seat them at each of the tables starting with the largest table
	for _, reservation := range w.reservations {
		//fmt.Printf("Reservation [%s] party of [%d] dislikes %v\n", reservation.name, reservation.size, reservation.dislikes)
		for j, table := range w.tables {
			//fmt.Printf("Table [%s] capacity [%d] dislikes %v\n", table.name, table.remainingSeats, reservation.dislikes)
			if table.remainingSeats >= reservation.size {
				if usePreference {
					skipSeatingAtCurrentTable := false
					for _, dislike := range reservation.dislikes {
						if _, found := table.seatMap[dislike]; found {
							//	fmt.Printf("dislike found at current table[%s], let's not use this table\n", table.name)
							skipSeatingAtCurrentTable = found
						}
					}
					if skipSeatingAtCurrentTable {
						// don't seat anyone at the current table because we've found someone sitting here that we don't like
						continue
					}
				}
				// lets seat these people at this table
				table.remainingSeats = table.remainingSeats - reservation.size
				table.seatMap[reservation.name] = reservation.size
				reservation.seated = true
				//fmt.Printf("Adding to table [%s], new capacity [%d]\n", table.name, table.remainingSeats)
				break
			}
			// if we don't have any more tables to sit the part at and they are not seated then error and return
			if j == len(w.tables)-1 && reservation.seated == false {
				return false, fmt.Errorf("no remaining table large enough for party [%s] found", reservation.name)
			}
		}
	}

	return true, nil
}

func (w *Wedding) printSeating() {
	fmt.Println("Printing seating arragement...")
	for i, table := range w.tables {
		if i != 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("Table %s: ", table.name)
		var counter int
		for k, v := range table.seatMap {
			counter++
			fmt.Printf("%s, party of %d", k, v)
			if counter != len(table.seatMap) {
				fmt.Printf(" & ")
			}
		}
	}
	fmt.Printf("\n")
}

func main() {
	w := NewWedding()

	// add sample test data
	// could also parse from file or parse from command argument
	// or read from stdin
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

	// attempt to seat guests
	var usePreferences = true
	ok, err := w.seatGuests(usePreferences)
	if err != nil {
		fmt.Printf("Unable to seat guests %s", err)
	}

	// if unsucessful seating guest and we tried the first time using preferenes
	// try seating guests again without using preferences
	if usePreferences && !ok {
		fmt.Println("Attempting to seat guests again without using their preferences")
		ok, err = w.seatGuests(usePreferences)
		if err != nil {
			fmt.Printf("Unable to seat guests using preferences, %s", err)
		}
	}

	if !ok {
		fmt.Println("Sorry, we tried, but we can't seat plan for this wedding, please call off wedding!")
		return
	}

	fmt.Println("Succesfully seated guests!")
	w.printSeating()
}
