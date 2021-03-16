//Package swagger This is a GPIO service implementation
package swagger

// State : The state
type State string

// List of State
const (
	LOW     State = "LOW"
	HIGH    State = "HIGH"
	UNKNOWN State = "UNKNOWN"
)

//MapState state map
var MapState = map[string]State{
	"LOW":     LOW,
	"HIGH":    HIGH,
	"UNKNOWN": UNKNOWN,
}
