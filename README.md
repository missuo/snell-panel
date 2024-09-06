# Snell Panel
Snell Panel for Surge

## Usage

1. Start the server

```bash
./snell-panel -token=your_token
```

2. Install Snell Server on your server

```bash
bash <(curl -Ls https://raw.githubusercontent.com/missuo/snell-panel/main/snell-install.sh) install your_panel_url your_token

or

bash <(curl -Ls https://snell.owo.nz/install/snell-install.sh) install your_panel_url your_token
```

3. Get the subscription info from the panel

```bash
http://your_panel_url/subscribe?token=your_token
```

## TODO

- [x] Add web UI

## License

GPL-3.0