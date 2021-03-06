package service

import (
	"context"
	"testing"
)

var sshPrivateKey = `-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----`

func TestCloneCode(t *testing.T) {
	err := PullCode(context.Background(), "house_scrapper", "git@github.com:ulysseskk/house.git", sshPrivateKey, "9a23220ff2597dbea933fd70e43f3a3fc1669cf9")
	if err != nil {
		panic(err)
	}
}
