package main

import (
	"errors"
	"sort"
)

type AnnouncementBoard interface {
	getAnnouncementIDsAt(int, int) []int
}

type AnnouncementPaper interface {
	addTo(AnnouncementBoard, int, int) error
	removeAndGetIDsOnTop() []int
}

type Board struct {
	rows         int
	cols         int
	currentTime  int
	installation []*installation
}

type Paper struct {
	id     int
	width  int
	height int
	board  *Board
}

type installation struct {
	paper *Paper
	row   int
	col   int
	time  int
}

func NewBoard(row int, col int) AnnouncementBoard {
	return &Board{rows: row, cols: col}
}

func (b *Board) getAnnouncementIDsAt(row, col int) []int {
	var ids []int
	for _, installation := range b.installation {
		if installation.contains(row, col) {
			ids = append(ids, installation.paper.id)
		}
	}
	if ids == nil {
		return []int{}
	}
	return ids
}

func (p *Paper) addTo(ab AnnouncementBoard, row, col int) error {
	b, ok := ab.(*Board)
	if !ok {
		return errors.New("invalid announcement board")
	}
	startCol := col - (p.width / 2)
	if row < 0 || startCol < 0 || row+p.height > b.rows || startCol+p.width > b.cols {
		if row < 0 || row >= b.rows || col < 0 || col >= b.cols {
			return errors.New("punez out of bound")
		}
		return errors.New("paper out of bound")
	}
	inst := &installation{paper: p, row: row, col: startCol, time: b.currentTime}
	b.currentTime++
	b.installation = append(b.installation, inst)
	p.board = b
	return nil
}

func (p *Paper) removeAndGetIDsOnTop() []int {
	if p.board == nil {
		return []int{}
	}
	var removed *installation
	b := p.board
	for i, inst := range b.installation {
		if inst.paper.id == p.id {
			removed = inst
			b.installation = append(b.installation[:i], b.installation[i+1:]...)
			break
		}
	}
	if removed == nil {
		return []int{}
	}
	overlapping := getInstallationOnTop(removed, b)
	sort.SliceStable(overlapping, func(i, j int) bool {
		return overlapping[i].time > overlapping[j].time
	})
	seen := make(map[int]bool)
	var ids []int
	for _, inst := range overlapping {
		if !seen[inst.paper.id] {
			ids = append(ids, inst.paper.id)
			seen[inst.paper.id] = true
		}
	}
	if ids == nil {
		return []int{}
	}
	p.board = nil
	return ids
}

func getInstallationOnTop(base *installation, b *Board) []*installation {
	var result []*installation
	for _, inst := range b.installation {
		if inst.time <= base.time {
			continue
		}
		if isOverLapping(base, inst) {
			continue
		}
		result = append(result, inst)
		result = append(result, getInstallationOnTop(inst, b)...)
	}
	return result

}

func isOverLapping(a, b *installation) bool {
	return a.col+a.paper.width <= b.col ||
		b.col+b.paper.width < a.col ||
		a.row+a.paper.height <= b.row ||
		b.row+b.paper.height < a.row
}

func NewPaper(width int, height int, ID int) AnnouncementPaper {
	return &Paper{id: ID, width: width, height: height}
}

func (inst *installation) contains(row, col int) bool {
	return (row >= inst.row) && (row < inst.row+inst.paper.height) && (col >= inst.col) && (col < inst.col+inst.paper.width)
}
