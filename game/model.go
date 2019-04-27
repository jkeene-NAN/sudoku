package game

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
	"encoding/json"

)

const numCandidates int = 9
const numRows int = 9
const numColumns int = 9
const numSubSquares int = 9
const NotSet = -1
const DefaultMaxIterations = 1000000000
const CheckChildrenDepthThreshold int = 50

func init() {
	rand.Seed(time.Now().Unix())

	log.SetOutput(os.Stdout)
}

type candidate struct {
	value  int
	row    int
	column int
}

type candidateList []*candidate

func (c candidate) clone() *candidate {
	var ret *candidate = &candidate{
		value:  c.value,
		row:    c.row,
		column: c.column,
	}

	return ret
}

func (c *candidate) equals(other *candidate) bool {
	return c.row == other.row && c.column == other.column && c.value == other.value
}

/*
Sub Grids
0|1|2
3|4|5
6|7|8
*/

type Game struct {
	Grid [][]int
}

func gridToString(grid [][]int) string {
	var buf bytes.Buffer
	buf.WriteString("\n")
	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			var value int = grid[row][column]
			if value == NotSet {
				buf.WriteString("-")
			} else {
				buf.WriteString(fmt.Sprintf("%d", value))
			}

			if column != numColumns-1 {
				buf.WriteString("|")
			}
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

func (gs *Game) String() string {
	return gridToString(gs.Grid)
}

func NewGame() *Game {
	var game *Game = &Game{
		Grid: make([][]int, numRows),
	}

	for row := 0; row < numRows; row++ {
		game.Grid[row] = make([]int, numRows)
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
	Grid               [][]int
	initialGameState *Game
	moves              candidateList
	GamePlayStatistics *GamePlayStatistics
}

func (gs *gameState) isMutable(can *candidate) bool {
	var row, column int

	row = can.row
	column = can.column
	return gs.initialGameState.Grid[row][column] == NotSet
}

func (gs *gameState) movesRemaining() int {
	var count = gs.setCount()
	return (numRows * numColumns) - count
}

func (gs *gameState) clone() *gameState {
	var ret *gameState = &gameState{
		Grid: make([][]int, numRows),
		initialGameState: gs.initialGameState,
	}


	for row := 0; row < numRows; row++ {
		ret.Grid[row] = make([]int, numColumns)
		for column := 0; column < numColumns; column++ {
			ret.Grid[row][column] = gs.Grid[row][column]
		}
	}

	return ret
}

func resetGameState(gs *gameState) {
	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			gs.Grid[row][column] = NotSet
		}
	}
}

func (gs *gameState) String() string {
	return gridToString(gs.Grid)
}

func (gs *gameState) Json() string {
	jsonBytes, _ := json.MarshalIndent(gs, "", "\t")
	return string(jsonBytes)
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
		return tree[length-1]
	}
}

func (tree searchTree) popBack() searchTree {
	var length int = len(tree)
	if length == 0 {
		return nil
	} else {
		return tree[0 : length-1]
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
		err = validateColumn(gs, n)
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
		Grid: make([][]int, numRows),
		moves: make(candidateList, 0, 9*9),
		GamePlayStatistics: &GamePlayStatistics{
			BackTracks: 0,
		},
		initialGameState: game,
	}

	for row := 0; row < numRows; row++ {
		gs.Grid[row] = make([]int, numColumns)
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
		return cl[length-1]
	}
}

func (cl candidateList) popBack() candidateList {
	var length int = len(cl)
	if length == 0 {
		return nil
	} else {
		return cl[0 : length-1]
	}
}

func createAllCandidatesList() candidateList {
	var ret candidateList = make(candidateList, 0, 9*9*9)
	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			for value := 0; value < numCandidates; value++ {
				var candidate *candidate = &candidate{
					value:  value,
					row:    row,
					column: column,
				}

				ret = append(ret, candidate)
			}
		}
	}

	return ret
}

func createValidCandidateListAsync(gs *gameState, allCandidate candidateList, checkChild bool) candidateList {
	var wait sync.WaitGroup
	wait.Add(len(allCandidate))
	var candidateChannel chan *candidate = make(chan *candidate, len(allCandidate))
	var validCandidates candidateList = make(candidateList, 0)

	var doCheck = func(clone *gameState, c *candidate) {
		defer wait.Done()
		var toAdd *candidate
		clone.addCandidate(c)
		var err error = validateGameState(clone)
		if err == nil {
			if checkChild {
				var childValidCandidates candidateList = createValidCandidateList(clone, allCandidate, false)
				if len(childValidCandidates) >= clone.movesRemaining() {
					toAdd = c
				} else {
					toAdd = nil
				}
			} else {
				toAdd = c
			}
		} else {
			toAdd = nil
		}

		candidateChannel <- toAdd
	}

	for _, candidate := range allCandidate {
		go doCheck(gs.clone(), candidate)
	}

	wait.Wait()

	for i := 0; i < len(allCandidate); i++ {
		var c *candidate
		c = <-candidateChannel
		if c != nil {
			validCandidates = append(validCandidates, c)
		}
	}

	return validCandidates
}

func createValidCandidateList(gs *gameState, allCandidates candidateList, checkChild bool) candidateList {
	var ret candidateList = make(candidateList, 0, len(allCandidates))
	var err error
	for _, c := range allCandidates {
		if !gs.isSet(c.row, c.column) {
			gs.addCandidate(c)
			err = validateGameState(gs)
			if err == nil {
				if checkChild {
					var childCandidates candidateList = createValidCandidateList(gs, allCandidates, false)
					if len(childCandidates) >= gs.movesRemaining() {
						ret = append(ret, c)
					}
				} else {
					ret = append(ret, c)
				}
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

type searchSnapShot struct {
	candidateCountHistogram []int
	moves                   candidateList
}

func createSearchSnapShot(gs *gameState, moves candidateList, tree searchTree) *searchSnapShot {
	var snapShot *searchSnapShot = &searchSnapShot{
		candidateCountHistogram: make([]int, len(tree)),
		moves: make(candidateList, len(moves)),
	}

	for i := 0; i < len(moves); i++ {
		snapShot.moves[i] = moves[i]
	}

	for i := 0; i < len(tree); i++ {
		snapShot.candidateCountHistogram[i] = len(tree[i])
	}

	return snapShot
}

func backTrack(gs *gameState, moves candidateList, tree searchTree) (*gameState, candidateList, searchTree, error) {
	if len(tree) == 0 {
		return nil, nil, nil, errors.New("tree is empty on call to back track")
	}

	//log.Print("backtracking...")
	var c *candidate = moves.back()
	moves = moves.popBack()
	gs.removeCandidate(c)

	var candidates candidateList = tree.back()
	tree = tree.popBack()
	if len(candidates) == 0 {
		return backTrack(gs, moves, tree)
	} else {
		c = candidates.back()
		candidates = candidates.popBack()
		tree = append(tree, candidates)
		gs.addCandidate(c)
		moves = append(moves, c)
		return gs, moves, tree, nil
	}

}

func printGameState(gs *gameState) {
	log.Print(gs.String())
}

func printTreeHistograms(tree searchTree, after int) {
	var treeLen int = len(tree)
	if treeLen < after {
		return
	} else {
		var buf bytes.Buffer
		for i := after; i < treeLen; i++ {
			var candidates candidateList = tree[i]
			var candidatesLen = len(candidates)
			buf.WriteString(fmt.Sprintf("%d:%d, ", i, candidatesLen))
		}

		log.Print(buf.String())
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
	var snapShotModulo int = 10000
	var checkChildren bool = false

	for playing && (iteration < maxIterations) {
		iteration++
		if (iteration % snapShotModulo) == 0 {
			log.Printf("iteration: %d", iteration)
		}
		if len(moves) >= CheckChildrenDepthThreshold {
			checkChildren = true
		} else {
			checkChildren = false
		}

		var candidates candidateList = createValidCandidateListAsync(gs, allCandidates, checkChildren)
		shuffleCandidates(candidates)
		if len(candidates) == 0 {
			gs, moves, tree, err = backTrack(gs, moves, tree)
			if err != nil {
				return nil, iteration, err
			} else {
				printTreeHistograms(tree, 60)
			}
		} else {
			var candidate *candidate = candidates.back()
			candidates = candidates.popBack()
			tree = append(tree, candidates)
			moves = append(moves, candidate)
			gs.addCandidate(candidate)
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
