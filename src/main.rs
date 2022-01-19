use pwa_server::server::serve::serve;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    serve()
}
