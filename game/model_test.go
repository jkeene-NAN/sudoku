package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func Test_creatAllCandidateList(t *testing.T) {
	var allCandidateList candidateList = createAllCandidatesList()
	assert.Equal(t, (9 * 9 * 9), len(allCandidateList))
}

func Test_shuffleCandidates(t *testing.T) {
	var allCandidateList candidateList = createAllCandidatesList()
	shuffleCandidates(allCandidateList)
	assert.Equal(t, (9 * 9 * 9), len(allCandidateList))

	for _, c := range allCandidateList {
		assert.NotNil(t, c)
	}
}

func Test_createValidCandidateList(t *testing.T) {
	var c *candidate = &candidate{
		row: 0,
		column: 0,
		value: 4,
	}

	var gs *gameState
	var err error
	gs, err = createGame(NewGame())
	assert.Nil(t, err)

	var allCandidates candidateList = createAllCandidatesList()
	allCandidates = candidateList{&candidate{}, &candidate{}, &candidate{}}
	allCandidates[0].value = 0
	allCandidates[1].value = 1


	allCandidates[2].value = 0
	allCandidates[2].row = 0
	allCandidates[2].column = 4

	var potentialCandidates = createValidCandidateList(gs, allCandidates)

	assert.Equal(t, len(allCandidates), len(potentialCandidates))
	c = allCandidates[0]
	gs.addCandidate(c)
	assert.True(t, gs.isSet(c.row, c.column))
	potentialCandidates = createValidCandidateList(gs, allCandidates)
	assert.NotEqual(t, len(allCandidates), len(potentialCandidates))
	//assert.Equal(t, (len(allCandidates) - (9 + 8 + 6 + 6)), len(potentialCandidates))
	assert.Equal(t, 0, len(potentialCandidates))
}

func Test_gameState_setCandidate(t *testing.T) {
	var c *candidate = &candidate{
		row: 0,
		column: 0,
		value: 4,
	}

	var gs *gameState
	var err error
	gs, err = createGame(NewGame())
	assert.Nil(t, err)

	gs.addCandidate(c)
	assert.Equal(t, 4, gs.Grid[c.row][c.column])
	assert.Equal(t, 1, gs.setCount())
	assert.True(t, gs.isSet(c.row, c.column))

	c = &candidate{
		row: 1,
		column: 1,
		value: 5,
	}
	gs.addCandidate(c)
	assert.Equal(t, 5, gs.Grid[c.row][c.column])
	assert.Equal(t, 2, gs.setCount())
	assert.True(t, gs.isSet(c.row, c.column))


	gs.removeCandidate(c)
	assert.Equal(t, NotSet, gs.Grid[c.row][c.column])
	assert.Equal(t, 1, gs.setCount())
	assert.False(t, gs.isSet(c.row, c.column))

}

func Test_computeBaseRow(t *testing.T) {
	var baseRow int
	var err error

	baseRow, err = computeBaseRow(-1)
	assert.NotNil(t, err)

	baseRow, err = computeBaseRow(10)
	assert.NotNil(t, err)

	var expect int = 0

	baseRow, err = computeBaseRow(0)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseRow)

	baseRow, err = computeBaseRow(1)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseRow)

	baseRow, err = computeBaseRow(2)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseRow)

	expect = 3
	baseRow, err = computeBaseRow(3)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseRow)

	baseRow, err = computeBaseRow(4)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseRow)

	baseRow, err = computeBaseRow(5)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseRow)

	expect = 6
	baseRow, err = computeBaseRow(6)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseRow)

	baseRow, err = computeBaseRow(7)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseRow)

	baseRow, err = computeBaseRow(8)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseRow)
}

func Test_computeBaseColumn(t *testing.T) {
	var baseColumn, expect int
	var err error

	baseColumn, err = computeBaseColumn(-1)
	assert.NotNil(t, err)

	baseColumn, err = computeBaseColumn(10)
	assert.NotNil(t, err)

	expect = 0
	baseColumn, err = computeBaseColumn(0)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseColumn)

	baseColumn, err = computeBaseColumn(3)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseColumn)

	baseColumn, err = computeBaseColumn(6)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseColumn)


	expect = 3
	baseColumn, err = computeBaseColumn(1)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseColumn)

	baseColumn, err = computeBaseColumn(4)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseColumn)

	baseColumn, err = computeBaseColumn(7)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseColumn)

	expect = 6
	baseColumn, err = computeBaseColumn(2)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseColumn)

	baseColumn, err = computeBaseColumn(5)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseColumn)

	baseColumn, err = computeBaseColumn(8)
	assert.Nil(t, err)
	assert.Equal(t, expect, baseColumn)
}




func Test_validateSubGrid(t *testing.T) {
	var gameState *gameState
	var err error
	gameState, err = createGame(NewGame())

	/*Valid Test Cases */
	err = validateSubGrid(gameState, 0)
	assert.Nil(t, err)

	err = validateSubGrid(gameState, 1)
	assert.Nil(t, err)

	err = validateSubGrid(gameState, 2)
	assert.Nil(t, err)

	err = validateSubGrid(gameState, 3)
	assert.Nil(t, err)

	err = validateSubGrid(gameState, 4)
	assert.Nil(t, err)

	err = validateSubGrid(gameState, 5)
	assert.Nil(t, err)

	err = validateSubGrid(gameState, 6)
	assert.Nil(t, err)

	err = validateSubGrid(gameState, 7)
	assert.Nil(t, err)

	err = validateSubGrid(gameState, 8)
	assert.Nil(t, err)

	/* should err on sub grid */
	err = validateSubGrid(gameState, -29)
	assert.NotNil(t, err)

	err = validateSubGrid(gameState, 345)
	assert.NotNil(t, err)


	gameState.Grid[0][0] = 8
	err = validateSubGrid(gameState, 0)
	assert.Nil(t, err)

	gameState.Grid[0][1] = 8
	err = validateSubGrid(gameState, 0)
	assert.NotNil(t, err)

	gameState.Grid[1][0] = 8
	err = validateSubGrid(gameState, 0)
	assert.NotNil(t, err)


	err = validateSubGrid(gameState, 3)
	assert.Nil(t, err)


	resetGameState(gameState)
	err = validateSubGrid(gameState, 0)
	assert.Nil(t, err)


	gameState.Grid[3][3] = 2
	gameState.Grid[4][4] = 2
	err = validateSubGrid(gameState, 4)
	assert.NotNil(t, err)

	gameState.Grid[4][4] = 0
	gameState.Grid[0][0] = 2
	err = validateSubGrid(gameState, 4)
	assert.Nil(t, err)

	resetGameState(gameState)
	err = validateSubGrid(gameState, 0)
	assert.Nil(t, err)

	gameState.Grid[0][0] = 0
	gameState.Grid[0][1] = 1
	gameState.Grid[0][2] = 2
	gameState.Grid[1][0] = 3
	gameState.Grid[1][1] = 4
	gameState.Grid[1][2] = 5
	gameState.Grid[2][0] = 6
	gameState.Grid[2][1] = 7
	gameState.Grid[2][2] = 8
	err = validateSubGrid(gameState, 0)
	assert.Nil(t, err)
}

func Test_testValidateGameState(t *testing.T) {
	var gameState *gameState
	var err error

	gameState, err = createGame(NewGame())
	assert.Nil(t, err)

	err = validateGameState(gameState)
	assert.Nil(t, err)

	/*
	Should fail on Sub Grid Validation
	 */
	gameState.Grid[0][0] = 3
	gameState.Grid[1][0] = 3
	err = validateGameState(gameState)
	assert.NotNil(t, err)

	resetGameState(gameState)

	row := 3
	gameState.Grid[row][0] = 3
	gameState.Grid[row][1] = 3
	err = validateGameState(gameState)
	assert.NotNil(t, err)

	resetGameState(gameState)
	column := 3
	gameState.Grid[0][column] = 3
	gameState.Grid[1][column] = 3
	err = validateGameState(gameState)
	assert.NotNil(t, err)
}

func Test_validateRow(t *testing.T) {


	var gameState *gameState
	var err error

	gameState, err = createGame(NewGame())
	assert.Nil(t, err)

	err = validateGameState(gameState)
	assert.Nil(t, err)


	err = validateRow(gameState, -234)
	assert.NotNil(t, err)

	err = validateRow(gameState, numRows)
	assert.NotNil(t, err)

	err = validateRow(gameState, numRows - 1)
	assert.Nil(t, err)

	row := 0

	gameState.Grid[row][0] = 3
	gameState.Grid[row][1] = 3
	err = validateRow(gameState, row)
	assert.NotNil(t, err)

	err = validateRow(gameState, row + 1)
	assert.Nil(t, err)

	resetGameState(gameState)
	row = 4
	for i := 0; i < numColumns; i++ {
		gameState.Grid[row][i] = i
	}

	err = validateRow(gameState, row)
	assert.Nil(t, err)
}

func Test_validateColumn(t *testing.T) {
	var gs *gameState
	var err error

	gs, err = createGame(NewGame())
	assert.Nil(t, err)


	err = validateColumn(gs, -234)
	assert.NotNil(t, err)

	err = validateColumn(gs, numColumns)
	assert.NotNil(t, err)

	err = validateColumn(gs, numColumns - 1)
	assert.Nil(t, err)

	resetGameState(gs)

	column := 4

	gs.Grid[0][column] = 3
	gs.Grid[1][column] = 3
	gs.Grid[3][column] = 3
	err = validateColumn(gs, column)
	assert.NotNil(t, err)

	err = validateColumn(gs, column + 1)
	assert.Nil(t, err)

	resetGameState(gs)

	for n := 0; n < numRows; n++ {
		gs.Grid[n][column] = n
	}

	err = validateColumn(gs, column)
	assert.Nil(t, err)
}

func Test_countSelected(t *testing.T) {
	var gs *gameState
	var err error

	gs, err = createGame(NewGame())
	assert.Nil(t, err)

	var count, expect int

	count = 0
	expect = 0

	count = countSelected(gs)
	assert.Equal(t, expect, count)


	gs.Grid[0][0] = 4
	expect += 1
	count = countSelected(gs)
	assert.Equal(t, expect, count)
}

func Test_isFinished(t *testing.T) {
	var expect, actual bool
	var gs *gameState
	var err error

	gs, err = createGame(NewGame())
	assert.Nil(t, err)

	expect = false
	actual = isFinished(gs)

	assert.Equal(t, expect, actual)

	gs.Grid[4][3] = 2
	actual = isFinished(gs)
	assert.Equal(t, expect, actual)

	expect = true
	for a := 0; a < numRows; a++ {
		for b := 0; b < numColumns; b++ {
			gs.Grid[a][b] = a
		}
	}

	/*
	Semantically this is not a valid terminal state for the game Sudoku, but the
	isFinished function does not evaluate constraints.  That will be handled by the
	function validategs
	 */
	actual = isFinished(gs)
	assert.Equal(t, expect, actual)

}

func Test_minMax(t *testing.T) {
	var minExpect int
	var maxExpect int
	var minActual int
	var maxActual int

	minExpect = 0
	maxExpect = 3
	minActual, maxActual = snapToGrid(0)
	assert.Equal(t, minExpect, minActual)
	assert.Equal(t, maxExpect, maxActual)

	minActual, maxActual = snapToGrid(1)
	assert.Equal(t, minExpect, minActual)
	assert.Equal(t, maxExpect, maxActual)

	minActual, maxActual = snapToGrid(2)
	assert.Equal(t, minExpect, minActual)
	assert.Equal(t, maxExpect, maxActual)

	minExpect = 3
	maxExpect = 6
	minActual, maxActual = snapToGrid(3)
	assert.Equal(t, minExpect, minActual)
	assert.Equal(t, maxExpect, maxActual)

	minActual, maxActual = snapToGrid(4)
	assert.Equal(t, minExpect, minActual)
	assert.Equal(t, maxExpect, maxActual)

	minActual, maxActual = snapToGrid(5)
	assert.Equal(t, minExpect, minActual)
	assert.Equal(t, maxExpect, maxActual)

	minExpect = 6
	maxExpect = 9
	minActual, maxActual = snapToGrid(6)
	assert.Equal(t, minExpect, minActual)
	assert.Equal(t, maxExpect, maxActual)

	minActual, maxActual = snapToGrid(7)
	assert.Equal(t, minExpect, minActual)
	assert.Equal(t, maxExpect, maxActual)

	minActual, maxActual = snapToGrid(8)
	assert.Equal(t, minExpect, minActual)
	assert.Equal(t, maxExpect, maxActual)

}



