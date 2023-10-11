package models

import (
	_ "image/png"
	"math/rand"
	"planeGame/entities"
	"time"
)

func InitGameObjects() ([]entities.Obstacle, []entities.Coin, entities.Airplane) {
	obstacles := make([]entities.Obstacle, 10)
	coins := make([]entities.Coin, 10)
	airplane := entities.Airplane{X: 50, Y: 300}
	for i := range obstacles {
		for {
			obstacle := entities.Obstacle{X: rand.Float64()*740 + 30, Y: rand.Float64()*540 + 30}
			if (obstacle.X-airplane.X)*(obstacle.X-airplane.X)+(obstacle.Y-airplane.Y)*(obstacle.Y-airplane.Y) >= 2500 {
				obstacles[i] = obstacle
				break
			}
		}
		coins[i] = entities.Coin{X: rand.Float64()*740 + 30, Y: rand.Float64()*540 + 30}
	}
	return obstacles, coins, airplane
}

func MoveObstacles(obstacles *[]entities.Obstacle) {
	for {
		time.Sleep(time.Millisecond * 20)
		entities.Mu.Lock()
		for i := range *obstacles {
			newX := (*obstacles)[i].X + rand.Float64()*20 - 10
			newY := (*obstacles)[i].Y + rand.Float64()*20 - 10
			if newX >= 25 && newX <= 775 {
				(*obstacles)[i].X = newX
			}
			if newY >= 25 && newY <= 575 {
				(*obstacles)[i].Y = newY
			}
		}
		entities.Mu.Unlock()
	}
}

func MoveCoins(coins *[]entities.Coin) {
	for {
		time.Sleep(time.Millisecond * 20)
		entities.Mu.Lock()
		for i := range *coins {
			newX := (*coins)[i].X + rand.Float64()*15 - 7.5
			newY := (*coins)[i].Y + rand.Float64()*15 - 7.5
			if newX >= 25 && newX <= 775 {
				(*coins)[i].X = newX
			}
			if newY >= 25 && newY <= 575 {
				(*coins)[i].Y = newY
			}
		}
		entities.Mu.Unlock()
	}
}

func ManageGame(gameStatus *string, score *int, lives *int, obstacles *[]entities.Obstacle, airplane *entities.Airplane, showExplosion *bool, explosionTime *time.Time) {
	for {
		time.Sleep(time.Millisecond * 50)
		entities.Mu.Lock()
		for _, obstacle := range *obstacles {
			if *gameStatus == "running" && time.Since(entities.LastCollisionTime) > time.Second && (obstacle.X-airplane.X)*(obstacle.X-airplane.X)+(obstacle.Y-airplane.Y)*(obstacle.Y-airplane.Y) < 2500 {
				*lives--
				entities.LastCollisionTime = time.Now()
				*showExplosion = true
				*explosionTime = time.Now()
				if *lives <= 0 {
					*gameStatus = "lost"
				}
			}
		}
		if *score >= 10 {
			*gameStatus = "won"
		}
		entities.Mu.Unlock()
	}
}
