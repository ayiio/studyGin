CREATE TABLE `user` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(20) DEFAULT '',
  `age` INT(11) DEFAULT '0',
  PRIMARY KEY(`id`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;


insert into user (name, age) values('test1', 10);
insert into user (name, age) values('test2', 20);
