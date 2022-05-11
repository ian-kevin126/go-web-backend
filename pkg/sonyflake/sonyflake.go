package main

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake"
)

var (
	sonyFlake     *sonyflake.Sonyflake
	sonyMachineID uint16
)

func getMachineID() (uint16, error) {
	return sonyMachineID, nil
}

// Init 需传入当前的机器ID
func Init(startTime string, machineId uint16) (err error) {
	sonyMachineID = machineId

	var st time.Time

	st, err = time.Parse("2006-01-02", startTime)

	if err != nil {
		return err
	}

	settings := sonyflake.Settings{
		StartTime: st,
		MachineID: getMachineID,
	}

	sonyFlake = sonyflake.NewSonyflake(settings)

	return
}

// GenID 生成id
func GenID() (id uint64, err error) {
	if sonyFlake == nil {
		err = fmt.Errorf("sonyflake not initted")
		return
	}

	id, err = sonyFlake.NextID()

	return
}

func main() {
	if err := Init("2020-07-01", 1); err != nil {
		fmt.Printf("Init failed, err: %v\n", err)
		return
	}

	id, _ := GenID()

	fmt.Println("id: ", id)
}
