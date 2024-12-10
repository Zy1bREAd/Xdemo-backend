package api

func ConvertToCmdSlice(bs []byte) []string {
	var ss []string // StringSlice
	var tempSlice []byte
	for k, v := range bs {
		if v != 32 {
			tempSlice = append(tempSlice, v)
			if k == len(bs)-1 {
				ss = append(ss, string(tempSlice))
			}
			continue
		}
		if len(tempSlice) != 0 {
			ss = append(ss, string(tempSlice))
			tempSlice = []byte{}
		}
	}
	return ss
}
