import mysql.connector
import os
from dotenv import load_dotenv
import csv
from mysql.connector import Error

# TODO: this code has problem that when we have
# exception between lines 10 to 35 it raises below exception
# 'UnboundLocalError: local variable 'connection' referenced before assignment'
# I dont find any solution yet.
def createSchema():
    try:
        sqlCmds = None
        with open("schema.sql", 'r') as file:
            sqlFile = file.read()
            sqlCmds = sqlFile.split(';')

        connection = mysql.connector.connect(host=os.getenv('DB_IP'),
                                             port=os.getenv('DB_PORT'),
                                         database=os.getenv('DB_NAME'),
                                         user=os.getenv('DB_USER'),
                                         password=os.getenv('DB_PASSWORD'))
        if connection.is_connected():
            cursor = connection.cursor()
            for cmd in sqlCmds:
                cursor.execute(cmd)

    except Error as e:
        print("Error while connecting to MySQL", e)
    finally:
        if connection.is_connected():
            cursor.close()
            connection.close()


def main():
    createSchema()
    
if __name__ == "__main__":
    main()
    print("Schema added successfully.")