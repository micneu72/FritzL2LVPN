package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/sethvargo/go-password/password"
)

/*func printSetup() {
	fmt.Printf("%s", alipsec_conf)
	fmt.Println("========")
	fmt.Printf("%s", bfipsec_conf)
	fmt.Printf("======== %s ipsec.secret", *nameSeiteAPtr)
	fmt.Printf("%s", ipsec_secret)
}*/

func printSetup(a string, b string, c string) {
	//Ausgabe Seite A
	fmt.Printf("Name der VPN:\t%s\n", a)
	fmt.Printf("\t  IP:\t%s\n", b)
	fmt.Printf("      DYNDNS:\t%s\n", c)
	fmt.Println("=========")
}

func printPSK(a string, b string, c string) {
	fmt.Println("ipsec.secret")
	fmt.Printf("@%s @%s : \"%s\"\n", a, b, c)
	fmt.Println("=========")
}
func genPSK() string {
	//PSK Generieren
	PSK, err := password.Generate(36, 10, 10, false, true)
	if err != nil {
		log.Fatal(err)
	}
	return PSK
}

func main() {

	//Angaben Seite A
	nameSeiteAPtr := flag.String("aname", "Seite A", "Name der Seite A")
	ipSeiteAPtr := flag.String("aip", "192.168.178.0", "IP Seite A")
	DYNDNSA := flag.String("adns", "a.myfritz.net", "Dyndns Seite A")

	//Angaben Seite B
	nameSeiteBPtr := flag.String("bname", "Seite B", "Name der Seite B")
	ipSeiteBPtr := flag.String("bip", "192.168.180.0", "IP Seite A")
	DYNDNSB := flag.String("bdns", "b.myfritz.net", "Dyndns Seite B")
	pconf := flag.Bool("P", false, "VPN Konfig anzeigen")

	PSK := genPSK()
	printSetup(*nameSeiteAPtr, *ipSeiteAPtr, *DYNDNSA)
	printSetup(*nameSeiteBPtr, *ipSeiteBPtr, *DYNDNSB)

	flag.Parse()

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

	ipsec_secret := `
	   @###DYNDNSA### @###DYNDNSB### : PSK "###PSK###"
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
	bfipsec_conf = strings.Replace(bfipsec_conf, "###IPSEITEA###", *ipSeiteAPtr, -1)
	bfipsec_conf = strings.Replace(bfipsec_conf, "###IPSEITEB###", *ipSeiteBPtr, -1)

	ipsec_secret = strings.Replace(ipsec_secret, "###PSK###", PSK, -1)
	ipsec_secret = strings.Replace(ipsec_secret, "###DYNDNSB###", *DYNDNSB, -1)
	ipsec_secret = strings.Replace(ipsec_secret, "###DYNDNSA###", *DYNDNSA, -1)

	if *pconf {
		printPSK(*DYNDNSA, *DYNDNSB, PSK)
		fmt.Printf("%s", alipsec_conf)
		fmt.Println("========")
		fmt.Printf("%s", bfipsec_conf)
	}

}
