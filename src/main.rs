use rxing;

fn main() {
    let file_name = "./samples/qr.png";

    let detect_results = rxing::helpers::detect_multiple_in_file(file_name);

    match detect_results {
        Ok(results) => print_results(results),
        Err(e) => println!("Error: {}", e),
    };
}

fn print_results(results: Vec<rxing::RXingResult>) {
    for result in results {
        println!("{} -> {}", result.getBarcodeFormat(), result.getText())
    }
}
