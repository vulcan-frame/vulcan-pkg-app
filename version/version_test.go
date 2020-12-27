package globalvars

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubVersion(t *testing.T) {
	cases := []struct {
		v         string
		az        string
		sv        []int64
		isRelease bool
	}{

		{
			"us-v1.20230428_153224", "us", []int64{1, 20230428153224}, true,
		},
		{
			"us-v1.0", "us", []int64{1, 0}, true,
		},
		{
			"us-v0.0", "us", []int64{0, 0}, true,
		},
		{
			"us-v0.1", "us", []int64{0, 1}, true,
		},
		{
			"us-v1.0.1", "us", []int64{0, 0}, false,
		},
		{
			"us-1.0", "us", []int64{0, 0}, false,
		},
		{
			"us-v1", "us", []int64{0, 0}, false,
		},
		{
			"v1.0", "", []int64{0, 0}, false,
		},
		{
			"v1", "", []int64{0, 0}, false,
		},
	}

	for _, c := range cases {
		az, sv, isRelease := GetSubVersion(c.v)
		assert.Equal(t, c.az, az)
		assert.Equal(t, c.sv, sv)
		assert.Equal(t, c.isRelease, isRelease)
	}

}
