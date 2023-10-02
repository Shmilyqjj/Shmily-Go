package main

func main() {
	m := map[string]string{"name": "John", "age": "30"}
	m1 := make(map[string]string, 1)
	m1["key"] = "val"
	for k, v := range m {
		println(k, v)
	}
}
