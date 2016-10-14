package template

func opt(s ...string) (ret string) {
	if len(s) == 1 {
		ret = s[0]
	}
	return
}
