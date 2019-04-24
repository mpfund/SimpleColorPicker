package main

import (
	"fmt"
	"github.com/kbinani/screenshot"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"image"
	"image/color"
	"log"
	"os"
	"time"
)

type Graphics struct{
	font *ttf.Font
	renderer *sdl.Renderer
}

type appInit struct{
	width int32
	height int32
	title string
	colorWindowSize int32
	textColor sdl.Color
}

type appSettings struct{
	running bool
	g *Graphics
	window *sdl.Window
	tColor colorful.Color
	texture *sdl.Texture
	init *appInit
	buttons []*Button
	logLabel *LogLabel
	hasFocus bool
}

func main() {
	os.Exit(run())
}

func run() int {
	appInit := &appInit{}
	appInit.colorWindowSize = int32(25)
	appInit.textColor = sdl.Color{140,140,140,0}
	appInit.width = 450
	appInit.height = 140
	appInit.title = "SimpleColorPicker"

	app,err := initSDL(appInit)
	app.buttons = []*Button{
		{
			&sdl.Rect{},
			"Copy",
			func(){ copyValue(app,"rgb")},
			appInit.textColor,
		},
		{
			&sdl.Rect{},
			"Copy",
			func(){ copyValue(app,"hex")},
			appInit.textColor,
		},
		{
			&sdl.Rect{},
			"Copy",
			func(){ copyValue(app,"hsl")},
			appInit.textColor,
		},
	}
	app.logLabel = &LogLabel{
		Rect:&sdl.Rect{10,appInit.height-20,170,20},
		Text:"hold CTRL to pick color",
		VisibleDuration:1*time.Second,
		ForeColor:sdl.Color{255,51,0,0},
	}

	app.logLabel.Show()

	if err != nil{
		log.Fatal(err)
	}

	app.running = true

	for app.running {
		event := sdl.PollEvent() // wait here until an event is in the event queue
		switch t:= event.(type) {
		case *sdl.QuitEvent:
			app.running = false
			return 0
		case *sdl.MouseButtonEvent:
			if t.Button==1{
				ForEachButton(app.buttons,func(b*Button){
					CheckClickEvent(t,b)})
			}
		case *sdl.WindowEvent:
			if t.Event == sdl.WINDOWEVENT_FOCUS_GAINED{
				app.hasFocus = true
			}
			if t.Event == sdl.WINDOWEVENT_FOCUS_LOST{
				app.hasFocus = false
			}
		}

		renderLoop(app)
	}

	return 0
}

func initSDL(appInit *appInit) (*appSettings,error) {
	var err error
	sdl.Init(sdl.INIT_EVERYTHING)

	app:=&appSettings{}
	if err := ttf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize TTF: %s\n", err)
		return nil,err
	}

	app.window, err = sdl.CreateWindow(appInit.title, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		appInit.width, appInit.height, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return nil,err
	}

	renderer, err := sdl.CreateRenderer(app.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return nil,err
	}

	font, err := ttf.OpenFont("test.ttf", 16);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open font: %s\n", err)
		return  nil,err
	}


	app.texture,_ =renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		2*appInit.colorWindowSize,2*appInit.colorWindowSize)

	app.g = &Graphics{font,renderer}
	app.tColor,_ = colorful.MakeColor(color.Black)
	app.init = appInit
	app.hasFocus = true
	return app,nil
}

func renderLoop(app *appSettings){
	g:= app.g
	g.renderer.Clear()
	wx:=app.init.colorWindowSize

	state := sdl.GetKeyboardState()
	keyTakeColorPressed := state[sdl.SCANCODE_LCTRL]==1||
		state[sdl.SCANCODE_RCTRL]==1

	if !app.hasFocus{
		app.logLabel.Text="focus window"
		app.logLabel.Show()
	}

	app.logLabel.Tick()

	ForEachButton(app.buttons,func(b*Button){b.Draw(g)})
	app.logLabel.Draw(g)

	mx,my,_:=sdl.GetGlobalMouseState()
	bounds := image.Rect(int(mx-app.init.colorWindowSize),
		int(my-app.init.colorWindowSize),int(mx+wx),int(my+wx))
	screenshot.GetDisplayBounds(0)
	img,err := screenshot.CaptureRect(bounds)
	if err != nil {
		log.Println(err)
	}


	src := sdl.Rect{0, 0, 2*wx, 2*wx}
	dst := sdl.Rect{10, 10, 100, 100}

	app.texture.Update(nil,img.Pix,img.Stride)
	g.renderer.Copy(app.texture, &src, &dst)

	if keyTakeColorPressed {
		color := getCenterColor(img)
		drawCrossair(g.renderer,dst)
		app.tColor,_ = colorful.MakeColor(color)
	}

	g.renderer.SetDrawColor(toRGBA(app.tColor))
	colorRect:=sdl.Rect{120,10,100,100}
	g.renderer.FillRect(&colorRect)
	g.renderer.SetDrawColor(0,0,0,0)

	colorRGB := toRGBString(app.tColor)
	colorHSL := toHSLString(app.tColor)
	colorHEX := app.tColor.Hex()

	renderText(g,colorRGB,app.init.textColor,&sdl.Point{290,10})
	renderText(g,colorHEX,app.init.textColor,&sdl.Point{290,40})
	renderText(g,colorHSL,app.init.textColor,&sdl.Point{290,70})

	app.buttons[0].Rect = &sdl.Rect{230,10,45,20}
	app.buttons[1].Rect = &sdl.Rect{230,40,45,20}
	app.buttons[2].Rect = &sdl.Rect{230,70,45,20}

	g.renderer.Present()
}



