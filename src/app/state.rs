use core::fmt;

pub struct App {
    pub path: String,
}

impl<'a> App {
    pub fn path() -> String {
        const DEFAULT_SERVE_DIRECTORY: &str = "public";

        match std::env::var("PWA_PATH") {
            Ok(p) => p,
            Err(e) => {
                // FIXME Fallback first to an argument passed in the command line
                // THEN default to a "public" folder
                println!("{}", e);
                println!("Defaulting to {}...", DEFAULT_SERVE_DIRECTORY);
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

    pub fn new() -> App {
        App { path: App::path() }
    }
}

impl fmt::Display for App {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
        write!(
            f,
            "\nApp name: {}\nApp version: {}\nApp path: {}",
            App::name(),
            App::version(),
            self.path
        )
    }
}
