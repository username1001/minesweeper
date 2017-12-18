package minesweeper

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

const (
	SAMPLE_GRID_WIDTH = 10
	SAMPLE_GRID_HEIGHT = 20
)

func newBlankGame() Minesweeper {
	return NewGame()
}

func newSampleGame() Minesweeper {
	return NewGame(Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
}

func TestGridMustNotBeSquaredForTheSakeOfTesting(t *testing.T) {
	assert.True(t, SAMPLE_GRID_WIDTH != SAMPLE_GRID_HEIGHT)
}

func TestBlock_SetBlock(t *testing.T) {
	block := new(Block)

	block.SetBlock(UNKNOWN)
	assert.Equal(t, block.Node, UNKNOWN)
}

func TestGame_SetGrid(t *testing.T) {
	minesweeper := newBlankGame()
	minesweeper.SetGrid(SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT)
	assert.Equal(t, minesweeper.(*game).Board.Grid, &Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
}

func TestGameWithGridArgument(t *testing.T) {
	minesweeper := newSampleGame()
	assert.Equal(t, minesweeper.(*game).Board.Grid, &Grid{SAMPLE_GRID_WIDTH, SAMPLE_GRID_HEIGHT})
}

func TestNewGridWhenStartedGame(t *testing.T) {
	minesweeper := newSampleGame()
	err := minesweeper.SetGrid(10, 20)
	assert.NotNil(t, err, "Must report an error upon setting a new grid from an already started game")
	assert.IsType(t, new(GameAlreadyStarted), err, "The error must be GameAlreadyStarted error type")
}

func TestFlaggedBlock(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.Flag(3, 6)
	assert.Equal(t, minesweeper.(*game).Blocks[3][6].Node, FLAGGED)
}

func TestGame_SetDifficulty(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	assert.Equal(t, minesweeper.(*game).Difficulty, EASY)
}

func TestShiftFromMaxPosition(t *testing.T) {
	grid := Grid{5, 5}
	x, y := shiftPosition(&grid, 4, 4)
	assert.Equal(t, struct{x int; y int}{ 0, 0}, struct{x int; y int}{x, y})
}

func TestBombsInPlace(t *testing.T) {

	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)

	numOfBombs := int(float32(game.width * game.height) * EASY_MULTIPLIER)
	countedBombs := 0
	for _, row := range game.Blocks {
		for _, block := range row {
			if block.Node == BOMB {
				countedBombs ++
			}
		}
	}
	assert.Equal(t, numOfBombs, countedBombs)
}

func TestTalliedBomb(t *testing.T) {
	minesweeper := newSampleGame()
	minesweeper.SetDifficulty(EASY)
	minesweeper.Play()

	game := minesweeper.(*game)
	width := game.width
	height := game.height

	count := func(blocks Blocks, x, y int) (has int) {
		if x >= 0 && y >= 0 &&
			x < width && y < height &&
				blocks[x][y].Node & BOMB == 1 {
					return 1
		}
		return
	}

	hasSurroundingTally := func(blocks Blocks, x, y int) int {
		if x >= 0 && y >= 0 &&
			x < width && y < height {
				switch blocks[x][y].Node {
				case NUMBER:
					return 1
				case BOMB:
					return -1
				default:
					return 0
				}
		}
		return -1
	}
	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == BOMB {
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x-1, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x-1, y))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x-1, y+1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x, y+1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x+1, y-1))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x+1, y))
				assert.NotEqual(t, 0, hasSurroundingTally(game.Blocks, x+1, y+1))
			}
		}
	}

	for x, row := range game.Blocks {
		for y, block := range row {
			if block.Node == NUMBER {
				var counted int
				counted = count(game.Blocks, x - 1, y - 1	) +
				count(game.Blocks, x - 1, y		) +
				count(game.Blocks, x - 1, y + 1	) +
				count(game.Blocks, x	, y - 1	) +
				count(game.Blocks, x	, y + 1	) +
				count(game.Blocks, x + 1, y - 1	) +
				count(game.Blocks, x + 1, y		) +
				count(game.Blocks, x + 1, y + 1	)
				assert.Equal(t, counted, block.value)
			}
		}
	}
}

func print(game *game) {
	for _, row := range game.Blocks {
		fmt.Println()
		for _, block := range row {
			if block.Node == BOMB {
				fmt.Print("* ")
			} else if block.Node == UNKNOWN {
				fmt.Print("  ")
			} else {
				fmt.Printf("%v ", block.value)
			}
		}
	}
}