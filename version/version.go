package globalvars

import (
	"strconv"
	"strings"
)

func GetSubVersion(v string) (az string, sv []int64, isRelease bool) {
	sv = make([]int64, 2)
	if len(v) == 0 {
		return
	}
	ss := strings.Split(v, "-")
	if len(ss) != 2 {
		return
	}

	az = ss[0]

	if strings.Index(ss[1], "v") != 0 {
		return
	}
	ss[1] = strings.Replace(ss[1], "v", "", 1)

	sss := strings.Split(ss[1], ".")
	if len(sss) != 2 {
		return
	}

	var err error

	sv[0], err = strconv.ParseInt(sss[0], 10, 64)
	if err != nil {
		return
	}

	sss[1] = strings.ReplaceAll(sss[1], "_", "")
	sv[1], err = strconv.ParseInt(sss[1], 10, 64)
	if err != nil {
		return
	}

	isRelease = true
	return
}
