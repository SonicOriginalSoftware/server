use actix_web::web;
use actix_web::{middleware, App, HttpServer};

use pwa_server::{
    api::route as api_route,
    app::{root, route as app_route, state::App as AppState},
    auth::route as auth_route,
};

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // TODO Server should be running on HTTP/2 over SSL/TLS

    const SCHEME: &str = "http";
    const HOST: &str = "localhost";
    const PORT: u32 = 8080;

    let app_state = AppState::new();

    println!("{}", app_state);

    println!("Starting server on {}://{}:{}...", SCHEME, HOST, PORT);

    // FIXME Logging not working

    // let app_data = web::Data::new(app_state);

    HttpServer::new(move || {
        App::new()
            .wrap(middleware::Logger::default())
            .wrap(middleware::DefaultHeaders::default())
            .wrap(middleware::NormalizePath::new(
                middleware::normalize::TrailingSlash::MergeOnly,
            ))
            .app_data(web::Data::new(app_state))
            .route("/", web::get().to(root))
            .route(&app_state.path, web::get().to(app_route))
            .route("/{filename:.*}", web::get().to(app_route))
            .service(web::resource("/api").to(api_route))
            .service(web::resource("/auth").to(auth_route))
    })
    .bind(format!("{}:{}", HOST, PORT))?
    .run()
    .await
}
