--first you need to create database  diamond default character set utf8(create database diamond);
--use diamond
--excute follow script to create table config_info

drop table if exists config_info cascade;

create table config_info (
         id bigint(64) unsigned not null auto_increment,
         data_id varchar(255) not null default '',
         group_id varchar(128) not null default '',
         content longtext not null,
         md5 varchar(32) not null default '',
         gmt_create datetime not null default '2020-12-30 17:56:00',
         gmt_modified datetime not null default '2020-12-30 17:56:00',
         primary key  (id),
         unique key uk_config_datagroup (data_id,group_id)
);