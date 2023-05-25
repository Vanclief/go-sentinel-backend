use rodio::{source::Source, Decoder, OutputStream};
use std::fs;
use std::fs::File;
use std::io::BufReader;
use std::sync::mpsc;
use std::thread;
use std::time::Duration;

use rxing;

fn main() {
    let dir_path = "./tmp";

    let (_stream, stream_handle) = OutputStream::try_default().unwrap();
    let file = BufReader::new(File::open("sounds/beep.mp3").unwrap());
    let source = Decoder::new(file).unwrap();
    stream_handle.play_raw(source.convert_samples()).unwrap();
    let (tx, rx) = mpsc::channel();

    println!("Scanning directory: {}", dir_path);

    loop {
        scan_directory(dir_path, tx.clone());
        thread::sleep(Duration::from_millis(500));

        for received in rx.try_iter() {
            println!("Received: {}", received);
        }
    }
}

fn scan_directory(dir_path: &str, tx: mpsc::Sender<String>) -> Result<(), std::io::Error> {
    let mut handles = vec![];
    let entries = fs::read_dir(dir_path)?;

    for entry in entries {
        let entry = entry?;
        let path = entry.path();

        if path.is_dir() {
            continue;
        }

        if let Some(extension) = path.extension() {
            if ["jpg", "jpeg", "png", "gif"].contains(&extension.to_str().unwrap_or_default()) {
                match path.to_str() {
                    Some(file_name) => {
                        let file_name = file_name.to_owned();
                        let thread_tx = tx.clone();

                        let handle = thread::spawn(move || {
                            detect_barcodes(&file_name);
                            fs::remove_file(&file_name).expect("Failed to remove file");
                            thread_tx.send(file_name).unwrap();
                        });
                        handles.push(handle);
                    }
                    _ => (),
                }
            }
        }
    }

    for handle in handles {
        handle.join().expect("Thread panicked");
    }

    Ok(())
}

fn detect_barcodes(file_name: &str) {
    let detect_results = rxing::helpers::detect_multiple_in_file(file_name);
    match detect_results {
        Ok(results) => print_results(results),
        Err(_) => (),
    };
}

fn print_results(results: Vec<rxing::RXingResult>) {
    for result in results {
        print_result(result);
    }
}

fn print_result(result: rxing::RXingResult) {
    println!("{} -> {}", result.getBarcodeFormat(), result.getText());
}
