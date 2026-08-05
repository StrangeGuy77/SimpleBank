DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS accounts; -- Last one to be dropped because entries and transfers references the primary key of accounts table.