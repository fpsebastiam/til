The recommended way to transfer large files into Snowflake tables is the `COPY INTO` command.
Here are the steps for a large csv load:
1. Compress your local csv file.
1. Ensure you have a `FILE FORMAT` for compressed csvs and a Snowflake `STAGE`.
1. Use the `PUT` command to send this file to a Snowflake STAGE.
1. Use the `COPY INTO` command to transfer the STAGE into your destination table.

A series of dummy commands to execute the steps above
```sql
CREATE OR REPLACE FILE FORMAT path_to_your_new_format TYPE='CSV' ... ;
CREATE OR REPLACE STAGE path_to_your_stage FILE_FORMAT=path_to_your_new_format;

PUT 'your_local_file_path'/your_file_name_and_extension@your_stage_name OVERWRITE=TRUE;

COPY INTO dst_table (columns,...)
FROM @path_to_your_stage/your_file_name_and_extension
FILE_FORMAT=path_to_your_new_format
(ADDITIONAL OPTIONS, e.g purge the staged file)
```

Relevant Docs:
- [COPY INTO](https://docs.snowflake.com/en/sql-reference/sql/copy-into-table)
- [FILE FORMATS](https://docs.snowflake.com/en/sql-reference/sql/create-file-format?utm_source=snowscope&utm_medium=serp&utm_term=file+format)
- [STAGES](https://docs.snowflake.com/en/sql-reference/sql/create-stage?utm_source=snowscope&utm_medium=serp&utm_term=stage)
