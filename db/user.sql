CREATE DATABASE IF NOT EXISTS db_wallet DEFAULT charset utf8 COLLATE utf8_general_ci;

DROP USER IF EXISTS wallet_user_cold;
CREATE USER 'wallet_user_cold'@'%' identified by 'cm1axC6K31j34m5';
GRANT SELECT,UPDATE,INSERT,DELETE ON db_wallet.* to 'wallet_user_cold'@'%' with grant option;
FLUSH PRIVILEGES;



--172.16.5.173
--  root
--  zH@NgZh1n2@e90517