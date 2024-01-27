INSERT INTO users (
    id, 
    name, 
    email, 
    password
) VALUES (UUID(), 'administrator', 'admin@my-wallet.com','$2y$10$pmU8V7pCRnrrVJMQG8fqYuDH92V2pE7HQo.b3BKhraA9s3CRur8Jy'), 
        (UUID(), 'example_user', 'user@my-wallet.com','$2y$10$pmU8V7pCRnrrVJMQG8fqYuDH92V2pE7HQo.b3BKhraA9s3CRur8Jy');
