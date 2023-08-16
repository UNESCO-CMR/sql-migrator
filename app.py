import os
import pymysql

# MySQL connection details
host = 'localhost'
username = 'root'
password = ''
database = 'myschoolonline'

# Path to the directory containing SQL files
sql_directory = '/Users/matt-klaus/Documents/UNESCO/Programming/sql-migrate-go/db_split/SQLDumpSplitterResult'

# Name of the DB schema file
schema_file = 'dbstructure.sql'

# Connect to MySQL
connection = pymysql.connect(host=host, user=username, password=password, database=database)

try:
    # Create a cursor object
    cursor = connection.cursor()

    # Read and execute the DB schema file first
    schema_file_path = os.path.join(sql_directory, schema_file)
    with open(schema_file_path, 'r') as schema_file:
        schema_sql = schema_file.read()
        cursor.execute(schema_sql)
    print(f'{schema_file_path} executed successfully.')

    # Get a list of SQL files in the directory
    sql_files = [file for file in os.listdir(sql_directory) if file.endswith('.sql') and file != schema_file]

    # Execute the remaining SQL files
    for sql_file in sql_files:
        file_path = os.path.join(sql_directory, sql_file)
        with open(file_path, 'r') as sql_file:
            sql = sql_file.read()
            cursor.execute(sql)
        print(f'{file_path} executed successfully.')

    # Commit the changes
    connection.commit()
    print('All SQL files executed successfully.')

except pymysql.Error as e:
    # Handle any errors that occur during execution
    print(f'Error executing SQL files: {e}')

finally:
    # Close the database connection
    connection.close()