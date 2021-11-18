package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
	customErrors "backendServer/pkg/errors"
)

type BoardUseCaseImpl struct {
	boardRepository    repositories.BoardRepository
	userRepository     repositories.UserRepository
	teamRepository     repositories.TeamRepository
	cardListRepository repositories.CardListRepository
	cardRepository     repositories.CardRepository
}

func CreateBoardUseCase(
	boardRepository repositories.BoardRepository,
	userRepository repositories.UserRepository,
	teamRepository repositories.TeamRepository,
	cardListRepository repositories.CardListRepository,
	cardRepository repositories.CardRepository,
) usecases.BoardUseCase {
	return &BoardUseCaseImpl{
		boardRepository:    boardRepository,
		userRepository:     userRepository,
		teamRepository:     teamRepository,
		cardListRepository: cardListRepository,
		cardRepository:     cardRepository,
	}
}

func (boardUseCase *BoardUseCaseImpl) GetUserBoards(uid uint) (teams *[]models.Team, err error) {
	teams, err = boardUseCase.userRepository.GetUserTeams(uid)
	if err != nil {
		return
	}

	for i, team := range *teams {
		boards, boardsErr := boardUseCase.teamRepository.GetTeamBoards(team.TID)
		if boardsErr != nil {
			err = boardsErr
			return
		}
		(*teams)[i].Boards = *boards
	}

	return
}

func (boardUseCase *BoardUseCaseImpl) CreateBoard(board *models.Board) (bid uint, err error) {
	err = boardUseCase.boardRepository.Create(board)
	if err != nil {
		return 0, err
	}
	return board.BID, nil
}

func (boardUseCase *BoardUseCaseImpl) GetBoard(uid, bid uint) (board *models.Board, err error) {
	isAccessed, err := boardUseCase.userRepository.IsBoardAccessed(uid, bid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	board, err = boardUseCase.boardRepository.GetByID(bid)
	if err != nil {
		return
	}

	lists, err := boardUseCase.boardRepository.GetBoardCardLists(bid)
	if err != nil {
		return nil, err
	}

	for i, list := range *lists {
		var cards *[]models.Card
		cards, err = boardUseCase.cardListRepository.GetCardListCards(list.CLID)
		if err != nil {
			return
		}
		for j, card := range *cards {
			var comments *[]models.Comment
			comments, err = boardUseCase.cardRepository.GetCardComments(card.CID)
			if err != nil {
				return
			}
			(*cards)[j].Comments = *comments
		}
		(*lists)[i].Cards = *cards
	}
	board.CardLists = *lists

	return
}

func (boardUseCase *BoardUseCaseImpl) UpdateBoard(uid uint, board *models.Board) (err error) {
	isAccessed, err := boardUseCase.userRepository.IsBoardAccessed(uid, board.BID)
	if err != nil {
		return err
	}
	if !isAccessed {
		return customErrors.ErrNoAccess
	}

	return boardUseCase.boardRepository.Update(board)
}

func (boardUseCase *BoardUseCaseImpl) DeleteBoard(uid, bid uint) (err error) {
	isAccessed, err := boardUseCase.userRepository.IsBoardAccessed(uid, bid)
	if err != nil {
		return err
	}
	if !isAccessed {
		return customErrors.ErrNoAccess
	}
	return boardUseCase.boardRepository.Delete(bid)
}
