use core::fmt;
use std::path::PathBuf;

pub struct App<'a> {
    pub name: &'a str,
    pub version: &'a str,
    pub path: PathBuf,
}

impl fmt::Display for App<'_> {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> Result<(), std::fmt::Error> {
        write!(
            f,
            "\nApp name: {}\nApp version: {}\nApp path: {}",
            self.name,
            self.version,
            self.path.display()
        )
    }
}
