package croot_test

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/go-hep/croot"
)

type Det struct {
	E float64
	T float64
}

type Event struct {
	I int64
	A Det
	B Det
}

type DataSlice struct {
	I     int64
	Data  float64
	Slice []float64
}

type DataArray struct {
	I     int64
	Data  float64
	Array [2]float64
}

type DataString struct {
	I      int64
	Data   float64
	String string
}

func TestTreeBuiltinsRW(t *testing.T) {
	const fname = "simple-event.root"
	const evtmax = 10000
	const splitlevel = 32
	const bufsiz = 32000
	const compress = 1
	const netopt = 0

	// write
	ref := make([]string, 0, 50)
	{
		add := func(str string) {
			ref = append(ref, str)
		}

		f, err := croot.OpenFile(fname, "recreate", "croot event file", compress, netopt)
		if err != nil {
			t.Errorf(err.Error())
		}

		// create a tree
		tree := croot.NewTree("tree", "tree", splitlevel)

		e := Event{}

		// create a branch with energy
		_, err = tree.Branch2("evt_i", &e.I, "evt_i/L", bufsiz)
		if err != nil {
			t.Errorf(err.Error())
		}

		_, err = tree.Branch2("evt_a_e", &e.A.E, "evt_a_e/D", bufsiz)
		if err != nil {
			t.Errorf(err.Error())
		}

		_, err = tree.Branch2("evt_a_t", &e.A.T, "evt_a_t/D", bufsiz)
		if err != nil {
			t.Errorf(err.Error())
		}

		_, err = tree.Branch2("evt_b_e", &e.B.E, "evt_b_e/D", bufsiz)
		if err != nil {
			t.Errorf(err.Error())
		}

		_, err = tree.Branch2("evt_b_t", &e.B.T, "evt_b_t/D", bufsiz)
		if err != nil {
			t.Errorf(err.Error())
		}

		// initialize our source of random numbers...
		src := rand.New(rand.NewSource(1))

		// fill some events with random numbers
		for iev := int64(0); iev != evtmax; iev++ {
			if iev%1000 == 0 {
				add(fmt.Sprintf(":: processing event %d...\n", iev))
			}

			e.I = iev
			// the two energies follow a gaussian distribution
			e.A.E = src.NormFloat64()
			e.B.E = src.NormFloat64()

			e.A.T = src.Float64()
			e.B.T = e.A.T * (src.NormFloat64()*1. + 0.)
			// e.A.Fs = []float64{e.A.E, e.A.T}
			// e.B.Fs = []float64{e.B.E, e.B.T}

			if iev%1000 == 0 {
				add(fmt.Sprintf("evt.i=   %8d\n", e.I))
				add(fmt.Sprintf("evt.a.e= %8.3f\n", e.A.E))
				add(fmt.Sprintf("evt.a.t= %8.3f\n", e.A.T))
				add(fmt.Sprintf("evt.b.e= %8.3f\n", e.B.E))
				add(fmt.Sprintf("evt.b.t= %8.3f\n", e.B.T))
			}
			_, err = tree.Fill()
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		f.Write("", 0, 0)
		f.Close("")
	}

	// read back
	chk := make([]string, 0, 50)
	{
		add := func(str string) {
			chk = append(chk, str)
		}

		f, err := croot.OpenFile(fname, "read", "croot event file", compress, netopt)
		if err != nil {
			t.Fatalf(err.Error())
		}
		tree := f.GetTree("tree")
		if tree.GetEntries() != evtmax {
			t.Fatalf("expected [%v] entries, got %v\n", evtmax, tree.GetEntries())
		}

		e := Event{}

		tree.SetBranchAddress("evt_i", &e.I)
		tree.SetBranchAddress("evt_a_e", &e.A.E)
		tree.SetBranchAddress("evt_a_t", &e.A.T)
		tree.SetBranchAddress("evt_b_e", &e.B.E)
		tree.SetBranchAddress("evt_b_t", &e.B.T)

		// read events
		for iev := int64(0); iev != evtmax; iev++ {
			if iev%1000 == 0 {
				add(fmt.Sprintf(":: processing event %d...\n", iev))
			}
			if tree.GetEntry(iev, 1) <= 0 {
				panic("error")
			}
			if iev%1000 == 0 {
				add(fmt.Sprintf("evt.i=   %8d\n", e.I))
				add(fmt.Sprintf("evt.a.e= %8.3f\n", e.A.E))
				add(fmt.Sprintf("evt.a.t= %8.3f\n", e.A.T))
				add(fmt.Sprintf("evt.b.e= %8.3f\n", e.B.E))
				add(fmt.Sprintf("evt.b.t= %8.3f\n", e.B.T))
			}

			if iev != e.I {
				t.Fatalf("invalid event number. expected %v, got %v", iev, e.I)
			}
		}
		f.Close("")
	}

	if !reflect.DeepEqual(ref, chk) {
		t.Errorf("log files do not match")
	}

	err := os.Remove(fname)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestTreeStructRW(t *testing.T) {
	const fname = "struct-event.root"
	const evtmax = 10000
	const splitlevel = 32
	const bufsiz = 32000
	const compress = 1
	const netopt = 0

	// write
	ref := make([]string, 0, 50)
	{
		add := func(str string) {
			ref = append(ref, str)
		}

		f, err := croot.OpenFile(fname, "recreate", "croot event file", compress, netopt)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// create a tree
		tree := croot.NewTree("tree", "tree", splitlevel)

		e := Event{}

		_, err = tree.Branch("evt", &e, bufsiz, 0)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// initialize our source of random numbers...
		src := rand.New(rand.NewSource(1))

		// fill some events with random numbers
		for iev := int64(0); iev != evtmax; iev++ {
			if iev%1000 == 0 {
				add(fmt.Sprintf(":: processing event %d...\n", iev))
			}

			e.I = iev
			// the two energies follow a gaussian distribution
			e.A.E = src.NormFloat64()
			e.B.E = src.NormFloat64()

			e.A.T = src.Float64()
			e.B.T = e.A.T * (src.NormFloat64()*1. + 0.)

			if iev%1000 == 0 {
				add(fmt.Sprintf("evt.i=   %8d\n", e.I))
				add(fmt.Sprintf("evt.a.e= %8.3f\n", e.A.E))
				add(fmt.Sprintf("evt.a.t= %8.3f\n", e.A.T))
				add(fmt.Sprintf("evt.b.e= %8.3f\n", e.B.E))
				add(fmt.Sprintf("evt.b.t= %8.3f\n", e.B.T))
			}
			_, err = tree.Fill()
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		f.Write("", 0, 0)
		f.Close("")
	}

	// read back
	chk := make([]string, 0, 50)
	{
		add := func(str string) {
			chk = append(chk, str)
		}

		f, err := croot.OpenFile(fname, "read", "croot event file", compress, netopt)
		if err != nil {
			t.Errorf(err.Error())
		}

		tree := f.GetTree("tree")
		if tree.GetEntries() != evtmax {
			t.Errorf("expected [%v] entries, got %v\n", evtmax, tree.GetEntries())
		}

		var e Event
		tree.SetBranchAddress("evt", &e)

		// read events
		for iev := int64(0); iev != evtmax; iev++ {
			if iev%1000 == 0 {
				add(fmt.Sprintf(":: processing event %d...\n", iev))
			}
			if tree.GetEntry(iev, 1) <= 0 {
				panic("error")
			}
			if iev%1000 == 0 {
				add(fmt.Sprintf("evt.i=   %8d\n", e.I))
				add(fmt.Sprintf("evt.a.e= %8.3f\n", e.A.E))
				add(fmt.Sprintf("evt.a.t= %8.3f\n", e.A.T))
				add(fmt.Sprintf("evt.b.e= %8.3f\n", e.B.E))
				add(fmt.Sprintf("evt.b.t= %8.3f\n", e.B.T))
			}

			if iev != e.I {
				t.Fatalf("invalid event number. expected %v, got %v", iev, e.I)
			}
		}
		f.Close("")
	}

	if !reflect.DeepEqual(ref, chk) {
		t.Errorf("log files do not match")
	}

	err := os.Remove(fname)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestTreeStructSlice(t *testing.T) {
	const fname = "struct-slice.root"
	const evtmax = 10
	const splitlevel = 32
	const bufsiz = 32000
	const compress = 1
	const netopt = 0

	// write
	ref := make([]string, 0, 50)
	{
		add := func(str string) {
			ref = append(ref, str)
		}

		f, err := croot.OpenFile(fname, "recreate", "croot event file", compress, netopt)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// create a tree
		tree := croot.NewTree("tree", "tree", splitlevel)

		e := DataSlice{}
		e.Slice = make([]float64, 0)

		_, err = tree.Branch("evt", &e, bufsiz, 0)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// initialize our source of random numbers...
		src := rand.New(rand.NewSource(1))

		// fill some events with random numbers
		for iev := int64(0); iev != evtmax; iev++ {
			if iev%1000 == 0 {
				add(fmt.Sprintf(":: processing event %d...\n", iev))
			}

			e.I = iev
			e.Data = src.NormFloat64()

			e.Slice = e.Slice[:0]
			e.Slice = append(e.Slice, e.Data, -e.Data)

			if len(e.Slice) != 2 {
				t.Errorf("invalid e.Slice size: %v (expected 2)", len(e.Slice))
			}

			if iev%1000 == 0 {
				add(fmt.Sprintf("evt.i=     %8d\n", e.I))
				add(fmt.Sprintf("evt.d=     %8.3f\n", e.Data))
				add(fmt.Sprintf("evt.slice= %8.3f %8.3f\n", e.Slice[0], e.Slice[1]))
			}
			_, err = tree.Fill()
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		f.Write("", 0, 0)
		f.Close("")
	}

	// read back
	chk := make([]string, 0, 50)
	{
		add := func(str string) {
			chk = append(chk, str)
		}

		f, err := croot.OpenFile(fname, "read", "croot event file", compress, netopt)
		if err != nil {
			t.Errorf(err.Error())
		}

		tree := f.GetTree("tree")
		if tree.GetEntries() != evtmax {
			t.Errorf("expected [%v] entries, got %v\n", evtmax, tree.GetEntries())
		}

		var e DataSlice
		e.Slice = make([]float64, 0, 2)
		tree.SetBranchAddress("evt", &e)

		// read events
		for iev := int64(0); iev != evtmax; iev++ {
			if iev%1000 == 0 {
				add(fmt.Sprintf(":: processing event %d...\n", iev))
			}
			if tree.GetEntry(iev, 1) <= 0 {
				panic("error")
			}
			if iev%1000 == 0 {
				add(fmt.Sprintf("evt.i=     %8d\n", e.I))
				add(fmt.Sprintf("evt.d=     %8.3f\n", e.Data))
				add(fmt.Sprintf("evt.slice= %8.3f %8.3f\n", e.Slice[0], e.Slice[1]))
			}

			if len(e.Slice) != 2 {
				t.Errorf("invalid e.Slice size: %v (expected 2)", len(e.Slice))
			}
			if e.Slice[0] != e.Data {
				t.Errorf("invalid e.Slice[0] value: %v (expected %v)",
					e.Slice[0], e.Data)
			}
			if e.Slice[1] != -e.Data {
				t.Errorf("invalid e.Slice[1] value: %v (expected %v)",
					e.Slice[1], -e.Data)
			}
			if iev != e.I {
				t.Fatalf("invalid event number. expected %v, got %v", iev, e.I)
			}
		}
		f.Close("")
	}

	if !reflect.DeepEqual(ref, chk) {
		t.Errorf("log files do not match")
	}

	err := os.Remove(fname)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestTreeStructArray(t *testing.T) {
	const fname = "struct-array.root"
	const evtmax = 10000
	const splitlevel = 32
	const bufsiz = 32000
	const compress = 1
	const netopt = 0

	// write
	ref := make([]string, 0, 50)
	{
		add := func(str string) {
			ref = append(ref, str)
		}

		f, err := croot.OpenFile(fname, "recreate", "croot event file", compress, netopt)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// create a tree
		tree := croot.NewTree("tree", "tree", splitlevel)

		e := DataArray{}

		_, err = tree.Branch("evt", &e, bufsiz, 0)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// initialize our source of random numbers...
		src := rand.New(rand.NewSource(1))

		// fill some events with random numbers
		for iev := int64(0); iev != evtmax; iev++ {
			if iev%1000 == 0 {
				add(fmt.Sprintf(":: processing event %d...\n", iev))
			}

			e.I = iev
			e.Data = src.NormFloat64()

			e.Array[0] = e.Data
			e.Array[1] = -e.Data

			if len(e.Array) != 2 {
				t.Errorf("invalid e.Array size: %v (expected 2)", len(e.Array))
			}

			if iev%1000 == 0 {
				add(fmt.Sprintf("evt.i=     %8d\n", e.I))
				add(fmt.Sprintf("evt.d=     %8.3f\n", e.Data))
				add(fmt.Sprintf("evt.array= %8.3f %8.3f\n", e.Array[0], e.Array[1]))
			}
			_, err = tree.Fill()
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		f.Write("", 0, 0)
		f.Close("")
	}

	// read back
	chk := make([]string, 0, 50)
	{
		add := func(str string) {
			chk = append(chk, str)
		}

		f, err := croot.OpenFile(fname, "read", "croot event file", compress, netopt)
		if err != nil {
			t.Errorf(err.Error())
		}

		tree := f.GetTree("tree")
		if tree.GetEntries() != evtmax {
			t.Errorf("expected [%v] entries, got %v\n", evtmax, tree.GetEntries())
		}

		var e DataArray
		tree.SetBranchAddress("evt", &e)

		// read events
		for iev := int64(0); iev != evtmax; iev++ {
			if iev%1000 == 0 {
				add(fmt.Sprintf(":: processing event %d...\n", iev))
			}
			if tree.GetEntry(iev, 1) <= 0 {
				panic("error")
			}
			if iev%1000 == 0 {
				add(fmt.Sprintf("evt.i=     %8d\n", e.I))
				add(fmt.Sprintf("evt.d=     %8.3f\n", e.Data))
				add(fmt.Sprintf("evt.array= %8.3f %8.3f\n", e.Array[0], e.Array[1]))
			}

			if len(e.Array) != 2 {
				t.Errorf("invalid e.Array size: %v (expected 2)", len(e.Array))
			}
			if e.Array[0] != e.Data {
				t.Errorf("invalid e.Array[0] value: %v (expected %v)",
					e.Array[0], e.Data)
			}
			if e.Array[1] != -e.Data {
				t.Errorf("invalid e.Array[0] value: %v (expected %v)",
					e.Array[1], -e.Data)
			}
			if iev != e.I {
				t.Fatalf("invalid event number. expected %v, got %v", iev, e.I)
			}
		}
		f.Close("")
	}

	if !reflect.DeepEqual(ref, chk) {
		t.Errorf("log files do not match")
	}

	err := os.Remove(fname)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestTreeStructString(t *testing.T) {
	const fname = "struct-string.root"
	const evtmax = 10000
	const splitlevel = 32
	const bufsiz = 32000
	const compress = 1
	const netopt = 0

	return // FIXME

	// write
	ref := make([]string, 0, 50)
	{
		add := func(str string) {
			ref = append(ref, str)
		}

		f, err := croot.OpenFile(fname, "recreate", "croot event file", compress, netopt)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// create a tree
		tree := croot.NewTree("tree", "tree", splitlevel)

		e := DataString{}

		_, err = tree.Branch("evt", &e, bufsiz, 0)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// initialize our source of random numbers...
		src := rand.New(rand.NewSource(1))

		// fill some events with random numbers
		for iev := int64(0); iev != evtmax; iev++ {
			if iev%1000 == 0 {
				add(fmt.Sprintf(":: processing event %d...\n", iev))
			}

			e.I = iev
			e.Data = src.NormFloat64()

			e.String = fmt.Sprintf("%v", e.Data)

			if iev%1000 == 0 {
				add(fmt.Sprintf("evt.i=     %8d\n", e.I))
				add(fmt.Sprintf("evt.d=     %v\n", e.Data))
				add(fmt.Sprintf("evt.s=     %s\n", e.String))
			}
			_, err = tree.Fill()
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		f.Write("", 0, 0)
		f.Close("")
	}

	// read back
	chk := make([]string, 0, 50)
	{
		add := func(str string) {
			chk = append(chk, str)
		}

		f, err := croot.OpenFile(fname, "read", "croot event file", compress, netopt)
		if err != nil {
			t.Errorf(err.Error())
		}

		tree := f.GetTree("tree")
		if tree.GetEntries() != evtmax {
			t.Errorf("expected [%v] entries, got %v\n", evtmax, tree.GetEntries())
		}

		var e DataString
		tree.SetBranchAddress("evt", &e)

		// read events
		for iev := int64(0); iev != evtmax; iev++ {
			fmt.Printf(":: processing event %d...\n", iev)
			if iev%1000 == 0 {
				add(fmt.Sprintf(":: processing event %d...\n", iev))
			}
			if tree.GetEntry(iev, 1) <= 0 {
				panic("error")
			}
			if iev%1000 == 0 {
				add(fmt.Sprintf("evt.i=     %8d\n", e.I))
				add(fmt.Sprintf("evt.d=     %vf\n", e.Data))
				add(fmt.Sprintf("evt.s=     %s\n", e.String))
			}

			if iev != e.I {
				t.Fatalf("invalid event number. expected %v, got %v", iev, e.I)
			}
		}
		f.Close("")
	}

	if !reflect.DeepEqual(ref, chk) {
		t.Fatalf("log files do not match\n==ref==\n%s\n==chk==\n%s\n", ref, chk)
	}

	err := os.Remove(fname)
	if err != nil {
		t.Errorf(err.Error())
	}
}

// EOF
