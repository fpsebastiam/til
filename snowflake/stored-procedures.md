CREATE OR REPLACE PROCEDURE mock_stored_procedure (
    user_id VARCHAR
)
  RETURNS table (some_feature NUMBER)
  LANGUAGE SQL
  AS
    $$
    DECLARE
        result RESULTSET;
        -- variables can be defined here
        transient_table VARCHAR DEFAULT 'my_db.my_schema.aggregated_feature_' || user_id;
    BEGIN
        -- or here, using let. Also note how '?' is later bound to a variable
        let aggregation_query := 'SELECT aggregated_feature FROM identifier(?) LIMIT 1;';
        BEGIN
            result := (EXECUTE IMMEDIATE :aggregation_query USING (transient_table));
            RETURN TABLE(result);
        EXCEPTION
            -- should be triggered on the first time one runs an aggregation, and
            -- the transient table is not there to be queried
            WHEN OTHER THEN
                -- in "EXECUTE IMMEDIATE" we have to use '?' to bind variables, but here
                -- we can use the :my_variable notation, which is clearer and less prone to bugs
                CREATE OR REPLACE TRANSIENT TABLE identifier(:transient_table) AS SELECT to_number(1230) AS aggregated_feature;
                result := (EXECUTE IMMEDIATE :aggregation_query USING(transient_table));
                RETURN TABLE(result);
        END;
    END;
    $$;
