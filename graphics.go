package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"image"
	"image/color"
	"time"
)

type Clickable interface{
	GetRect() *sdl.Rect
	GetOnClick() OnClick
}

type OnClick func()

type Button struct{
	Rect *sdl.Rect
	Text string
	OnClick OnClick
	ForeColor sdl.Color
}

func (b *Button) GetRect() *sdl.Rect{
	return b.Rect
}
func (b *Button) GetOnClick() OnClick{
	return b.OnClick
}

func (b *Button) Draw(g *Graphics){
	re,gr,bl,a,_ := g.renderer.GetDrawColor()
	g.renderer.SetDrawColor(240,240,240,0)
	g.renderer.DrawRect(b.Rect)
	renderText(g,b.Text,b.ForeColor,&sdl.Point{b.Rect.X,b.Rect.Y})
	g.renderer.SetDrawColor(re,gr,bl,a)
}


type LogLabel struct{
	Rect *sdl.Rect
	Text string
	VisibleDuration time.Duration
	visibleStartTime time.Time
	visible bool
	ForeColor sdl.Color
	BackColor sdl.Color
}

func (b *LogLabel)Show(){
	b.visible = true
	b.visibleStartTime = time.Now()
}


func (b *LogLabel)Tick () {
	if time.Now().Sub(b.visibleStartTime)>b.VisibleDuration{
		b.visible=false
	}
}

func (b *LogLabel) Draw(g *Graphics){
	if !b.visible{
		return
	}
	re,gr,bl,a,_ := g.renderer.GetDrawColor()
	renderText(g,b.Text,b.ForeColor,&sdl.Point{b.Rect.X,b.Rect.Y})
	g.renderer.SetDrawColor(re,gr,bl,a)
}

func ForEachButton(buttons []*Button, f func(b *Button)){
	for _,x:=range buttons{
		f(x)
	}
}

func CheckClickEvent(t *sdl.MouseButtonEvent,clickable Clickable){
	clickPoint := &sdl.Point{t.X,t.Y}
	if clickPoint.InRect(clickable.GetRect()) && t.State==1{
		clickable.GetOnClick()()
	}
}

func renderText(g *Graphics,text string,color sdl.Color,point *sdl.Point){
	surface,_ := g.font.RenderUTF8Blended(text,color)
	defer surface.Free()
	txtSurface,_ := g.renderer.CreateTextureFromSurface(surface)
	defer txtSurface.Destroy()
	g.renderer.Copy(txtSurface,nil,
		&sdl.Rect{point.X,point.Y,surface.W,surface.H})
}

func getCenterColor(img *image.RGBA)color.Color{
	c:= img.At(img.Rect.Dx()/2,img.Rect.Dy()/2)
	return c
}

