// Copyright (c) 2025 SeyedAli
// Licensed under the MIT License. See LICENSE file in the project root for details.

package option_test

import (
	"testing"

	"github.com/seyallius/goxide/rusty/option"
)

// largePayload exposes the copy cost of inline storage alongside its allocation
// savings. Its size is deliberately representative rather than a public API type.
type largePayload struct {
	data [256]byte
}

var (
	optionIntSink    option.Option[int]
	optionLargeSink  option.Option[largePayload]
	intSink          int
	largePayloadSink largePayload
)

func BenchmarkOption(b *testing.B) {
	b.Run("SomeInt", benchmarkOptionSomeInt)
	b.Run("SomeLarge", benchmarkOptionSomeLarge)
	b.Run("NoneInt", benchmarkOptionNoneInt)
	b.Run("UnwrapInt", benchmarkOptionUnwrapInt)
	b.Run("UnwrapLarge", benchmarkOptionUnwrapLarge)
	b.Run("UnwrapOrNone", benchmarkOptionUnwrapOrNone)
	b.Run("MapSomeInt", benchmarkOptionMapSomeInt)
	b.Run("MapNoneInt", benchmarkOptionMapNoneInt)
	b.Run("FlatMapSomeInt", benchmarkOptionFlatMapSomeInt)
	b.Run("FlatMapNoneInt", benchmarkOptionFlatMapNoneInt)
}

func benchmarkOptionSomeInt(b *testing.B) {
	b.ReportAllocs()
	for i := 0; b.Loop(); i++ {
		optionIntSink = option.Some(i)
	}
}

func benchmarkOptionSomeLarge(b *testing.B) {
	var payload largePayload
	b.ReportAllocs()
	for i := 0; b.Loop(); i++ {
		payload.data[0] = byte(i)
		optionLargeSink = option.Some(payload)
	}
}

func benchmarkOptionNoneInt(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		optionIntSink = option.None[int]()
	}
}

func benchmarkOptionUnwrapInt(b *testing.B) {
	value := option.Some(42)
	b.ReportAllocs()
	for b.Loop() {
		intSink = value.Unwrap()
	}
}

func benchmarkOptionUnwrapLarge(b *testing.B) {
	value := option.Some(largePayload{data: [256]byte{0: 42}})
	b.ReportAllocs()
	for b.Loop() {
		largePayloadSink = value.Unwrap()
	}
}

func benchmarkOptionUnwrapOrNone(b *testing.B) {
	value := option.None[int]()
	b.ReportAllocs()
	for b.Loop() {
		intSink = value.UnwrapOr(42)
	}
}

func benchmarkOptionMapSomeInt(b *testing.B) {
	value := option.Some(42)
	b.ReportAllocs()
	for b.Loop() {
		optionIntSink = option.Map(value, incrementOptionInt)
	}
}

func benchmarkOptionMapNoneInt(b *testing.B) {
	value := option.None[int]()
	b.ReportAllocs()
	for b.Loop() {
		optionIntSink = option.Map(value, incrementOptionInt)
	}
}

func benchmarkOptionFlatMapSomeInt(b *testing.B) {
	value := option.Some(42)
	b.ReportAllocs()
	for b.Loop() {
		optionIntSink = option.FlatMap(value, someIncrementedOptionInt)
	}
}

func benchmarkOptionFlatMapNoneInt(b *testing.B) {
	value := option.None[int]()
	b.ReportAllocs()
	for b.Loop() {
		optionIntSink = option.FlatMap(value, someIncrementedOptionInt)
	}
}

func incrementOptionInt(value int) int { return value + 1 }

func someIncrementedOptionInt(value int) option.Option[int] {
	return option.Some(value + 1)
}
