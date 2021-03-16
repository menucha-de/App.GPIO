package swagger

import (
	"errors"
	"strconv"

	capture "github.com/peramic/capture"
)

var count int = Count()

var config *capture.CycleConfiguration
var labels [8]string

func init() {
	labels = [...]string{"HS1", "HS2", "HS3", "HS4", "SWS1/SWD1", "SWS2/SWD2", "LS1", "LS2"}
	config = capture.GetConfig()
	if len(config.Devices) == 0 { //initialize
		dev := capture.Device{Name: "gpio"}
		dev.Fields = make(map[string]capture.Field)
		for i := 0; i < count; i++ {
			f := capture.Field{
				ID: strconv.Itoa(i + 1), Name: labels[i], Label: labels[i], Properties: map[string]string{"direction": "OUTPUT", "enabled": "false"}, Value: "LOW",
			}
			dev.Fields[f.ID] = f

			setDirection(i+1, "OUTPUT")
			setState(i+1, "LOW")
		}

		config.AddDevice(dev)
	} else { //set
		for _, device := range config.Devices {
			for _, field := range device.Fields {
				if len(field.Properties) > 0 {
					initialState, ok := field.Properties["initialState"]
					direction := field.Properties["direction"]
					id, _ := strconv.Atoi(field.ID)
					setDirection(id, MapDirection[direction])
					if direction == "OUTPUT" {
						if ok {
							setState(id, initialState)
							field.Value = initialState
						} else {
							setState(id, "LOW")
							field.Value = "LOW"
						}
						device.Fields[field.ID] = field
					}
				}
			}
			config.SetDevice(device)
		}
	}
	if len(config.Cycles) > 0 {
		for _, v := range config.Cycles {
			//err := capture.Add(v)
			//if err != nil {
			if v.Enabled {
				enableFields(v)
			}
			//}

		}
	}
}

func setDeviceLabel(device string, label string) error {

	dev, err := config.GetDevice(device)
	if err != nil {
		return err
	}
	dev.Label = label
	return config.SetDevice(dev)
}
func setDeviceProperty(device string, prop string, value string) error {

	dev, ok := config.Devices[device]
	if !ok {
		return errors.New("Device with id " + device + " does not exist")
	}
	if dev.Properties == nil {
		dev.Properties = make(map[string]string)
	}
	dev.Properties[prop] = value
	return config.SetDevice(dev)
}

func setDeviceFieldProperty(device string, field string, property string, value string) error {
	f, err := checkField(device, field, true)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(field)
	if err != nil {
		return errors.New("Wrong Pin" + field)
	}
	if f.Properties[property] == value {
		return nil
	}
	switch property {
	case "direction":
		direction, ok := MapDirection[value]
		if !ok {
			return errors.New("Unknown direction")

		}

		SetEnable(id, false)

		if direction == INPUT {
			delete(f.Properties, "initialState")
			setState(id, "LOW")
		}
		setDirection(id, direction)
		Clients.Mu.Lock()
		if direction == INPUT && len(Clients.Clients) != 0 {
			SetEnable(id, true)
		}
		Clients.Mu.Unlock()
		f.Value = GetState(id)

		f.Properties["direction"] = value

		config.SetDeviceField(device, f)

	case "initialState":
		s, ok := MapState[value]
		if !ok {
			return errors.New("Unknown state")

		}
		if s == UNKNOWN {
			delete(f.Properties, "initialState")
		} else {
			f.Properties["initialState"] = value
		}
		config.SetDeviceField(device, f)

	default:
		return errors.New("Unknown Property" + property)
	}
	return nil
}
func getDeviceFieldProperty(device string, field string, property string) (string, error) {
	f, err := checkField(device, field, false)
	if err != nil {
		return "", err
	}

	p, pok := f.Properties[property]
	if !pok {
		return "", errors.New("Unknown Property" + property)
	}
	return p, nil
}
func getDeviceFieldProperties(device string, field string) (map[string]string, error) {
	f, err := checkField(device, field, false)
	if err != nil {
		return nil, err
	}

	return f.Properties, nil
}
func setDeviceFieldProperties(device string, field string, properties map[string]string) error {
	f, err := checkField(device, field, true)
	if err != nil {
		return err
	}
	if properties == nil {
		return errors.New("Properties can't be null")
	}
	oldproperties := make(map[string]string)
	for k, v := range f.Properties {
		oldproperties[k] = v
	}
	for kk, vv := range properties {
		err := setDeviceFieldProperty(device, field, kk, vv)
		if err != nil {
			f.Properties = oldproperties
			config.SetDeviceField(device, f)
			return err
		}
	}
	return nil

}

func setDeviceFieldValue(device string, field string, value string) error {
	f, err := checkField(device, field, true)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(field)
	if err != nil {
		return errors.New("Wrong Pin" + field)
	}
	state, ok := MapState[value]
	if !ok || state == UNKNOWN {
		return errors.New("Unknown state")
	}

	if f.Properties["direction"] == "INPUT" {
		return errors.New("Cannot set INPUT State")

	}
	setState(id, value)
	f.Value = value
	err = config.SetDeviceField(device, f)
	if err != nil {
		return err
	}
	return nil
}
func getDeviceFieldValue(device string, field string) (interface{}, error) {
	_, err := checkField(device, field, false)
	if err != nil {
		return "", err
	}
	id, err := strconv.Atoi(field)
	if err != nil {
		return "", errors.New("Wrong Pin" + field)
	}
	return GetState(id), nil
}
func getDeviceFields(device string) (map[string]capture.Field, error) {
	dev, ok := config.Devices[device]
	if !ok {
		return nil, errors.New("Device with id " + device + " does not exist")
	}
	for idd, field := range dev.Fields {
		id := field.ID
		idx, err := strconv.Atoi(id)
		if err != nil {
			return nil, errors.New("Wrong Pin ID " + id)
		}
		field.Value = GetState(idx)
		dev.Fields[idd] = field
	}

	return dev.Fields, nil
}
func getDeviceField(device string, field string) (capture.Field, error) {
	f, err := checkField(device, field, false)
	if err != nil {
		return capture.Field{}, err
	}
	id, err := strconv.Atoi(field)
	if err != nil {
		return capture.Field{}, errors.New("Wrong Pin " + field)
	}
	f.Value = GetState(id)
	return f, nil
}
func setDeviceField(device string, field string, fieldvalue capture.Field) error {
	f, err := checkField(device, field, true)
	if err != nil {
		return err
	}
	if field != fieldvalue.ID {
		return errors.New("Field id doesn't match")
	}
	if fieldvalue.Properties == nil {
		return errors.New("Field properties can't be null")
	}
	p, ok := fieldvalue.Properties["direction"]
	if ok && p == "INPUT" {
		return errors.New("Can't set value to an INPUT pin")
	}
	if !ok && f.Properties["direction"] == "INPUT" {
		return errors.New("Can't set value to an INPUT pin")
	}
	str, ok := fieldvalue.Value.(string)
	if !ok || (str != "LOW" && str != "HIGH") {
		return errors.New("Wrong Field Value")
	}
	err = setDeviceFieldProperties(device, field, fieldvalue.Properties)
	if err != nil {
		return err
	}

	err = setDeviceFieldValue(device, field, str)
	if err != nil {
		return err
	}
	err = setDeviceFieldLabel(device, field, fieldvalue.Label)
	if err != nil {
		return err
	}

	return nil
}

func setDeviceFieldLabel(device string, field string, value string) error {

	f, err := checkField(device, field, true)
	if err != nil {
		return err
	}
	f.Label = value

	config.SetDeviceField(device, f)

	return nil
}
func enableInputs() {

	for i := 0; i < count; i++ {
		if GetDirection(i) == INPUT && getEnable(i) == false {
			SetEnable(i, true)
		}
	}
}

/////
func enableFields(report capture.CycleSpec) {
	for id, fields := range report.FieldSubscriptions {
		dev, ok := config.Devices[id]
		if ok {
			for _, p := range fields {
				f, fok := dev.Fields[p]
				if fok {
					if f.Properties["direction"] == "INPUT" {
						idd, _ := strconv.Atoi(f.ID)
						SetEnable(idd, true)
					}
				}
			}

		}

	}
}
func checkFields(cycle *capture.CycleSpec) error {
	if !cycle.Enabled {
		return nil
	}
	for d, f := range cycle.FieldSubscriptions {
		dev, ok := config.Devices[d]
		if !ok {
			return errors.New("Report device does not exist")
		}
		if len(f) == 0 {
			return errors.New("Cycle must have at least one field")
		}
		for _, fid := range f {
			field, fok := dev.Fields[fid]
			if !fok {
				return errors.New("Report field " + fid + " does not exist")
			}
			if field.Properties["direction"] != "INPUT" {
				return errors.New("Report field " + fid + " is not an INPUT field")
			}
		}

	}
	return nil
}
func checkField(device string, field string, checkEnabled bool) (capture.Field, error) {
	dev, ok := config.Devices[device]
	if !ok {
		return capture.Field{}, errors.New("Device with id " + device + " does not exist")
	}
	f, fok := dev.Fields[field]
	if !fok {
		return capture.Field{}, errors.New("Wrong Pin " + field)
	}

	if checkEnabled {
		p, ok := f.Properties["enabled"]

		if ok && p == "true" {

			return capture.Field{}, errors.New("Pin " + field + " is in use")
		}
	}
	return f, nil
}
