use core::fmt;

pub struct App<'a> {
    pub name: &'a str,
    pub version: &'a str,
    pub path: String,
}

impl<'a> App<'a> {
    fn path() -> String {
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

    pub fn new() -> App<'a> {
        App {
            name: env!("CARGO_PKG_NAME"),
            version: env!("CARGO_PKG_VERSION"),
            path: App::path(),
        }
    }
}

impl fmt::Display for App<'_> {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
        write!(
            f,
            "\nApp name: {}\nApp version: {}\nApp path: {}",
            self.name, self.version, self.path
        )
    }
}
