use actix_web::{HttpRequest, HttpResponse, Responder};

pub async fn route(_request: HttpRequest) -> impl Responder {
    println!("Received auth request!");
    HttpResponse::Ok().body("Should return with Auth response!")
}
