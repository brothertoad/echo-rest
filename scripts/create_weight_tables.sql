drop table if exists weightDaily;
drop table if exists weightSum;

create table weightDaily (
  date integer not null unique,
  weight integer not null
);

create table weightSum (
  month integer,
  year integer,
  count integer,
  total integer,
  avg integer,
  unique (month, year)
);
