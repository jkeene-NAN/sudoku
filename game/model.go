package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const numCandidates int = 9
const numRows int = 9
const numColumns int = 9
const numSubSquares int = 9
const NotSet = -1
const DefaultMaxIterations = 1000000000

func init() {
	rand.Seed(time.Now().Unix())
}

type candidate struct {
	value int
	row int
	column int
}

type candidateList []*candidate

func (c candidate) clone() (*candidate) {
	var ret *candidate = &candidate{
		value: c.value,
		row: c.row,
		column: c.column,
	}

	return ret
}

/*
Sub Grids
0|1|2
3|4|5
6|7|8
 */

type Game struct {
	Grid [numRows][numColumns] int
}

func NewGame() (*Game) {
	var game *Game = &Game{

	}

	for row := 0; row < numRows; row++{
		for column := 0; column < numColumns; column++ {
			game.Grid[row][column] = NotSet
		}
	}

	return game
}

type GamePlayStatistics struct {
	BackTracks int
	Iterations int
}

type gameState struct {
	Grid [numRows][numColumns]int
	moves candidateList
	GamePlayStatistics *GamePlayStatistics
}

func resetGameState(gs *gameState) {
	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			gs.Grid[row][column] = NotSet
		}
	}
}

func (gs *gameState) addCandidate(candidate *candidate) {
	gs.Grid[candidate.row][candidate.column] = candidate.value
}

func (gs *gameState) removeCandidate(candidate *candidate) {
	gs.Grid[candidate.row][candidate.column] = NotSet
}

func (gs *gameState) setCount() int {
	var count int = 0
	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			if gs.Grid[row][column] != NotSet {
				count++
			}
		}
	}

	return count
}

func (gs *gameState) isSet(row, column int) bool {
	return gs.Grid[row][column] != NotSet
}

type searchTree [](candidateList)

func (tree searchTree) back() candidateList {
	var length int = len(tree)
	if length == 0 {
		return nil
	} else {
		return tree[length - 1]
	}
}

func (tree searchTree) popBack() searchTree {
	var length int = len(tree)
	if length == 0 {
		return nil
	} else {
		return tree[0:length - 1]
	}
}


const subGridRows int = 3
const subGridColumns int = 3
/*
Grid 0
(0, 0), (0, 1), (0, 2)
(1, 0), (1, 1), (1, 2)
(2, 0), (2, 1), (2, 2)
 */
/*
Grid 1
(0, 3), (0, 4), (0, 5)
(1, 3), (1, 4), (1, 5)
(2, 3), (2, 4), (2, 5)
 */
/*
Grid 2
(0, 6), (0, 7), (0, 8)
(1, 6), (1, 7), (1, 8)
(2, 6), (2, 7), (2, 8)
*/
/*
Grid 3
(3, 0), (3, 1), (3, 2)
(4, 0), (4, 1), (4, 2)
(5, 0), (5, 1), (5, 2)
 */
/*
Grid 4
(3, 3), (3, 4), (3, 5)
(4, 3), (4, 4), (4, 5)
(5, 3), (5, 4), (5, 5)
 */
/*
Grid 5
(3, 6), (3, 7), (3, 8)
(4, 6), (4, 7), (4, 8)
(5, 6), (5, 7), (5, 8)
 */

/*
Grid 6
(6, 0), (6, 1), (6, 2)
(7, 0), (7, 1), (7, 2)
(8, 0), (8, 1), (8, 2)
 */
/*
Grid 7
(6, 3), (6, 4), (6, 5)
(7, 3), (7, 4), (7, 5)
(8, 3), (8, 4), (8, 5)
 */
/*
Grid 8
(6, 6), (6, 7), (6, 8)
(7, 6), (7, 7), (7, 8)
(8, 6), (8, 7), (8, 8)
 */


/*
Grid N
(
 */

func computeBaseRow(subSquare int) (int, error) {
	if subSquare < 0 || subSquare > numSubSquares {
		return -1, errors.New(fmt.Sprintf("INVALID SUB SQUARE: %d", subSquare))
	}

	subSquare = (subSquare / subGridRows) * subGridRows
	return subSquare, nil
}

func computeBaseColumn(subSquare int) (int, error) {
	if subSquare < 0 || subSquare > numSubSquares {
		if subSquare < 0 || subSquare > numSubSquares {
			return -1, errors.New(fmt.Sprintf("INVALID SUB SQUARE: %d", subSquare))
		}
	}

	subSquare = (subSquare % subGridColumns) * subGridColumns
	return subSquare, nil
}


func validateSubGrid(gameState *gameState, subSquare int) (err error) {
	var baseRow, baseColumn int
	baseRow, err = computeBaseRow(subSquare)
	if err != nil {
		return err
	}
	baseColumn, err = computeBaseColumn(subSquare)
	if err != nil {
		return err
	}


	maxRow := baseRow + 3
	maxColumn := baseColumn + 3

	var valueCounts map[int]int = make(map[int]int)
	var values []int = make([]int, 0)
	var errorValues []int = make([]int, 0)
	var valid bool = true
	var value int = NotSet

	for row := baseRow; row < maxRow; row++ {
		for column := baseColumn; column < maxColumn; column++ {
			value = gameState.Grid[row][column]
			if value != NotSet {
				count, ok := valueCounts[value]
				if !ok {
					values = append(values, value)
					count = 0
				}

				valueCounts[value] = count + 1
			}
		}
	}

	for _, value := range values {
		count, _ := valueCounts[value]
		if count > 1 || count == 0 {
			valid = false
			errorValues = append(errorValues, value)
		}
	}

	if !valid {
		return errors.New(fmt.Sprintf("the following values are too represented in the sub grid %d: %v",
			subSquare,
			errorValues))
	} else {
		return nil
	}


}

func validateRow(gameState *gameState, row int) (err error) {
	if row < 0 || row >= numRows {
		return errors.New(fmt.Sprintf("row passed into validate row is not valid: %d", row))
	}

	var valueCounts map[int]int = make(map[int]int)
	var values []int = make([]int, 0)
	var errorValues []int = make([]int, 0)
	var valid bool = true
	var value int = NotSet

	for n := 0; n < numColumns; n++ {
		value = gameState.Grid[row][n]
		if value != -1 {
			count, ok := valueCounts[value]
			if !ok {
				values = append(values, value)
				count = 0
			}

			valueCounts[value] = count + 1
		}
	}

	for _, val := range values {
		if valueCounts[val] != 1 {
			valid = false
			errorValues = append(errorValues, val)
		}
	}

	if !valid {
		return errors.New(fmt.Sprintf("Row %d has the following duplicated values: %v", row, errorValues))
	} else {
		return nil
	}

}

func validateColumn(gameState *gameState, column int) (err error) {
	if column < 0 || column >= numRows {
		return errors.New(fmt.Sprintf("column passed into validate column is not valid: %d", column))
	}

	var valueCounts map[int]int = make(map[int]int)
	var values []int = make([]int, 0)
	var errorValues []int = make([]int, 0)
	var valid bool = true
	var value int = NotSet

	for n := 0; n < numRows; n++ {
		value = gameState.Grid[n][column]
		if value != NotSet {
			count, ok := valueCounts[value]
			if !ok {
				values = append(values, value)
				count = 0
			}

			valueCounts[value] = count + 1
		}
	}

	for _, val := range values {
		if valueCounts[val] != 1 {
			valid = false
			errorValues = append(errorValues, val)
		}
	}

	if !valid {
		return errors.New(fmt.Sprintf("Column %d has the following duplicated values: %v", column, errorValues))
	} else {
		return nil
	}

}


func validateGameState(gs *gameState) (err error) {
	/*Validate Sub Grids*/

	for n := 0; n < numSubSquares; n++ {
		err = validateSubGrid(gs, n)
		if err != nil {
			return err
		}
	}


	/*Validate Rows*/
	for n := 0; n < numRows; n++ {
		err = validateRow(gs, n)
		if err != nil {
			return err
		}
	}

	/*Validate Columns*/
	for n := 0; n < numColumns; n++ {
		err = validateRow(gs, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func isFinished(gs *gameState) bool {
	count := gs.setCount()
	return count == (numRows * numColumns)
}

func countSelected(gameState *gameState) (count int) {
	count = 0
	for row := 0; row < numRows; row++ {
		for col := 0; col < numColumns; col++ {
			if gameState.Grid[row][col] >= 0 {
				count += 1
			}
		}
	}

	return count
}

func createGame(game *Game) (*gameState, error) {
	var err error
	var gs *gameState = &gameState{
		moves: make(candidateList, 0, 9 * 9),
		GamePlayStatistics: &GamePlayStatistics{
			BackTracks: 0,
		},
	}

	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			gs.Grid[row][column] = game.Grid[row][column]
		}
	}


	err = validateGameState(gs)
	if err != nil {
		return nil, err
	}


	return gs, nil
}

func snapToGrid(n int) (int, int) {
	var div int = n / 3
	var min int = (3 * div)
	var max int = min + 3

	return min, max
}




func (cl candidateList) back() *candidate {
	var length int = len(cl)
	if length == 0 {
		return nil
	} else {
		return cl[length - 1]
	}
}

func (cl candidateList) popBack() candidateList {
	var length int = len(cl)
	if length == 0 {
		return nil
	} else {
		return cl[0:length - 1]
	}
}

func createAllCandidatesList() candidateList {
	var ret candidateList = make(candidateList, 0, 9 * 9 * 9)
	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			for value := 0; value < numCandidates; value++ {
				var candidate *candidate = &candidate{
					value: value,
					row: row,
					column: column,
				}

				ret = append(ret, candidate)
			}
		}
	}

	return ret
}

func createValidCandidateList(gs *gameState, allCandidates candidateList) candidateList {
	var ret candidateList = make(candidateList, 0, len(allCandidates))
	var err error
	for _, c := range allCandidates {
		if !gs.isSet(c.row, c.column) {
			gs.addCandidate(c)
			err = validateGameState(gs)
			if err == nil {
				ret = append(ret, c)
			} else {
				//log.Printf("%s", err)
			}

			gs.removeCandidate(c)
		} else {
			//log.Printf("skipping candidate for being already set: %v", *c)
		}
	}
	//log.Printf("returning %d candidates", len(ret))
	return ret
}

func shuffleCandidates(candidates candidateList) {
	length := len(candidates)
	for i, c := range candidates {
		j := rand.Int() % length
		if j < 0 {
			j = -j
		}
		candidates[i] = candidates[j]
		candidates[j] = c
	}
}


func PlayGame(initialGameState *Game, maxIterations int) (*Game, int, error) {
	var gs *gameState
	var err error

	gs, err = createGame(initialGameState)
	if err != nil {
		return nil, NotSet, err
	}

	var allCandidates candidateList = createAllCandidatesList()
	var tree searchTree = make(searchTree, 0, 81)
	var moves candidateList = make(candidateList, 0, 81)
	var playing bool = !isFinished(gs)
	var iteration int = 0
	for playing && (iteration < maxIterations){
		iteration++

		var candidates candidateList = createValidCandidateList(gs, allCandidates)
		shuffleCandidates(candidates)
		if len(candidates) == 0 {
			//TODO
			playing = false
		} else {
			var candidate *candidate = candidates.back()
			candidates = candidates.popBack()
			tree = append(tree, candidates)
			moves = append(moves, candidate)
			gs.addCandidate(candidate)
			gs.Grid[candidate.row][candidate.column] = candidate.value
			//log.Printf("gs.setCount(): %d", gs.setCount())
			playing = !isFinished(gs)
		}
	}

	gs.GamePlayStatistics.Iterations = iteration

	var g *Game = NewGame()
	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			g.Grid[row][column] = gs.Grid[row][column]
		}
	}

	return g, gs.GamePlayStatistics.Iterations, nil
}

