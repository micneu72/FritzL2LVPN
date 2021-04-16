package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/sethvargo/go-password/password"
)

func main() {

	//Angaben Seite A
	nameSeiteAPtr := flag.String("aname", "Seite A", "Name der Seite A")
	ipSeiteAPtr := flag.String("aip", "192.168.178.0", "IP Seite A")
	DYNDNSA := flag.String("adns", "a.myfritz.net", "Dyndns Seite A")

	//Angaben Seite B
	nameSeiteBPtr := flag.String("bname", "Seite B", "Name der Seite A")
	ipSeiteBPtr := flag.String("bip", "192.168.180.0", "IP Seite A")
	DYNDNSB := flag.String("bdns", "b.myfritz.net", "Dyndns Seite B")
	//PSK Generieren
	PSK, err := password.Generate(36, 10, 10, false, true)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	//Ausgabe Seite A
	fmt.Printf("Name der VPN Seite A:\t%s\n", *nameSeiteAPtr)
	fmt.Printf("\t  IP Seite A:\t%s\n", *ipSeiteAPtr)
	fmt.Printf("      DYNDNS Seite A:\t%s\n", *DYNDNSA)
	fmt.Println("=========")

	//Ausgabe Seite B
	fmt.Printf("Name der VPN Seite B:\t%s\n", *nameSeiteBPtr)
	fmt.Printf("\t  IP Seite B:\t%s\n", *ipSeiteBPtr)
	fmt.Printf("      DYNDNS Seite B:\t%s\n", *DYNDNSB)
	fmt.Println("=========")
	fmt.Printf("PSK:\t\t\t\"%s\"\n", PSK)
	fmt.Println("=========")

	fmt.Println("ipsec.secret")
	fmt.Printf("@%s @%s : \"%s\"\n", *DYNDNSB, *DYNDNSB, PSK)
	fmt.Println("=========")

	alipsec_conf := `conn ###NAME###
    aggressive = yes
    fragmentation = yes
    keyexchange = ikev1
    mobike = yes
    reauth = yes
    rekey = yes
    forceencaps = no
    installpolicy = yes
    type = tunnel
    dpdaction = restart
    dpddelay = 10s
    dpdtimeout = 60s

    left=###DYNDNSA###
    #left = %any
    right = ###DYNDNSB###

    leftid = ###DYNDNSA###
    ikelifetime = 3600s
    lifetime = 3600s
    ike = aes256-sha512-modp1024!
    leftauth = psk
    rightauth = psk
    rightid = @###DYNDNSB###
    reqid = 1
    rightsubnet = ###IPSEITEB###/24
    leftsubnet = ###IPSEITEA###/24
    esp = aes256-sha512-modp1024!
    auto = route

`

	bfipsec_conf := `vpncfg {
	connections {
		enabled = yes;
		editable = yes;
		conn_type = conntype_lan;
		name = "###NAME###";
		boxuser_id = 0;
		always_renew = yes;
		reject_not_encrypted = no;
		dont_filter_netbios = yes;
		localip = 0.0.0.0;
		local_virtualip = 0.0.0.0;
		remotehostname = "###DYNDNSA###";
		remote_virtualip = 0.0.0.0;
		localid {
			fqdn = "###DYNDNSB###";
		}
		remoteid {
			fqdn = "###DYNDNSA###";
		}
		mode = phase1_mode_idp;
		phase1ss = "all/all/all";
		keytype = connkeytype_pre_shared;
		key = "###PSK###";
		cert_do_server_auth = no;
		use_nat_t = yes;
		use_xauth = no;
		use_cfgmode = no;
		phase2localid {
			ipnet {
				ipaddr = ###IPSEITEB###;
				mask = 255.255.255.0;
			}
		}
		phase2remoteid {
			ipnet {
				ipaddr = ###IPSEITEA###;
				mask = 255.255.255.0;
			}
		}
		phase2ss = "esp-all-all/ah-none/comp-all/pfs";
		accesslist = "permit ip any ###IPSEITEA### 255.255.255.0";
	}
	ike_forward_rules = "udp 0.0.0.0:500 0.0.0.0:500",
			"udp 0.0.0.0:4500 0.0.0.0:4500";
						
} 
`
	conNAME := "\"" + *nameSeiteAPtr + " => " + *nameSeiteBPtr + "\""
	// Replace
	alipsec_conf = strings.Replace(alipsec_conf, "###NAME###", conNAME, -1)
	alipsec_conf = strings.Replace(alipsec_conf, "###DYNDNSA###", *DYNDNSA, -1)
	alipsec_conf = strings.Replace(alipsec_conf, "###DYNDNSB###", *DYNDNSB, -1)
	alipsec_conf = strings.Replace(alipsec_conf, "###DYNDNSA###", *DYNDNSA, -1)
	alipsec_conf = strings.Replace(alipsec_conf, "###PSK###", PSK, -1)
	alipsec_conf = strings.Replace(alipsec_conf, "###IPSEITEA###", *ipSeiteAPtr, -1)
	alipsec_conf = strings.Replace(alipsec_conf, "###IPSEITEB###", *ipSeiteBPtr, -1)

	// Replace
	bfipsec_conf = strings.Replace(bfipsec_conf, "###NAME###", conNAME, -1)
	bfipsec_conf = strings.Replace(bfipsec_conf, "###DYNDNSA###", *DYNDNSA, -1)
	bfipsec_conf = strings.Replace(bfipsec_conf, "###DYNDNSB###", *DYNDNSB, -1)
	bfipsec_conf = strings.Replace(bfipsec_conf, "###DYNDNSA###", *DYNDNSA, -1)
	bfipsec_conf = strings.Replace(bfipsec_conf, "###PSK###", PSK, -1)
	bfipsec_conf = strings.Replace(bfipsec_conf, "###IPSEITEA###", *ipSeiteAPtr, -1)
	bfipsec_conf = strings.Replace(bfipsec_conf, "###IPSEITEB###", *ipSeiteBPtr, -1)

	fmt.Printf("%s", alipsec_conf)
	fmt.Printf("%s", bfipsec_conf)
}
