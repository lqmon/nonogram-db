package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const cellSize = 40
const hintArea = 120

func parseInput() (int, int, []string, []string, string) {

	scanner := bufio.NewScanner(os.Stdin)

	var height, width int
	var rows []string
	var cols []string
	var goal string

	mode := ""

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "height") {
			height, _ = strconv.Atoi(strings.Split(line, " ")[1])
		} else if strings.HasPrefix(line, "width") {
			width, _ = strconv.Atoi(strings.Split(line, " ")[1])
		} else if line == "rows" {
			mode = "rows"
		} else if line == "columns" {
			mode = "cols"
		} else if strings.HasPrefix(line, "goal") {
			goal = strings.Trim(line[5:], "\"")
		} else {

			if mode == "rows" {
				rows = append(rows, line)
			} else if mode == "cols" {
				cols = append(cols, line)
			}
		}
	}

	return height, width, rows, cols, goal
}

func drawRect(img *image.RGBA, x, y, w, h int, c color.Color) {

	r := image.Rect(x, y, x+w, y+h)
	draw.Draw(img, r, &image.Uniform{c}, image.Point{}, draw.Src)
}

func drawText(img *image.RGBA, x, y int, label string) {

	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString(label)
}

func drawHints(img *image.RGBA, rows, cols []string) {

	for i, r := range rows {

		hints := strings.Split(r, ",")

		for j, h := range hints {

			drawText(
				img,
				hintArea-20*(len(hints)-j),
				hintArea+i*cellSize+25,
				h,
			)
		}
	}

	for i, c := range cols {

		hints := strings.Split(c, ",")

		for j, h := range hints {

			drawText(
				img,
				hintArea+i*cellSize+10,
				hintArea-20*(len(hints)-j),
				h,
			)
		}
	}
}

func drawGrid(img *image.RGBA, height, width int) {

	for y := 0; y <= height; y++ {
		for x := 0; x < width*cellSize; x++ {
			img.Set(hintArea+x, hintArea+y*cellSize, color.Black)
		}
	}

	for x := 0; x <= width; x++ {
		for y := 0; y < height*cellSize; y++ {
			img.Set(hintArea+x*cellSize, hintArea+y, color.Black)
		}
	}
}

func drawSolution(img *image.RGBA, height, width int, goal string) {

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			idx := y*width + x

			if goal[idx] == '1' {

				drawRect(
					img,
					hintArea+x*cellSize,
					hintArea+y*cellSize,
					cellSize,
					cellSize,
					color.Black,
				)
			}
		}
	}
}

func createImage(height, width int) *image.RGBA {

	imgWidth := width*cellSize + hintArea
	imgHeight := height*cellSize + hintArea

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	drawRect(img, 0, 0, imgWidth, imgHeight, color.White)

	return img
}

func saveImage(img *image.RGBA, filename string) {

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}

func main() {

	height, width, rows, cols, goal := parseInput()

	// Problem image (no solution cells)
	problem := createImage(height, width)
	drawHints(problem, rows, cols)
	drawGrid(problem, height, width)

	saveImage(problem, "nonogram_problem.png")

	// Solution image
	solution := createImage(height, width)
	drawHints(solution, rows, cols)
	drawSolution(solution, height, width, goal)
	drawGrid(solution, height, width)

	saveImage(solution, "nonogram_solution.png")

	fmt.Println("Generated:")
	fmt.Println("nonogram_problem.png")
	fmt.Println("nonogram_solution.png")
}