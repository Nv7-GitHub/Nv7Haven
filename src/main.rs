#[macro_use] extern crate rocket;

mod crufst;

#[get("/")]
fn index() -> rocket::response::Redirect {
    rocket::response::Redirect::to("https://api.nv7haven.com")
}

#[launch]
fn rocket() -> _ {
    println!("Started!");
    rocket::build()
        .mount("/", routes![index])
        .mount("/crufst", routes![crufst::newgame, crufst::join, crufst::events, crufst::place])
        .manage(crufst::state::Crufst::new())
}