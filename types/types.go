package types

import "time"

type Store struct {
	ID uint
	Chain string
}

type StoreWindow struct {
	Start time.Time
	End   time.Time
	Price uint
	Store Store
}

type Schedule struct {
	Start  time.Time
	End    time.Time
	Chains map[string]bool
}
