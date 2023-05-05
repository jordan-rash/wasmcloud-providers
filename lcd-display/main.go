package main

import (
	"errors"

	lcd "github.com/wjessop/lcm1602_lcd"
	"golang.org/x/exp/io/i2c"

	display "github.com/jordan-rash/wasmcloud-interfaces/lcd-display"
	provider "github.com/wasmCloud/provider-sdk-go"
	core "github.com/wasmcloud/interfaces/core/tinygo"
	msgpack "github.com/wasmcloud/tinygo-msgpack"
)

var (
	p          *provider.WasmcloudProvider
	lcdDisplay *lcd.LCM1602LCD
	lcdDevice  *i2c.Device

	ErrInvalidOperation error = errors.New("invalid operation provided")
)

func main() {
	var err error

	p, err = provider.New(
		"jordan-rash:display",
		provider.WithNewLinkFunc(handleNewLink),
		provider.WithDelLinkFunc(handleDelLink),
		provider.WithHealthCheckMsg(healthCheckMsg),
		provider.WithProviderActionFunc(providerAction),
		provider.WithShutdownFunc(shutdown),
	)
	if err != nil {
		panic(err)
	}

	err = p.Start()
	if err != nil {
		panic(err)
	}
}

func healthCheckMsg() string {
	return ""
}

func providerAction(action provider.ProviderAction) (*provider.ProviderResponse, error) {
	resp := provider.ProviderResponse{}

	switch action.Operation {
	case "Display.DisplayLine":
		d := msgpack.NewDecoder(action.Msg)
		msg, err := display.MDecodeLine(&d)
		if err != nil {
			p.Logger.V(1).Info("failed to decode line: %s", err)
			return nil, err
		}

		p.Logger.V(1).Info("printing %s to line %d", msg.Text, msg.LineNumber)

		// text, row, position
		if err := lcdDisplay.WriteString(msg.Text, int(msg.LineNumber), 0); err != nil {
			return nil, err
		}

		var sizer msgpack.Sizer
		sizer.WriteBool(true)
		buf := make([]byte, sizer.Len())
		encoder := msgpack.NewEncoder(buf)
		encoder.WriteBool(true)

		resp.Msg = buf

		return &resp, nil
	case "Display.Clear":
		if err := lcdDisplay.Clear(); err != nil {
			return nil, err
		}

		var sizer msgpack.Sizer
		sizer.WriteBool(true)
		buf := make([]byte, sizer.Len())
		encoder := msgpack.NewEncoder(buf)
		encoder.WriteBool(true)

		resp.Msg = buf
		return &resp, nil
	default:
		return nil, ErrInvalidOperation
	}
}

func handleNewLink(linkdef core.LinkDefinition) error {
	var err error

	lcdDevice, err = i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x27)
	if err != nil {
		return err
	}

	lcdDisplay, err = lcd.NewLCM1602LCD(lcdDevice)
	if err != nil {
		return err
	}

	return nil
}

func handleDelLink(_ core.LinkDefinition) error {
	defer lcdDevice.Close()
	return nil
}

func shutdown() error {
	defer lcdDevice.Close()
	return nil
}
