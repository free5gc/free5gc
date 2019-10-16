.PHONY: install uninstall clean

HOMEASN1DIR = $(HOME)/asn1c

ASN1FILE = ../support/15.2.0/38413.asn

COMPILELANG = -lang=Go
FLAGS = -pdu=all -fcompound-names -findirect-choice -fno-include-deps $(COMPILELANG) 

all:
		$(HOMEASN1DIR)/asn1c/asn1c $(FLAGS) $(ASN1FILE)
install:
		cp -r ../../../asn1c-s1ap $(HOMEASN1DIR)
		cd $(HOMEASN1DIR) && test -f configure || autoreconf -iv && ./configure && make
clean:
		rm -f *.go

uninstall: 
		rm -rf $(HOMEASN1DIR)
