package swagger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/menucha-de/capture"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/menucha-de/transport"
	"github.com/menucha-de/logging"
	"github.com/menucha-de/utils"
)

type wbSocketClients struct {
	Clients map[*websocket.Conn]bool
	Mu      sync.RWMutex
}

//Clients websocket opened connections
var Clients = wbSocketClients{Clients: make(map[*websocket.Conn]bool)}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var lg *logging.Logger = logging.GetLogger("gpio")

//GetConfiguration Gets configuration
func getDevices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	devices := capture.GetDevices()
	newdev := make(map[string]capture.Device)
	for m, device := range devices {

		ff := make(map[string]capture.Field)
		for n, field := range device.Fields {

			if field.Properties["direction"] == "INPUT" {
				id, _ := strconv.Atoi(field.ID)
				field.Value = GetState(id)
				//	newdev[m].Fields[n] = field
			}
			ff[n] = field
		}
		device.Fields = ff
		newdev[m] = device
	}
	var err = json.NewEncoder(w).Encode(newdev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	device := vars["device"]
	dev, _ := capture.GetDevice(device)
	var err = json.NewEncoder(w).Encode(dev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	dev, _ := capture.GetDevice(device)
	w.Write([]byte(dev.Label))

}
func setLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	var b bytes.Buffer
	n, err := b.ReadFrom(r.Body)
	if err != nil || n == 0 {
		http.Error(w, "Could not read label value", http.StatusBadRequest)
		return
	}
	err = setDeviceLabel(device, b.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
func deleteLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]

	err := setDeviceLabel(device, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func getProperties(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	dev, err := capture.GetDevice(device)
	if err != nil {
		http.Error(w, "Device with ID "+device+" does not exist", http.StatusInternalServerError)
		return
	}
	val := dev.Properties
	err = json.NewEncoder(w).Encode(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
func getProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	name := vars["name"]
	dev, err := capture.GetDevice(device)
	if err != nil {
		http.Error(w, "Device with ID "+device+" does not exist", http.StatusInternalServerError)
		return
	}
	props := dev.Properties
	if props == nil {
		http.Error(w, "Device with ID "+device+" has no properties", http.StatusInternalServerError)
		return
	}
	val, ok := props[name]
	if !ok {
		http.Error(w, "Device with ID "+device+" has no property "+name, http.StatusBadRequest)
		return
	}

	w.Write([]byte(val))

}
func setProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]

	name := vars["name"]

	var b bytes.Buffer
	n, err := b.ReadFrom(r.Body)
	if err != nil || n == 0 {
		http.Error(w, "Could not read property value", http.StatusBadRequest)
		return
	}
	err = setDeviceProperty(device, name, b.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func setFieldProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]
	name := vars["name"]

	var b bytes.Buffer
	n, err := b.ReadFrom(r.Body)
	if err != nil || n == 0 {
		http.Error(w, "Could not read property value", http.StatusBadRequest)
		return
	}
	err = setDeviceFieldProperty(device, field, name, b.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	dev, _ := capture.GetDevice(device)
	msg := WsMessage{device, dev.Fields[field]}
	HandleMessages(msg)
}
func getFieldProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]
	name := vars["name"]
	val, err := getDeviceFieldProperty(device, field, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(val))
}
func getFieldProperties(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]

	val, err := getDeviceFieldProperties(device, field)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getField(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]

	val, err := getDeviceField(device, field)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getFields(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]

	val, err := getDeviceFields(device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func setField(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]
	var fieldValue *capture.Field
	err := utils.DecodeJSONBody(w, r, &fieldValue)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			lg.WithError(err).Error("Failed to get properties")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	err = setDeviceField(device, field, *fieldValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	dev, _ := capture.GetDevice(device)
	msg := WsMessage{device, dev.Fields[field]}
	HandleMessages(msg)
}

func setFieldProperties(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]
	properties := make(map[string]string)
	err := utils.DecodeJSONBody(w, r, properties)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			lg.WithError(err).Error("Failed to get properties")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	err = setDeviceFieldProperties(device, field, properties)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	dev, _ := capture.GetDevice(device)
	msg := WsMessage{device, dev.Fields[field]}
	HandleMessages(msg)
}

func setFieldValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]

	var b bytes.Buffer
	n, err := b.ReadFrom(r.Body)
	if err != nil || n == 0 {
		http.Error(w, "Could not read property value", http.StatusBadRequest)
		return
	}
	err = setDeviceFieldValue(device, field, b.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	dev, _ := capture.GetDevice(device)
	msg := WsMessage{device, dev.Fields[field]}
	HandleMessages(msg)
}
func getFieldValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]

	v, err := getDeviceFieldValue(device, field)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	val := fmt.Sprintf("%v", v)
	w.Write([]byte(val))

}
func getFieldLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]

	v, err := getDeviceField(device, field)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(v.Label))

}

func setFieldLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]

	var b bytes.Buffer
	n, err := b.ReadFrom(r.Body)
	if err != nil || n == 0 {
		http.Error(w, "Could not read property value", http.StatusBadRequest)
		return
	}
	err = setDeviceFieldLabel(device, field, b.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	dev, _ := capture.GetDevice(device)
	msg := WsMessage{device, dev.Fields[field]}
	HandleMessages(msg)
}
func deleteFieldLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device := vars["device"]
	field := vars["field"]

	err := setDeviceFieldLabel(device, field, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	dev, _ := capture.GetDevice(device)
	msg := WsMessage{device, dev.Fields[field]}
	HandleMessages(msg)
}

/*func getKeepAliveConfiguration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	enc := json.NewEncoder(w)
	err := enc.Encode(config.KeepAliveConfiguration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func setKeepAliveConfiguration(w http.ResponseWriter, r *http.Request) {
	var conf KeepAliveConfiguration

	err := utils.DecodeJSONBody(w, r, &conf)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			lg.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	config.KeepAliveConfiguration = &conf
	w.WriteHeader(http.StatusNoContent)
}*/

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Make sure we close the connection when the function returns
	//defer ws.Close()
	// Register our new client
	Clients.Mu.Lock()
	defer Clients.Mu.Unlock()
	Clients.Clients[ws] = true
	if len(Clients.Clients) == 1 { //only at first connection
		enableInputs()
	}

}

//HandleMessages ws message handler
func HandleMessages(msg WsMessage) {
	Clients.Mu.Lock()
	defer Clients.Mu.Unlock()
	for client := range Clients.Clients {

		err := client.WriteJSON(msg)
		if err != nil {
			lg.Debug("WebSocket client connection lost", err.Error())
			client.Close()
			delete(Clients.Clients, client)
		}
	}
}

//HandleMessages1 ws message handler
func HandleMessages1(topic string, data capture.CaptureData) {

	field := capture.Field{}
	switch data.Value {
	case "0":
		field.Value = "LOW"
	case "1":
		field.Value = "HIGH"
	default:
		field.Value = "UNKNOWN"
	}
	field.ID = topic
	msg := WsMessage{DeviceID: "gpio", Field: field}

	HandleMessages(msg)
}

//reports
func getReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var err = json.NewEncoder(w).Encode(capture.GetCycles())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func addReport(w http.ResponseWriter, r *http.Request) {
	var report *capture.CycleSpec
	err := utils.DecodeJSONBody(w, r, &report)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			lg.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	err = checkFields(report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	id, err := capture.Define(*report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(id))
}
func getReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	s, err := capture.GetCycle(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewEncoder(w).Encode(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func setReport(w http.ResponseWriter, r *http.Request) {
	var report *capture.CycleSpec
	vars := mux.Vars(r)
	err := utils.DecodeJSONBody(w, r, &report)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			lg.WithError(err).Error("Failed to get report")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	id := vars["id"]

	if report.ID != id {
		http.Error(w, "ID of report does not match ", http.StatusBadRequest)
		return
	}
	err = checkFields(report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	err = capture.Update(*report)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	if report.Enabled {
		enableFields(*report)
	}
	w.WriteHeader(http.StatusNoContent)
}
func deleteReport(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	capture.Undefine(id)
	w.WriteHeader(http.StatusNoContent)
}

func getSubscriptors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	s, err := capture.GetSubscriptors(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = json.NewEncoder(w).Encode(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func addSubscriptor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var subscriptor *transport.Subscriptor

	err := utils.DecodeJSONBody(w, r, &subscriptor)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			lg.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	idd, err := capture.DefineSubscriptor(id, *subscriptor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(idd))
}

func getSubscriptor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	subID := vars["subscriptorId"]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	sub, err := capture.GetSubscriptor(id, subID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func setSubscriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "plain/text; charset=UTF-8")
	var subscr *transport.Subscriptor
	vars := mux.Vars(r)
	err := utils.DecodeJSONBody(w, r, &subscr)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			lg.WithError(err).Error("Failed to get report")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	id := vars["id"]
	subID := vars["subscriptorId"]

	err = capture.UpdateSubscriptor(id, subID, *subscr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func deleteSubscriptor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "plain/text; charset=UTF-8")

	vars := mux.Vars(r)
	id := vars["id"]
	subID := vars["subscriptorId"]
	err := capture.UndefineSubscriptor(id, subID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func getAPI(w http.ResponseWriter, r *http.Request) {
	file := "./www/openapi.yaml"
	http.ServeFile(w, r, file)
}
func mockPinValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pin := vars["id"]
	value := vars["value"]

	m := capture.CaptureData{Date: time.Now(), Device: "gpio", Field: pin, Value: value}
	capture.Pub.Publish(pin, m)
	HandleMessages1(pin, m)
	w.WriteHeader(http.StatusNoContent)
}
