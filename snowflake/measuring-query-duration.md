Here is a simple way to measure query duration for a Snoflake warehouse:

```sql
select
    count(*) as total_queries
    , avg(total_elapsed_time)
    , median(total_elapsed_time)
    , percentile_cont(0.90) within group (order by total_elapsed_time)
    , percentile_cont(0.95) within group (order by total_elapsed_time)
    , percentile_cont(0.99) within group (order by total_elapsed_time)
from table(information_schema.query_history_by_warehouse(
    'YOUR_WAREHOUSE',
    result_limit => 1000,
    END_TIME_RANGE_START=>to_timestamp_ltz('2023-11-30 15:00:00.000 -0300')
))
where 1=1
and query_text like '%insert_some_pattern_here%'
and database_name='YOUR_DB';
```
