# DNSUpdater - Like Dynamic DNS for Cloudflare DNS

DNSUpdater updates Cloudflare DNS with you NAT'ed (external ip) address. It's like Dynamic DNS for Cloudflare.

DNSUpdater uses API Tokens and not API Keys. The API Token you setup for this requires only the following permissions:

* Zone - Zone Settings - Read
* Zone - Zone - Read
* Zone - DNS - Edit

## Example Usage
```shell
dnsupdater -config myconfig.json
```

## Cron Example
```shell
#IP Monitoring & Updateing
0 */6 * * * /home/user/jobs/dnsupdater/dnsupdater -config /home/user/jobs/dnsupdater/config.json
```

## License

This project is licensed under the BSD 3-clause license