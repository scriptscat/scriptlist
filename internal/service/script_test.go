package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_script_requestSyncUrl(t *testing.T) {
	s := &script{}

	code, err := s.requestSyncUrl("https://raw.githubusercontent.com/CodFrm/StudyGit/master/src/test.js")
	assert.Nil(t, err)
	assert.NotEmpty(t, code)

}
