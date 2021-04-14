// +build integration

package test

import (
	"testing"
)

// Run runs a test with system keyring.
// nolint: thelper
func Run(
	t *testing.T,
	service, key string,
	expect RunExpect,
	test func(t *testing.T),
) {
	runTest(t, service, key, expect, test)
}
