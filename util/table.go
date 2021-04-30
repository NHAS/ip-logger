package util

import (
	"fmt"

	"github.com/NHAS/ip-logger/models"
)

type table struct {
	name          string
	rowNames      []string
	values        [][]string
	valueMaxSizes []int
}

func (t *table) updateMaxs(rn ...string) error {
	if len(rn) > len(t.rowNames) {
		return fmt.Errorf("Wrong size guy")
	}

	if t.valueMaxSizes == nil {
		t.valueMaxSizes = make([]int, len(t.rowNames))
	}

	for i, v := range rn {
		if t.valueMaxSizes[i] < len(v) {
			t.valueMaxSizes[i] = len(v)
		}
	}

	return nil
}

func (t *table) AddRow(rn ...string) error {
	t.rowNames = append(t.rowNames, rn...)

	err := t.updateMaxs(rn...)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) AddValues(vals ...string) error {
	if len(vals) > len(t.rowNames) {
		return fmt.Errorf("Error more values than exist in the row name")
	}

	t.values = append(t.values, vals)

	err := t.updateMaxs(vals...)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Print() {

	top := "|"
	for i, rh := range t.rowNames {
		top += fmt.Sprintf(" %-"+fmt.Sprintf("%d", t.valueMaxSizes[i])+"s |", rh)
	}

	fmt.Printf("%"+fmt.Sprintf("%d", len(top)/2-len(t.name))+"s\n", t.name)

	seperator(len(top))
	fmt.Println(top)
	seperator(len(top))

	for _, row := range t.values {
		line := "|"
		for i, v := range row {
			line += fmt.Sprintf(" %-"+fmt.Sprintf("%d", t.valueMaxSizes[i])+"s |", v)

		}
		fmt.Println(line)
		seperator(len(line))
	}
}

func PrintTable(domain string, u []models.Url) {

	var t table
	t.name = "URLs"

	t.AddRow("Label", "Short URL", "Destination", "Number Visits")
	for _, url := range u {
		t.AddValues(url.Label, domain+"/a/"+url.Identifier, url.Destination, fmt.Sprintf("%d", len(url.Vists)))
	}

	t.Print()

}

func PrintVisits(u models.Url) {
	var t table
	t.name = "Visits"

	t.AddRow("Time", "IP", "UA")
	for _, visit := range u.Vists {
		t.AddValues(visit.CreatedAt.Format("02 Jan 06 15:04"), visit.IP, visit.UA)
	}

	t.Print()
}

func seperator(i int) {
	for n := 0; n < i; n++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
}
