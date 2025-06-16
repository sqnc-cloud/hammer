use clap::Parser;
use std::env;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    #[arg(short, long)]
    pub uri: String,

    #[arg(short, long)]
    pub database: String,

    #[clap(short, long, value_parser, num_args = 1.., value_delimiter = ' ')]
    pub collections: Vec<String>,
}

fn main() {
    let args = Args::parse();

    println!("URI {}!", args.uri);
    println!("DB {}!", args.database);
    println!("collections {:?}!", args.collections);
    println!("{:?}", env::current_dir());
}
