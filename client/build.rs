use walkdir::{WalkDir};
use std::path::PathBuf;
use std::io;

fn get_files_in_dir(path: &str) -> Result<Vec<PathBuf>, io::Error> {
    let mut files = Vec::new();

    for entry in WalkDir::new(path) {
        let entry = entry?;
        if entry.file_type().is_file() {
            files.push(entry.path().to_owned());
        }
    }

    Ok(files)
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let proto_files = get_files_in_dir("proto")?;

    let proto_files_refs: Vec<&std::path::Path> = proto_files
        .iter()
        .map(|p| p.as_path())
        .filter(|p| p.extension().unwrap_or_default() == "proto")
        .collect();

    tonic_build::configure()
        .build_server(false)
        .compile(&proto_files_refs, &["proto"])?;
    Ok(())
}
