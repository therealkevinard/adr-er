package view

// focusState tracks currently-focused child model. valid states are defined in the iota below.
// this type exposes funcs for wraparound cycling of focusState state.
type focusState int

const (
	focusList focusState = iota
	focusViewer
	focusMax // provides upper boundary for wraparound-increments
)

// Next increments focusState, restarting at 0 if we escape the boundary.
func (f focusState) Next(from focusState) focusState {
	// visualizing 3 plus 1 with focusMax == 3:
	// (3 + 1) % 3 => 4 % 3 => 1
	return (from + 1) % focusMax
}

// Prev decrements focusState, restarting at the end if we escape the boundary.
func (f focusState) Prev(from focusState) focusState {
	// visualizing 0 minus 1 with focusMax == 3:
	// (0 - 1 + 3) % 3 => 2 % 3 => 2
	return (from - 1 + focusMax) % focusMax
}
