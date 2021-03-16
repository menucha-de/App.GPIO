//Package swagger This is a GPIO service implementation
package swagger

// Direction : The direction
type Direction string

// List of Direction
const (
	INPUT  Direction = "INPUT"
	OUTPUT Direction = "OUTPUT"
)

// MapDirection :Direction Map
var MapDirection = map[string]Direction{
	"INPUT":  INPUT,
	"OUTPUT": OUTPUT,
}
