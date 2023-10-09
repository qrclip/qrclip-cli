![QRClip Logo](https://cdn.qrclip.io/images/qrclip-github4.png)
# QRClip - Secure CLI File & Text Sharing
Secure, fast, and end-to-end encrypted file & text sharing right from your terminal.

---

## 🌟 Features

### 📲 Interactive QR Code Terminal Sharing

With QRClip's CLI, the power of QR codes meets the flexibility of the terminal.

- **Terminal to Phone**: Create a QR code in your terminal, give it a quick scan with your phone, and voilà! The data is sent from the terminal straight to your device.

- **Phone to Terminal**: Craft a receiver QR code right in your terminal. A simple phone scan and your data zips from your phone right into the terminal.

### 🔒 End-to-End Encryption

Ensure your data's integrity and privacy.

- **XChaCha20-Poly1305**: State-of-the-art encryption to secure your data.

- **Zero Knowledge**: Your data remains confidential, even from us.

### 💾 Flexible Data Transfer

Empower your sharing with QRClip's diverse options.

- **Files & Text**: Share both textual data and files directly from your terminal.

- **Expiration Control**: Set how long your shared data remains accessible.

- **Transfer Limits**: Define how many times your data can be accessed.

### ⚙️ Utility Features

Additional tools to enhance your QRClip experience.

- **Storage Selection**: Choose where your QRClips are stored.

- **Check Limits**: Quickly view your QRClip usage limits.

### 🌍 Cross-Platform Compatibility

QRClip ensures a seamless experience across a spectrum of devices and platforms.

- **Pre-built Binaries**: No matter if you're on Linux, MacOS, or Windows, QRClip's CLI is ready for you. Pre-built binaries are available for various platforms under releases.
- **Build from Source**: Want a tailored experience? Customize and compile your own version of the CLI with minimal fuss.
- **Mobile Applications**: On the go? QRClip's dedicated apps for both [Android](https://play.google.com/store/apps/details?id=io.qrclip.app&pcampaignid=pcampaignidMKT-Other-global-all-co-prtnr-py-PartBadge-Mar2515-1) and [iOS](https://apps.apple.com/us/app/qrclip/id6446234709) make sure you can share and receive data wherever you are.
- **Web Application**: QRClip's [web app](https://app.qrclip.io) allows for easy data sharing without the need for any installation. Check it out at QRClip Web.

## 🚀 Getting Started

### Prerequisites

- Ensure you have [Go](https://golang.org/) installed (version 1.17 or later).

### Installation

#### Using Pre-built Binaries

For ease of use, we provide pre-built binaries for various platforms. You can download the appropriate version for your system from our [Releases](https://github.com/qrclip/qrclip-cli/releases) page.

#### Building from Source

If you prefer building from the source:

```bash
git clone https://github.com/qrclip/qrclip-cli.git
cd qrclip-cli/cmd/qrclip
go build
```

## 🔐 Authentication

QRClip offers a seamless authentication process tailored for both QR code enthusiasts and traditionalists.

### Logging In

You have the flexibility to log in using a QR code or with your username and password.

#### 1. Using QR Code:

Simply type in the following command, and a QR code will be generated for you to scan:

```bash
qrclip login
```

Once the QR code appears, grab your phone and scan it using the QRClip app. Voilà, you're logged in without typing a password!

#### 2. Using Username and Password:
You can directly provide your username and password:

```bash
qrclip login -u myemail@email.com -p "MySecretPassword"
```

Alternatively, just provide your username, and you'll be prompted to enter your password:

```bash
qrclip login -u myemail@email.com
```

### Logging Out

To ensure your security, you can log out and clear your credentials anytime with:

```bash
qrclip logout
```

## 📤 Data Transfer

Sharing data becomes incredibly easy and secure with QRClip. Here's how you can send and receive data right from your terminal:

### Sending Data

QRClip offers flexibility in sending both messages and files.

#### 1. Sending a Message:

To quickly send a text message, use:

```bash
qrclip send -m "Message to Send"
```

#### 2. Sending a File:

For sending a file, provide the path to the desired file:

```bash
qrclip send -f /path/to/fileToSend
```

#### 3. Sending Both Message and File:

Combine both options to send a message and file simultaneously:

```bash
qrclip send -m "Message to Send" -f /path/to/fileToSend
```

---

Once the data transfer completes, both a QR code and a direct link will be presented in your terminal. You can either:

-**Scan the QR Code:** Use your phone to swiftly download the data.

-**Share the Link:** Copy the provided link and distribute it as you see fit, granting others access to the data.

---

#### Additional Send Options:
**Expiry:** Set the lifespan of your shared data (default is 2880 mins).
```bash
-e 2880
```
**Max Transfers:** Determine the number of times your data can be accessed (default is 5).
```bash
-mt 5
```
**Allow Deletion:** Decide if the receiver can delete the shared data (default is true).
```bash
-ad true
```

### Receiving Data

Retrieve shared data effortlessly with QRClip.

#### 1. Generate a Receiver:
Use the following command to start receiving data:

```bash
qrclip receive
```

After generating the receiver, simply scan the QR code with your phone. This action will launch the QRClip app, allowing you to draft a message or select files. Once you've made your choices, just hit "send" in the app. When the QRClip transfer completes, press 'Enter' in your terminal to seamlessly download the data right there.

#### 2. Fetch a Specific QRClip:
You might need the QRClip ID, SubID, and Encryption Key:

```bash
qrclip receive -i QRClipID -s QRClipSubID -k 32CharactersEncryptionKeyEncodedInBase64Url
```

Executing this will showcase the conveyed message and initiate the file download process.

#### 3. Fetch by URL:
To retrieve data using a specific URL:

```bash
qrclip receive -u "QRClipURL"
```

Running this command will unveil the shared message and commence the file transfer.

## ⚙️ Utility Commands

QRClip provides utilities to help you understand and control your usage.

### 📊 Check Usage Limits

To understand the restrictions on your account or transfer limits:

- **Command Variants**:
    - `check`
    - `c`
    - `limits`
    - `limit`

```bash
qrclip check
```

### 🌐 Select Storage

To choose where your data gets stored or to switch between different storage options:

```bash
qrclip storage
```

## 📘 Help & Info

Whenever you're in doubt or need a refresher on how to use QRClip commands, these built-in tools are here to assist:

### 🆘 Help Guide

Feeling lost? The CLI has got your back with an onboard guide. Access it via:

- **Command Variants**:
  - `help`
  - `h`
  - `--help`

```bash
qrclip help
```

### 🏷️ Version Info

Stay up-to-date and check the version of your QRClip tool with:

- **Command Variants**:
    - `version`
    - `v`
    - `--version`

```bash
qrclip --version
```

## 🌐 Dive Deeper

Thanks for exploring QRClip's CLI documentation. For a more visual experience and to understand the diverse use-cases, visit our official [QRClip Landing Page](https://www.qrclip.io). To dive deeper into insights, updates, and intriguing articles, don't miss out on our [QRClip Blog](https://www.qrclip.io/blog).

Stay connected, stay secure. Happy sharing!