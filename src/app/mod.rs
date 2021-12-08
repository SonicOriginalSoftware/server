use crate::app::state::App;
use std::path::Path;

use actix_files::NamedFile;
use actix_web::{http, web, HttpRequest, HttpResponse, Responder, Result};

pub mod state;

pub async fn root() -> impl Responder {
    HttpResponse::MovedPermanently()
        .set_header(http::header::LOCATION, "/app")
        .finish()
}

pub async fn app(request: HttpRequest, app_state: web::Data<App<'_>>) -> Result<NamedFile> {
    let path = request.match_info().query("filename");
    Ok(NamedFile::open(
        Path::new(&app_state.path).join(if path == "" { "index.html" } else { path }),
    )?)
}
