use std::io::Write;
use color_eyre::eyre::Result;
use whoami::realname;
fn main() -> std::result::Result<(), color_eyre::Report> {
    println!("Welcome {} to crabscript!", realname());

    start()?;
    Ok(())
}

static PROMPT: &str = "->>";

fn start() -> Result<()> {
    loop {
        print!("{} ",PROMPT);
        std::io::stdout().flush()?;
    
        // take a line from stdin and process it
        let mut line = String::new();
        std::io::stdin().read_line(&mut line)?;

        let lex = ast::lexer::Lexer::new(line.into());

        for tok in lex.into_iter() {
            println!("{:?}", tok);
        }
    }
    Ok(())
}