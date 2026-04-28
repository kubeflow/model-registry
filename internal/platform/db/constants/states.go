package constants

// ArtifactStateMapping maps artifact state strings to their integer values
// These correspond to the database integer values for artifact states
var ArtifactStateMapping = map[string]int32{
	"UNKNOWN":             0,
	"PENDING":             1,
	"LIVE":                2,
	"MARKED_FOR_DELETION": 3,
	"DELETED":             4,
	"ABANDONED":           5,
	"REFERENCE":           6,
}

// ArtifactStateNames maps artifact state integers to their string names
// These correspond to the database integer values for artifact states
var ArtifactStateNames = map[int32]string{
	0: "UNKNOWN",
	1: "PENDING",
	2: "LIVE",
	3: "MARKED_FOR_DELETION",
	4: "DELETED",
	5: "ABANDONED",
	6: "REFERENCE",
}

// ExecutionStateMapping maps execution state strings to their integer values
// These correspond to the database integer values for execution states
var ExecutionStateMapping = map[string]int32{
	"UNKNOWN":  0,
	"NEW":      1,
	"RUNNING":  2,
	"COMPLETE": 3,
	"FAILED":   4,
	"CACHED":   5,
	"CANCELED": 6,
}

// ExecutionStateNames maps execution state integers to their string names
// These correspond to the database integer values for execution states
var ExecutionStateNames = map[int32]string{
	0: "UNKNOWN",
	1: "NEW",
	2: "RUNNING",
	3: "COMPLETE",
	4: "FAILED",
	5: "CACHED",
	6: "CANCELED",
}
