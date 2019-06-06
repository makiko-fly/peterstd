package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeSlot(t *testing.T) {
	slot := NewTimeSlot("09:25:00-10:25:00")
	assert.Equal(t, true, slot.IsValid(time.Date(0, 0, 0, 9, 25, 0, 0, TZShanghai())))
	assert.Equal(t, true, slot.IsValid(time.Date(0, 0, 0, 9, 30, 0, 0, TZShanghai())))
	assert.Equal(t, false, slot.IsValid(time.Date(0, 0, 0, 10, 25, 0, 0, TZShanghai())))
	assert.Equal(t, false, slot.IsValid(time.Date(0, 0, 0, 10, 30, 0, 0, TZShanghai())))
}

func TestTimeSlotGroup(t *testing.T) {
	group := NewTimeSlotGroup("09:00:00-09:30:00", "10:00:00-10:30:00")
	assert.Equal(t, true, group.IsValid(time.Date(0, 0, 0, 9, 0, 0, 0, TZShanghai())))
	assert.Equal(t, true, group.IsValid(time.Date(0, 0, 0, 9, 15, 0, 0, TZShanghai())))
	assert.Equal(t, false, group.IsValid(time.Date(0, 0, 0, 9, 30, 0, 0, TZShanghai())))
	assert.Equal(t, false, group.IsValid(time.Date(0, 0, 0, 10, 30, 0, 0, TZShanghai())))
	assert.Equal(t, false, group.IsValid(time.Date(0, 0, 0, 11, 30, 0, 0, TZShanghai())))
}
