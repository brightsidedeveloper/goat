package goat

import (
	"context"
	"io"
)

type Component interface {
	Render(ctx context.Context, w io.Writer, props any) error
}
