-- test database
create database test;
create user test with encrypted password 'test123';
grant all privileges on database test to test;

-- development database
create database development;
create user development with encrypted password 'development123';
grant all privileges on database development to development;
revoke connect on database development from public;

-- production database
create database production;
create user production with encrypted password 'production123';
grant all privileges on database production to production;
revoke connect on database production from public;
