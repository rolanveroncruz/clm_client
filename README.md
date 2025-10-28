# CLM Client
The CLM Client is one of three components that compose the Certificate Lifecycle Management system.
It is the server application that will run on the server whose certificate is being managed.
It implements an http server to respond to CLM Acme's attempt to setup the HTTP-01 solution.
For now, the acme component of the client is separated, for ease of testing and development.
Once we have the Global Sign code working, we can merge the acme code into this client.





The clm client aims to have the following functions:

* Respond to ACME challenges:
  1) HTTP01 Challenge (what we are working on now, Oct 28, 2025)
  2) TLS (still have to figure out)
* Generate a Certificate Signing Request (CSR) - (currently with the acme_client, clm_acme)
* Accept a new certificate - (currently with the acme_client, clm_acme)
* Install the certificate (next step)
  * for nginx (next step)
  * apache
  * other web servers
 
### Oct 28, 2025
We currently have 3 APIs to support the HTTP-01 Challenge:
* `/acme/login` - login to the clm_acme server. This is needed to get a JWT token to be allowed to
   put the token-authorization pair using the next API.
* `/acme/.well-known/acme-challenge/put-pair` - save the token-authorization pair
* `/acme/acme-challenge/{token}` - get the authorization string for the token


## HTTP-01 Challenge
In the Http-01 Challenge the ACME server will povide:
* a token, and
* an authorization string
These are sent to clm_acme.

Then the ACME server will make a GET call to
`http:<DOMAIN>/.well-known/acme-challenge/<tokan>` in which it expects to 
receive the authorization string as a text/plain response.

To implement this, this client could either:
* create a file and store it in the server's file system mapped to `/.well-known/acme-challange`.
  *  this will probably require running this client in sudo (probably not a good idea)
  
* expose the `/.well-known/acme-challenge/<token>` as a valid URL path that can be called.
  * this will necessitate this client running at `PORT :80` because the ACME server requires it so (also not a good idea)

* expose `/acme-challenge/<token>`, but have nginx map the path `/.well-known` to this client.
  * this seems the best alternative, as it doesn't require running as `sudo` in case the file path needs is unowned, and only requires
a small configuration change for nginx.

