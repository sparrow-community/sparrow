package processing

import (
	"time"

	"github.com/sony/sonyflake/v2"
)

var sf *sonyflake.Sonyflake

func InitSonyflake(startTime time.Time, machineID func() (int, error)) {
	settings := sonyflake.Settings{
		StartTime: startTime,
		MachineID: machineID,
	}
	var err error
	sf, err = sonyflake.New(settings)
	if err != nil {
		panic(err)
	}
	if sf == nil {
		panic("sonyflake not created")
	}
}

// NextID generates a new unique ID using Sonyflake
func NextID() int64 {
	if sf == nil {
		panic("Sonyflake not initialized")
	}
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return id
}
