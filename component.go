package goat

import (
	"context"
)

type Component interface {
	Render(ctx context.Context, props any) VNode
}
