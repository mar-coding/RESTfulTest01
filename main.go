package main

func main() {
	a := App{}
	dsn := LoadDB()
	a.Initialize("mysql", dsn)
	a.Run(":8010")
}
