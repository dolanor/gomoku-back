package ia

import (
  "./../protocol"
  "./../referee"
  "fmt"
 // "strconv"
)

const (
  TWO_ALIGN = 1
  THREE_ALIGN = 5
  FOUR_ALIGN = 10
  // Compute : base + pawn taken
  BASE_PAWN_TAKEN = 4
  // Most important, wins over the rest every time
  FIVE_ALIGN = 500
  MAX_INIT = -42000
  MIN_INIT = 42000
)

var mapCopies uintptr = 0
var smallestIndex int
var highestIndex int

type minMaxStruct struct {
  M      []protocol.MapData
  Player int8
  Depth  int
  End    bool
}

func getOtherPlayer(player int8) int8 {
  if (player == 0) {
    return 1
  }
  return 0
}

/**
 * Plays a pawn for player at index idx if possible, otherwise returns false
 * @param  m      map
 * @param  idx    index to play
 * @param  player player to play for
 * @return bool   true if played
 */
func playIdx(m []protocol.MapData, idx int, player int8) bool {
  cell := m[idx]
  if (cell.Empty && cell.Playable) {
    m[idx].Empty = false
    m[idx].Playable = false
    m[idx].Player = player

    return true
  }
  return false
}

func eval(data *minMaxStruct) int {
  val := 0
  val += TWO_ALIGN * CountSequences(data.M, data.Player, 2)
  val += THREE_ALIGN * CountSequences(data.M, data.Player, 3)
  val += FOUR_ALIGN * CountSequences(data.M, data.Player, 4)
  val += FIVE_ALIGN * CountSequences(data.M, data.Player, 5)
  return val
}

func max(data *minMaxStruct) int {
  if (data.Depth == 0 || data.End) {
    data.Player = getOtherPlayer(data.Player)
    return eval(data)
  }
  max := MAX_INIT

  mapcp := make([]protocol.MapData, len(data.M))
  for i := smallestIndex; i < highestIndex; i++ {
    copy(mapcp, data.M)
    mapCopies += 1
    if playIdx(mapcp, i, data.Player) {
      _, end, valid := referee.Exec(mapcp, i)
      if (valid) {
        tmp := min(&minMaxStruct{mapcp, getOtherPlayer(data.Player), data.Depth - 1, end})
        if (tmp > max) {
          max = tmp
        }
      }
    }
  }
  return max
}

func min(data *minMaxStruct) int {
  if (data.Depth == 0 || data.End) {
    data.Player = getOtherPlayer(data.Player)
    return eval(data)
  }
  min := MIN_INIT

  mapcp := make([]protocol.MapData, len(data.M))
  for i := smallestIndex; i < highestIndex; i++ {
    copy(mapcp, data.M)
    mapCopies += 1
    if playIdx(mapcp, i, data.Player) {
      _, end, valid := referee.Exec(mapcp, i)
      if (valid) {
        tmp := max(&minMaxStruct{mapcp, getOtherPlayer(data.Player), data.Depth - 1, end})
        if (tmp < min) {
          min = tmp
        }
      }
    }
  }
  return min
}
func initSmallMax(m []protocol.MapData) {
  smallestIndex = -1
  for i := 0; i < protocol.MAP_SIZE; i++ {
    if (!m[i].Empty) {
      if (smallestIndex == -1) {
        smallestIndex = i
      }
      highestIndex = i
    }
  }
  if (smallestIndex > (19 * 2)) {
    smallestIndex -= (19 * 2)
  }
  if ((highestIndex + (19 * 2)) < protocol.MAP_SIZE) {
    highestIndex += (19 * 2)
  }
  // fmt.Println("iteration window size", highestIndex - smallestIndex, "(" + strconv.Itoa(smallestIndex) + "->" + strconv.Itoa(highestIndex) + ")")
}

func MinMax(m []protocol.MapData, player int8, depth int) (int) {
  initSmallMax(m)
  max := MAX_INIT
  maxIdx := 0
  mapCopies = 0

  mapcp := make([]protocol.MapData, len(m))
  for i := smallestIndex; i < highestIndex; i++ {
    copy(mapcp, m)
    mapCopies += 1
    if playIdx(mapcp, i, player) {
      _, end, valid := referee.Exec(mapcp, i)
      if (valid) {
        tmp := min(&minMaxStruct{mapcp, player, depth - 1, end})
        if (tmp > max) {
          fmt.Println(tmp, i)
          max = tmp
          maxIdx = i
        }
      }
    }
  }
  // fmt.Println("Total map cp", mapCopies)
  return maxIdx
  //  fmt.Println(m)
  //  fmt.Println(map)
}
