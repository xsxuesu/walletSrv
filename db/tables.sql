CREATE TABLE IF NOT EXISTS hd_count(
	`id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `type` varchar(50) NOT NULL COMMENT 'coin type',
	hdcount varchar(100) NOT NULL COMMENT 'hd path count',
	updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'hd_count create and update time',
	PRIMARY KEY (`id`)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='hd wallet address subpath record';

CREATE TABLE IF NOT EXISTS btc_addr (
	address varchar(100) NOT NULL COMMENT 'bitcoin address',
	private varchar(200) NOT NULL COMMENT 'bitcoin encrypted private key',
	PRIMARY KEY  (address)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='bitcoin address and private key records';

	CREATE TABLE IF NOT EXISTS eth_addr (
	address varchar(100) NOT NULL COMMENT 'ethereum address',
	private varchar(200) NOT NULL COMMENT 'ethereum encrypted private key',
	PRIMARY KEY  (address)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='ethereum address and private key records';


	CREATE TABLE IF NOT EXISTS usdt_addr (
	address varchar(100) NOT NULL COMMENT 'usdt address',
	private varchar(200) NOT NULL COMMENT 'usdt encrypted private key',
	PRIMARY KEY  (address)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='usdt address and private key records';


	CREATE TABLE IF NOT EXISTS btc_hdaddr (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  address varchar(100) NOT NULL COMMENT 'btc hd address',
	private varchar(200) NOT NULL COMMENT 'btc encrypted private key',
	PRIMARY KEY  (address)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='btc address and private key records';

	CREATE TABLE IF NOT EXISTS eth_hdaddr (
	  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    address varchar(100) NOT NULL COMMENT 'ethereum hd address',
    private varchar(200) NOT NULL COMMENT 'ethereum encrypted private key',
    created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',
    PRIMARY KEY  (`id`),
    UNIQUE KEY  ix_eth_hdaddr(`address`)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='ethereum hd address and private key records';


	CREATE TABLE IF NOT EXISTS usdt_hdaddr (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    address varchar(100) NOT NULL COMMENT 'usdt hd address',
    private varchar(200) NOT NULL COMMENT 'usdt encrypted private key',
    PRIMARY KEY  (`id`),
    UNIQUE KEY  ix_usdt_hdaddr(`address`)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='usdt hd address and private key records';

CREATE TABLE IF NOT EXISTS serial_btc (
	`id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  serial varchar(100) NOT NULL COMMENT 'serial no',
	type varchar(10) NOT NULL COMMENT 'coin  type',
	f varchar(100) NOT NULL COMMENT 'transfer from address',
	t varchar(100) NOT NULL COMMENT 'transfer to address',
	value varchar(300) NOT NULL COMMENT 'transfer value',
	fee double NOT NULL COMMENT 'tansfer fee',
	time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'order create and update time',
	txid varchar(200) NOT NULL COMMENT 'transaction hash',
	status varchar(20) NOT NULL COMMENT 'transaction status',
	PRIMARY KEY  (serial),
	UNIQUE KEY  ix_serial_btc(`serial`)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='bitcoin transaction records';


CREATE TABLE IF NOT EXISTS serial_eth (
	serial varchar(100) NOT NULL COMMENT 'serial no',
	type varchar(10) NOT NULL COMMENT 'coin  type',
	f varchar(100) NOT NULL COMMENT 'transfer from address',
	t varchar(100) NOT NULL COMMENT 'transfer to address',
	value varchar(300) NOT NULL COMMENT 'transfer value',
	fee double NOT NULL COMMENT 'tansfer fee',
	time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'order create and update time',
	txid varchar(200) NOT NULL COMMENT 'transaction hash',
	status varchar(20) NOT NULL COMMENT 'transaction status',
	PRIMARY KEY  (serial)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='ethereum transaction records';


CREATE TABLE IF NOT EXISTS serial_usdt (
	serial varchar(100) NOT NULL COMMENT 'serial no',
	type varchar(10) NOT NULL COMMENT 'coin  type',
	f varchar(100) NOT NULL COMMENT 'transfer from address',
	t varchar(100) NOT NULL COMMENT 'transfer to address',
	value varchar(300) NOT NULL COMMENT 'transfer value',
	fee double NOT NULL COMMENT 'tansfer fee',
	time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'order create and update time',
	txid varchar(200) NOT NULL COMMENT 'transaction hash',
	status varchar(20) NOT NULL COMMENT 'transaction status',
	PRIMARY KEY  (serial)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='usdt transaction records';


CREATE TABLE IF NOT EXISTS decrypt_btc (
	serial varchar(100) NOT NULL COMMENT 'serial no',
	type varchar(10) NOT NULL COMMENT 'coin type',
	f varchar(100) NOT NULL COMMENT 'transfer from address',
	t varchar(100) NOT NULL COMMENT 'transfer to address',
	hash varchar(10) NOT NULL COMMENT 'hash function',
	time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'order create and update time',
	PRIMARY KEY  (serial)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='decrypt private key records';

CREATE TABLE IF NOT EXISTS decrypt_eth (
	serial varchar(100) NOT NULL COMMENT 'serial no',
	type varchar(10) NOT NULL COMMENT 'coin type',
	f varchar(100) NOT NULL COMMENT 'transfer from address',
	t varchar(100) NOT NULL COMMENT 'transfer to address',
	hash varchar(10) NOT NULL COMMENT 'hash function',
	time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'order create and update time',
	PRIMARY KEY  (serial)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='decrypt private key records';

CREATE TABLE IF NOT EXISTS decrypt_usdt (
	serial varchar(100) NOT NULL COMMENT 'serial no',
	type varchar(10) NOT NULL COMMENT 'coin type',
	f varchar(100) NOT NULL COMMENT 'transfer from address',
	t varchar(100) NOT NULL COMMENT 'transfer to address',
	hash varchar(10) NOT NULL COMMENT 'hash function',
	time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'order create and update time',
	PRIMARY KEY  (serial)
	)ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='decrypt private key records';