package qdm

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizeWithFailed(t *testing.T) {
	assert := assert.New(t)

	svc := &service{
		client: resty.New().
			SetBaseURL("https://ecapis.qdm.cloud/api/v1"),
	}

	_, err := svc.Authorize("", "")
	if assert.Error(err) {
		assert.Equal(err.Error(), "Authentication failed")
	}
}
