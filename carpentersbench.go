/*
  To the extent possible under law, Ian Davis has waived all copyright
  and related or neighboring rights to this Amberfell Source Code file.
  This work is published from the United Kingdom. 
*/
package main

type CarpentersBench struct {
	pos Vectori
}

func NewCarpentersBench(pos Vectori) *CarpentersBench {
	obj := CarpentersBench{pos: pos}
	return &obj
}

func (self *CarpentersBench) Label() string {
	return "Carpenter's Bench"
}

func (self *CarpentersBench) Recipes() []Recipe {
	return carpenterRecipes
}
