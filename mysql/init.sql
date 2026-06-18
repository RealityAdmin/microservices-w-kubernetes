CREATE TABLE IF NOT EXISTS Users (
    id int AUTO_INCREMENT PRIMARY KEY, 
    email VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS Products (
    id int AUTO_INCREMENT PRIMARY KEY, 
    name VARCHAR(255), 
    price FLOAT
);

CREATE TABLE IF NOT EXISTS Orders (
    id int AUTO_INCREMENT PRIMARY KEY, 
    userId int, productId int, 
    FOREIGN KEY (userId) REFERENCES Users(id),
    FOREIGN KEY (productId) REFERENCES Products(id)
);