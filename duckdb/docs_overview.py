import pandas as pd
from faker import Faker

import duckdb
from duckdb.typing import *


def generate_random_name():
    fake = Faker()
    return fake.name()


# udf support
duckdb.create_function("random_name", generate_random_name, [], VARCHAR)
# duckdb.remove_function("random_name")

if __name__ == "__main__":
    any_relation_name = duckdb.sql(
        "select 42 as i"
    )  # this queries an in-memory database
    print(type(any_relation_name))  # this should be a duckdb relation
    duckdb.sql("select i*2 AS k from any_relation_name").show()

    # reading data from pandas dataframes (polars and arrow also supported)
    pandas_df = pd.DataFrame({"my_col": [42]})
    duckdb.sql("SELECT * FROM pandas_df")

    # writing data to disk
    duckdb.sql("COPY (SELECT 42) TO 'out.parquet'")

    # handling connections to db file
    with duckdb.connect("file.db") as con:
        con.sql("CREATE OR REPLACE TABLE test_table(i integer)")
        con.sql("INSERT INTO test_table VALUES (42)")
        # prepared statements
        con.execute("INSERT INTO test_table VALUES(?)", [49])
        con.executemany("INSERT INTO test_table VALUES(?)", [[43], [44], [45]])
        con.execute("SELECT * FROM test_table WHERE i > ? ORDER BY i ASC", [44])
        print(con.fetchall())

        con.table("test_table").show()

    conn = duckdb.connect("file.db")
    conn.table("test_table").show()
    conn.close()

    duckdb.read_csv("sample.csv", header=True, sep=",")
    # also supports parquet, json, dataframes, arrow objects, numpy objects
    duckdb.sql("select * from 'sample.csv'").show()

    # function API
    res = duckdb.sql("SELECT random_name()").fetchall()
    print(res)
