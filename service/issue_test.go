package service

import (
	"github.com/stretchr/testify/assert"
	pb "github.com/ynishi/gdean/pb"
	"testing"
)

func TestShouldFetchIssue(t *testing.T) {
	issue := pb.Issue{}
	assert.NotNil(t, issue)
}
