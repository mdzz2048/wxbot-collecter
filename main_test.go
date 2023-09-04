package main

import (
	"testing"
)

func TestSaveToCubox(t *testing.T) {
	url := "https://cloud.tencent.com/developer/article/1849807"
	SaveToCubox(url)
}
