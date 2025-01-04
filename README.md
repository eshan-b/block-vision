<img src=".github/logo.svg" alt="block-vision Logo" width="350">
<p>

**block-vision** is a sleek terminal-based interface for real-time cryptocurrency prices, powered by the CoinMarketCap API. Built with Go, it's perfect for crypto enthusiasts who love working from the command line.

## Showcase

![App Showcase](.github/demo.gif)

## How to Run

1. Clone the repository:
   ```bash
   git clone https://github.com/eshan-b/block-vision.git
   cd block-vision
   ```
2. Run the executable from your terminal:
   ```bash
   ./block-vision.exe
   ```
> [!NOTE]
> If you modify the source code, you have to rebuild the executable:
> ```bash
> go mod tidy
> go build .
> ```

## Roadmap

This is just the beginning! **block-vision** has a lot of room to grow and evolve.

- [x] View trending cryptocurrencies
- [x] Search cryptocurrencies and view info
- [ ] View portfolio with Metamask sign-in

## Technologies Used

| Technology    | Purpose                                       |
|---------------|-----------------------------------------------|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea)    | MVU framework for building terminal UIs      |
| [Lipgloss](https://github.com/charmbracelet/lipgloss)       | Styling for text and colors in the terminal   |
| [Bubbles](https://github.com/charmbracelet/bubbles)         | Component library for terminal UIs         |
| [CoinGecko API](https://www.coingecko.com/en/api)           | Fetches real-time cryptocurrency data        |

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE.txt) file for details.