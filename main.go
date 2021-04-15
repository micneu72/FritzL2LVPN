package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/sethvargo/go-password/password"
)

func main() {

	//Angaben Seite A
	nameSeiteAPtr := flag.String("aname", "Seite A", "Name der Seite A")
	ipSeiteAPtr := flag.String("aip", "192.168.178.0", "IP Seite A")
	maskSeiteAPtr := flag.String("amask", "255.255.255.0", "Netzmaske Seite A (255.255.255.0)")
	DYNDNSA := flag.String("adns", "a.myfritz.net", "Dyndns Seite A")

	//Angaben Seite B
	nameSeiteBPtr := flag.String("bname", "Seite B", "Name der Seite A")
	ipSeiteBPtr := flag.String("bip", "192.168.180.0", "IP Seite A")
	maskSeiteBPtr := flag.String("bmask", "255.255.255.0", "Netzmaske Seite B (255.255.255.0)")
	DYNDNSB := flag.String("bdns", "b.myfritz.net", "Dyndns Seite B")
	//PSK Generieren
	res, err := password.Generate(36, 10, 10, false, true)
	if err != nil {
		log.Fatal(err)
	}
	PSK := flag.String("psk", res, "PSK Getestet 36 Zeichen ohne angaben wird ein PSK generiert")

	flag.Parse()

	//Ausgabe Seite A
	fmt.Printf("Name der VPN Seite A:\t%s\n", *nameSeiteAPtr)
	fmt.Printf("\t  IP Seite A:\t%s\n", *ipSeiteAPtr)
	fmt.Printf("   Netzmaske Seite A:\t%s\n", *maskSeiteAPtr)
	fmt.Printf("      DYNDNS Seite A:\t%s\n", *DYNDNSA)
	fmt.Println("=========")

	//Ausgabe Seite B
	fmt.Printf("Name der VPN Seite B:\t%s\n", *nameSeiteBPtr)
	fmt.Printf("\t  IP Seite B:\t%s\n", *ipSeiteBPtr)
	fmt.Printf("   Netzmaske Seite B:\t%s\n", *maskSeiteBPtr)
	fmt.Printf("      DYNDNS Seite B:\t%s\n", *DYNDNSB)
	fmt.Println("=========")
	fmt.Printf("PSK:\t\t\t\"%s\"\n", *PSK)
	fmt.Println("=========")

	fmt.Println("ipsec.secret")
	fmt.Printf("@%s @%s : \"%s\"\n", *DYNDNSB, *DYNDNSB, *PSK)
	fmt.Println("=========")
	lipsecA := "conn \"" + *nameSeiteAPtr + "=>" + *nameSeiteBPtr + "\"\n" +
		"\taggressive = yes\n" +
		"\tfragmentation = yes\n" +
		"\tkeyexchange = ikev1\n" +
		"\tmobike = yes\n" +
		"\treauth = yes\n" +
		"\trekey = yes\n" +
		"\tforceencaps = no\n" +
		"\tinstallpolicy = yes\n" +
		"\ttype = tunnel\n" +
		"\tdpdaction = restart\n" +
		"\tdpddelay = 10s\n" +
		"\tdpdtimeout = 60s\n" +
		"\n" +
		"\tleft = " + *DYNDNSA + "\n" +
		"\t#left = %any\n" +
		"\tright = " + *DYNDNSB + "\n" +
		"\n" +
		"\tleftid = " + *DYNDNSA + "\n" +
		"\tikelifetime = 3600s\n" +
		"\tlifetime = 3600s\n" +
		"\tike = aes256-sha512-modp1024!\n" +
		"\tleftauth = psk\n" +
		"\trightauth = psk\n" +
		"\trightid = " + *DYNDNSB + "\n" +
		"\treqid = 1\n" +
		"\trightsubnet = " + *ipSeiteBPtr + "/24\n" +
		"\tleftsubnet = " + *ipSeiteAPtr + "/24\n" +
		"\tesp = aes256-sha512-modp1024!\n" +
		"\tauto = route\n"

	fmt.Printf("%s", lipsecA)
	fmt.Println("=========")
	lipsecB := "conn \"" + *nameSeiteBPtr + "=>" + *nameSeiteAPtr + "\"\n" +
		"\taggressive = yes\n" +
		"\tfragmentation = yes\n" +
		"\tkeyexchange = ikev1\n" +
		"\tmobike = yes\n" +
		"\treauth = yes\n" +
		"\trekey = yes\n" +
		"\tforceencaps = no\n" +
		"\tinstallpolicy = yes\n" +
		"\ttype = tunnel\n" +
		"\tdpdaction = restart\n" +
		"\tdpddelay = 10s\n" +
		"\tdpdtimeout = 60s\n" +
		"\n" +
		"\tleft = " + *DYNDNSB + "\n" +
		"\t#left = %any\n" +
		"\tright = " + *DYNDNSA + "\n" +
		"\n" +
		"\tleftid = " + *DYNDNSB + "\n" +
		"\tikelifetime = 3600s\n" +
		"\tlifetime = 3600s\n" +
		"\tike = aes256-sha512-modp1024!\n" +
		"\tleftauth = psk\n" +
		"\trightauth = psk\n" +
		"\trightid = " + *DYNDNSA + "\n" +
		"\treqid = 1\n" +
		"\trightsubnet = " + *ipSeiteAPtr + "/24\n" +
		"\tleftsubnet = " + *ipSeiteBPtr + "/24\n" +
		"\tesp = aes256-sha512-modp1024!\n" +
		"\tauto = route\n"

	fmt.Printf("%s", lipsecB)
	fmt.Println("=========")
	fipsecA := "vpncfg {\n" +
		"\tconnections {\n" +
		"\t\tenabled = yes;\n" +
		"\t\teditable = yes;\n" +
		"\t\tconn_type = conntype_lan;\n" +
		"\t\tname = " + "\"" + *nameSeiteBPtr + "=>" + *nameSeiteAPtr + "\"\n" +
		"\t\tboxuser_id = 0;\n" +
		"\t\talways_renew = yes;\n" +
		"\t\treject_not_encrypted = no;\n" +
		"\t\tdont_filter_netbios = yes;\n" +
		"\t\tlocalip = 0.0.0.0;\n" +
		"\t\tlocal_virtualip = 0.0.0.0;\n" +
		"\t\tremotehostname = " + *DYNDNSA + ";\n" +
		"\t\tremote_virtualip = 0.0.0.0;\n" +
		"\t\tlocalid {\n" +
		"\t\t\tfqdn = " + *DYNDNSB + ";\n" +
		"\t\t}\n" +
		"\t\tremoteid {\n" +
		"\t\t\tfqdn = " + *DYNDNSA + ";\n" +
		"\t\t}\n" +
		"\t\tmode = phase1_mode_idp;\n" +
		"\t\tphase1ss = \"all/all/all\";\n" +
		"\t\tkeytype = connkeytype_pre_shared;\n" +
		"\t\tkey = \"" + *PSK + "\";\n" +
		"\t\tcert_do_server_auth = no;\n" +
		"\t\tuse_nat_t = yes;\n" +
		"\t\tuse_xauth = no;\n" +
		"\t\tuse_cfgmode = no;\n" +
		"\t\tphase2localid {\n" +
		"\t\t\tipnet {\n" +
		"\t\t\t\tipaddr = " + *ipSeiteBPtr + ";\n" +
		"\t\t\t\tmask = " + *maskSeiteBPtr + ";\n" +
		"\t\t\t}\n" +
		"\t\t}\n" +
		"\t\tphase2remoteid {\n" +
		"\t\t\tipnet {\n" +
		"\t\t\t\tipaddr = " + *ipSeiteAPtr + ";\n" +
		"\t\t\t\tmask = " + *maskSeiteAPtr + ";\n" +
		"\t\t\t}\n" +
		"\t\t}\n" +
		"\t\tphase2ss = \"esp-all-all/ah-none/comp-all/pfs\";\n" +
		"\t\taccesslist = \"permit ip any " + *ipSeiteAPtr + " " + *maskSeiteAPtr + "\";\n" +
		"\t}\n" +
		"\tike_forward_rules = \"udp 0.0.0.0:500 0.0.0.0:500\",\n" +
		"\t\t\t\t\"udp 0.0.0.0:4500 0.0.0.0:4500\";\n" +
		"}\n"
	fmt.Printf("%s", fipsecA)
	fmt.Println("=========")
}
