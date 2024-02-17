package main

func PrintString(object MalObject) (string, error) {
	return object.String(), nil
	// switch object.(type) {
	// case MalAtom:
	// 	return object.(MalAtom).Token(), nil

	// case MalList:
	// 	var ret []string
	// 	for _, form := range object.(MalList).Objects() {
	// 		s, err := PrintString(form)
	// 		if err != nil {
	// 			return "", err
	// 		}
	// 		ret = append(ret, s)
	// 	}
	// 	return fmt.Sprintf("(%s)", strings.Join(ret, " ")), nil

	// default:
	// 	return "", fmt.Errorf("unsupported MalObject type: %v", object)
	// }
}
