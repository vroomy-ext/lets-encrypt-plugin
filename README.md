# Let's Encrypt Plugin
Let's Encrypt Plugin is a Vroomy plugin for auto-generation of SSL certificates

## Usage
Add `github.com/vroomy/lets-encrypt-plugin as letsEncrypt` to the `plugins` section of your Vroomy configuration. Be sure to run `vpn update letsEncrypt` if you've just added the plugin to your project.

### Example output
```
● Let's Encrypt :: Certificate is expired (or expiring soon), executing renewal process
● Let's Encrypt :: Client created
2020/07/15 15:29:54 [INFO] acme: Registering account for [User email]
● Let's Encrypt :: User registered
2020/07/15 15:29:54 [INFO] [your-domain.com] acme: Obtaining bundled SAN certificate
2020/07/15 15:29:55 [INFO] [your-domain.com] AuthURL: https://acme-v02.api.letsencrypt.org/acme/authz-v3/[Cert Authorization URL]
2020/07/15 15:29:55 [INFO] [your-domain.com] acme: use tls-alpn-01 solver
2020/07/15 15:29:55 [INFO] [your-domain.com] acme: Trying to solve TLS-ALPN-01
2020/07/15 15:30:01 [INFO] [your-domain.com] The server validated our request
2020/07/15 15:30:01 [INFO] [your-domain.com] acme: Validations succeeded; requesting certificates
2020/07/15 15:30:02 [INFO] [your-domain.com] Server responded with a certificate.
● Let's Encrypt :: Certificates obtained
● Let's Encrypt :: Certificate renewal process complete
```