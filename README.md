# Net-Cat project

## Overview

Net-Cat is a simple terminal-based chat application written in Go. It allows multiple clients to connect to a server and exchange messages in real-time. Clients can also change their nicknames and text colors.

## Features

- **Real-time Messaging**: Clients can send and receive messages instantly.
- **Nickname Management**: Users can set and change their nicknames.
- **Color Customization**: Users can customize their text color.
- **Command Help**: Users can request a list of available commands.
- **Chat History**: The server maintains a history of chat messages and sends it to new clients upon joining.
- **Log File**: All chat messages are logged to a file named with the current date and time.

## Installation

1. **Clone the Repository**:

    ```sh
    git clone https://zone01normandie.org/git/gkopoin/net-cat.git
    ```

2. **Build the Project**:

    ```sh
    go build -o TCPChat
    ```

3. **Run the Server**:

    ```sh
    ./TCPChat
    ```

## Usage

1. **Start the Server**:

    Run the built binary to start the server.

    ```sh
    ./TCPChat
    ```

2. **Connect as a Client**:

    Use `telnet` or `netcat` (nc) to connect to the server. Replace `localhost` and `port` with the server's address and port number.

    ```sh
    telnet localhost 12345
    ```

    or

    ```sh
    nc localhost 8989
    ```

3. **Commands**:

    - **/nick `<new_name>`**: Change your nickname.
    - **/color `<color_name>`**: Change your text color. Supported colors include `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`.
    - **/help**: Show available commands.

## Logging

The server logs all messages to a file. The log file is named based on the current date and time and is created in the same directory as the server executable.