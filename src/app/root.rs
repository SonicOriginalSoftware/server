use actix_web::{http, HttpResponse, Responder};

pub async fn root() -> impl Responder {
    HttpResponse::MovedPermanently()
        .set_header(http::header::LOCATION, "/app")
        .finish()
}
