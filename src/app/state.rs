pub struct App {
    pub path: String,
}

impl<'a> App {
    // FIXME This probably does not belong here
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

    pub fn name() -> &'static str {
        env!("CARGO_PKG_NAME")
    }

    pub fn version() -> &'static str {
        env!("CARGO_PKG_VERSION")
    }

    pub fn new(path: &str) -> App {
        App {
            path: path.to_string(),
        }
    }
}
