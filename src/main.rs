use clap::Parser;
use std::env;
use std::process::Command;

#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
struct Args {
    #[arg(short, long)]
    pub uri: String,

    #[arg(short, long)]
    pub database: String,

    #[clap(short, long, value_parser, num_args = 1.., value_delimiter = ' ')]
    pub collections: Vec<String>,

    #[arg(short, long)]
    pub output_dir: String,
}

fn main() {
    let args = Args::parse();

    println!("URI {}", args.uri);
    println!("DB {}", args.database);
    println!("collections {:?}", args.collections);

    let current_dir = env::current_dir().unwrap().display().to_string();
    for collection in args.collections {
        let output_dir = format!("{}/backup/{}.json", current_dir.clone(), collection.clone(),);

        println!("{output_dir}");

        let output = Command::new("mongoexport")
            .arg("--uri")
            .arg(format!("mongodb://{}", args.uri.clone()))
            .arg("--db")
            .arg(args.database.clone())
            .arg("--collection")
            .arg(collection.clone())
            .arg("--out")
            .arg(output_dir)
            .arg("--jsonFormat")
            .arg("canonical")
            .arg("--type")
            .arg("json")
            .output()
            .expect("failed to execute process");
        println!("{:?}", output);
    }
}
