package main

import (
	_ "image/png"
	"math/rand"
	"planeGame/scenes"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func main() { rand.Seed(time.Now().UnixNano()); pixelgl.Run(scenes.Run) }
