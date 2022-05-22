use random_string;
use rocket::http::Status;
use rocket::response::stream;
use rocket::Either;
use rocket::serde::{Deserialize};

pub mod state;
use state::*;

mod sim;

/*
View in browser by:
1. Go to /crufst/new
2. Execute 
new EventSource("http://127.0.0.1:49154/crufst/events/" + document.body.innerText)
3. Execute
await (await fetch("http://127.0.0.1:49154/crufst/join/" + document.body.innerText)).text()
4. View messages in Network tab
*/

const CODE_LENGTH: usize = 5;

#[get("/new")]
pub fn newgame(state: &rocket::State<Crufst>) -> String {
  loop {
    let code: String = random_string::generate(CODE_LENGTH, "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890");
    if !state.games.lock().unwrap().contains_key(&code) {
      state.games.lock().unwrap().insert(code.clone(), Game::new(&code));
      return code;
    }
  }
}

#[get("/join/<code>")]
pub fn join(state: &rocket::State<Crufst>, code: String) -> (Status, String) {
  if !state.games.lock().unwrap().contains_key(&code) {
    return (Status::InternalServerError, "Invalid game code.".to_string());
  }

  let mut lock = state.games.lock().unwrap();
  let game = lock.get_mut(&code).unwrap();
  if game.state != GameState::Waiting {
    return (Status::InternalServerError, "Game already full.".to_string());
  }

  // Send join event
  game.state = GameState::Ready;
  if let Err(e) = game.msgs.send(Message::Join) {
    return (Status::InternalServerError, e.to_string());
  }

  (Status::Ok, "Successfully joined game.".to_string())
}

#[get("/events/<code>")]
pub async fn events(state: &rocket::State<Crufst>, code: String) -> Either<stream::EventStream![], (Status, &'static str)> {
  if !state.games.lock().unwrap().contains_key(&code) {
    return Either::Right((Status::InternalServerError, "Invalid game code."));
  }

  let lock = state.games.lock().unwrap();
  let game = lock.get(&code).unwrap();
  let mut sub = game.msgs.subscribe();

  
  Either::Left(stream::EventStream! {
    loop {
      match sub.recv().await {
        Ok(n) => {
          yield stream::Event::json(&n);
    
          if n == Message::Close {
            break;
          }
        }
        Err(_) => break
      }
    }
  })
}

#[derive(Deserialize)]
#[serde(crate = "rocket::serde")]
pub struct PlaceEvent {
  player: Player,
  col: usize,
}

#[post("/place/<code>", data = "<ev>")]
pub fn place(state: &rocket::State<Crufst>, code: String, ev: rocket::serde::json::Json<PlaceEvent>) -> (Status, String) {
  if !state.games.lock().unwrap().contains_key(&code) {
    return (Status::InternalServerError, "Invalid game code.".to_string());
  }

  let mut lock = state.games.lock().unwrap();
  let game = lock.get_mut(&code).unwrap();

  // Try placing
  let res = game.sim.place(ev.player, ev.col);
  if res == sim::SimResult::Failed {
    return (Status::BadRequest, "Cannot place there.".to_string());
  }

  // Update boards
  if let Err(e) = game.msgs.send(Message::Place{player: ev.player, col: ev.col}) {
    return (Status::InternalServerError, e.to_string());
  }
  match res {
    sim::SimResult::Continue => {
      (Status::Ok, "Successfully placed.".to_string())
    },
    sim::SimResult::Win => {
      // Send Win message
      if let Err(e) = game.msgs.send(Message::Win(ev.player)) {
        return (Status::InternalServerError, e.to_string());
      }
      (Status::Ok, "Game completed.".to_string())
    },
    sim::SimResult::Tie => {
      // Send Win message
      if let Err(e) = game.msgs.send(Message::Tie) {
        return (Status::InternalServerError, e.to_string());
      }
      (Status::Ok, "Game tied.".to_string())
    },
    sim::SimResult::Failed => unreachable!(),
  }
}