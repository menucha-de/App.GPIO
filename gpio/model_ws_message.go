package swagger

import capture "github.com/menucha-de/capture"

// WsMessage The websocket message
type WsMessage struct {
	// The keep-alive enable state
	DeviceID string `json:"deviceId,omitempty"`
	// The interval of the keep-alive event in milliseconds
	Field capture.Field `json:"field,omitempty"`
}
