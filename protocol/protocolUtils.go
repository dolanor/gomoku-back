package protocol

type MapData struct {
	Empty, Playable bool
	Player          int
}

const IDLE = "IDLE"
const START_OF_GAME = "START_OF_GAME"
const PLAY_TURN = "PLAY_TURN"
const END_OF_GAME = "END_OF_GAME"
const ENTER_ROOM = "ENTER_ROOM"
const REFRESH = "REFRESH"

type MessageIdle struct {
	Type string
}

type MessageStartOfGame struct {
	Type         string
	PlayerNumber int
}

type MessagePlayTurn struct {
	Type           string
	Map            []MapData
	AvailablePawns [2]int
	CapturedPawns  [2]int
}

type MessageEndOfGame struct {
	Type           string
	Map            []MapData
	AvailablePawns [2]int
	CapturedPawns  [2]int
	Winner         int
}

type MessageEnterRoom struct {
	Type string
	Room int
}

type MessageRefresh struct {
	Type           string
	Map            []MapData
	AvailablePawns [2]int
	CapturedPawns  [2]int
}

func SendEndOfGame(m []MapData, availablePawns [2]int, capturedPawns [2]int, winner int) (*MessageEndOfGame) {
	return &MessageEndOfGame{
		END_OF_GAME,
		m,
		availablePawns,
		capturedPawns,
		winner}
}

func SendPlayTurn(m []MapData, availablePawns [2]int, capturedPawns [2]int) (*MessagePlayTurn) {
	return &MessagePlayTurn{
		PLAY_TURN,
		m,
		availablePawns,
		capturedPawns}
}

func SendStartOfGame(number int) (*MessageStartOfGame) {
	return &MessageStartOfGame{
		START_OF_GAME,
		number}
}

func SendIdle() (*MessageIdle) {
	return &MessageIdle{
		IDLE}
}

func SendRefresh(m []MapData, availablePawns [2]int, capturedPawns [2]int) (*MessageRefresh) {
	return &MessageRefresh{
		REFRESH,
		m,
		availablePawns,
		capturedPawns}
}

func InitGameData() ([]MapData, [2]int, [2]int) {
	myMap := make([]MapData, 19 * 19)
	for x := 0; x < 19 * 19; x++ {
		myMap[x].Empty = true
		myMap[x].Playable = true
		myMap[x].Player = -1
	}
	var availablePawns [2]int
	availablePawns[0] = 60
	availablePawns[1] = 60

	var capturedPawns [2]int
	capturedPawns[0] = 0
	capturedPawns[1] = 0

	return myMap, availablePawns, capturedPawns
}
