package dispatcher

import (
	"fmt"
	"time"
)

type Job struct {
	ID      int
	Message string
}

// Overwrite here
func (j *Job) Do() error {
	time.Sleep(2 * time.Second)
	fmt.Println(j)
	return nil
}
