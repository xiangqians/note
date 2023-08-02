// validate test
// @author xiangqian
// @date `00:23 2023/08/03`
package validate

import "testing"

func TestUserName(t *testing.T) {
	userName := "1213-_ASD"
	err := UserName(userName)
	if err != nil {
		panic(err)
	}
}
