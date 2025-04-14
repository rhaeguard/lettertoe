package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type JsonPayload struct {
	Three []string `json:"3"`
	Four  []string `json:"4"`
	Five  []string `json:"5"`
}

const GAME_BOARD_SIZE = 3
const LETTER_BOARD_HEIGHT = 2
const LETTER_BOARD_WIDTH = 5

var GAME_BOARD = [GAME_BOARD_SIZE][GAME_BOARD_SIZE]uint8{
	{' ', ' ', ' '},
	{' ', ' ', ' '},
	{' ', ' ', ' '},
}

var LETTER_BOARD = [LETTER_BOARD_HEIGHT][LETTER_BOARD_WIDTH]uint8{
	{' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', ' '},
}

var WordsData JsonPayload

func determineLetters() {
	words := WordsData.Three

	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })

	seenLetters := []rune{
		' ',
	}

	index := 0

	for index < 3 {
		word := words[index]

		for _, ch := range strings.ToUpper(word) {
			seenLetters = append(seenLetters, ch)
		}

		index += 1
	}

	rand.Shuffle(len(seenLetters), func(i, j int) { seenLetters[i], seenLetters[j] = seenLetters[j], seenLetters[i] })

	index = 0

	for i := range LETTER_BOARD_HEIGHT {
		for j := range LETTER_BOARD_WIDTH {
			LETTER_BOARD[i][j] = uint8(seenLetters[index])
			index += 1
		}
	}
}

func main() {

	jsonFile, err := os.Open("words.json")
	if err != nil {
		panic("could not open the words.json file")
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		panic("could not read the words.json file bytes")
	}

	if err := json.Unmarshal(jsonBytes, &WordsData); err != nil {
		panic("could not unmarshall the words.json file")
	}

	determineLetters()

	ScreenWidth, ScreenHeight := int32(700), int32(1000)
	FontSize := int32(60)

	BoardSize := float32(ScreenWidth) * 0.8
	CellSize := BoardSize / 3
	BoardPosition := rl.NewVector2(
		(float32(ScreenWidth)-BoardSize)/2,
		(float32(ScreenHeight)-BoardSize)/3,
	)
	Board := rl.NewRectangle(
		BoardPosition.X, BoardPosition.Y, BoardSize, BoardSize,
	)

	BoardPositonEnd := BoardPosition.Y + BoardSize

	InputAreaTotal := rl.NewRectangle(
		0, BoardPositonEnd, float32(ScreenWidth),
		(float32(ScreenHeight)-Board.Height)-Board.Y,
	)

	LetterAreaHeight := InputAreaTotal.Height * 0.8

	LettersArea := rl.NewRectangle(
		(float32(ScreenWidth)-LetterAreaHeight*2.5)/2,
		InputAreaTotal.Y+(float32(InputAreaTotal.Height)-LetterAreaHeight)/2,
		LetterAreaHeight*2.5,
		LetterAreaHeight,
	)

	LetterCellSize := LetterAreaHeight / 2

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(ScreenWidth, ScreenHeight, "lettertoe")
	defer rl.CloseWindow()

	font := rl.GetFontDefault()

	var userSelectedBoardCell rl.Vector2 = rl.NewVector2(-1, -1)
	var userSelectedLetterCell rl.Vector2 = rl.NewVector2(-1, -1)

	var draggingLetterPos rl.Vector2 = rl.NewVector2(-1, -1)

	for !rl.WindowShouldClose() {

		{ // determine what cell user is hovering over
			mp := rl.GetMousePosition()

			if rl.CheckCollisionPointRec(mp, Board) {
				deltaX := mp.X - Board.X
				deltaY := mp.Y - Board.Y

				cellX := float64(deltaX / CellSize)
				cellY := float64(deltaY / CellSize)

				if cellX >= 0 && cellY >= 0 {
					userSelectedBoardCell.X = float32(math.Floor(cellX))
					userSelectedBoardCell.Y = float32(math.Floor(cellY))
				} else {
					userSelectedBoardCell.X = -1
					userSelectedBoardCell.Y = -1
				}
			} else {
				userSelectedBoardCell.X = -1
				userSelectedBoardCell.Y = -1
			}

			if rl.CheckCollisionPointRec(mp, LettersArea) {
				deltaX := mp.X - LettersArea.X
				deltaY := mp.Y - LettersArea.Y

				cellX := float64(deltaX / LetterCellSize)
				cellY := float64(deltaY / LetterCellSize)

				if cellX >= 0 && cellY >= 0 {
					xx := float32(math.Floor(cellX))
					yy := float32(math.Floor(cellY))

					if LETTER_BOARD[int(yy)][int(xx)] != ' ' {
						userSelectedLetterCell.X = xx
						userSelectedLetterCell.Y = yy
					} else {
						userSelectedLetterCell.X = -1
						userSelectedLetterCell.Y = -1
					}

				} else {
					userSelectedLetterCell.X = -1
					userSelectedLetterCell.Y = -1
				}
			} else {
				userSelectedLetterCell.X = -1
				userSelectedLetterCell.Y = -1
			}
		}

		{ // drag mechanic

			if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				if draggingLetterPos.X >= 0 {

				}

				sx, sy := int(userSelectedLetterCell.X), int(userSelectedLetterCell.Y)

				if userSelectedLetterCell.X >= 0 && LETTER_BOARD[sy][sx] != ' ' {
					draggingLetterPos = userSelectedLetterCell
				}
			}

			if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
				if draggingLetterPos.X >= 0 {

					if userSelectedBoardCell.X >= 0 {
						sx, sy := int(draggingLetterPos.X), int(draggingLetterPos.Y)
						dx, dy := int(userSelectedBoardCell.X), int(userSelectedBoardCell.Y)

						GAME_BOARD[dy][dx] = LETTER_BOARD[sy][sx]
						LETTER_BOARD[sy][sx] = ' '
					}

					draggingLetterPos.X = -1
					draggingLetterPos.Y = -1
				}
			}

		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.SkyBlue)

		bannerSize := rl.MeasureTextEx(font, "lettertoe", float32(FontSize), float32(FontSize/5))
		rl.DrawTextEx(font, "lettertoe", rl.NewVector2(float32(ScreenWidth-int32(bannerSize.X))/2, 0), float32(FontSize), float32(FontSize/5), rl.DarkBlue)

		rl.DrawRectangleV(
			BoardPosition, rl.NewVector2(BoardSize, BoardSize), rl.RayWhite,
		)

		if userSelectedBoardCell.X >= 0 {
			rl.DrawRectangleRec(
				rl.NewRectangle(
					BoardPosition.X+userSelectedBoardCell.X*CellSize,
					BoardPosition.Y+userSelectedBoardCell.Y*CellSize,
					CellSize,
					CellSize,
				),
				rl.Yellow,
			)
		}

		for y := range 3 {
			for x := range 3 {
				rl.DrawRectangleLinesEx(
					rl.NewRectangle(
						BoardPosition.X+float32(x)*CellSize,
						BoardPosition.Y+float32(y)*CellSize,
						CellSize,
						CellSize,
					),
					1,
					rl.DarkBlue,
				)

				letter := fmt.Sprintf("%c", GAME_BOARD[y][x])

				singleLetterSize := rl.MeasureTextEx(
					font,
					letter,
					float32(FontSize*2),
					0,
				)

				rl.DrawTextEx(
					font,
					letter,
					rl.NewVector2(
						BoardPosition.X+float32(x)*CellSize+(CellSize-singleLetterSize.X)/2,
						BoardPosition.Y+float32(y)*CellSize+(CellSize-singleLetterSize.Y)/2,
					),
					float32(FontSize*2),
					0,
					rl.DarkBlue,
				)
			}
		}

		// letters area

		for y := range 2 {
			for x := range 5 {
				if LETTER_BOARD[y][x] == ' ' {
					continue
				}

				cellColor := rl.RayWhite

				if userSelectedLetterCell.X >= 0 && (int(userSelectedLetterCell.X) == x && int(userSelectedLetterCell.Y) == y) {
					cellColor = rl.Red
				}

				rl.DrawRectangleRec(
					rl.NewRectangle(
						LettersArea.X+float32(x)*LetterCellSize,
						LettersArea.Y+float32(y)*LetterCellSize,
						LetterCellSize,
						LetterCellSize,
					),
					cellColor,
				)

				rl.DrawRectangleLinesEx(
					rl.NewRectangle(
						LettersArea.X+float32(x)*LetterCellSize,
						LettersArea.Y+float32(y)*LetterCellSize,
						LetterCellSize,
						LetterCellSize,
					),
					1,
					rl.DarkBlue,
				)

				letter := fmt.Sprintf("%c", LETTER_BOARD[y][x])

				singleLetterSize := rl.MeasureTextEx(
					font,
					letter,
					float32(FontSize),
					0,
				)

				rl.DrawTextEx(
					font,
					letter,
					rl.NewVector2(
						LettersArea.X+float32(x)*LetterCellSize+(LetterCellSize-singleLetterSize.X)/2,
						LettersArea.Y+float32(y)*LetterCellSize+(LetterCellSize-singleLetterSize.Y)/2,
					),
					float32(FontSize),
					0,
					rl.DarkBlue,
				)
			}
		}

		{ // draw the dragged letter
			if draggingLetterPos.X >= 0 {
				pos := rl.GetMousePosition()

				rl.DrawRectangleV(
					pos,
					rl.NewVector2(LetterCellSize, LetterCellSize),
					rl.LightGray,
				)

				x, y := int(draggingLetterPos.X), int(draggingLetterPos.Y)

				letter := fmt.Sprintf("%c", LETTER_BOARD[y][x])

				singleLetterSize := rl.MeasureTextEx(
					font,
					letter,
					float32(FontSize),
					0,
				)

				rl.DrawTextEx(
					font,
					letter,
					rl.NewVector2(
						pos.X+(LetterCellSize-singleLetterSize.X)/2,
						pos.Y+(LetterCellSize-singleLetterSize.Y)/2,
					),
					float32(FontSize),
					0,
					rl.DarkBlue,
				)
			}
		}

		rl.EndDrawing()
	}

}
