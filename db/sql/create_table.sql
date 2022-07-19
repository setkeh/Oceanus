CREATE TABLE Images (
	id VARCHAR(36) PRIMARY KEY,
	imagePath VARCHAR(255),
	imageName VARCHAR(255)
);

ALTER TABLE Images ADD imageName VARCHAR(255) AFTER imagePath;

SELECT * FROM Images;
