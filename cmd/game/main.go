package main

import (
	"fmt"
	"time"

	"gorl/fw/core/gem"
	input "gorl/fw/core/input/input_handling"
	"gorl/fw/core/logging"
	"gorl/fw/core/render"
	"gorl/fw/core/settings"
	"gorl/fw/core/store"
	"gorl/fw/physics"
	"gorl/game"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v3.3-core/gl"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	// PRE-INIT
	go func() {
		err := http.ListenAndServe("localhost:6969", nil)
		if err != nil {
			panic(err)
		}
	}()

	// settings
	settings_path := "settings.json"
	err := settings.LoadSettings(settings_path)
	if err != nil {
		fmt.Println("Error loading settings:", err)
		fmt.Println("Using fallback settings.")
		settings.UseFallbackSettings()
	}

	// logging
	logging.Init(settings.CurrentSettings().LogPath)
	logging.Info("Logging initialized")
	if err == nil {
		logging.Info("Settings loaded successfully.")
	} else {
		logging.Warning("Settings loading unsuccessful, using fallback.")
	}

	// INITIALIZATION
	// raylib window
	rl.InitWindow(
		int32(settings.CurrentSettings().ScreenWidth),
		int32(settings.CurrentSettings().ScreenHeight),
		settings.CurrentSettings().Title)
	defer rl.CloseWindow()
	rl.SetTargetFPS(int32(settings.CurrentSettings().TargetFps))

	// raw gl
	if err := gl.Init(); err != nil {
		logging.Fatal("Failed to initialize OpenGL:", err)
	}

	// rendering
	render.Init(rl.NewVector2(
		float32(settings.CurrentSettings().ScreenWidth),
		float32(settings.CurrentSettings().ScreenHeight)))
	defer render.Deinit()

	logging.Info("Rendering initialized.")

	// initialize audio
	//audio.InitAudio()
	//defer audio.DeinitAudio()

	// collision
	//collision.InitCollision()
	//defer collision.DeinitCollision()

	// physics
	physics.InitPhysics((1.0 / 60.0), rl.Vector2Zero(), (1.0 / 32.0))
	defer physics.DeinitPhysics()

	gem.Init()
	defer gem.Deinit()

	// lighting
	//lighting.InitLighting()
	//defer lighting.DeinitLighting()

	// animtion (premades need init and update)
	//animation.InitPremades(render.Rs.CurrentStage.Camera, render.GetWSCameraOffset())

	// register audio tracks
	//audio.RegisterMusic("aza-tumbleweeds", "audio/music/azakaela/azaFMP2_field7_Tumbleweeds.ogg")
	//audio.RegisterMusic("aza-outwest", "audio/music/azakaela/azaFMP2_scene1_OutWest.ogg")
	//audio.RegisterMusic("aza-frontier", "audio/music/azakaela/azaFMP2_town_Frontier.ogg")
	//audio.CreatePlaylist("main-menu", []string{"aza-tumbleweeds", "aza-outwest", "aza-frontier"})
	//audio.SetGlobalVolume(0.9)
	//audio.SetMusicVolume(0.7)
	//audio.SetSFXVolume(0.9)

	// gui
	//gui.InitBackend()

	// cursor
	//rl.HideCursor()

	// scenes
	//scenes.Sm.RegisterScene("dev", &scenes.DevScene{})
	//scenes.Sm.EnableScene("dev")

	//scenes.RegisterScene("some_name", &uscenes.TemplateScene{})
	//scenes.EnableScene("some_name")
	//scenes.DisableScene("some_name")

	//rl.DisableCursor()
	game.Init()

	// GAME LOOP
	//rl.SetExitKey(rl.KeyEnd) // Set a key to exit the game
	shouldExit := false

	// frame time measurement stuff
	//	frameStart := time.Now()
	//	var frameTime time.Duration = 0

	debugTex := rl.LoadRenderTexture(int32(settings.CurrentSettings().RenderWidth), int32(settings.CurrentSettings().RenderHeight))

	startTimer := 0.0

	for !shouldExit {
		//frameStart = time.Now()

		if startTimer < 5 {
			startTimer += float64(rl.GetFrameTime())
			rl.BeginDrawing()
			rl.EndDrawing()
			continue
		} else {
			startTimer = 6
		}

		rl.BeginTextureMode(debugTex)
		rl.ClearBackground(rl.Blank)

		shouldFixedUpdate := physics.Update()
		drawables, inputReceivers := gem.Traverse(shouldFixedUpdate)

		rl.EndTextureMode()

		//scenes.UpdateScenes() // TODO: rework scenes to be more clear
		//scenes.FixedUpdateScenes()

		rl.BeginDrawing()

		render.Draw(drawables)

		// input is processed at the end of the frame, because here we know in
		// what order the entities were drawn, and can be sure whatever the
		// user clicked was really visible at the front.
		//inputEventReceivers := append(inputReceivers, drawableInputReceivers...)
		input.HandleInputEvents(inputReceivers)

		// Draw Debug Info
		//DrawDebugInfo(frameTime)
		rl.DrawTexturePro(
			debugTex.Texture,
			rl.NewRectangle(0, 0, float32(debugTex.Texture.Width), -float32(debugTex.Texture.Height)),
			rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
			rl.Vector2Zero(),
			0, rl.White,
		)

		rl.EndDrawing()

		//audio.Update()
		//frameTime = time.Since(frameStart) // calculate after rl.EndDrawing() to include rendering time

		appState, ok := store.Get[*store.AppState]()
		shouldExit = rl.WindowShouldClose() || (!ok || appState.ShouldQuit)
	}

	//scenes.Sm.DisableAllScenes()
}

func DrawDebugInfo(frameTime time.Duration) {
	rl.DrawFPS(10, 10)
	rl.DrawText("dt: "+frameTime.String(), 10, 30, 20, rl.Lime)
	//physics.DrawColliders(true, true, true)
	//render.DebugDrawStageViewports(
	//	rl.NewVector2(10, 10), 4, render,
	//	[]*render.RenderStage{defaultRenderStage},
	//)
	//gem.DebugDrawEntities(rl.NewVector2(10, 50), 12)
	gem.DebugDrawHierarchy(rl.NewVector2(10, 50), 8)
}
