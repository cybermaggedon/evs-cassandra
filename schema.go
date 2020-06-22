package main

func (l *Loader) InitSchema() error {

	qry := l.session.Query(`
            CREATE KEYSPACE IF NOT EXISTS cyberprobe WITH REPLICATION = {
                'class': 'SimpleStrategy', 'replication_factor': '1'
            }
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	qry = l.session.Query(`
            CREATE TYPE IF NOT EXISTS cyberprobe.dns_query ( 
                name text, type text, cls text 
            )
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	qry = l.session.Query(`
           CREATE TYPE IF NOT EXISTS cyberprobe.dns_answer (
                name text, type text, cls text, address inet
            )
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	qry = l.session.Query(`
            CREATE TYPE IF NOT EXISTS cyberprobe.http_request ( method text )
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	qry = l.session.Query(`
            CREATE TYPE IF NOT EXISTS cyberprobe.http_response (
                status text, code int 
            )
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	qry = l.session.Query(`
            CREATE TYPE IF NOT EXISTS cyberprobe.indicator (
                id text, type text, value text, category text,
                source text, author text, description text, probability float
            )
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	qry = l.session.Query(`
            CREATE TABLE IF NOT EXISTS cyberprobe.event (
                id text, time timestamp, action text, device text,
                network text,
                srcip inet, destip inet, srcport int, destport int,
                protocol text, dns_type text, 
                url text,
                dns_query list<frozen<dns_query>>,
                dns_answer list<frozen<dns_answer>>, 
                http_request frozen<http_request>,
                http_response frozen<http_response>,
                header map<text, text>,
                indicator list<frozen<indicator>>,
                primary key(device, action, time, id)
            )
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	qry = l.session.Query(`
            CREATE INDEX IF NOT EXISTS event_srcip ON
                cyberprobe.event (srcip)
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	qry = l.session.Query(`
            CREATE INDEX IF NOT EXISTS event_destip ON
                cyberprobe.event (destip)
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	qry = l.session.Query(`
            CREATE INDEX IF NOT EXISTS event_srcport ON
                cyberprobe.event (srcport)
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	qry = l.session.Query(`
            CREATE INDEX IF NOT EXISTS event_destport ON
            cyberprobe.event (destport)
        `)
	if err := qry.Exec(); err != nil {
		return err
	}

	l.insert = l.session.Query(`
            INSERT INTO cyberprobe.event (
                id, time, action, device, network,
                srcip, destip, srcport, destport, protocol, 
                dns_type, dns_query, dns_answer, http_request,
                http_response, url, header, indicator)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
       `)

	return nil

}
