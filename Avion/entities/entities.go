package entities

import (
	"sync"
	"time"
)

type Coin struct{ X, Y float64 }
type Obstacle struct{ X, Y float64 }
type Airplane struct{ X, Y float64 }

var Mu sync.Mutex
var LastCollisionTime time.Time
