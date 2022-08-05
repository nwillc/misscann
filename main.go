package main

import (
	"fmt"
	"github.com/nwillc/genfuncs"
	"github.com/nwillc/genfuncs/container"
	"github.com/nwillc/genfuncs/container/sequences"
	"golang.org/x/exp/slices"
)

const (
	MISSIONARY person = 1
	CANNIBAL   person = -1
	LEFT       side   = 0
	RIGHT      side   = 1
)

type (
	person int
	side   int
	group  container.GSlice[person]
	river  struct {
		left  group
		right group
		boat  side
		prior *river
	}
)

func main() {
	start := newRiver(
		group{MISSIONARY, MISSIONARY, MISSIONARY, CANNIBAL, CANNIBAL, CANNIBAL},
		group{},
		LEFT,
		nil,
	)

	possible := container.NewDeque[*river](start)
	considering := container.GMap[string, *river]{start.String(): start}
	for possible.Len() > 0 {
		possibility := possible.Remove()
		if possibility.success() {
			fmt.Println("Success:")
			for {
				fmt.Println("       ", possibility)
				if possibility.prior == nil {
					break
				}
				possibility = possibility.prior
			}
			return
		}
		fmt.Println("Considering:", possibility)
		fmt.Println("  New Options:")
		possibility.possibleNext().ForEach(func(_ int, o *river) {
			if considering.Contains(o.String()) {
				return
			}
			fmt.Println("        ", o)
			considering[o.String()] = o
			possible.Add(o)
		})
	}
}

func newRiver(left, right group, boat side, prior *river) *river {
	return &river{
		left:  left.sort(),
		right: right.sort(),
		boat:  boat,
		prior: prior,
	}
}

func (s side) String() string {
	switch s {
	case LEFT:
		return "< "
	case RIGHT:
		return " >"
	default:
		return "??"
	}
}

func (p person) String() string {
	switch p {
	case MISSIONARY:
		return "M"
	case CANNIBAL:
		return "C"
	default:
		return "?"
	}
}

func (b group) allowable() bool {
	asSlice := container.GSlice[person](b)
	if !sequences.Any[person](asSlice, func(p person) bool { return p == MISSIONARY }) {
		return true
	}
	return sequences.Fold[person, int](asSlice, 0, func(acc int, p person) int { return acc + int(p) }) >= 0
}

func (b group) String() string {
	return sequences.JoinToString[person](container.GSlice[person](b), genfuncs.StringerToString[person](), " ", "[", "]")
}

func (b group) sort() group {
	sorted := container.GSlice[person](b).SortBy(func(a, b person) bool { return int(a) < int(b) })
	return group(sorted)
}

func (b group) remove(i int) (person, group) {
	cp := slices.Clone(b)
	p := cp[i]
	return p, slices.Delete(cp, i, i+1)
}

func (b group) add(p person) group {
	cp := slices.Clone(b)
	return append(cp, p)
}

func (r river) allowable() bool {
	return r.left.allowable() && r.right.allowable()
}

func (r river) String() string {
	return r.left.String() + r.boat.String() + r.right.String()
}

func (r river) success() bool {
	return container.GSlice[person](r.left).Len() == 0
}

func (r river) ferryFrom() group {
	if r.boat == LEFT {
		return r.left
	}
	return r.right
}

func (r river) ferryTo() group {
	if r.boat == LEFT {
		return r.right
	}
	return r.left
}

func (r river) ferryOne(i int) (result *river, ok bool) {
	from := r.ferryFrom()
	to := r.ferryTo()
	boat := group{}
	var p person

	p, from = from.remove(i)
	boat = boat.add(p)
	to = to.add(p)

	if !boat.allowable() {
		return nil, false
	}

	if r.boat == LEFT {
		result = newRiver(from, to, RIGHT, &r)
	} else {
		result = newRiver(to, from, LEFT, &r)
	}
	if !result.allowable() {
		return nil, false
	}
	return result, true
}

func (r river) ferryTwo(i, j int) (result *river, ok bool) {
	from := r.ferryFrom()
	to := r.ferryTo()
	boat := group{}
	var p person

	// first
	first := genfuncs.Min(i, j)
	p, from = from.remove(first)
	boat = boat.add(p)
	to = to.add(p)

	// second
	second := genfuncs.Max(i, j) - 1
	p, from = from.remove(second)
	boat = boat.add(p)
	to = to.add(p)

	if !boat.allowable() {
		return nil, false
	}

	if r.boat == LEFT {
		result = newRiver(from, to, RIGHT, &r)
	} else {
		result = newRiver(to, from, LEFT, &r)
	}
	if !result.allowable() {
		return nil, false
	}
	return result, true
}

func (r river) possibleNext() container.GSlice[*river] {
	allowed := container.GMap[string, *river]{}
	for i := 0; i < len(r.ferryFrom()); i++ {
		// try one person
		one, ok := r.ferryOne(i)
		if ok {
			allowed[one.String()] = one
		}
		for j := 0; j < len(r.ferryFrom()); j++ {
			if i == j {
				continue
			}
			// try two
			two, ok := r.ferryTwo(i, j)
			if ok {
				allowed[two.String()] = two
			}
		}
	}
	return allowed.Values()
}
