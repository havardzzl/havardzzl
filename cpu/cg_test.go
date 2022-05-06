package main

import "testing"

func TestC(t *testing.T) {
	res := podRegexp.FindStringSubmatch("kubepods-burstable-pod503aa307_2ead_4099_be3a_6e824c92ab09.slice")
	for _, v := range res {
		t.Log(v)
	}
}
