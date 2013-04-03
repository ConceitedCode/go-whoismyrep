package whoismyrep

import (
	"flag"
	"testing"
)

var (
	wimr *WIMR
)

func init() {
	flag.Parse()
	wimr = Open()
}

func TestRepByZip(t *testing.T) {
	res, err := wimr.RepsByZip("60660", "5601")
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res {
		t.Logf("%+v", r)
	}
}
