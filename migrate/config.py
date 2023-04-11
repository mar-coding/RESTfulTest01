import mysql.connector
import os
from dotenv import load_dotenv
import csv

def test_connection():
    load_dotenv()
    connection = None
    try:
        connection = mysql.connector.connect(host='127.0.0.1',
                                             port='3306',
                                         database='myDB',
                                         user='root',
                                         password='132egfM*m134A13tgqk@', use_pure=True, connection_timeout=180)

        if connection.is_connected():
            db_Info = connection.get_server_info()
            print("Connected to MySQL database... MySQL Server version on ", db_Info)
    except Error as e:
        print("Error while connecting to MySQL", e)
    finally:
        if connection.is_connected():
            connection.close()
            print("MySQL connection is closed")

def add_schema():
    sqlCmds = None
    with open("schema.sql", 'r') as file:
        sqlFile = file.read()
        sqlCmds = sqlFile.split(';')

    try:
        connection = mysql.connector.connect(host='127.0.0.1',
                                             port='3306',
                                         database='myDB',
                                         user='root',
                                         password='132egfM*m134A13tgqk@', use_pure=True, connection_timeout=180)

        if connection.is_connected():
            for cmd in sqlCmds:
                connection._execute_query(cmd)
    except Exception as e:
        print("Error while running SQL", e)
    finally:
        if connection.is_connected():
            connection.close()
def main():
    add_schema()
    
if __name__ == "__main__":
    main()
    print("Schema added successfully.")