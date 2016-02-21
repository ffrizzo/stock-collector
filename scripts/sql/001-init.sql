create database stockCollector;
create table stock (id serial not null, ticker varchar(255), account varchar(255),
                    username varchar(255), sell numeric(14,2), rate numeric(14,2),
                    buy numeric(14,2), time timestamp, primary key(id));
