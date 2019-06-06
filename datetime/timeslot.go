package datetime

import (
	"strings"
	"time"
)

const (
	layoutTimeSlot   = "15:04:05"
	SecondsPerHour   = 3600
	SecondsPerMinute = 60
)

type TimeSlot struct {
	Start int
	End   int
}

func NewTimeSlot(pattern string) *TimeSlot {
	parts := strings.SplitN(pattern, "-", 2)
	start, _ := time.Parse(layoutTimeSlot, parts[0])
	end, _ := time.Parse(layoutTimeSlot, parts[1])
	return &TimeSlot{
		Start: seconds(start),
		End:   seconds(end),
	}
}

// seconds 无视年月日，算出秒数
func seconds(t time.Time) int {
	return t.Hour()*SecondsPerHour + t.Minute()*SecondsPerMinute + t.Second()
}

func (s *TimeSlot) isValid(ts int) bool {
	return s.Start <= ts && ts < s.End
}

func (s *TimeSlot) IsValid(t time.Time) bool {
	return s.isValid(seconds(t))
}

type TimeSlotGroup struct {
	slots []*TimeSlot
}

func NewTimeSlotGroup(patterns ...string) *TimeSlotGroup {
	group := new(TimeSlotGroup)
	for _, pat := range patterns {
		group.slots = append(group.slots, NewTimeSlot(pat))
	}
	return group
}

func (g *TimeSlotGroup) IsValid(t time.Time) bool {
	ts := seconds(t)
	for _, slot := range g.slots {
		if slot.isValid(ts) {
			return true
		}
	}
	return false
}
