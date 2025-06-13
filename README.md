# clm_client
The clm_client has the following functions:

* Respond to ACME challenges:
  1) HTTP01 Challenge
  2) TLS
* Generate a Certificate Signing Request (CSR)
* Accept a new certificate
* Install the certificate
  * for nginx
  

## HTTP-01 Challenge
In the Http-01 Challenge the ACME server will povide:
* a token, and
* an authorization string

Then the ACME server will make a GET call to
`http:<DOMAIN>/.well-known/acme-challenge/<tokan>` upon which it expects to 
receive the authorization string as a text/plain response.

To implement this, this client could either:
* create a file and store it in the server's file system mapped to `/.well-known/acme-challange`.
  *  this will probably require running this client in sudo (probably not a good idea)
  
* expose the `/.well-known/acme-challenge/<token>` as a valid URL path that can be called.
  * this will necessitate this client running at `PORT :80` because the ACME server requires it so (also not a good idea)

* expose `/acme-challenge/<token>`, but have nginx map the path `/.well-known` to this client.
  * this seems the best alternative, as it doesn't require running as `sudo` in case the file path needs is unowned, and only requires
a small configuration change for nginx.

