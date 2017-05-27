package game

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"errors"
)


func TestSolver_Solve(t *testing.T) {
	var game *Game
	var err error
	var solver = CreateSolver()
	game = NewGame()
	game, _, err = solver.Solve(game)

	assert.Nil(t, err)
}


func Test_cast(t *testing.T) {
	doCast := func(ob interface{}) error {
		switch ob.(type) {
		case int:
			return nil
		default:
			return errors.New("not an int")
		}
	 }


	err := doCast(4)
	assert.Nil(t, err)

	err = doCast("I am not an integer")
	assert.NotNil(t, err)

	doCast2 := func(ob interface{}) (int, error) {
		val, ok := ob.(int)
		if !ok {
			return -1, errors.New("not an int")
		} else {
			return val, nil
		}
	}
	var val int
	val, err = doCast2(4)
	assert.Nil(t, err)
	assert.Equal(t, 4, val)

	val, err = doCast2("I am not an integer")
	assert.NotNil(t, err)
	assert.Equal(t, -1, val)
}


func Test_AsyncCandidateListCreator_createCandidateList(t *testing.T) {
	var candidateListCreator AsyncCandidateListCreator = AsyncCandidateListCreator{
		CheckChildDepth: 50,
	}

	gs, err := createGame(NewGame())
	assert.Nil(t, err)
	var allCandidates candidateList = createAllCandidatesList()
	var candidates candidateList = candidateListCreator.createCandidates(gs, allCandidates)
	assert.NotEmpty(t, candidates)
}