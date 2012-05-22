/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

type Forge struct {
	pos Vectori
}

func NewForge(pos Vectori) *Forge {
	obj := Forge{pos: pos}
	return &obj
}

func (self *Forge) Label() string {
	return "Forge"
}

func (self *Forge) Recipes() []Recipe {
	return forgeRecipes
}
