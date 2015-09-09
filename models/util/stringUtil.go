package util


func StrJoin(strs []string, joint string)string{
	ret := ""
	for idx, s := range strs{
		if idx > 0{
			ret = ret + joint
		}
		ret = ret + s
	}
	return ret
}