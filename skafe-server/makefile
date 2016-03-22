DESTDIR=/usr/local
CREATEMODE=0644

all: skafe-server

skafe-server: link.go main.go rbScriptWorker.go ruleEngine.go scriptEngine.go setup.go syncwriter.go queue.go
	go build

install: skafe-server
	# Install the server bin
	install -D skafe-server $(DESTDIR)/usr/bin/skafe-server
	# Install the server conf
	install -m $(CREATEMODE) -D config/skafe-server.conf $(DESTDIR)/etc/skafe/skafe-server.conf
	# Create the tls directories
	install -d $(DESTDIR)/etc/skafe/tls
	install -m $(CREATEMODE) -d $(DESTDIR)/etc/skafe/tls/private
	# Create the rules directory
	install -m $(CREATEMODE) -d $(DESTDIR)/etc/skafe/rules
	# Create the lib directory
	install -d $(DESTDIR)/usr/lib/skafe
	# Install library files
	install -m $(CREATEMODE) scripts/ruby.rb $(DESTDIR)/usr/lib/skafe/ruby.rb
	# Create the scripts directory
	install -d $(DESTDIR)/etc/skafe/scripts
