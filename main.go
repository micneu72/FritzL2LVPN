package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

func genPSK() string {
	//PSK Generieren
	PSK, err := password.Generate(36, 10, 10, false, true)
	if err != nil {
		log.Fatal(err)
	}
	return PSK
}

func create_ipsec_secret(a string, b string, c string, d string, e string) string {
	z := strings.Replace(a, "###PSK###", b, -1)
	z = strings.Replace(z, "###DYNDNSA###", c, -1)
	z = strings.Replace(z, "###DYNDNSB###", d, -1)
	z = strings.Replace(z, "###IP###", e, -1)
	return z
}

func create_ipsec_conf(a string, b string, c string, d string, e string, f string, g string) string {
	// Replace
	z := strings.Replace(a, "###NAME###", b, -1)
	z = strings.Replace(z, "###DYNDNSA###", c, -1)
	z = strings.Replace(z, "###DYNDNSB###", d, -1)
	z = strings.Replace(z, "###IPSEITEA###", e, -1)
	z = strings.Replace(z, "###IPSEITEB###", f, -1)
	z = strings.Replace(z, "###PSK###", g, -1)
	return z
}

func createDir(a string, b string, c string, d string) string {
	DirPATH := "./" + a + "_" + b + "/" + c + "/" + d
	os.MkdirAll(DirPATH, os.ModePerm)
	return DirPATH
}

func wConfig(path string, filename string, config string) string {
	path = path + "/" + filename
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.WriteString(config)
	if err2 != nil {
		log.Fatal(err2)
	}
	return path
}

func createNftables(nftables string, ip string) string {
	nft_conf := strings.Replace(nftables, "###IPSEITEB###", ip, -1)
	return nft_conf
}
func main() {
	PSK := genPSK()

	//Angaben Seite A
	nameSeiteAPtr := flag.String("aname", "Seite A", "Name der Seite A")
	ipSeiteAPtr := flag.String("aip", "192.168.178.0", "IP Seite A")
	DYNDNSA := flag.String("adns", "a.myfritz.net", "Dyndns Seite A")

	//Angaben Seite B
	nameSeiteBPtr := flag.String("bname", "Seite B", "Name der Seite B")
	ipSeiteBPtr := flag.String("bip", "192.168.180.0", "IP Seite A")
	DYNDNSB := flag.String("bdns", "b.myfritz.net", "Dyndns Seite B")
	pconf := flag.Bool("P", false, "VPN Konfig anzeigen")
	wconf := flag.Bool("W", false, "VPN Konfig schreiben")

	flag.Parse()
	//vpnsite := [][]string{{*nameSeiteAPtr, *ipSeiteAPtr, *DYNDNSA}, {*nameSeiteBPtr, *ipSeiteBPtr, *DYNDNSB}}

	lipsec_conf := `conn "###NAME###"
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

	fipsec_conf := `vpncfg {
        connections {
                enabled = yes;
                editable = yes;
                conn_type = conntype_lan;
                name = "###NAME###";
                always_renew = no;
                reject_not_encrypted = no;
                dont_filter_netbios = yes;
                localip = 0.0.0.0;
                local_virtualip = 0.0.0.0;
                remoteip = 0.0.0.0;
                remote_virtualip = 0.0.0.0;
                remotehostname = "###DYNDNSA###";
                localid {
                        fqdn = "###DYNDNSB###";
                }
                remoteid {
                        fqdn = "###DYNDNSA###";
                }
                mode = phase1_mode_aggressive;
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
}`

	nftables := `### nftables.conf
	nft insert rule nat postrouting oifname <LAN> ip daddr ###IPSEITEB###/24 accept
	nft add rule filter forward iifname <LAN> oifname <WAN> ip saddr ###IPSEITEB###/24 ct state new accept`

	ipsec_secret := `
#### ipsec.secret ###IP###
@###DYNDNSA### @###DYNDNSB### : PSK "###PSK###"`

	conNAME := *nameSeiteAPtr + " <=> " + *nameSeiteBPtr

	alipsec_conf := create_ipsec_conf(lipsec_conf, conNAME, *DYNDNSA, *DYNDNSB, *ipSeiteAPtr, *ipSeiteBPtr, PSK)
	blipsec_conf := create_ipsec_conf(lipsec_conf, conNAME, *DYNDNSB, *DYNDNSA, *ipSeiteBPtr, *ipSeiteAPtr, PSK)
	afipsec_conf := create_ipsec_conf(fipsec_conf, conNAME, *DYNDNSB, *DYNDNSA, *ipSeiteBPtr, *ipSeiteAPtr, PSK)
	bfipsec_conf := create_ipsec_conf(fipsec_conf, conNAME, *DYNDNSA, *DYNDNSB, *ipSeiteAPtr, *ipSeiteBPtr, PSK)

	ipsec_secret_a := create_ipsec_secret(ipsec_secret, PSK, *DYNDNSA, *DYNDNSB, *ipSeiteAPtr)
	ipsec_secret_b := create_ipsec_secret(ipsec_secret, PSK, *DYNDNSB, *DYNDNSA, *ipSeiteBPtr)

	//nftables = strings.Replace(nftables, "###IPSEITEB###", *ipSeiteBPtr, -1)
	nftables_a := createNftables(nftables, *ipSeiteBPtr)
	nftables_b := createNftables(nftables, *ipSeiteAPtr)

	printSetup(*nameSeiteAPtr, *ipSeiteAPtr, *DYNDNSA)
	printSetup(*nameSeiteBPtr, *ipSeiteBPtr, *DYNDNSB)

	if *pconf {
		fmt.Printf("%s\n", ipsec_secret_a)
		fmt.Println("========")
		fmt.Printf("%s\n", ipsec_secret_b)
		fmt.Println("========")
		fmt.Printf("%s", alipsec_conf)
		fmt.Println("========")
		fmt.Printf("%s", blipsec_conf)
		fmt.Println("========")
		fmt.Printf("%s", afipsec_conf)
		fmt.Println("\n========")
		fmt.Printf("%s", bfipsec_conf)
		fmt.Println("\n========")
		fmt.Printf("%s\n", nftables_a)
		fmt.Println("========")
		fmt.Printf("%s\n", nftables_b)

	}

	if *wconf {
		// Verzeichnisse Seite A
		pFa := createDir(*nameSeiteAPtr, *nameSeiteBPtr, *ipSeiteAPtr, "Fritzbox")
		pLa := createDir(*nameSeiteAPtr, *nameSeiteBPtr, *ipSeiteAPtr, "strongSwan")
		// Verzeichnisse Seite B
		pFb := createDir(*nameSeiteAPtr, *nameSeiteBPtr, *ipSeiteBPtr, "Fritzbox")
		pLb := createDir(*nameSeiteAPtr, *nameSeiteBPtr, *ipSeiteBPtr, "strongSwan")

		cpL_nft_a := wConfig(pLa, "nft.conf", nftables_a)
		cpL_nft_b := wConfig(pLb, "nft.conf", nftables_b)

		cpLa_secret := wConfig(pLa, "ipsec.secret", ipsec_secret_a)
		cpLb_secret := wConfig(pLb, "ipsec.secret", ipsec_secret_a)

		cpLa_conf := wConfig(pLa, "ipsec.conf", alipsec_conf)
		cpLb_conf := wConfig(pLb, "ipsec.conf", blipsec_conf)

		cpFa_cfg := wConfig(pFa, *ipSeiteAPtr+".cfg", afipsec_conf)
		cpFb_cfg := wConfig(pFb, *ipSeiteBPtr+".cfg", bfipsec_conf)

		fmt.Printf("\nDatei erstellt:\n\t%s \n\t%s \n\t%s \n\t%s \n\t%s \n\t%s \n\t%s \n\t%s\n", cpLa_secret, cpLb_secret, cpLa_conf, cpLb_conf, cpFa_cfg, cpFb_cfg, cpL_nft_a, cpL_nft_b)

	}
}
