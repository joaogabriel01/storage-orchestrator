package strategies

import (
	"github.com/joaogabriel01/storage-orchestrator/pkg/protocols"
	unit_test "github.com/joaogabriel01/storage-orchestrator/pkg/test"
)

var units map[string]protocols.StorageUnit[string, string]
var mock1 *unit_test.UnitMock
var mock2 *unit_test.UnitMock
var targets []string

func initialSetup() {
	mock1 = unit_test.NewUnitMock()
	mock2 = unit_test.NewUnitMock()
	units = make(map[string]protocols.StorageUnit[string, string])
	units["mock1"] = mock1
	units["mock2"] = mock2
	targets = []string{"mock1", "mock2"}
}
