package main

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/veandco/go-sdl2/sdl"
)

func toRGBString(c colorful.Color)string{
	r,g,b:=c.RGB255()
	return fmt.Sprintf("rgb(%v,%v,%v)",r,g,b)
}

func toHSLString(c colorful.Color)string{
	h,s,l:=c.Hsl()
	return fmt.Sprintf("hsl(%.0f,%.0f,%.0f)",h,s*100,l*100)
}

func toRGBA(c colorful.Color)(r,g,b,a uint8){
	r,g,b = c.RGB255()
	return r,g,b,a
}

func drawCrossair(r *sdl.Renderer,rect sdl.Rect){
	middleOffset := int32(5)

	// top-middle
	r.DrawLine(rect.X+rect.W/2,rect.Y,
		rect.X+rect.W/2,rect.Y+rect.H/2-middleOffset)
	// middle-bottom
	r.DrawLine(rect.X+rect.W/2,rect.Y+rect.H/2+middleOffset,
		rect.X+rect.W/2,rect.Y+rect.H)
	//left-middle
	r.DrawLine(rect.X,rect.Y+rect.H/2,
		rect.X+rect.W/2-middleOffset,rect.Y+rect.H/2)
	//middle-right
	r.DrawLine(rect.X+rect.W/2+middleOffset,rect.Y+rect.H/2,
		rect.X+rect.W,rect.Y+rect.H/2)
}

func copyValue(app *appSettings, what string){
	app.logLabel.Text = "copied to clipboard"
	switch what{
	case "hex":
		sdl.SetClipboardText(app.tColor.Hex())
		app.logLabel.Show()
		break;
	case "rgb":
		sdl.SetClipboardText(toRGBString(app.tColor))
		app.logLabel.Show()
		break;
	case "hsl":
		sdl.SetClipboardText(toHSLString(app.tColor))
		app.logLabel.Show()
		break;
	}
}