BEGIN;

DO $$
DECLARE
    fkey_record RECORD;
    drop_statement TEXT;
    create_statement TEXT;
BEGIN
    FOR fkey_record IN
        SELECT
            con.conname AS constraint_name,
            conrelid::regclass AS referencing_table,
            confrelid::regclass AS referenced_table,
            (SELECT string_agg(att.attname, ', ') FROM unnest(con.conkey) AS col(attnum) JOIN pg_attribute AS att ON att.attrelid = con.conrelid AND att.attnum = col.attnum) AS referencing_columns,
            (SELECT string_agg(att.attname, ', ') FROM unnest(con.confkey) AS col(attnum) JOIN pg_attribute AS att ON att.attrelid = con.confrelid AND att.attnum = col.attnum) AS referenced_columns
        FROM
            pg_constraint AS con
        JOIN
            pg_class AS conf ON con.confrelid = conf.oid
        WHERE
            con.contype = 'f' AND conf.relname = 'users'
            AND con.confupdtype = 'c'
    LOOP
        drop_statement := 'ALTER TABLE ' || quote_ident(fkey_record.referencing_table::text) || ' DROP CONSTRAINT ' || quote_ident(fkey_record.constraint_name) || ';';
        RAISE NOTICE 'Executing: %', drop_statement;
        EXECUTE drop_statement;

        create_statement := 'ALTER TABLE ' || quote_ident(fkey_record.referencing_table::text) || ' ADD CONSTRAINT ' || quote_ident(fkey_record.constraint_name) ||
                            ' FOREIGN KEY (' || fkey_record.referencing_columns || ')' ||
                            ' REFERENCES ' || quote_ident(fkey_record.referenced_table::text) || ' (' || fkey_record.referenced_columns || ');';
        RAISE NOTICE 'Executing: %', create_statement;
        EXECUTE create_statement;
    END LOOP;
END;
$$;

COMMIT;
