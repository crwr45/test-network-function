// Copyright (C) 2021 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package clusterrolebinding_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/test-network-function/test-network-function/pkg/tnf"
	crb "github.com/test-network-function/test-network-function/pkg/tnf/handlers/clusterrolebinding"
)

func Test_NewClusterRoleBinding(t *testing.T) {
	newCrb := crb.NewClusterRoleBinding(testTimeoutDuration, testServiceAccount, testPodNamespace)
	assert.NotNil(t, newCrb)
	assert.Equal(t, testTimeoutDuration, newCrb.Timeout())
	assert.Equal(t, newCrb.Result(), tnf.ERROR)
}

func Test_ReelFirstPositive(t *testing.T) {
	newCrb := crb.NewClusterRoleBinding(testTimeoutDuration, testServiceAccount, testPodNamespace)
	assert.NotNil(t, newCrb)
	firstStep := newCrb.ReelFirst()
	re := regexp.MustCompile(firstStep.Expect[0])
	matches := re.FindStringSubmatch(testInputFail)
	assert.Len(t, matches, 1)
	assert.Equal(t, testInputFail, matches[0])
}

func Test_ReelFirstPositiveEmpty(t *testing.T) {
	newCrb := crb.NewClusterRoleBinding(testTimeoutDuration, testServiceAccount, testPodNamespace)
	assert.NotNil(t, newCrb)
	firstStep := newCrb.ReelFirst()
	re := regexp.MustCompile(firstStep.Expect[0])
	matches := re.FindStringSubmatch(testInputSuccess)
	assert.Len(t, matches, 1)
	assert.Equal(t, testInputSuccess, matches[0])
}

func Test_ReelFirstNegative(t *testing.T) {
	newCrb := crb.NewClusterRoleBinding(testTimeoutDuration, testServiceAccount, testPodNamespace)
	assert.NotNil(t, newCrb)
	firstStep := newCrb.ReelFirst()
	re := regexp.MustCompile(firstStep.Expect[0])
	matches := re.FindStringSubmatch(testInputError)
	assert.Len(t, matches, 0)
}

func Test_ReelMatchSuccess(t *testing.T) {
	newCrb := crb.NewClusterRoleBinding(testTimeoutDuration, testServiceAccount, testPodNamespace)
	assert.NotNil(t, newCrb)
	step := newCrb.ReelMatch("", "", testInputSuccess)
	assert.Nil(t, step)
	assert.Equal(t, tnf.SUCCESS, newCrb.Result())
	assert.Len(t, newCrb.GetClusterRoleBindings(), 0)
}

func Test_ReelMatchFail(t *testing.T) {
	newCrb := crb.NewClusterRoleBinding(testTimeoutDuration, testServiceAccount, testPodNamespace)
	assert.NotNil(t, newCrb)
	step := newCrb.ReelMatch("", "", testInputFail)
	assert.Nil(t, step)
	assert.Equal(t, tnf.FAILURE, newCrb.Result())
	assert.Len(t, newCrb.GetClusterRoleBindings(), 3)
}

// Just ensure there are no panics.
func Test_ReelEof(t *testing.T) {
	newCrb := crb.NewClusterRoleBinding(testTimeoutDuration, testServiceAccount, testPodNamespace)
	assert.NotNil(t, newCrb)
	newCrb.ReelEOF()
}

const (
	testTimeoutDuration = time.Second * 2
	testPodNamespace    = "testPodNamespace"
	testServiceAccount  = "testServiceAccount"
	testInputError      = ""
	testInputSuccess    = "NAME\tSERVICE_ACCOUNTS"
	testInputFail       = `NAME	 SERVICE_ACCOUNTS
	test-builder                                                     map[kind:ServiceAccount name:testServiceAccount namespace:testPodNamespace]
	test-default                                                     map[kind:ServiceAccount name:testServiceAccount namespace:testPodNamespace]
	test-deployer                                                    map[kind:ServiceAccount name:testServiceAccount namespace:testPodNamespace]
	`
)
