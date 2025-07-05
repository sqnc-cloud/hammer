use std::default;

use clap::{Command, arg};

use crate::blockchain::Blockchain;

pub struct Cli {
    bc: Blockchain,
}

impl Cli {
    pub fn new() -> Result<Cli, failure::Error> {
        Ok(Cli {
            bc: Blockchain::new()?,
        })
    }

    pub fn run(&mut self) -> Result<(), failure::Error> {
        let matches = Command::new("blockchain")
            .version("0.1")
            .author("Iago Ferreira <iago-ferreira@outlook.com>")
            .about("a proof of concept blockchain")
            .subcommand(Command::new("print"))
            .about("print all the blocks")
            .subcommand(Command::new("clear"))
            .about("clear the blockchain")
            .subcommand(Command::new("generate"))
            .arg(arg!(<DATA>"block data"))
            .get_matches();

        if let Some(matches) = matches.subcommand_matches("generate") {
            if let Some(data) = matches.get_one::<String>("DATA") {
                self.add_block(String::from(data))?;
            }
        }

        if matches.subcommand_matches("print").is_some() {
            self.print_chain()?;
        }

        if matches.subcommand_matches("print").is_some() {
            self.clear_blockchain()?;
        }

        Ok(())
    }

    pub fn add_block(&mut self, data: String) -> Result<(), failure::Error> {
        println!("Adding block with data {data}");
        &self.bc.add_block(data)?;
        Ok(())
    }
    pub fn print_chain(&self) -> Result<(), failure::Error> {
        for block in self.bc.iter() {
            println!("{:?}", block)
        }
        Ok(())
    }
    pub fn clear_blockchain(&self) -> Result<(), failure::Error> {
        let default_database = sled::open("/dev/null");
        let old_database = mem::replace(&mut self.bc.db, default_database);
        Ok(())
    }
}
