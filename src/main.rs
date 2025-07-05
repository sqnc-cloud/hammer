use clap::Parser;
use std::env;
use std::fs::File;
use std::io::{Read, Seek, Write};
use std::path::Path;
use std::process::Command;
use walkdir::{DirEntry, WalkDir};
use zip::write::SimpleFileOptions;

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

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::parse();

    println!("URI {}", args.uri);
    println!("DB {}", args.database);
    println!("collections {:?}", args.collections);

    let current_dir = env::current_dir().unwrap().display().to_string();
    let output_dir = format!("{}/backup/", current_dir.clone());
    let output_zip = format!("{}/backup.zip", current_dir.clone());
    for collection in &args.collections {
        let output_file = format!("{}{}.json", output_dir.clone(), collection.clone());

        println!("{output_file}");

        let output = Command::new("mongoexport")
            .arg("--uri")
            .arg(format!("mongodb://{}", args.uri.clone()))
            .arg("--db")
            .arg(args.database.clone())
            .arg("--collection")
            .arg(collection.clone())
            .arg("--out")
            .arg(output_file.clone())
            .arg("--jsonFormat")
            .arg("canonical")
            .arg("--type")
            .arg("json")
            .output()
            .expect("failed to execute process");

        println!("{:?}", output);
    }

    if !Path::new(&output_dir).is_dir() {
        panic!("Output directory does not exist");
    }
    let output_zip_path = Path::new(&output_zip);
    let file = File::create(output_zip_path).unwrap();
    let walkdir = WalkDir::new(output_dir.clone());
    let it = walkdir.into_iter();

    zip_dir(
        &mut it.filter_map(|e| e.ok()),
        output_dir.clone().as_ref(),
        file,
        zip::CompressionMethod::Stored,
    )?;
    return Ok(());
}
fn zip_dir<T>(
    it: &mut dyn Iterator<Item = DirEntry>,
    prefix: &Path,
    writer: T,
    method: zip::CompressionMethod,
) -> anyhow::Result<()>
where
    T: Write + Seek,
{
    let mut zip = zip::ZipWriter::new(writer);
    let options = SimpleFileOptions::default()
        .compression_method(method)
        .unix_permissions(0o755);

    let mut buffer = Vec::new();
    for entry in it {
        let path = entry.path();
        let name = path.strip_prefix(prefix)?;
        let path_str = name
            .to_str()
            .ok_or_else(|| anyhow::anyhow!("{:?} is not a valid UTF-8 path", name))?;

        if path.is_file() {
            println!("adding file {path:?} as {name:?} ...");
            zip.start_file(path_str, options)?;
            File::open(path)?.read_to_end(&mut buffer)?;
            zip.write_all(&buffer)?;
            buffer.clear();
        } else if !name.as_os_str().is_empty() {
            println!("adding dir {path_str:?} as {name:?} ...");
            zip.add_directory(path_str, options)?;
        }
    }
    zip.finish()?;
    Ok(())
}

