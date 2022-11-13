package main

import (
	eb "github.com/IgneousRed/EBitEn"
	m "github.com/IgneousRed/gomisc"
	ebt "github.com/hajimehoshi/ebiten/v2"
)

var windowSize = m.Vec2I(600, 600)
var font eb.Font

type tileState int

const (
	tSEmpty tileState = iota
	tSX
	tSO
)

type finish int

const (
	fNil finish = iota
	fDraw
	fHumanWin
	fBotWin
)

type game struct {
	rng    m.PCG32
	finish finish
	tiles  [3][3]tileState
	filled int
	cursor m.Vec[int]
}

func gameInit() game {
	g := game{}
	g.rng = m.PCG32Init()
	g.cursor = m.Vec2I(1, 1)
	return g
}
func (g *game) cursorUpdate(b bool, vec m.Vec[int]) {
	if b {
		g.cursor = g.cursor.Add(vec).Wrap1(3)
	}
}
func (g *game) tile(pos m.Vec[int]) *tileState {
	return &g.tiles[pos[0]][pos[1]]
}
func tStoF(tile tileState) finish {
	if tile == tSX {
		return fHumanWin
	} else if tile == tSO {
		return fBotWin
	}
	return fNil
}
func (g *game) checkTiles(x0, y0, x1, y1, x2, y2 int) bool {
	if g.tiles[x0][y0] == g.tiles[x1][y1] &&
		g.tiles[x1][y1] == g.tiles[x2][y2] &&
		g.tiles[x0][y0] != tSEmpty {
		g.finish = tStoF(g.tiles[x0][y0])
		return true
	}
	return false
}
func (g *game) updateVictory() bool {
	for x := range g.tiles {
		if g.checkTiles(x, 0, x, 1, x, 2) {
			return true
		}
	}
	for y := range g.tiles {
		if g.checkTiles(0, y, 1, y, 2, y) {
			return true
		}
	}
	for i := 0; i < 2; i++ {
		if g.checkTiles(0, i*2, 1, 1, 2, 2-i*2) {
			return true
		}
	}
	if g.filled == 9 {
		g.finish = fDraw
		return true
	}
	return false
}
func (g *game) Update() {
	if g.finish != fNil {
		if eb.KeysDown(ebt.KeySpace) {
			*g = gameInit()
		}
		return
	}
	g.cursorUpdate(eb.KeysDown(ebt.KeyArrowRight, ebt.KeyD), m.Vec2I(1, 0))
	g.cursorUpdate(eb.KeysDown(ebt.KeyArrowUp, ebt.KeyW), m.Vec2I(0, 1))
	g.cursorUpdate(eb.KeysDown(ebt.KeyArrowLeft, ebt.KeyA), m.Vec2I(-1, 0))
	g.cursorUpdate(eb.KeysDown(ebt.KeyArrowDown, ebt.KeyS), m.Vec2I(0, -1))
	if eb.KeysDown(ebt.KeySpace) && *g.tile(g.cursor) == tSEmpty {
		*g.tile(g.cursor) = tSX
		g.filled++
		if g.updateVictory() {
			return
		}
		for {
			pos := m.Vec2I(g.rng.Range(3), g.rng.Range(3))
			if *g.tile(pos) == tSEmpty {
				*g.tile(pos) = tSO
				g.filled++
				break
			}
		}
		g.updateVictory()
	}
}
func tilePos(pos m.Vec[int]) m.Vec[int] {
	return m.Vec2I(200, 200).Mul(pos)
}
func (g *game) Draw() {
	eb.DrawLineI(m.Vec2I(200, 0), m.Vec2I(200, 600), 1, eb.White)
	eb.DrawLineI(m.Vec2I(400, 0), m.Vec2I(400, 600), 1, eb.White)
	eb.DrawLineI(m.Vec2I(0, 200), m.Vec2I(600, 200), 1, eb.White)
	eb.DrawLineI(m.Vec2I(0, 400), m.Vec2I(600, 400), 1, eb.White)
	cursorStart := tilePos(g.cursor).Add(m.Vec2I(10, 10))
	cursorEnd := cursorStart.Add(m.Vec2I(180, 180))
	eb.DrawLineI(cursorStart, m.Vec2I(cursorStart[0], cursorEnd[1]), 1, eb.Red)
	eb.DrawLineI(cursorStart, m.Vec2I(cursorEnd[0], cursorStart[1]), 1, eb.Red)
	eb.DrawLineI(cursorEnd, m.Vec2I(cursorEnd[0], cursorStart[1]), 1, eb.Red)
	eb.DrawLineI(cursorEnd, m.Vec2I(cursorStart[0], cursorEnd[1]), 1, eb.Red)
	for x := range g.tiles {
		for y, tile := range g.tiles[x] {
			pos := tilePos(m.Vec2I(x, y)).Add(m.Vec2I(30, 20))
			if tile == tSX {
				eb.DrawTextI(font, 200, pos, "X", eb.Green)
			} else if tile == tSO {
				eb.DrawTextI(font, 200, pos, "O", eb.Blue)
			}
		}
	}
	if g.finish != fNil {
		if g.finish == fHumanWin {
			eb.DrawTextI(font, 20, m.Vec2I(250, 310), "Human Win", eb.Magenta)
		} else if g.finish == fBotWin {
			eb.DrawTextI(font, 20, m.Vec2I(250, 310), "Bot Win", eb.Magenta)
		} else if g.finish == fDraw {
			eb.DrawTextI(font, 20, m.Vec2I(250, 310), "DRAW", eb.Magenta)
		}
		eb.DrawTextI(font, 20, m.Vec2I(250, 290), "Press Space", eb.Magenta)
	}
}
func main() {
	f, err := eb.FontNew("FiraCode-Medium.ttf")
	m.FatalErr("", err)
	font = f
	g := gameInit()
	eb.InitGame("ttt", windowSize, &g)
}
