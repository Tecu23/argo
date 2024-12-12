package uci

import (
	"fmt"
	"strconv"
)

type Option interface {
	UciName() string
	UciString() string
	Set(s string) error
}

type BoolOption struct {
	Name  string
	Value *bool
}

func (opt *BoolOption) UciName() string {
	return opt.Name
}

func (opt *BoolOption) UciString() string {
	return fmt.Sprintf("option name %v type %v default %v",
		opt.Name, "check", *opt.Value)
}

func (opt *BoolOption) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*opt.Value = v
	return nil
}
