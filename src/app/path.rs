pub fn path() -> String {
    match std::env::var("PWA_PATH") {
        Ok(p) => p,
        Err(_) => {
            // println!("{}", e);
            const DEFAULT_SERVE_DIRECTORY: &str = "public";

            println!(
                "PWA_PATH not set. Defaulting to [{}]...",
                DEFAULT_SERVE_DIRECTORY
            );
            DEFAULT_SERVE_DIRECTORY.to_string()
        }
    }
}
