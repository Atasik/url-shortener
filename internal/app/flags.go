package app

import "flag"

type Flags struct {
	Postgres bool
}

func (f *Flags) getFlags() {
	flag.BoolVar(&f.Postgres, "db", false, "flag for Postgres database")
	flag.Parse()
}
