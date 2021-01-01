# go-ipaddr-server

Web server for getting your client IP Address and hosting server's egress IP Address.

## Deploy

```bash
docker run irotoris/ipaddr-server:latest
```

or deploy image to other PaaS.

## Usage

- `GET /`: Get your client ip address.
- `GET /egress`: Get egress ip address of hosting server. Optional `CheckIpUrl` GET parameter (default: <http://checkip.amazonaws.com/>)

## License

MIT
