package model

type ProcessStatus string

const (
	IN_QUEUE   ProcessStatus = "IN_QUEUE"
	PROCESSING ProcessStatus = "PROCESSING"
	SUCCESS    ProcessStatus = "SUCCESS"
	FAIL       ProcessStatus = "FAIL"
)
