use crate::app::state::App;

use actix_files::NamedFile;
use actix_web::{web, HttpRequest, Result};

use std::path::Path;

pub async fn route(request: HttpRequest, _app_state: web::Data<App>) -> Result<NamedFile> {
    let path = request.match_info().query("filename");
    Ok(NamedFile::open(
        Path::new(&super::path::path()).join(if path == "" { "index.html" } else { path }),
    )?)
}
