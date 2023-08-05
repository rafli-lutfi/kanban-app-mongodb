# Kanban App MongoDB

simple kanban app using golang programming language and mongoDB database.

steps to run this application:

1. Clone this repository into your local computer

2. Copy a file called **.env.example** then paste it in the same directory and rename it with **.env**

3. Fill in all variables

    ```md
    # PORT
    PORT = 

    #MONGODB URI
    MONGO_URI = 

    #MONGODB COLLECTION
    DATABASE_NAME = 

    # JWT KEY
    JWT_SECRET_KEY = 
    ```

4. Run `go build` in your terminal and a file with .exe extension will appear

5. Run the `.exe` file

6. Open a browser and go to <http://localhost:port/> **(match *port* with the one in the env file)**
