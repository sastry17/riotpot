## Trudy

Trudy is a transparent proxy that can modify and drop traffic for arbitrary TCP connections. Trudy can be used to programmatically modify TCP traffic for proxy-unaware clients. Trudy creates a 2-way "pipe" for each connection it proxies. The device you are proxying (the "client") connects to Trudy (but doesn't know this) and Trudy connects to the client's intended destination (the "server"). Traffic is then passed between these pipes. Users can create Go functions to mangle data between pipes.  [See it in action!](https://asciinema.org/a/7zkywm0biuz1wa64az3tmox8v) For a practical overview, check out [@tsusanka](https://twitter.com/tsusanka)'s very good [blog post](https://blog.susanka.eu/how-to-modify-general-tcp-ip-traffic-on-the-fly-with-trudy) on using Trudy to analyze Telegram's MTProto. 

Trudy can also proxy TLS connections. Obviously, you will need a valid certificate or a client that does not validate certificates.

Trudy was designed for monitoring and modifying proxy-unaware devices that use non-HTTP protocols. If you want to monitor, intercept, and modify HTTP traffic, Burp Suite is probably the better option. 

## Author

Written by Kelby Ludwig ([@kelbyludwig](https://twitter.com/kelbyludwig))

### Why I Built This
I have done security research that involved sitting between a embedded device and a server and modifying some custom binary protocol on the fly. This usually is a slow process that involves sniffing legitimate traffic, and then rebuilding packets programmatically. Trudy enables Burp-like features for generalized TCP traffic.

### Simple Setup

0. Configure a virtual machine (Trudy has been tested on a 64-bit Debian 8 VM) to shove all traffic through Trudy. I personally use a Vagrant VM that sets this up for me. The Vagrant VM is available [here](https://github.com/praetorian-inc/mitm-vm). If you would like to use different `--to-ports` values, you can use Trudy's command line flags to change Trudy's listening ports.

    `iptables -t nat -A PREROUTING -i eth1 -p tcp --dport 8888 -m tcp -j REDIRECT --to-ports 8080`

    `iptables -t nat -A PREROUTING -i eth1 -p tcp --dport 443 -m tcp -j REDIRECT --to-ports 6443`

    `iptables -t nat -A PREROUTING -i eth1 -p tcp -m tcp -j REDIRECT --to-ports 6666`

    `ip route del 0/0`

    `route add default gw 192.168.1.1 dev eth1`

    `sysctl -w net.ipv4.ip_forward=1`

1. Clone the repo on the virtual machine and build the Trudy binary.

    `git clone https://github.com/kelbyludwig/trudy.git`

    `cd trudy`

    `go install`

2. Run the Trudy binary as root. This starts the listeners. If you ran the `iptables` commands above, `iptables` will forward traffic destined for port 443 to port 6443. Trudy listens on this port and expects traffic coming into this port to be TLS. All other TCP connections will be forwarded through port 6666. 

    `sudo $GOPATH/bin/trudy`

3. Setup your host machine to use the virtual machine as its router. You should see connections being made in Trudy's console but not notice any traffic issues on the host machine (except TLS errors).

4. In order to manipulate data, just implement whatever functions you might need within the `module` package. The default implementations for these functions are hands-off, so if they do not make sense for your situation, feel free to leave them as they are. More detailed documentation is in the `module` package and the data flow is detailed below.


5. To access the interceptor, visit `http://<IP ADDRESS OF VM>:8888/` in your web browser. The only gotcha here is you must visit the interceptor after starting Trudy but before Trudy receives a packet that it wants to intercept. 

## Data Flow

Module methods are called in this order. Downward arrows indicate a branch if the `Do*` function returns true.

```
Deserialize -> Drop -> DoMangle ->  DoIntercept -> DoPrint -> Serialize -> BeforeWriteTo(Server|Client) -> AfterWriteTo(Server|Client)
                         |         /                |        /
                         |_> Mangle                 |_> PrettyPrint
```
