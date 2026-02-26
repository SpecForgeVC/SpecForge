-- Add password_hash column to users table
ALTER TABLE users ADD COLUMN password_hash TEXT;

-- Update superadmin user with Argon2id hashed password "password123"
UPDATE users 
SET password_hash = '$argon2id$v=19$m=65536,t=3,p=4$tUwPJBnFIGLXMekKuVT/jw$Hd+XzvIemnOi0Stx0ewn1JtitdCPfXj0ghYjdoDdxnI8'
WHERE email = 'admin@specforge.io';
