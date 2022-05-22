use rocket::tokio::sync::broadcast::{Sender, channel};
use std::{sync::Mutex, collections::HashMap};
use rocket::serde::{Serialize, Deserialize};
use super::sim;

pub struct Crufst {
  pub games: Mutex<HashMap<String, Game>>,
}

#[derive(Copy, Clone, Serialize, Deserialize, PartialEq)]
#[serde(crate = "rocket::serde")]
pub enum Player {
  Owner,
  Joiner
}

#[derive(Clone, Serialize, Deserialize, PartialEq)]
#[serde(crate = "rocket::serde")]
pub enum Message {
  Join,
  Place{player: Player, col: usize},
  Win(Player),
  Tie,
  Close,
}

pub struct Game {
  pub state: GameState,
  pub code: String,
  pub msgs: Sender<Message>,
  pub sim: sim::Board,
}

#[derive(PartialEq, Debug)]
pub enum GameState {
  Waiting,
  Ready
}

impl Game {
  pub fn new(code: &String) -> Self {
    Self {
      state: GameState::Waiting,
      code: code.clone(),
      msgs: channel(1024).0,
      sim: sim::Board::new(),
    }
  }
}

impl Crufst {
  pub fn new() -> Self {
    Self {
      games: Mutex::new(HashMap::new()),
    }
  }
}