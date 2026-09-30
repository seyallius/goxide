package result_test

import "testing"

// These grouped benchmarks keep each comparison under one workload. The
// sub-benchmark names are the variants consumed by benchbook, so a report can
// compare traditional and Result implementations directly.

func BenchmarkSuccess(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalSuccess)
	b.Run("result", benchmarkResultSuccess)
	b.Run("result_unwrap_or", benchmarkResultSuccessUnwrapOr)
}

func BenchmarkError(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalError)
	b.Run("result", benchmarkResultError)
}

func BenchmarkChainedSuccess(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalChainedSuccess)
	b.Run("result", benchmarkResultChainedSuccess)
}

func BenchmarkErrorHandling(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalErrorHandling)
	b.Run("result_with_try", benchmarkResultWithTry)
	b.Run("result_with_and_then", benchmarkResultWithAndThen)
}

func BenchmarkDBCreateUser(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalDBCreateUser)
	b.Run("result", benchmarkResultDBCreateUser)
}

func BenchmarkDBFindUser(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalDBFindUser)
	b.Run("result", benchmarkResultDBFindUser)
}

func BenchmarkDBFindUserNotFound(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalDBFindUserNotFound)
	b.Run("result", benchmarkResultDBFindUserNotFound)
}

func BenchmarkDBUpdateUser(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalDBUpdateUser)
	b.Run("result", benchmarkResultDBUpdateUser)
}

func BenchmarkDBGetOrCreateUser(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalDBGetOrCreateUser)
	b.Run("result", benchmarkResultDBGetOrCreateUser)
}

func BenchmarkDBChainedOperations(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalDBChainedOperations)
	b.Run("result", benchmarkResultDBChainedOperations)
	b.Run("result_bubble_up", benchmarkResultDBChainedOperationsBubbleUp)
}

func BenchmarkDBErrorHandlingWithFallback(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalDBErrorHandlingWithFallback)
	b.Run("result", benchmarkResultDBErrorHandlingWithFallback)
}

func BenchmarkDBCreateUserAllocs(b *testing.B) {
	b.Run("traditional", benchmarkTraditionalDBCreateUserAllocs)
	b.Run("result", benchmarkResultDBCreateUserAllocs)
}
