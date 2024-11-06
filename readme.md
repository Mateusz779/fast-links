# Fast Links - Quickly Open Links on Multiple Devices

## Project Goal

Fast Links was created to streamline teamwork and increase productivity when working with multiple devices. The application allows for instant sharing of links in real-time by creating virtual rooms where all participants see the same links.

## Attention!
**The project is mainly intended for use within an intranet. Exposing it on the internet may result in unintended links being added to the room. The author is not responsible for any unlawful use of the tool.**


## Use Cases

- **Online Meetings**: Quickly share materials, presentations, and documents.
- **Teamwork**: Jointly browse websites, codes, and projects.
- **Education**: Use for interactive presentations and collaborative problem-solving.

## Features

- **Creating Virtual Rooms**: Each user can create their own room and invite others.
- **Adding Links**: Users add links that are immediately displayed on the screens of all room participants.
- **Scanning QR Codes**: Add links by scanning a QR code.
- **Scalability**: The application can handle any number of users and rooms.

## How to Use

1. **Installation**: Only Docker is required.
2. **Running**: Execute the `run_container.sh` script, optionally specifying the port using the `-p` parameter.
3. **Access**: Open the specified address in your browser.
4. **Adding Links**: Enter the link in the input field and confirm.
5. **Listening Mode**: Go to the `/app` page or specify a specific room ID using the `?id=` parameter.

## Technologies

- **Golang**: Efficient backend programming language.
- **html5-qrcode**: Library for scanning QR codes in the browser.

## Additional Suggestions

- **Documentation**: Consider creating a more detailed guide to explain advanced configuration and usage options.
- **User Interface**: Improving the UI could make the application more user-friendly.
- **Testing**: Conduct comprehensive tests to ensure application stability and reliability.
- **Security**: Ensure the security of the application, especially if it will be used to share sensitive data.