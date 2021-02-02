
drop database if exists imgupdb;

create database if not exists imgupdb;

use imgupdb;

CREATE TABLE IF NOT EXISTS img (
  id int(20) AUTO_INCREMENT,
  filename VARCHAR(45) NULL,
  filepath VARCHAR(300) NULL,
  imgdata LONGBLOB,
  CONSTRAINT PRIMARY KEY (id)
);