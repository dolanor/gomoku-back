package ia

import (
  "./../referee"
  "./../protocol"
  "fmt"
  "unsafe"
)

var totalMapCopies uintptr = 0
/*
 * We start off by guiding ai to build up sequences, and then emphasis on taking pawns
 */
const (
  TWO_ALIGN = 1
  THREE_ALIGN = 2
  FOUR_ALIGN = 3
  // Compute : base + pawn taken
  BASE_PAWN_TAKEN = 4
  // Most important, wins over the rest every time
  FIVE_ALIGN = 500
  NON_INIT = -42
)
/**
 * Plays a pawn for player at index idx if possible, otherwise returns false
 * @param  m      map
 * @param  idx    index to play
 * @param  player player to play for
 * @return bool   true if played
 */
func playableIdx(m []protocol.MapData, idx int, player int) bool {
  cell := m[idx]
  if (cell.Empty && cell.Playable) {
    m[idx].Empty = false
    m[idx].Playable = false
    m[idx].Player = player

    return true
  }
  return false
}

func eval(m []protocol.MapData, player int, capture int) (int) {
  val := 0
  val += TWO_ALIGN * CountSequences(m, player, 2)
  val += THREE_ALIGN * CountSequences(m, player, 3)
  val += FOUR_ALIGN * CountSequences(m, player, 4)
  val += FIVE_ALIGN * CountSequences(m, player, 5)
  return val
}

func min(m []protocol.MapData, player int, depth int, capture int) (int, []protocol.MapData, int) {
  if (depth == 0) {
    return eval(m, player, capture), m, capture
  }

  min_val := NON_INIT
  var ret []protocol.MapData = nil
  ncaptured := 0
  ok := false
  var end bool

  for i := 0; i < protocol.MAP_SIZE; i++ {
    totalMapCopies += 1
    tmpMap := make([]protocol.MapData, len(m))
    copy(tmpMap, m)
    if playableIdx(tmpMap, i, player) {
      tmpMap, capture, _, ok = referee.Exec(tmpMap, i)
      if (end) {
        fmt.Println("end from min")
        // TODO : handle end
        return eval(tmpMap, player, capture), tmpMap, capture
      }
      if (ok) {
        val, nMap, ncap := max(tmpMap, player, depth - 1, capture)
        if (val < min_val || min_val == NON_INIT) {
          ret = nMap
          min_val = val
          ncaptured = ncap
        }
      }
    }
  }
  return min_val, ret, ncaptured
}

func max(m []protocol.MapData, player int, depth int, capture int) (int, []protocol.MapData, int) {
  if (depth == 0) {
    return eval(m, player, capture), m, capture
  }

  max_val := NON_INIT
  var ret []protocol.MapData = nil
  ncaptured := 0
  ok := false
  var end bool

  for i := 0; i < protocol.MAP_SIZE; i++ {
    totalMapCopies += 1
    tmpMap := make([]protocol.MapData, len(m))
    copy(tmpMap, m)
    if playableIdx(tmpMap, i, player) {
      tmpMap, capture, end, ok = referee.Exec(tmpMap, i)
      if (end) {
        fmt.Println("end from max")
        // TODO : handle end
        return eval(tmpMap, player, capture), tmpMap, capture
      }
      if (ok) {
        val, nMap, ncap := min(tmpMap, player, depth - 1, capture)
        if (val > max_val || max_val == NON_INIT) {
          ret = nMap
          max_val = val
          ncaptured = ncap
        }
      }
    }
  }
  return max_val, ret, ncaptured
}

func MinMax(m []protocol.MapData, player int, depth int) ([]protocol.MapData, int) {
  //  fmt.Println(m)
  totalMapCopies = 0
  ret, nmap, captured := max(m, player, depth, 0)
  fmt.Println("MinMax for player", player, ret)
  fmt.Println("Total map copies", totalMapCopies)
  fmt.Print("Total byte allocated by operation ")
  fmt.Println((totalMapCopies * (uintptr(len(m)) * unsafe.Sizeof(m[0]))) / 1000000, "mo")
  //  fmt.Println(nmap)

  return nmap, captured
}
