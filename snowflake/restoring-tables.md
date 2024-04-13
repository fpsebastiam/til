Sometimes Snowflake time travel feature is not enough to restore a table to a previous state.


```sql
-- restoring Snowflake tables that have either been dropped or replaced

create table db.schema.test_table as
select 'v1' as version, 1 as id;


-- replace the table with a new version
create or replace table db.schema.test_table as
select 'v2' as version, 2 as id;


-- replace the table with a new version
create or replace table db.schema.test_table as
select 'v3' as version, 3 as id;


-- drop the tablje
drop table db.schema.test_table;

-- restore to v3
undrop table db.schema.test_table;
select * from db.schema.test_table; -- >> 'v3', 3

-- restore to v2
alter table db.schema.test_table rename to test_table_v3;
undrop table db.schema.test_table;
select * from db.schema.test_table; -- >> 'v2', 2

-- restore to v1
alter table db.schema.test_table rename to test_table_v2;
undrop table db.schema.test_table;
select * from db.schema.test_table; -- >> 'v1', 1
```