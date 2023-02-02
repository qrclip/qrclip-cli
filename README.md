<p align="center">
  <img src="https://cdn.qrclip.io/images/qrclip-github4.png" alt="qrclip" />
</p>
<p></p>
<h2 align="center">Command Line Interface (CLI)</h2>


<h2>Easily Share Files And Texts With Any Device</h2>
<p>Experience secure and private data transfer with self‑destructing QR codes and links</p>
<a href="https://www.qrclip.io">QRClip.io</a>
<br>

```
LOGIN
 To login using qr code
  qrclip l
 With username and password
  qrclip l -u myemail@email.com -p "MySecretPassword"
 With username (password will be asked)
  qrclip l -u myemail@email.com

LOGOUT
 To clear the credentials
  qrclip logout

CHECK LIMITS
 Check QRClip limits of the current user
  qrclip c

SEND
 qrclip s -m "Message to Send" -f fileToSend
 qrclip s -m "Message to Send"
 qrclip s -f fileToSend
 Other Options:
  -e 15      ( Expiration Time in minutes - default 15 )
  -mt 2      ( Max transfers - default 2 )
  -ad true   ( Allow delete - default true )

RECEIVE
 Receive Mode:
  qrclip r
 Get QRClip:
  qrclip r -i QRClipID -s QRClipSubID -k 32CharactersEncryptionKeyEncodedInBase64Url
  qrclip r -u "QRClipURL"

SELECT STORAGE
 Select storage:
  qrclip storage

HELP
 qrclip h
```