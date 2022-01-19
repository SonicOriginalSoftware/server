use actix_web::{
    middleware,
    web::{get, resource, Data},
    App, HttpServer,
};

use {
    crate::api::route as api_route,
    crate::app::{path as app_path, root as app_root, route as app_route, state::App as AppState},
    crate::auth::route as auth_route,
};

#[actix_web::main]
pub async fn serve() -> std::io::Result<()> {
    // TODO Server should be running on HTTP/2 over SSL/TLS

    println!("{}", AppState::name());
    println!("{}", AppState::version());

    const SCHEME: &str = "http";
    const HOST: &str = "localhost";
    const PORT: u32 = 8080;

    let serve_path = app_path::path();

    println!("  Starting server on {}://{}:{}...", SCHEME, HOST, PORT);
    println!("  Serving from [{}]...", serve_path);

    // FIXME Logging not working

    HttpServer::new(move || {
        App::new()
            .wrap(middleware::Logger::default())
            .wrap(middleware::DefaultHeaders::default())
            .wrap(middleware::NormalizePath::new(
                middleware::normalize::TrailingSlash::MergeOnly,
            ))
            .app_data(Data::new(AppState::new()))
            .route("/", get().to(app_root::root))
            .route(&serve_path, get().to(app_route::route))
            .route("/{filename:.*}", get().to(app_route::route))
            .service(resource("/api").to(api_route))
            .service(resource("/auth").to(auth_route))
    })
    .bind(format!("{}:{}", HOST, PORT))?
    .run()
    .await
}
