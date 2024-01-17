package memorydb

import "github.com/unicornultrafoundation/go-helios/u2udb"

type Mod func(store u2udb.Store) u2udb.Store

type Producer struct {
	fs   *fakeFS
	mods []Mod
}

// NewProducer of memory db.
func NewProducer(namespace string, mods ...Mod) u2udb.IterableDBProducer {
	return &Producer{
		fs:   newFakeFS(namespace),
		mods: mods,
	}
}

// Names of existing databases.
func (p *Producer) Names() []string {
	return p.fs.ListFakeDBs()
}

// OpenDB or create db with name.
func (p *Producer) OpenDB(name string) (u2udb.Store, error) {
	db := p.fs.OpenFakeDB(name)

	for _, mod := range p.mods {
		db = mod(db)
	}

	return db, nil
}
