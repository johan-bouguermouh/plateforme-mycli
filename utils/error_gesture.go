package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/urfave/cli"
)

func OnUsageError(ctx *cli.Context, err error, _ bool) error {
	type subCommandHelp struct {
		flagName string
		usage    string
	}

	// Calculate the maximum width of the flag name field
	// for a good looking printing
	vflags := visibleFlags(ctx.Command.Flags)
	help := make([]subCommandHelp, len(vflags))
	maxWidth := 0
	for i, f := range vflags {
		s := strings.Split(f.String(), "\t")
		if len(s[0]) > maxWidth {
			maxWidth = len(s[0])
		}

		help[i] = subCommandHelp{flagName: s[0], usage: s[1]}
	}
	maxWidth += 2

	var errMsg strings.Builder

	// Do the good-looking printing now
	fmt.Fprintln(&errMsg, "Invalid command usage,", err.Error())
	if len(help) > 0 {
		fmt.Fprintln(&errMsg, "\nSUPPORTED FLAGS:")
		for _, h := range help {
			spaces := string(bytes.Repeat([]byte{' '}, maxWidth-len(h.flagName)))
			fmt.Fprintf(&errMsg, "   %s%s%s\n", h.flagName, spaces, h.usage)
		}
	}
	//console.Fatal(errMsg.String())
	return err
}

func visibleFlags(fl []cli.Flag) []cli.Flag {
	visible := []cli.Flag{}
	for _, flag := range fl {
		field := flagValue(flag).FieldByName("Hidden")
		if !field.IsValid() || !field.Bool() {
			visible = append(visible, flag)
		}
	}
	return visible
}

func flagValue(f cli.Flag) reflect.Value {
	fv := reflect.ValueOf(f)
	for fv.Kind() == reflect.Ptr {
		fv = reflect.Indirect(fv)
	}
	return fv
}