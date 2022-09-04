<p align="center">
  <img src="https://cdn.qrclip.io/images/qrclip-github.png" alt="qrclip" />
</p>
<p></p>
<h2 align="center">Command Line Interface (CLI)</h2>
<br>

<h2>Transfer any file to any device via QR Code</h2>


<a href="https://www.qrclip.io">QRClip.io</a> 
Send and receive encrypted data to and from any device without using your personal cloud or email accounts by simply scanning QR codes.
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
  qrclip r -i QRClipID -s QRClipSubID -k 32CharactersEncryptionKeyEncodedInBase64
  qrclip r -u "QRClipURL"

SELECT STORAGE
 Select storage:
  qrclip storage

ENCRYPT (OFFLINE)
 Encrypt with automatic generated key:
  qrclip e -f fileToEncrypt
 Encrypt with a specified key:
  qrclip e -f fileToEncrypt -k 32CharactersEncryptionKeyEncodedInBase64

DECRYPT (OFFLINE)
 qrclip d -f fileToDecrypt -k 32CharactersEncryptionKeyEncodedInBase64

GENERATE KEY ENCODED IN BASE64
 Generate a random key:
  qrclip g
 Generate a key with phrase:
  qrclip g -p TheMountainFlyingOverTheRedRiver
   < 32 characters, X's are appended
   > 32 characters, the text is shortened to 32 characters

HELP
 qrclip h
```