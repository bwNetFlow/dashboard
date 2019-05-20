package counters

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLabelHash(t *testing.T) {
	labelMap1 := NewLabel(map[string]string{
		"topic":     "flow-messages-enriched",
		"partition": "0",
	})
	labelMap2 := NewLabel(map[string]string{
		"topic":     "flow-messages-enriched",
		"partition": "0",
	})
	labelMap3 := NewLabel(map[string]string{
		"partition": "0",
		"topic":     "flow-messages-enriched",
	})

	Convey("Should hash labels correctly", t, func() {
		So(labelMap1.hash(), ShouldEqual, uint32(2844051426))
		So(labelMap2.hash(), ShouldEqual, uint32(2844051426))
		So(labelMap3.hash(), ShouldEqual, uint32(2844051426))
	})
}
