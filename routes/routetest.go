package routes

import "github.com/ISBachvarov21/Rune-GO/runehttp"

// rune get /
func Root(c *runehttp.Context) (runehttp.Response, error) {
	return runehttp.Text("ok "), nil
}

// rune post /
func Post(c *runehttp.Context) (runehttp.Response, error) {
	return runehttp.Text("ok but post"), nil
}
