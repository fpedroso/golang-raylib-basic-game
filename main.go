package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	windowWidth  = 1024
	windowHeight = 768
	fps          = 60
	camZoom      = 3
	playerSpeed  = 3
	windowTitle  = "Go Animal Crossing Tutorial"
)

var (
	grassTexture                                  rl.Texture2D
	fencesTexture                                 rl.Texture2D
	tex                                           rl.Texture2D
	playerTexture                                 rl.Texture2D
	playerSrcRec                                  rl.Rectangle
	playerDestRec                                 rl.Rectangle
	music                                         rl.Music
	cam                                           rl.Camera2D
	gamePaused                                    bool
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerRight, playerLeft bool
	playerFrame                                   int
	frameCount                                    int
	tileDest                                      rl.Rectangle
	tileSrc                                       rl.Rectangle
	tileMap                                       []int
	srcMap                                        []string
	mapW, mapH                                    int
	camTarget                                     rl.Vector2

	backgroundColor = rl.NewColor(147, 211, 196, 255)
	running         = true
)

func drawScene() {
	for i := 0; i < len(tileMap); i++ {
		if tileMap[i] != 0 {
			tileDest.X = tileDest.Width * float32(i%mapW)
			tileDest.Y = tileDest.Height * float32(i/mapH)
			if srcMap[i] == "g" {
				tex = grassTexture
			}
			if srcMap[i] == "f" {
				tex = fencesTexture
				tileSrc.X = 16
				tileSrc.Y = 16
				rl.DrawTexturePro(grassTexture, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
			}
			tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))
			tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Width)))
			rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(0, 0), 0, rl.White)
		}
	}

	rl.DrawTexturePro(playerTexture, playerSrcRec, playerDestRec, rl.NewVector2(0, 0), 0, rl.White)

	if gamePaused {
		text := "PAUSED!"
		fontSize := 48
		measures := rl.MeasureTextEx(rl.GetFontDefault(), text, float32(fontSize), 1)
		rl.DrawText(text, int32(cam.Target.X)-int32(measures.X/2), int32(cam.Target.Y)-int32(measures.Y/2), int32(fontSize), rl.Black)
	}
}

func input() {
	// pause
	if rl.IsKeyPressed(rl.KeyP) {
		gamePaused = !gamePaused
	}

	if gamePaused {
		return
	}

	// character movement
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = 0
		playerDown = true
	}
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = 1
		playerUp = true
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = 2
		playerLeft = true
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = 3
		playerRight = true
	}
}

func update() {
	running = !rl.WindowShouldClose()

	playerSrcRec.X = playerSrcRec.Width * float32(playerFrame)

	if playerMoving {
		if playerDown {
			playerDestRec.Y += playerSpeed
		}
		if playerUp {
			playerDestRec.Y -= playerSpeed
		}
		if playerLeft {
			playerDestRec.X -= playerSpeed
		}
		if playerRight {
			playerDestRec.X += playerSpeed
		}
		if frameCount%8 == 1 {
			playerFrame++
		}
	} else if frameCount%45 == 1 {
		playerFrame++
	}

	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}
	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
	}

	playerSrcRec.X = playerSrcRec.Width * float32(playerFrame)
	playerSrcRec.Y = playerSrcRec.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if gamePaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	camTarget = rl.NewVector2(float32(playerDestRec.X+(playerDestRec.Width/2)), float32(playerDestRec.Y+(playerDestRec.Height/2)))
	cam.Target = camTarget

	playerMoving = false
	playerUp, playerLeft, playerRight, playerDown = false, false, false, false
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(backgroundColor)
	rl.BeginMode2D(cam)
	drawScene()
	rl.EndMode2D()
	rl.EndDrawing()
}

func loadMap(mapFile string) {
	file, err := os.ReadFile(mapFile)

	if err != nil {
		log.Fatal(err)
	}

	remNewLines := strings.ReplaceAll(strings.ReplaceAll(string(file), "\r\n", "\n"), "\n", " ")
	sliced := strings.Split(remNewLines, " ")
	mapW = -1
	mapH = -1

	for i := 0; i < len(sliced); i++ {
		s, err := strconv.ParseInt(sliced[i], 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		m := int(s)

		if mapW == -1 {
			mapW = m
		} else if mapH == -1 {
			mapH = m
		} else if i < mapW*mapH+2 {
			tileMap = append(tileMap, m)
		} else {
			srcMap = append(srcMap, sliced[i])
		}
	}

	if len(tileMap) > mapW*mapH {
		tileMap = tileMap[:len(tileMap)-1]
	}
}

func start() {
	rl.InitWindow(windowWidth, windowHeight, windowTitle)
	rl.SetTargetFPS(fps)

	grassTexture = rl.LoadTexture("resources/sprites/tilesets/grass.png")
	fencesTexture = rl.LoadTexture("resources/sprites/tilesets/fences.png")
	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	playerTexture = rl.LoadTexture("resources/sprites/characters/basic_character_spritesheet.png")
	playerSrcRec = rl.NewRectangle(0, 0, 48, 48)
	playerDestRec = rl.NewRectangle(0, 0, 48, 48)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("resources/audio/silly_fun.mp3")
	gamePaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(
		rl.NewVector2(float32(windowWidth/2), float32(windowHeight/2)),
		camTarget,
		0,
		camZoom,
	)

	loadMap("one.map")
}

func quit() {
	rl.UnloadTexture(grassTexture)
	rl.UnloadTexture(fencesTexture)
	rl.UnloadTexture(playerTexture)

	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()

	rl.CloseWindow()
}

func centerWindow() {
	monitor := rl.GetCurrentMonitor()
	monitorWidth := rl.GetMonitorWidth(monitor)
	monitorHeight := rl.GetMonitorHeight(monitor)
	rl.SetWindowPosition((monitorWidth/2)-(windowWidth/2), (monitorHeight/2)-(windowHeight/2))
}

func main() {
	start()
	centerWindow()
	defer quit()

	for running {
		input()
		update()
		render()
	}
}
