package game

import (
	"log"
	"os"
	"errors"
	"sync"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "Solver", log.LstdFlags)
}


type CandidateListCreator interface {
	createCandidates(gs *gameState, allCandidates candidateList) candidateList
}

type AsyncCandidateListCreator struct {
	CheckChildDepth int
}

func (creator *AsyncCandidateListCreator) createCandidates(gs *gameState,
	allCandidate candidateList) candidateList {


	var wait sync.WaitGroup
	wait.Add(len(allCandidate))
	var candidateChannel chan *candidate = make(chan *candidate, len(allCandidate))
	var validCandidates candidateList = make(candidateList, 0)
	var checkChild bool = gs.setCount() >= creator.CheckChildDepth

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

/*
Backtracker interface that is actually a struct. Make this an interface in the future if we care.
 */
type Backtracker struct {

}

func doBackTrack (backTracks int, gs *gameState, moves candidateList, tree searchTree) (*gameState, candidateList, searchTree, error) {
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
		return doBackTrack(backTracks + 1, gs, moves, tree)
	} else {
		c = candidates.back()
		candidates = candidates.popBack()
		tree = append(tree, candidates)
		gs.addCandidate(c)
		moves = append(moves, c)
		return gs, moves, tree, nil
	}

}

func (b Backtracker) backTrack(gs *gameState, moves candidateList, tree searchTree) (*gameState, candidateList, searchTree, error) {
	return doBackTrack(0, gs, moves, tree)
}

/**
Solver with configuration to turn off logging and adjust log output levels.
 */
type Solver struct {
	Log bool
	IterationReportInterval int
	MaximumIterations int
	ChildCreator CandidateListCreator
}

/**
Creates a Solver with some default configurations.
The configuration options turn off logging and set both the IterationReportingInterval
and the MaximumIterations to entirely arbitrary values.
CheckChildren is also set to True.  Setting this to false will result in a considerably higher run time.
 */
func CreateSolver() *Solver {
	var ret *Solver = &Solver {
		Log: false,
		IterationReportInterval: 100,
		MaximumIterations: 1000000,
		ChildCreator: &AsyncCandidateListCreator{
			CheckChildDepth: 50,
		},
	}

	return ret
}


/*
Solves a sudoku puzzle.
Returns a game state representing the best the solver could do (should solve anything that is solvable,
some game play statistics that are not well implemented, and an error that is nil if everything went well.
 */
func (solver *Solver) Solve(game *Game) (*Game, *GamePlayStatistics, error) {
	if game == nil {
		return nil, nil, errors.New("game is nil on call to Solve.")
	}

	var (
		gs *gameState
		err error
	)

	gs, err = createGame(game)
	if err != nil {
		return nil, nil, err
	}

	var allCandidates candidateList = createAllCandidatesList()
	var tree searchTree = make(searchTree, 0, 81)
	var moves candidateList = make(candidateList, 0, 81)
	var playing bool = !isFinished(gs)
	var iteration int = 0
	var snapShotModulo int = solver.IterationReportInterval
	var maxIterations int = solver.MaximumIterations
	var logging bool = solver.Log
	var gamePlayStatistics *GamePlayStatistics = &GamePlayStatistics{}
	var childCreator CandidateListCreator = solver.ChildCreator

	if logging {
		logger.Printf("solver.Solve initiating")
	}

	for playing && (gamePlayStatistics.Iterations < maxIterations) {
		gamePlayStatistics.Iterations++
		if (iteration % snapShotModulo) == 0 {
			if logging {
				log.Printf("iteration: %d", iteration)
				log.Printf("game state set count: %d, len(moves): %d, len(tree): %d, len(tree.back()): %d",
					gs.setCount(), len(moves), len(tree), len(tree.back()))
				log.Printf("gs.setCount(): %d", gs.setCount())
			}
		}

		var candidates candidateList = childCreator.createCandidates(gs, allCandidates)
		shuffleCandidates(candidates)
		if len(candidates) == 0 {
			gs, moves, tree, err = backTrack(gs, moves, tree)
			if err != nil {
				return nil,  gamePlayStatistics, err
			} else {
				if logging {
					printTreeHistograms(tree, 60)
				}

			}
		} else {
			var candidate *candidate = candidates.back()
			candidates = candidates.popBack()
			tree = append(tree, candidates)
			moves = append(moves, candidate)
			gs.addCandidate(candidate)
			playing = !isFinished(gs)
			if logging {
				printTreeHistograms(tree, 60)
			}
		}
	}

	var g *Game = NewGame()
	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			g.Grid[row][column] = gs.Grid[row][column]
		}
	}

	return g, gamePlayStatistics, nil
}




