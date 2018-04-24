package object

func FullyQ(ns, name string) string {
	if ns == "" {
		return name
	}
	return ns + "\\" + name
}
