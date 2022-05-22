use super::state::Player;

#[derive(Copy, Clone, PartialEq)]
pub enum Cell {
  Empty,
  Filled(Player)
}

const BOARD_COLS: usize = 7;
const BOARD_ROWS: usize = 6;

pub struct Board {
  cells: [[Cell; BOARD_ROWS]; BOARD_COLS], // Access by doing cells[col][row]
  last_turn: Player,
}

#[derive(PartialEq)]
pub enum SimResult {
  Failed,
  Continue,
  Win,
  Tie
}

impl Board {
  pub fn new() -> Self {
    Self {cells: [[Cell::Empty; BOARD_ROWS]; BOARD_COLS], last_turn: Player::Joiner}
  }

  pub fn place(&mut self, pl: Player, col: usize) -> SimResult {
    if pl == self.last_turn { // Going twice in a row
      return SimResult::Failed;
    }

    if col > BOARD_COLS - 1 { // Too far to the right
      return SimResult::Failed;
    }

    let c = self.cells[col];
    let mut last = BOARD_COLS;
    for (i, val) in c.iter().enumerate().rev() {
      if *val == Cell::Empty {
        last = i;
        break;
      }
    }

    if last == BOARD_COLS { // Column is full
      return SimResult::Failed;
    }

    // Update board
    self.cells[col][last] = Cell::Filled(pl);
    self.last_turn = pl;

    // Check for connect 4
    for offset_col in -1..2 {
      for offset_row in -1..2 {
        if offset_col == 0 && offset_row == 0 {
          continue
        }

        // In bounds
        let (valid, rowval, colval) = calcoffset(last, col, offset_row, offset_col);
        if !valid {
          continue;
        }

        // Is your own
        let isfilled = self.cells[colval][rowval] == Cell::Filled(pl);
        if !isfilled {
          continue
        }

        // Check for four
        let mut currcol = col;
        let mut currrow = last;
        let mut works = true;
        for _ in 0..4 {
          (works, currrow, currcol) = calcoffset(currrow, currcol, offset_row, offset_col);
          if !works {
            break;
          }

          if self.cells[currcol][currrow] != Cell::Filled(pl) {
            works = false;
            break;
          }
        }

        if works { // They won!
          return SimResult::Win;
        }
      }
    }

    // Check for tie
    let mut is_tie = true;
    'outer: for col in self.cells {
      for cell in col {
        if cell == Cell::Empty {
          is_tie = false;
          break 'outer;
        }
      }
    }

    if is_tie {
      return SimResult::Tie;
    }

    SimResult::Continue
  }
}

fn calcoffset(row: usize, col: usize, rowoff: i8, coloff: i8) -> (bool, usize, usize) {
  let rowv = (row as i8) + rowoff;
  let colv = (col as i8) + coloff;
  if (rowv < 0) || (rowv >= BOARD_ROWS as i8) || (colv < 0) || (colv >= BOARD_COLS as i8) {
    return (false, 0, 0);
  }

  (true, rowv as usize, colv as usize)
}