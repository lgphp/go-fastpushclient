package fastpushclient

import (
	_ "github.com/chentaihan/aesCbc"
	"github.com/rogpeppe/fastuuid"
	"testing"
)

func TestClient(t *testing.T) {

	println(fastuuid.MustNewGenerator().Hex128())
}
