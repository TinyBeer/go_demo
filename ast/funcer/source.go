package main

import (
	"context"
	"strings"
)

type Foo struct {
}

func (*Foo) NeedContext(ctx context.Context) {

}

func (*Foo) NotNeedContext() {

}

func ContextWanted(name string) string {
	return strings.TrimSpace(name)
}
