# WireCrab 🦀

A lightweight network traffic analyzer built with Go. Still under active development, but already functional for basic packet capturing and analysis.

## What's this?

WireCrab is my attempt at creating a simpler, more modern network traffic analyzer. Think Wireshark, but more focused on today's debugging needs. I started this project because I got tired of the complexity of existing tools when all I needed was to quickly check what's going on with my network traffic.

## Current State

Right now it can:
- Capture live traffic from network interfaces
- Display basic packet info (protocols, IPs, sizes)
- Show detailed packet analysis
- Works on Windows and Linux (macOS support coming soon)

**Warning**: This is a work in progress! Expect bugs and missing features. I'm actively working on it, but use at your own risk.

## Requirements

- Go 1.24 or newer
- Tshark (yeah, ironically we still need it - for now)
- Admin/root privileges (for packet capture)

## Roadmap

- [ ] Packet filtering
- [ ] Packet diff (byte by byte)
- [ ] Resend packet
- [ ] Generate the source code in different languages to send the same selected packet

## Contributing

Feel free to jump in! I'm still figuring out the direction of this project, but if you find it interesting:

1. Check out the issues
2. Fork & create a branch
3. Submit a PR

Just keep in mind things might change a lot as the project evolves.

## License

Apache License 2.0 - See [LICENSE](LICENSE) file

## Why Another Packet Analyzer?

Look, I know there are tons of packet analyzers out there. But most feel either too complex or too simple. I wanted something in between - powerful enough for daily development work, but without the learning curve of Wireshark. Plus, I thought it would be fun to build one from scratch.
