package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	_ "image/png"
	"planeGame/models"
	"planeGame/views"
	"sync"
	"time"
)

type Coin struct{ X, Y float64 }
type Obstacle struct{ X, Y float64 }
type Airplane struct{ X, Y float64 }

var LastCollisionTime time.Time
var Mu sync.Mutex

func Run() {
	cfg := pixelgl.WindowConfig{Title: "Aventura Aérea: Cazador de Monedas", Bounds: pixel.R(0, 0, 800, 600)}
	win, _ := pixelgl.NewWindow(cfg)
	skyPic, _ := views.LoadPicture("Assets/sky.png")
	skySprite := pixel.NewSprite(skyPic, skyPic.Bounds())
	airplanePicRight, _ := views.LoadPicture("Assets/JetRight.png")
	airplanePicLeft, _ := views.LoadPicture("Assets/JetLeft.png")
	airplaneSpriteRight := pixel.NewSprite(airplanePicRight, airplanePicRight.Bounds())
	airplaneSpriteLeft := pixel.NewSprite(airplanePicLeft, airplanePicLeft.Bounds())
	rockPic, _ := views.LoadPicture("Assets/rock.png")
	rockSprite := pixel.NewSprite(rockPic, rockPic.Bounds())
	coinPic, _ := views.LoadPicture("Assets/coin.png")
	coinSprite := pixel.NewSprite(coinPic, coinPic.Bounds())
	explosionPic, _ := views.LoadPicture("Assets/explosion.png")
	explosionSprite := pixel.NewSprite(explosionPic, explosionPic.Bounds())
	gameStatus := "running"
	obstacles, coins, airplane := models.InitGameObjects()
	score := 0
	lives := 3
	showExplosion := false
	explosionTime := time.Time{}
	direction := "right"
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	go models.MoveObstacles(&obstacles)
	go models.MoveCoins(&coins)
	go models.ManageGame(&gameStatus, &score, &lives, &obstacles, &airplane, &showExplosion, &explosionTime)

	for !win.Closed() {
		win.Clear(colornames.Skyblue)
		skySprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 1.5).Moved(pixel.V(400, 300)))
		Mu.Lock()
		for _, obstacle := range obstacles {
			rockSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.02).Moved(pixel.V(obstacle.X, obstacle.Y)))
		}
		for i, coin := range coins {
			if gameStatus == "running" && (coin.X-airplane.X)*(coin.X-airplane.X)+(coin.Y-airplane.Y)*(coin.Y-airplane.Y) < 2500 {
				coins = append(coins[:i], coins[i+1:]...)
				score++
				break
			}
			coinSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.1).Moved(pixel.V(coin.X, coin.Y)))
		}
		if direction == "right" {
			airplaneSpriteRight.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.05).Moved(pixel.V(airplane.X, airplane.Y)))
		}
		if direction == "left" {
			airplaneSpriteLeft.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.05).Moved(pixel.V(airplane.X, airplane.Y)))
		}
		if showExplosion && time.Since(explosionTime).Seconds() < 0.5 {
			explosionSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.05).Moved(pixel.V(airplane.X, airplane.Y)))
		} else {
			showExplosion = false
		}
		Mu.Unlock()
		txt := text.New(pixel.V(50, 550), atlas)
		txt.Color = colornames.Black
		fmt.Fprintf(txt, "Monedas obtenidas: %d Vidas: %d", score, lives)
		txt.Draw(win, pixel.IM)
		if gameStatus == "won" {
			txt := text.New(pixel.V(100, 300), atlas)
			txt.Color = colornames.Black
			fmt.Fprintf(txt, "Haz ganado! presiona P para reiniciar")
			txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2))
		}
		if gameStatus == "lost" {
			txt := text.New(pixel.V(100, 300), atlas)
			txt.Color = colornames.Black
			fmt.Fprintf(txt, "Juego terminado. Presiona P para reiniciar")
			txt.Draw(win, pixel.IM.Scaled(txt.Orig, 2))
		}
		win.Update()
		if win.JustPressed(pixelgl.KeyP) {
			gameStatus = "running"
			obstacles, coins, airplane = models.InitGameObjects()
			score = 0
			lives = 3
		}
		if gameStatus == "running" {
			if win.Pressed(pixelgl.KeyUp) && airplane.Y <= 575 {
				airplane.Y += 10
			}
			if win.Pressed(pixelgl.KeyDown) && airplane.Y >= 25 {
				airplane.Y -= 10
			}
			if win.Pressed(pixelgl.KeyLeft) && airplane.X >= 25 {
				airplane.X -= 10
				direction = "left"
			}
			if win.Pressed(pixelgl.KeyRight) && airplane.X <= 775 {
				airplane.X += 10
				direction = "right"
			}
		}
		time.Sleep(time.Millisecond * 16)
	}
}
