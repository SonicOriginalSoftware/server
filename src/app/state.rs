pub struct App {}

impl App {
    // FIXME This probably does not belong here
    pub fn name() -> &'static str {
        env!("CARGO_PKG_NAME")
    }

    pub fn version() -> &'static str {
        env!("CARGO_PKG_VERSION")
    }

    pub fn new() -> App {
        App {}
    }
}
