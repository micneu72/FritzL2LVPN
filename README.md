# FritzL2LVPN
Erstellt eine VPN Konfiguration zwischen AVM Fritzbox und Linux/Unix (stongSwan)
die PSK wird generiert mit einer Länge von 36 Zeichen.
Getestet mit einer Fritzbox 7490 FRITZ!OS: 07.21

Aufruf:
```
 micneu@mne-mbp  /tmp  ./FritzL2LVPN -h                                                        ✔  3697  19:50:26
Usage of ./FritzL2LVPN:
  -P	VPN Konfig anzeigen
  -W	VPN Konfig schreiben
  -adns string
    	Dyndns Seite A (default "a.myfritz.net")
  -aip string
    	IP Seite A (default "192.168.178.0")
  -aname string
    	Name der Seite A (default "Seite A")
  -bdns string
    	Dyndns Seite B (default "b.myfritz.net")
  -bip string
    	IP Seite A (default "192.168.180.0")
  -bname string
    	Name der Seite B (default "Seite B")



