package main

func main() {
	sg := NewSVGGrid()
	//sg.RenderGrid(&CircularGradient{})
	ic, err := NewImageContent("images/ada-lovelace.jpg")
	if err != nil {
		panic(err)
	}
	sg.RenderGrid(ic)
}
