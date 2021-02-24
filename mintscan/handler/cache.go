package handler

import (
	"context"
	"fmt"
)

var (
	BondDenom = ""
)

func SetBondDenom() error {
	if s == nil {
		fmt.Println("s is nil")
		return nil
	}
	denom, err := s.Client.GRPC.GetBondDenom(context.Background())
	if err != nil {
		return err
	}

	BondDenom = denom

	return nil
}
